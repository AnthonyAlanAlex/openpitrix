// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package qingcloud

import (
	"context"
	"fmt"

	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/stringutil"
)

type Provider struct {
	ctx context.Context
}

func NewProvider(ctx context.Context) *Provider {
	return &Provider{
		ctx,
	}
}

func (p *Provider) ParseClusterConf(versionId, runtimeId, conf string, clusterWrapper *models.ClusterWrapper) error {
	frameInterface, err := vmbased.GetFrameInterface(p.ctx, nil)
	if err != nil {
		return err
	}
	return frameInterface.ParseClusterConf(versionId, runtimeId, conf, clusterWrapper)
}

func (p *Provider) SplitJobIntoTasks(job *models.Job) (*models.TaskLayer, error) {
	return vmbased.SplitJobIntoTasks(p.ctx, job)
}

func (p *Provider) HandleSubtask(task *models.Task) error {
	handler := GetProviderHandler(p.ctx)
	return vmbased.HandleSubtask(p.ctx, task, handler)
}

func (p *Provider) WaitSubtask(task *models.Task) error {
	handler := GetProviderHandler(p.ctx)
	return vmbased.WaitSubtask(p.ctx, task, handler)
}

func (p *Provider) DescribeSubnets(ctx context.Context, req *pb.DescribeSubnetsRequest) (*pb.DescribeSubnetsResponse, error) {
	handler := GetProviderHandler(p.ctx)
	return handler.DescribeSubnets(ctx, req)
}

func (p *Provider) CheckResource(ctx context.Context, clusterWrapper *models.ClusterWrapper) error {
	handler := GetProviderHandler(p.ctx)
	return handler.CheckResourceQuotas(ctx, clusterWrapper)
}

func (p *Provider) DescribeVpc(runtimeId, vpcId string) (*models.Vpc, error) {
	handler := GetProviderHandler(p.ctx)
	return handler.DescribeVpc(runtimeId, vpcId)
}

func (p *Provider) ValidateCredential(url, credential, zone string) error {
	handler := GetProviderHandler(p.ctx)
	zones, err := handler.DescribeZones(url, credential)
	if err != nil {
		return err
	}
	if zone == "" {
		return nil
	}
	if !stringutil.StringIn(zone, zones) {
		return fmt.Errorf("cannot access zone [%s]", zone)
	}
	return nil
}

func (p *Provider) UpdateClusterStatus(job *models.Job) error {
	return nil
}

func (p *Provider) DescribeRuntimeProviderZones(url, credential string) ([]string, error) {
	handler := GetProviderHandler(p.ctx)
	return handler.DescribeZones(url, credential)
}

func (p *Provider) DescribeClusterDetails(ctx context.Context, cluster *models.ClusterWrapper) error {
	return nil
}
