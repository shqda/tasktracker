package handlers

import (
	"TaskTracker_/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskHandlerInterface interface {
	GetLastTask(c *gin.Context)
	PostTask(c *gin.Context)
}

type TaskHandler struct {
	taskService *services.TaskService
}

var (
	ErrInvalidJSON = errors.New("invalid JSON")
)

func NewTaskHandler(ts *services.TaskService) *TaskHandler {
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
	c.JSON(http.StatusCreated, ts.taskService.CreateTask(input.Task))
}

func (ts *TaskHandler) GetLastTask(c *gin.Context) {
	task, err := ts.taskService.LastTask()
	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
	} else {
		c.JSON(200, task)
	}
}
