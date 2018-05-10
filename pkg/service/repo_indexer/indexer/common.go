package indexer

import (
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/client"
	appclient "openpitrix.io/openpitrix/pkg/client/app"
	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Indexer interface {
	IndexRepo() error
}

func GetIndexer(repo *pb.Repo) Indexer {
	var i Indexer
	providers := repo.GetProviders()
	if stringutil.StringIn(constants.ProviderKubernetes, providers) {
		i = NewHelmIndexer(repo)
	} else {
		i = NewDevkitIndexer(repo)
	}
	return i
}

type indexer struct {
	repo *pb.Repo
}
type appInterface interface {
	GetName() string
	GetDescription() string
	GetIcon() string
	GetHome() string
	GetSources() []string
	GetKeywords() []string
}
type versionInterface interface {
	GetVersion() string
	GetAppVersion() string
	GetDescription() string
	GetUrls() []string
}

func (i *indexer) syncAppInfo(app appInterface) (string, error) {
	chartName := app.GetName()
	repoId := i.repo.GetRepoId().GetValue()
	owner := i.repo.GetOwner().GetValue()

	var appId string
	ctx := client.GetSystemUserContext()
	appManagerClient, err := appclient.NewAppManagerClient(ctx)
	if err != nil {
		return appId, err
	}
	req := pb.DescribeAppsRequest{}
	req.RepoId = []string{repoId}
	req.Owner = []string{owner}
	req.ChartName = []string{chartName}
	res, err := appManagerClient.DescribeApps(ctx, &req)
	if err != nil {
		return appId, err
	}
	description := pbutil.ToProtoString(app.GetDescription())
	icon := pbutil.ToProtoString(app.GetIcon())
	home := pbutil.ToProtoString(app.GetHome())
	sources := pbutil.ToProtoString(strings.Join(app.GetSources(), ","))
	keywords := pbutil.ToProtoString(strings.Join(app.GetKeywords(), ","))
	if res.TotalCount == 0 {
		createReq := pb.CreateAppRequest{}
		createReq.RepoId = pbutil.ToProtoString(repoId)
		createReq.ChartName = pbutil.ToProtoString(chartName)
		createReq.Name = pbutil.ToProtoString(chartName)
		createReq.Description = description
		createReq.Icon = icon
		createReq.Home = home
		createReq.Sources = sources
		createReq.Keywords = keywords

		createRes, err := appManagerClient.CreateApp(ctx, &createReq)
		if err != nil {
			return appId, err
		}
		appId = createRes.GetApp().GetAppId().GetValue()
		return appId, err

	} else {
		modifyReq := pb.ModifyAppRequest{}
		modifyReq.AppId = res.AppSet[0].AppId
		modifyReq.Name = pbutil.ToProtoString(chartName)
		modifyReq.ChartName = pbutil.ToProtoString(chartName)
		modifyReq.Description = description
		modifyReq.Icon = icon
		modifyReq.Home = home
		modifyReq.Sources = sources
		modifyReq.Keywords = keywords

		modifyRes, err := appManagerClient.ModifyApp(ctx, &modifyReq)
		if err != nil {
			return appId, err
		}
		appId = modifyRes.GetApp().GetAppId().GetValue()
		return appId, err
	}
}

func (i *indexer) syncAppVersionInfo(appId string, version versionInterface, index int) (string, error) {
	owner := i.repo.GetOwner().GetValue()

	var versionId string
	ctx := client.GetSystemUserContext()
	appManagerClient, err := appclient.NewAppManagerClient(ctx)
	if err != nil {
		return versionId, err
	}
	appVersionName := version.GetVersion()
	if version.GetAppVersion() != "" {
		appVersionName += fmt.Sprintf(" [%s]", version.GetAppVersion())
	}
	packageName := version.GetUrls()[0]
	description := version.GetDescription()
	req := pb.DescribeAppVersionsRequest{}
	req.AppId = []string{appId}
	req.Owner = []string{owner}
	req.Name = []string{appVersionName}
	res, err := appManagerClient.DescribeAppVersions(ctx, &req)
	if err != nil {
		return versionId, err
	}
	if res.TotalCount == 0 {
		createReq := pb.CreateAppVersionRequest{}
		createReq.AppId = pbutil.ToProtoString(appId)
		createReq.Owner = pbutil.ToProtoString(owner)
		createReq.Name = pbutil.ToProtoString(appVersionName)
		createReq.PackageName = pbutil.ToProtoString(packageName)
		createReq.Description = pbutil.ToProtoString(description)
		createReq.Sequence = pbutil.ToProtoUInt32(uint32(index))

		createRes, err := appManagerClient.CreateAppVersion(ctx, &createReq)
		if err != nil {
			return versionId, err
		}
		versionId = createRes.GetAppVersion().GetVersionId().GetValue()
		return versionId, err
	} else {
		modifyReq := pb.ModifyAppVersionRequest{}
		modifyReq.VersionId = res.AppVersionSet[0].VersionId
		modifyReq.PackageName = pbutil.ToProtoString(packageName)
		modifyReq.Description = pbutil.ToProtoString(description)
		modifyReq.Sequence = pbutil.ToProtoUInt32(uint32(index))

		modifyRes, err := appManagerClient.ModifyAppVersion(ctx, &modifyReq)
		if err != nil {
			return versionId, err
		}
		versionId = modifyRes.GetAppVersion().GetVersionId().GetValue()
		return versionId, err
	}
}
