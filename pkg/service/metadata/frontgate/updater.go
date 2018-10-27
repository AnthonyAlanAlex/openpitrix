package frontgate

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"google.golang.org/grpc"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/metadata/drone"
	"openpitrix.io/openpitrix/pkg/pb/metadata/pilot"
	"openpitrix.io/openpitrix/pkg/pb/metadata/types"
	"openpitrix.io/openpitrix/pkg/service/metadata/drone/droneutil"
	"openpitrix.io/openpitrix/pkg/util/gziputil"
	"openpitrix.io/openpitrix/pkg/util/httputil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/retryutil"
	"openpitrix.io/openpitrix/pkg/version"
)

var (
	Version                     = version.ShortVersion
	CheckInterval               = 10 * time.Second
	RetryInterval               = 3 * time.Second
	RetryCount                  = 5
	OpenPitrixReleaseUrlPattern = "https://github.com/openpitrix/openpitrix/releases/download/%s/openpitrix-%s-bin.tar.gz"
	HttpServePath               = "/opt/openpitrix/bin"
	DowloadPathPattern          = "/opt/openpitrix/bin/%s"
	DowloadFilePathPattern      = "/opt/openpitrix/bin/%s/%s"
	PilotVersionFilePath        = "/opt/openpitrix/conf/pilot-version"
	KeyPrefix                   = "/"
	KeyRegexp                   = regexp.MustCompile(`^\/\d+\.\d+\.\d+\.\d+\/host\/ip$`)
	EtcdEndpoints               = []string{"127.0.0.1:2379"}
)

type Updater struct {
	conn       *grpc.ClientConn
	connClosed <-chan struct{}

	etcd *EtcdClientManager
	cfg  *pbtypes.FrontgateConfig
}

func NewUpdater(conn *grpc.ClientConn, connClosed <-chan struct{}, cfg *pbtypes.FrontgateConfig) *Updater {
	return &Updater{
		conn:       conn,
		connClosed: connClosed,
		etcd:       NewEtcdClientManager(),
		cfg:        cfg,
	}
}

func (u *Updater) checkPilotVersionDiff() (bool, string) {
	pilotClient := pbpilot.NewPilotServiceForFrontgateClient(u.conn)

	ctx := context.Background()
	input := &pbtypes.Empty{}
	pilotVersion, err := pilotClient.GetPilotVersion(ctx, input)
	if err != nil {
		logger.Warn(ctx, "Get pilot version failed, %+v", err)
		return false, ""
	}

	logger.Debug(ctx, "Get pilot version [%s]", pilotVersion.ShortVersion)
	logger.Debug(ctx, "Get self  version [%s]", Version)

	var short string
	v := strings.SplitN(strings.Trim(pilotVersion.ShortVersion, "\""), "-", 2)
	if len(v) != 0 {
		short = v[0]
	}

	return pilotVersion.ShortVersion != Version, short
}

func (u *Updater) createPilotVersionFile(pilotVersion string) error {
	f, err := os.Create(PilotVersionFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(pilotVersion)
	if err != nil {
		return err
	}

	return nil
}

func (u *Updater) downloadNewRelease(pilotVersion string) error {
	err := os.MkdirAll(fmt.Sprintf(DowloadPathPattern, pilotVersion), os.ModeDir|os.ModePerm)
	if err != nil {
		return err
	}

	url := fmt.Sprintf(OpenPitrixReleaseUrlPattern, pilotVersion, pilotVersion)
	logger.Info(nil, "Trying to download new release from url [%s]", url)

	err = retryutil.Retry(RetryCount, RetryInterval, func() error {
		resp, err := httputil.HttpGet(url)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("download new release from url [%s] failed, status %s", url, resp.Status)
		}

		archiveFiles, err := gziputil.LoadArchive(resp.Body)
		if err != nil {
			return err
		}

		for fileName, fileBytes := range archiveFiles {
			filePath := fmt.Sprintf(DowloadFilePathPattern, pilotVersion, fileName)

			logger.Info(nil, "Write downloaded file [%s] to [%s]", fileName, filePath)
			f, err := os.Create(filePath)
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = f.Write(fileBytes)
			if err != nil {
				return err
			}

			err = os.Chmod(filePath, os.ModePerm)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	err = u.createPilotVersionFile(pilotVersion)
	if err != nil {
		return err
	}

	return nil
}

func (u *Updater) getDroneList() ([]string, error) {
	etcdConfig := u.cfg.GetEtcdConfig()
	etcdClient, err := u.etcd.GetClient(
		EtcdEndpoints,
		time.Duration(etcdConfig.GetTimeoutSeconds())*time.Second,
		int(etcdConfig.GetMaxTxnOps()),
	)
	if err != nil {
		return nil, err
	}

	vs, err := etcdClient.GetValuesByPrefix(KeyPrefix)
	if err != nil {
		return nil, err
	}

	drones := []string{}
	for k, v := range vs {
		if KeyRegexp.MatchString(k) {
			logger.Debug(nil, "Matched key [%s] value [%s]", k, v)
			drones = append(drones, v)
		}
	}

	return drones, nil
}

func (u *Updater) getDroneVersion(ctx context.Context, client pbdrone.DroneServiceClient) (string, error) {
	droneVersion, err := client.GetDroneVersion(ctx, &pbtypes.DroneEndpoint{})
	if err != nil {
		return "", err
	}

	var short string
	v := strings.SplitN(strings.Trim(droneVersion.ShortVersion, "\""), "-", 2)
	if len(v) != 0 {
		short = v[0]
	}

	return short, nil
}

func (u *Updater) distributeDrone(drone string, pilotVersion string) error {
	ctx := context.Background()

	client, conn, err := droneutil.DialDroneService(ctx, drone, constants.DroneServicePort)
	if err != nil {
		return err
	}
	defer conn.Close()

	droneVersion, err := u.getDroneVersion(ctx, client)
	if err != nil {
		return err
	}

	logger.Debug(ctx, "Pilot version [%s]", pilotVersion)
	logger.Debug(ctx, "Drone version [%s]", droneVersion)

	if pilotVersion != droneVersion {
		logger.Info(nil, "Trying to distribute drone with version [%s] from frontgate[%s] to drone[%s]", pilotVersion, u.cfg.Host, drone)
		req := &pbtypes.DistributeDroneRequest{
			PilotVersion:     pbutil.ToProtoString(pilotVersion),
			FrontgateAddress: pbutil.ToProtoString(u.cfg.Host),
		}

		_, err = client.DistributeDrone(ctx, req)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Updater) distributeDrones(pilotVersion string) error {
	drones, err := u.getDroneList()
	if err != nil {
		return err
	}

	logger.Debug(nil, "Get drone list %+v", drones)

	for _, drone := range drones {
		err := u.distributeDrone(drone, pilotVersion)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u *Updater) Serve() {
	ticker := time.NewTicker(CheckInterval)
	defer ticker.Stop()

	for t := range ticker.C {
		logger.Debug(nil, "Tick at [%s]", t)

		select {
		case <-u.connClosed:
			return
		default:
			diff, pilotVersion := u.checkPilotVersionDiff()
			if diff {
				err := u.downloadNewRelease(pilotVersion)
				if err != nil {
					logger.Warn(nil, "Download new release failed, %+v", err)
					continue
				}
				logger.Info(nil, "Frontgate exit")
				os.Exit(0)
			}

			err := u.distributeDrones(pilotVersion)
			if err != nil {
				logger.Warn(nil, "Distribute drone failed, %+v", err)
			}
		}
	}
}
