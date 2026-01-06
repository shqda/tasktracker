package services

import (
	"TaskTracker_/internal/models"
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestTaskService_CreateTask(t *testing.T) {
	tests := []struct {
		name  string
		ts    TaskService
		title string
		want  TaskService
	}{
		{
			name: "empty tasks list",
			ts: TaskService{
				tasks:  make([]models.Task, 0),
				nextID: 0,
			},
			title: "simple task",
			want: TaskService{
				tasks: []models.Task{
					models.Task{Title: "simple task"},
				},
				nextID: 1,
			},
		},
		{
			name: "non-empty tasks list",
			ts: TaskService{
				tasks: []models.Task{
					models.Task{Title: "1st task", ID: 0},
				},
				nextID: 1,
			},
			title: "2nd task",
			want: TaskService{
				tasks: []models.Task{
					models.Task{Title: "1st task", ID: 0},
					models.Task{Title: "2nd task", ID: 1},
				},
				nextID: 2,
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.ts.CreateTask(tc.title)
			assert.Equal(t, tc.ts, tc.want)
		})
	}
}

func TestTaskService_LastTask(t *testing.T) {
	tests := []struct {
		name      string
		ts        TaskService
		want      models.Task
		WantError error
	}{
		{
			name: "empty tasks list",
			ts: TaskService{
				tasks:  make([]models.Task, 0),
				nextID: 0,
			},
			want:      models.Task{},
			WantError: ErrTaskNotFound,
		},
		{
			name: "non-empty tasks list",
			ts: TaskService{
				tasks: []models.Task{
					models.Task{Title: "1st task", ID: 0},
					models.Task{Title: "2nd task", ID: 1},
				},
				nextID: 2,
			},
			want:      models.Task{Title: "2nd task", ID: 1},
			WantError: nil,
		},
		{
			name: "last task empty title",
			ts: TaskService{
				tasks: []models.Task{
					models.Task{Title: "", ID: 0},
				},
				nextID: 2,
			},
			want:      models.Task{Title: "", ID: 0},
			WantError: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.ts.LastTask()
			assert.Equal(t, err, tc.WantError)
			assert.Equal(t, result, tc.want)
		})
	}
}
