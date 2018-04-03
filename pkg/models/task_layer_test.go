// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package models

import (
	"strings"
	"testing"
)

func TestWalkTree(t *testing.T) {
	var taskLayer0, taskLayer1, taskLayer2, taskLayer3 TaskLayer
	taskLayer0 = TaskLayer{
		Tasks: nil,
		Child: nil,
	}
	taskLayer1 = TaskLayer{
		Tasks: []*Task{{TaskId: "0"}},
		Child: &taskLayer0,
	}
	taskLayer2 = TaskLayer{
		Tasks: []*Task{{TaskId: "1"}, {TaskId: "2"}},
		Child: &taskLayer1,
	}
	taskLayer3 = TaskLayer{
		Tasks: []*Task{{TaskId: "3"}},
		Child: &taskLayer2,
	}

	expectTasks := []string{"3", "1", "2", "0"}

	var waitTasks, execTasks []string

	err := taskLayer3.WalkTree(func(parent *TaskLayer, current *TaskLayer) error {
		if parent != nil {
			for _, parentTask := range parent.Tasks {
				waitTasks = append(waitTasks, parentTask.TaskId)
			}
		}

		if current != nil {
			for _, task := range current.Tasks {
				execTasks = append(execTasks, task.TaskId)
			}
		}
		return nil
	})

	if err != nil {
		t.Errorf("Error is %+v", err)
	}

	if strings.Join(waitTasks, ",") != strings.Join(expectTasks, ",") ||
		strings.Join(execTasks, ",") != strings.Join(expectTasks, ",") {
		t.Errorf("Wrong task order")
	}
}
