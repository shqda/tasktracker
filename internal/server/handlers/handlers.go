package handlers

import (
	"TaskTracker_/internal/services"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskHandler struct {
	taskService *services.TaskService
}

var (
	ErrInvalidJSON = errors.New("invalid JSON")
)

func NewTaskHandler() *TaskHandler {
	return &TaskHandler{
		taskService: services.NewTaskService(),
	}
}

func (ts *TaskHandler) PostTask(c *gin.Context) {
	var input struct {
		Task string `json:"task" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrInvalidJSON)
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
