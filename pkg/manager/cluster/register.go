// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package cluster

import (
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pi"
)

type Register struct {
	*pi.Pi
	ClusterId      string
	SubnetId       string
	VpcId          string
	FrontgateId    string
	Owner          string
	ClusterType    int32
	ClusterWrapper *models.ClusterWrapper
	Runtime        *Runtime
}

func (r *Register) RegisterClusterWrapper() error {
	// register cluster
	if r.ClusterWrapper.Cluster != nil {
		r.ClusterWrapper.Cluster.ClusterId = r.ClusterId
		r.ClusterWrapper.Cluster.RuntimeEnvId = r.Runtime.RuntimeEnvId
		r.ClusterWrapper.Cluster.FrontgateId = r.FrontgateId
		r.ClusterWrapper.Cluster.VpcId = r.VpcId
		r.ClusterWrapper.Cluster.Owner = r.Owner
		r.ClusterWrapper.Cluster.ClusterType = r.ClusterType
		_, err := r.Db.
			InsertInto(models.ClusterTableName).
			Columns(models.ClusterColumns...).
			Record(r.ClusterWrapper.Cluster).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterTableName, r.ClusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster node
	for _, clusterNode := range r.ClusterWrapper.ClusterNodes {
		clusterNode.ClusterId = r.ClusterId
		clusterNode.NodeId = models.NewClusterNodeId()
		clusterNode.Owner = r.Owner
		_, err := r.Db.
			InsertInto(models.ClusterNodeTableName).
			Columns(models.ClusterNodeColumns...).
			Record(clusterNode).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterNodeTableName, r.ClusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster common
	for _, clusterCommon := range r.ClusterWrapper.ClusterCommons {
		clusterCommon.ClusterId = r.ClusterId
		_, err := r.Db.
			InsertInto(models.ClusterCommonTableName).
			Columns(models.ClusterCommonColumns...).
			Record(clusterCommon).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterCommonTableName, r.ClusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster link
	for _, clusterLink := range r.ClusterWrapper.ClusterLinks {
		clusterLink.ClusterId = r.ClusterId
		clusterLink.Owner = r.Owner
		_, err := r.Db.
			InsertInto(models.ClusterLinkTableName).
			Columns(models.ClusterLinkColumns...).
			Record(clusterLink).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterLinkTableName, r.ClusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster role
	for _, clusterRole := range r.ClusterWrapper.ClusterRoles {
		clusterRole.ClusterId = r.ClusterId
		_, err := r.Db.
			InsertInto(models.ClusterRoleTableName).
			Columns(models.ClusterRoleColumns...).
			Record(clusterRole).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterRoleTableName, r.ClusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	// register cluster loadbalancer
	for _, clusterLoadbalancer := range r.ClusterWrapper.ClusterLoadbalancers {
		clusterLoadbalancer.ClusterId = r.ClusterId
		_, err := r.Db.
			InsertInto(models.ClusterLoadbalancerTableName).
			Columns(models.ClusterLoadbalancerColumns...).
			Record(clusterLoadbalancer).
			Exec()
		if err != nil {
			logger.Errorf("Failed to insert table [%s] with cluster id [%s]: %+v",
				models.ClusterLoadbalancerTableName, r.ClusterWrapper.Cluster.ClusterId, err)
			return err
		}
	}

	return nil
}
