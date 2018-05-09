// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"openpitrix.io/openpitrix/pkg/client"
	clusterclient "openpitrix.io/openpitrix/pkg/client/cluster"
	"openpitrix.io/openpitrix/pkg/logger"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/plugins/vmbased"
	"openpitrix.io/openpitrix/pkg/util/jsonutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

type Processor struct {
	Task *models.Task
}

func NewProcessor(task *models.Task) *Processor {
	return &Processor{
		Task: task,
	}
}

// Post process when task is start
func (t *Processor) Pre() error {
	if t.Task.Directive == "" {
		logger.Warn("Skip empty task [%s] directive", t.Task.TaskId)
		return nil
	}
	var err error
	ctx := client.GetSystemUserContext()
	clusterClient, err := clusterclient.NewClient(ctx)
	if err != nil {
		logger.Error("Executing task [%s] post processor failed: %+v", t.Task.TaskId, err)
		return err
	}
	switch t.Task.TaskAction {
	case vmbased.ActionRunInstances:
		// volume created before instance, so need to change RunInstances task directive
		instance, err := models.NewInstance(t.Task.Directive)
		if err == nil {
			clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{instance.NodeId})
			if err == nil {
				instance.VolumeId = clusterNodes[0].GetVolumeId().GetValue()
				// write back
				t.Task.Directive, err = instance.ToString()
				logger.Debug("Task [%s] new directive: %s", t.Task.TaskId, t.Task.Directive)
			}
		}
	case vmbased.ActionFormatAndMountVolume:
		meta, err := models.NewMeta(t.Task.Directive)
		if err == nil {
			clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{meta.NodeId})
			if err == nil {
				clusterNode := clusterNodes[0]
				clusterRole := clusterNode.GetClusterRole()
				cmd := vmbased.FormatAndMountVolumeCmd(
					clusterNode.GetDevice().GetValue(),
					clusterRole.GetMountPoint().GetValue(),
					clusterRole.GetFileSystem().GetValue(),
					clusterRole.GetMountOptions().GetValue(),
				)
				meta.Cnodes = jsonutil.ToString(vmbased.GetCmdCnodes(
					clusterNode.GetPrivateIp().GetValue(),
					&models.Cmd{
						Cmd:     cmd,
						Timeout: meta.Timeout,
						Id:      t.Task.TaskId,
					},
				))
				meta.DroneIp = clusterNode.GetPrivateIp().GetValue()

				t.Task.TaskAction = vmbased.ActionRegisterCmd
				// write back
				t.Task.Directive, err = meta.ToString()

				logger.Debug("Task [%s] new directive: %s", t.Task.TaskId, t.Task.Directive)
			}
		}
	case vmbased.ActionRegisterMetadata:
		meta, err := models.NewMeta(t.Task.Directive)
		if err == nil {
			pbClusterWrappers, err := clusterClient.GetClusterWrappers(ctx, []string{meta.ClusterId})
			if err == nil {
				metadata := &vmbased.MetadataV1{
					ClusterWrapper: pbClusterWrappers[0],
				}
				meta.Cnodes = jsonutil.ToString(metadata.GetClusterCnodes())

				// write back
				t.Task.Directive, err = meta.ToString()

				logger.Debug("Task [%s] new directive: %s", t.Task.TaskId, t.Task.Directive)
			}
		}
	case vmbased.ActionRegisterCmd:
		// when CreateCluster need to reload ip
		meta, err := models.NewMeta(t.Task.Directive)
		if err == nil {
			if meta.DroneIp == "" {
				clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{meta.NodeId})
				if err == nil {
					meta.DroneIp = clusterNodes[0].GetPrivateIp().GetValue()
					cnodes, err := models.NewCmd(meta.Cnodes)
					if err == nil {
						cnodes.Id = t.Task.TaskId
						meta.Cnodes = jsonutil.ToString(vmbased.GetCmdCnodes(meta.DroneIp, cnodes))
						// write back
						t.Task.Directive, err = meta.ToString()

						logger.Debug("Task [%s] new directive: %s", t.Task.TaskId, t.Task.Directive)
					}
				}
			}
		}
	case vmbased.ActionStartConfd:
		// when CreateCluster need to reload ip
		meta, err := models.NewMeta(t.Task.Directive)
		if err == nil {
			if meta.DroneIp == "" {
				clusterNodes, err := clusterClient.GetClusterNodes(ctx, []string{meta.NodeId})
				if err == nil {
					meta.DroneIp = clusterNodes[0].GetPrivateIp().GetValue()

					// write back
					t.Task.Directive, err = meta.ToString()

					logger.Debug("Task [%s] new directive: %s", t.Task.TaskId, t.Task.Directive)
				}
			}
		}
	case vmbased.ActionCreateVolumes:

	case vmbased.ActionWaitFrontgateAvailable:
	default:
		logger.Info("Nothing to do with task [%s] pre processor", t.Task.TaskId)
	}
	if err != nil {
		logger.Error("Executing task [%s] pre processor failed: %+v", t.Task.TaskId, err)
	}
	return err
}

// Post process when task is done
func (t *Processor) Post() error {
	var err error
	ctx := client.GetSystemUserContext()
	clusterClient, err := clusterclient.NewClient(ctx)
	if err != nil {
		logger.Error("Executing task [%s] post processor failed: %+v", t.Task.TaskId, err)
		return err
	}
	switch t.Task.TaskAction {
	case vmbased.ActionRunInstances:
		if t.Task.Directive == "" {
			logger.Warn("Skip empty task [%s] directive", t.Task.TaskId)
		}
		instance, err := models.NewInstance(t.Task.Directive)
		if err == nil {
			_, err = clusterClient.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
				ClusterNode: &pb.ClusterNode{
					NodeId:     pbutil.ToProtoString(instance.NodeId),
					InstanceId: pbutil.ToProtoString(instance.InstanceId),
					Device:     pbutil.ToProtoString(instance.Device),
					PrivateIp:  pbutil.ToProtoString(instance.PrivateIp),
				},
			})
		}
	case vmbased.ActionCreateVolumes:
		if t.Task.Directive == "" {
			logger.Warn("Skip empty task [%s] directive", t.Task.TaskId)
		}
		volume, err := models.NewVolume(t.Task.Directive)
		if err == nil {
			_, err = clusterClient.ModifyClusterNode(ctx, &pb.ModifyClusterNodeRequest{
				ClusterNode: &pb.ClusterNode{
					NodeId:   pbutil.ToProtoString(t.Task.NodeId),
					VolumeId: pbutil.ToProtoString(volume.VolumeId),
				},
			})
		}
	default:
		logger.Info("Nothing to do with task [%s] post processor", t.Task.TaskId)
	}
	if err != nil {
		logger.Error("Executing task [%s] post processor failed: %+v", t.Task.TaskId, err)
	}
	return err
}
