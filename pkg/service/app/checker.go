// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package app

import (
	"context"

	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/pb"
)

func (p *Server) Checker(ctx context.Context, req interface{}) error {
	switch r := req.(type) {
	case *pb.CreateAppRequest:
		return manager.NewChecker(ctx, r).
			Required("name", "repo_id").
			Exec()
	case *pb.ModifyAppRequest:
		return manager.NewChecker(ctx, r).
			Required("app_id").
			Exec()
	case *pb.DeleteAppsRequest:
		return manager.NewChecker(ctx, r).
			Required("app_id").
			Exec()
	case *pb.CreateAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Required("app_id", "name", "repo_id").
			Exec()
	case *pb.ModifyAppVersionRequest:
		return manager.NewChecker(ctx, r).
			Required("version_id").
			Exec()
	case *pb.DeleteAppVersionsRequest:
		return manager.NewChecker(ctx, r).
			Required("version_id").
			Exec()
	case *pb.GetAppVersionPackageRequest:
		return manager.NewChecker(ctx, r).
			Required("version_id").
			Exec()
	case *pb.GetAppVersionPackageFilesRequest:
		return manager.NewChecker(ctx, r).
			Required("version_id").
			Exec()
	}
	return nil
}
