package service

import (
	"errors"
	"tasktracker/internal/model"
)

var (
	ErrNoTasks         = errors.New("no tasks found")
	ErrCreatingFailure = errors.New("task creating failure")
)

type TaskStorage interface {
	InsertTask(title string) (int, error)
	GetLastTask() (*model.Task, error)
}

type TaskService struct {
	Storage TaskStorage
}

func NewTaskService(s TaskStorage) *TaskService {
	return &TaskService{Storage: s}
}

func (ts *TaskService) CreateTask(title string) (*model.Task, error) {
	id, err := ts.Storage.InsertTask(title)
	if err != nil {
		return nil, ErrCreatingFailure
	}
	return &model.Task{
		ID:    int32(id),
		Title: title,
	}, nil
}

func (ts *TaskService) LastTask() (*model.Task, error) {
	task, err := ts.Storage.GetLastTask()
	if err != nil {
		return nil, ErrNoTasks
	}
	return task, nil
}
