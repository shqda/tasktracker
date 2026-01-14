package storage

import "tasktracker/internal/models"

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TaskStorage --structname=MockTaskStorage --case=underscore
type TaskStorage interface {
	InsertTask(title string) (int, error)
	GetLastTask() (*models.Task, error)
}
