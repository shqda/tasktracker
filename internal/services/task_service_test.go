package services

import (
	"github.com/stretchr/testify/assert"
	"tasktracker/internal/models"
	"tasktracker/internal/storage/mocks"
	"testing"
)

func TestTaskService_CreateTask(t *testing.T) {
	mockStorage := mocks.NewMockTaskStorage(t)
	svc := TaskService{Storage: mockStorage}

	methodName := "InsertTask"

	t.Run("step1: create task1", func(t *testing.T) {
		mockStorage.On(methodName, "task1").Return(1, nil)
		task, err := svc.CreateTask("task1")
		assert.NoError(t, err)
		assert.Equal(t, int32(1), task.ID)
	})

	t.Run("step2: create task2", func(t *testing.T) {
		mockStorage.On(methodName, "task2").Return(2, nil)
		task, err := svc.CreateTask("task2")
		assert.NoError(t, err)
		assert.Equal(t, int32(2), task.ID)
	})

	mockStorage.AssertExpectations(t)
}

func TestTaskService_LastTask(t *testing.T) {
	methodName := "GetLastTask"

	tests := []struct {
		name      string
		setupMock func(*mocks.MockTaskStorage)
		wantTask  *models.Task
		wantErr   error
	}{
		{
			name: "empty database",
			setupMock: func(m *mocks.MockTaskStorage) {
				m.On(methodName).Return(nil, ErrNoTasks)
			},
			wantTask: nil,
			wantErr:  ErrNoTasks,
		},
		{
			name: "last task exists",
			setupMock: func(m *mocks.MockTaskStorage) {
				m.On(methodName).Return(&models.Task{ID: 1, Title: "2nd task"}, nil)
			},
			wantTask: &models.Task{ID: 1, Title: "2nd task"},
			wantErr:  nil,
		},
		{
			name: "last task with empty title",
			setupMock: func(m *mocks.MockTaskStorage) {
				m.On(methodName).Return(&models.Task{ID: 0, Title: ""}, nil)
			},
			wantTask: &models.Task{ID: 0, Title: ""},
			wantErr:  nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			mockStorage := mocks.NewMockTaskStorage(t)
			tc.setupMock(mockStorage)

			svc := TaskService{Storage: mockStorage}

			result, err := svc.LastTask()
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantTask, result)

			mockStorage.AssertExpectations(t)
		})
	}
}
