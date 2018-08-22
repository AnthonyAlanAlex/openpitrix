// Copyright 2017 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package senderutil

import (
	"context"
	"encoding/json"
	"net/http"

	context2 "golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

const senderKey = "sender"

type Info struct {
	UserId string `json:"user_id"`
}

func GetSystemUser() *Info {
	return &Info{UserId: "system"}
}

func (info *Info) ToJson() string {
	ret, _ := json.Marshal(info)
	return string(ret)
}

func GetSenderFromContext(ctx context.Context) *Info {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	//logger.Debug(nil, "%+v", md[senderKey])
	if len(md[senderKey]) == 0 {
		return nil
	}
	sender := Info{}
	err := json.Unmarshal([]byte(md[senderKey][0]), &sender)
	if err != nil {
		panic(err)
	}
	return &sender
}

func AuthUserInfo(ctx context.Context, authKey string) *Info {
	//logger.Debug(ctx, "got auth key: %+v", authKey)
	// TODO: validate auth key && get user info from db
	return GetSystemUser()
}

func ServeMuxSetSender(ctx context2.Context, request *http.Request) metadata.MD {
	md := metadata.MD{}
	authKey := request.Header.Get("X-Auth-Key")
	user := AuthUserInfo(ctx, authKey)
	md[senderKey] = []string{user.ToJson()}
	return md
}

func NewContext(ctx context.Context, user *Info) context.Context {
	md := metadata.MD{}
	md[senderKey] = []string{user.ToJson()}
	return metadata.NewOutgoingContext(ctx, md)
}
