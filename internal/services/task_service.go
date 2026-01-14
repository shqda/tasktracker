package services

import (
	"errors"
	"tasktracker/internal/models"
	"tasktracker/internal/storage"
)

var (
	ErrNoTasks         = errors.New("no tasks found")
	ErrCreatingFailure = errors.New("task creating failure")
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TaskServiceInterface --structname=MockTaskService --case=underscore
type TaskServiceInterface interface {
	CreateTask(title string) (*models.Task, error)
	LastTask() (*models.Task, error)
}

type TaskService struct {
	Storage storage.TaskStorage
}

func NewTaskService() *TaskService {
	return new(TaskService)
}

func (ts *TaskService) CreateTask(title string) (*models.Task, error) {
	id, err := ts.Storage.InsertTask(title)
	if err != nil {
		return nil, ErrCreatingFailure
	}
	return &models.Task{
		ID:    int32(id),
		Title: title,
	}, nil
}

func (ts *TaskService) LastTask() (*models.Task, error) {
	task, err := ts.Storage.GetLastTask()
	if err != nil {
		return nil, ErrNoTasks
	}
	return task, nil
}
