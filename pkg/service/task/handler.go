// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package task

import (
	"context"
	"fmt"
	"strings"

	"openpitrix.io/openpitrix/pkg/constants"
	"openpitrix.io/openpitrix/pkg/db"
	"openpitrix.io/openpitrix/pkg/gerr"
	"openpitrix.io/openpitrix/pkg/manager"
	"openpitrix.io/openpitrix/pkg/models"
	"openpitrix.io/openpitrix/pkg/pb"
	"openpitrix.io/openpitrix/pkg/util/pbutil"
	"openpitrix.io/openpitrix/pkg/util/senderutil"
)

func (p *Server) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	newTask := models.NewTask(
		"",
		req.GetJobId().GetValue(),
		req.GetNodeId().GetValue(),
		req.GetTarget().GetValue(),
		req.GetTaskAction().GetValue(),
		req.GetDirective().GetValue(),
		s.UserId,
		req.GetFailureAllowed().GetValue(),
	)

	if req.GetStatus().GetValue() == constants.StatusFailed {
		newTask.Status = req.GetStatus().GetValue()
	}

	_, err := p.Db.
		InsertInto(models.TaskTableName).
		Columns(models.TaskColumns...).
		Record(newTask).
		Exec()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
	}

	if newTask.Status != constants.StatusFailed {
		err = p.controller.queue.Enqueue(newTask.TaskId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorCreateResourcesFailed)
		}
	}

	res := &pb.CreateTaskResponse{
		TaskId: pbutil.ToProtoString(newTask.TaskId),
		JobId:  pbutil.ToProtoString(newTask.JobId),
	}
	return res, nil
}

func (p *Server) DescribeTasks(ctx context.Context, req *pb.DescribeTasksRequest) (*pb.DescribeTasksResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)
	var tasks []*models.Task
	offset := pbutil.GetOffsetFromRequest(req)
	limit := pbutil.GetLimitFromRequest(req)

	query := p.Db.
		Select(models.TaskColumns...).
		From(models.TaskTableName).
		Offset(offset).
		Limit(limit).
		Where(manager.BuildFilterConditions(req, models.TaskTableName)).
		Where(db.Eq("owner", s.UserId)).
		OrderDir("create_time", true)

	_, err := query.Load(&tasks)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}
	count, err := query.Count()
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorDescribeResourcesFailed)
	}

	res := &pb.DescribeTasksResponse{
		TaskSet:    models.TasksToPbs(tasks),
		TotalCount: uint32(count),
	}
	return res, nil
}

func (p *Server) RetryTasks(ctx context.Context, req *pb.RetryTasksRequest) (*pb.RetryTasksResponse, error) {
	s := senderutil.GetSenderFromContext(ctx)

	taskIds := req.GetTaskId()
	var tasks []*models.Task
	query := p.Db.
		Select(models.TaskColumns...).
		From(models.TaskTableName).
		Where(db.Eq("task_id", taskIds)).
		Where(db.Eq("owner", s.UserId))

	_, err := query.Load(&tasks)
	if err != nil {
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorRetryTaskFailed, strings.Join(taskIds, ","))
	}

	if len(tasks) != len(taskIds) {
		err = fmt.Errorf("retryTasks [%s] with count [%d]", strings.Join(taskIds, ","), len(tasks))
		return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorRetryTaskFailed, strings.Join(taskIds, ","))
	}

	for _, taskId := range taskIds {
		err = p.controller.queue.Enqueue(taskId)
		if err != nil {
			return nil, gerr.NewWithDetail(gerr.Internal, err, gerr.ErrorRetryTaskFailed, strings.Join(taskIds, ","))
		}
	}

	res := &pb.RetryTasksResponse{
		TaskSet: models.TasksToPbs(tasks),
	}
	return res, nil
}
