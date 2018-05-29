// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package frontgate

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/pb/frontgate"
	"openpitrix.io/openpitrix/pkg/pb/types"
	"openpitrix.io/openpitrix/pkg/service/pilot/pilotutil"
)

type Server struct {
	cfg  *ConfigManager
	etcd *EtcdClientManager

	ch   *pilotutil.FrameChannel
	conn *grpc.ClientConn
	err  error
}

func Serve(cfg *ConfigManager) {
	p := &Server{
		cfg:  cfg,
		etcd: NewEtcdClientManager(),
	}

	go ServeReverseRpcServerForPilot(cfg.Get(), p)
	go pbfrontgate.ListenAndServeFrontgateService("tcp",
		fmt.Sprintf(":%d", constants.FrontgateServicePort),
		p,
	)

	<-make(chan bool)
}

func ServeReverseRpcServerForPilot(
	cfg *pbtypes.FrontgateConfig,
	service pbfrontgate.FrontgateService,
) {
	logger.Info("ReverseRpcServerForPilot beign")
	defer logger.Info("ReverseRpcServerForPilot end")

	var lastErrCode = codes.OK

	for {
		ch, conn, err := pilotutil.DialFrontgateChannel(
			context.Background(), fmt.Sprintf("%s:%d", cfg.PilotHost, cfg.PilotPort),
			grpc.WithInsecure(),
		)
		if err != nil {
			gerr, ok := status.FromError(err)
			if !ok {
				logger.Error("err shoule be grpc error type")
				time.Sleep(time.Second)
				continue
			}

			if gerr.Code() != codes.Unavailable || gerr.Code() != lastErrCode {
				logger.Error("did not connect: %v", gerr.Err())
			}

			lastErrCode = gerr.Code()
			continue
		} else {
			if lastErrCode == codes.Unavailable {
				logger.Info("pilot connect ok")
			}

			lastErrCode = codes.OK
		}

		pbfrontgate.ServeFrontgateService(ch, service)
		conn.Close()
	}
}
