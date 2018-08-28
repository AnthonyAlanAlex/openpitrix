// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package repoiface

import (
	"context"
	"fmt"
	neturl "net/url"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/logger"
)

type Err error

var (
	ErrGetIndexYamlFailed   Err = fmt.Errorf("get index.yaml failed")
	ErrParseUrlFailed       Err = fmt.Errorf("parse url failed")
	ErrDecodeJsonFailed     Err = fmt.Errorf("decode json failed")
	ErrEmptyAccessKeyId     Err = fmt.Errorf("access key id is empty")
	ErrEmptySecretAccessKey Err = fmt.Errorf("secret access key is empty")
	ErrSchemeNotMatched     Err = fmt.Errorf("scheme not matched")
	ErrInvalidType          Err = fmt.Errorf("invalid repo type")
)

type RepoInterface interface {
	ReadFile(ctx context.Context, filename string) ([]byte, error)
	WriteFile(ctx context.Context, filename string, data []byte) error
}

func New(ctx context.Context, repoType, url, credential string) (RepoInterface, error) {
	u, err := neturl.ParseRequestURI(url)
	if err != nil {
		logger.Error(ctx, "Parse url [%s] failed, error: %+v", url, err)
		return nil, ErrParseUrlFailed
	}

	switch repoType {
	case constants.TypeS3:
		if u.Scheme != constants.TypeS3 {
			return nil, ErrSchemeNotMatched
		}
		return NewS3Interface(ctx, u, credential)
	case constants.TypeHttp:
		if u.Scheme != constants.TypeHttp {
			return nil, ErrSchemeNotMatched
		}
		return NewHttpInterface(ctx, u)
	case constants.TypeHttps:
		if u.Scheme != constants.TypeHttps {
			return nil, ErrSchemeNotMatched
		}
		return NewHttpInterface(ctx, u)
	default:
		return nil, ErrInvalidType
	}
}
