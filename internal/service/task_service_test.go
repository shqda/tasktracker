package service

import (
	"errors"
	"tasktracker/internal/model"
	"tasktracker/internal/service/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTaskService_CreateTask(t *testing.T) {
	methodName := "InsertTask"

	tests := []struct {
		name      string
		title     string
		setupMock func(*service.MockTaskStorage)
		wantTask  *model.Task
		wantErr   error
	}{
		{
			name:  "valid",
			title: "task1",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName, "task1").Once().Return(1, nil)
			},
			wantTask: &model.Task{ID: 1, Title: "task1"},
			wantErr:  nil,
		},
		{
			name:  "storage error",
			title: "task1",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName, "task1").Once().Return(0, errors.New("db error"))
			},
			wantTask: nil,
			wantErr:  ErrCreatingFailure,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := service.NewMockTaskStorage(t)
			tt.setupMock(mockStorage)

			svc := TaskService{Storage: mockStorage}

			task, err := svc.CreateTask(tt.title)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantTask, task)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestTaskService_GetLastTask(t *testing.T) {
	methodName := "GetLastTask"

	tests := []struct {
		name      string
		setupMock func(*service.MockTaskStorage)
		wantTask  *model.Task
		wantErr   error
	}{
		{
			name: "empty database",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName).Once().Return(nil, errors.New("db error"))
			},
			wantTask: nil,
			wantErr:  ErrNoTasks,
		},
		{
			name: "last task exists",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName).Once().Return(&model.Task{ID: 1, Title: "2nd task"}, nil)
			},
			wantTask: &model.Task{ID: 1, Title: "2nd task"},
			wantErr:  nil,
		},
		{
			name: "last task with empty title",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName).Once().Return(&model.Task{ID: 0, Title: ""}, nil)
			},
			wantTask: &model.Task{ID: 0, Title: ""},
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := service.NewMockTaskStorage(t)
			tt.setupMock(mockStorage)

			svc := TaskService{Storage: mockStorage}

			result, err := svc.GetLastTask()
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.wantTask, result)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestTaskService_GetAllTasks(t *testing.T) {
	methodName := "GetAllTasks"

	tests := []struct {
		name       string
		setupMock  func(*service.MockTaskStorage)
		wantTasks  []model.Task
		wantErrMsg string
	}{
		{
			name: "valid",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName).Once().Return([]model.Task{
					{ID: 1, Title: "task1"},
					{ID: 2, Title: "task2"},
				}, nil)
			},
			wantTasks: []model.Task{
				{ID: 1, Title: "task1"},
				{ID: 2, Title: "task2"},
			},
		},
		{
			name: "empty database",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName).Once().Return([]model.Task{}, nil)
			},
			wantTasks: []model.Task{},
		},
		{
			name: "storage error",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName).Once().Return(nil, errors.New("db error"))
			},
			wantTasks:  nil,
			wantErrMsg: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := service.NewMockTaskStorage(t)
			tt.setupMock(mockStorage)

			svc := TaskService{Storage: mockStorage}

			result, err := svc.GetAllTasks()
			if tt.wantErrMsg != "" {
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantTasks, result)

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestTaskService_RenameTask(t *testing.T) {
	methodName := "UpdateTask"

	tests := []struct {
		name       string
		id         int
		title      string
		setupMock  func(*service.MockTaskStorage)
		wantErrMsg string
	}{
		{
			name:  "valid",
			id:    1,
			title: "new name",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName, 1, "new name").Once().Return(nil)
			},
		},
		{
			name:  "storage error",
			id:    1,
			title: "new name",
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName, 1, "new name").Once().Return(errors.New("db error"))
			},
			wantErrMsg: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := service.NewMockTaskStorage(t)
			tt.setupMock(mockStorage)

			svc := TaskService{Storage: mockStorage}

			err := svc.RenameTask(tt.id, tt.title)
			if tt.wantErrMsg != "" {
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	methodName := "DeleteTask"

	tests := []struct {
		name       string
		id         int
		setupMock  func(*service.MockTaskStorage)
		wantErrMsg string
	}{
		{
			name: "valid",
			id:   1,
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName, 1).Once().Return(nil)
			},
		},
		{
			name: "invalid id",
			id:   999,
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName, 999).Once().Return(ErrNoTasks)
			},
			wantErrMsg: ErrNoTasks.Error(),
		},
		{
			name: "storage error",
			id:   1,
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName, 1).Once().Return(errors.New("db error"))
			},
			wantErrMsg: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := service.NewMockTaskStorage(t)
			tt.setupMock(mockStorage)

			svc := TaskService{Storage: mockStorage}

			err := svc.DeleteTask(tt.id)
			if tt.wantErrMsg != "" {
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}

			mockStorage.AssertExpectations(t)
		})
	}
}

func TestTaskService_GetTaskByID(t *testing.T) {
	methodName := "GetTaskByID"

	tests := []struct {
		name       string
		id         int
		setupMock  func(*service.MockTaskStorage)
		wantTask   *model.Task
		wantErrMsg string
	}{
		{
			name: "valid",
			id:   1,
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName, 1).Once().Return(&model.Task{ID: 1, Title: "task1"}, nil)
			},
			wantTask: &model.Task{ID: 1, Title: "task1"},
		},
		{
			name: "storage error",
			id:   1,
			setupMock: func(m *service.MockTaskStorage) {
				m.On(methodName, 1).Once().Return(nil, errors.New("db error"))
			},
			wantTask:   nil,
			wantErrMsg: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockStorage := service.NewMockTaskStorage(t)
			tt.setupMock(mockStorage)

			svc := TaskService{Storage: mockStorage}

			result, err := svc.GetTaskByID(tt.id)
			if tt.wantErrMsg != "" {
				assert.EqualError(t, err, tt.wantErrMsg)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.wantTask, result)

			mockStorage.AssertExpectations(t)
		})
	}
}
