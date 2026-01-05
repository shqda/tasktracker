package services

import (
	"TaskTracker_/internal/models"
	"errors"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

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
