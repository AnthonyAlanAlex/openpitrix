// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"time"

	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/utils"
)

const RepoSelectorTableName = "repo_selector"

func NewRepoSelectorId() string {
	return utils.GetUuid("repos-")
}

type RepoSelector struct {
	RepoSelectorId string
	RepoId         string
	SelectorKey    string
	SelectorValue  string

	CreateTime time.Time
}

var RepoSelectorColumns = GetColumnsFromStruct(&RepoSelector{})

func NewRepoSelector(repoId, selectorKey, selectorValue string) *RepoSelector {
	return &RepoSelector{
		RepoSelectorId: NewRepoSelectorId(),
		RepoId:         repoId,
		SelectorKey:    selectorKey,
		SelectorValue:  selectorValue,

		CreateTime: time.Now(),
	}
}

func RepoSelectorToPb(repoSelector *RepoSelector) *pb.RepoSelector {
	pbRepoSelector := pb.RepoSelector{}
	pbRepoSelector.RepoId = utils.ToProtoString(repoSelector.RepoId)
	pbRepoSelector.RepoSelectorId = utils.ToProtoString(repoSelector.RepoSelectorId)
	pbRepoSelector.SelectorKey = utils.ToProtoString(repoSelector.SelectorKey)
	pbRepoSelector.SelectorValue = utils.ToProtoString(repoSelector.SelectorValue)

	pbRepoSelector.CreateTime = utils.ToProtoTimestamp(repoSelector.CreateTime)
	return &pbRepoSelector
}

func RepoSelectorsToPbs(repoSelectors []*RepoSelector) (pbRepoSelectors []*pb.RepoSelector) {
	for _, repoSelector := range repoSelectors {
		pbRepoSelectors = append(pbRepoSelectors, RepoSelectorToPb(repoSelector))
	}
	return
}

func RepoSelectorsMap(repoSelectors []*RepoSelector) map[string][]*RepoSelector {
	selectorsMap := make(map[string][]*RepoSelector)
	for _, l := range repoSelectors {
		repoId := l.RepoId
		selectorsMap[repoId] = append(selectorsMap[repoId], l)
	}
	return selectorsMap
}
