// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package client

import (
	"context"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func SetSystemUserToContext(ctx context.Context) context.Context {
	return senderutil.ContextWithSender(ctx, senderutil.GetSystemSender())
}

func SetUserToContext(ctx context.Context, user *pb.User) context.Context {
	return senderutil.ContextWithSender(ctx, senderutil.GetSender(user))
}
