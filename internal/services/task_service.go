package services

import (
	"TaskTracker_/internal/models"
	"errors"
	"tasktracker/internal/models"
	"tasktracker/internal/storage"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TaskServiceInterface --structname=MockTaskService
type TaskServiceInterface interface {
	CreateTask(title string) models.Task
	LastTask() (models.Task, error)
}

type TaskService struct {
	tasks  []models.Task
	nextID int32
}

func NewTaskService() *TaskService {
	return &TaskService{}
}

func (ts *TaskService) CreateTask(title string) models.Task {
	task := models.Task{
		ID:    ts.nextID,
		Title: title,
	}
	ts.nextID++
	ts.tasks = append(ts.tasks, task)
	return task
}

func (ts *TaskService) LastTask() (models.Task, error) {
	if len(ts.tasks) != 0 {
		return ts.tasks[len(ts.tasks)-1], nil
	}
	return models.Task{}, ErrTaskNotFound
}
