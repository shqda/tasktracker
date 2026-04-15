package handler

import (
	"errors"
	"net/http"
	"strconv"
	"tasktracker/internal/model"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidJSON = errors.New("invalid JSON")
)

type TaskServiceInterface interface {
	CreateTask(title string) (*model.Task, error)
	GetLastTask() (*model.Task, error)
	GetTaskByID(id int) (*model.Task, error)
	GetAllTasks() ([]model.Task, error)
	RenameTask(id int, title string) error
	DeleteTask(id int) error
}

type TaskHandler struct {
	taskService TaskServiceInterface
}

func NewTaskHandler(ts TaskServiceInterface) *TaskHandler {
	return &TaskHandler{taskService: ts}
}

func parseID(c *gin.Context) (int, error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
	}
	return id, err
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, task)
}

func (ts *TaskHandler) GetLastTask(c *gin.Context) {
	task, err := ts.taskService.GetLastTask()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (ts *TaskHandler) GetTaskByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	task, err := ts.taskService.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (ts *TaskHandler) GetAllTasks(c *gin.Context) {
	task, err := ts.taskService.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (ts *TaskHandler) RenameTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	var input struct {
		Title string `json:"title" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": ErrInvalidJSON.Error()})
		return
	}
	err = ts.taskService.RenameTask(id, input.Title)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, input)
}

func (ts *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	err = ts.taskService.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}
