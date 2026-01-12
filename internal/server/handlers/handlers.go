package handlers

import (
	"TaskTracker_/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=TaskHandlerInterface --structname=MockTaskHandler --case=underscore
type TaskHandlerInterface interface {
	GetLastTask(c *gin.Context)
	PostTask(c *gin.Context)
}

type TaskHandler struct {
	taskService services.TaskServiceInterface
}

var (
	ErrInvalidJSON = errors.New("invalid JSON")
)

func NewTaskHandler(ts services.TaskServiceInterface) *TaskHandler {
	if ts == nil {
		ts = services.NewTaskService()
	}
	return &TaskHandler{taskService: ts}
}

func (ts *TaskHandler) PostTask(c *gin.Context) {
	var input struct {
		Task string `json:"task" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON.Error()})
		return
	}
	task, err := ts.taskService.CreateTask(input.Task)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (ts *TaskHandler) GetLastTask(c *gin.Context) {
	task, err := ts.taskService.LastTask()
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, task)
	}
}
