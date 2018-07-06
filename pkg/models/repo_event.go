// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/topic"
	"openpitrix.io/openpitrix/pkg/util/idutil"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
)

const RepoEventTableName = "repo_event"

func NewRepoEventId() string {
	return idutil.GetUuid("repoe-")
}

type RepoEvent struct {
	RepoEventId string
	RepoId      string
	Status      string
	Result      string
	Owner       string
	CreateTime  time.Time
	StatusTime  time.Time
}

func (r RepoEvent) GetTopicResource() topic.Resource {
	return topic.NewResource(RepoEventTableName, r.RepoEventId).SetStatus(r.Status)
}

var RepoEventColumns = GetColumnsFromStruct(&RepoEvent{})

func NewRepoEvent(repoId, owner string) *RepoEvent {
	return &RepoEvent{
		RepoEventId: NewRepoEventId(),
		RepoId:      repoId,
		Owner:       owner,
		Status:      constants.StatusPending,
		Result:      "",
		CreateTime:  time.Now(),
		StatusTime:  time.Now(),
	}
}

func RepoEventToPb(repoTask *RepoEvent) *pb.RepoEvent {
	pbRepoTask := pb.RepoEvent{}
	pbRepoTask.RepoEventId = pbutil.ToProtoString(repoTask.RepoEventId)
	pbRepoTask.RepoId = pbutil.ToProtoString(repoTask.RepoId)
	pbRepoTask.Status = pbutil.ToProtoString(repoTask.Status)
	pbRepoTask.Result = pbutil.ToProtoString(repoTask.Result)
	pbRepoTask.Owner = pbutil.ToProtoString(repoTask.Owner)
	pbRepoTask.CreateTime = pbutil.ToProtoTimestamp(repoTask.CreateTime)
	pbRepoTask.StatusTime = pbutil.ToProtoTimestamp(repoTask.StatusTime)
	return &pbRepoTask
}

func RepoEventsToPbs(repoTasks []*RepoEvent) (pbRepoTasks []*pb.RepoEvent) {
	for _, repoTask := range repoTasks {
		pbRepoTasks = append(pbRepoTasks, RepoEventToPb(repoTask))
	}
	return
}
