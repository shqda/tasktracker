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
	GetTaskByID(id int) (*model.Task, error)
	GetLastTask() (*model.Task, error)
	GetAllTasks() ([]model.Task, error)
	DeleteTask(id int) error
	UpdateTask(id int, title string) error
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

func (ts *TaskService) GetLastTask() (*model.Task, error) {
	task, err := ts.Storage.GetLastTask()
	if err != nil {
		return nil, ErrNoTasks
	}
	return task, nil
}

func (ts *TaskService) GetTaskByID(id int) (*model.Task, error) {
	task, err := ts.Storage.GetTaskByID(id)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (ts *TaskService) GetAllTasks() ([]model.Task, error) {
	tasks, err := ts.Storage.GetAllTasks()
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (ts *TaskService) RenameTask(id int, title string) error {
	return ts.Storage.UpdateTask(id, title)
}

func (ts *TaskService) DeleteTask(id int) error {
	return ts.Storage.DeleteTask(id)
}
