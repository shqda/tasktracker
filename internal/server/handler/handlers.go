package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"tasktracker/internal/errs"
	"tasktracker/internal/model"

	"github.com/gin-gonic/gin"
)

const (
	internalErrorMsg = "internal error"
	invalidJSONMsg   = "invalid JSON"
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
		slog.Warn("failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSONMsg})
		return
	}
	task, err := ts.taskService.CreateTask(input.Task)
	if err != nil {
		slog.Error("failed to create task", "title", input.Task, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrorMsg})
		return
	}
	slog.Info("task created", "id", task.ID)
	c.JSON(http.StatusCreated, task)
}

func (ts *TaskHandler) GetLastTask(c *gin.Context) {
	task, err := ts.taskService.GetLastTask()
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Error("failed to get task", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrorMsg})
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
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Error("failed to get task", "id", id, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrorMsg})
		return
	}
	c.JSON(http.StatusOK, task)
}

func (ts *TaskHandler) GetAllTasks(c *gin.Context) {
	task, err := ts.taskService.GetAllTasks()
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Error("failed to get tasks", "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrorMsg})
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
		slog.Warn("failed to bind JSON", "err", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidJSONMsg})
		return
	}
	err = ts.taskService.RenameTask(id, input.Title)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Error("failed to rename task", "id", id, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrorMsg})
		return
	}
	slog.Info("task renamed", "id", id)
	c.JSON(http.StatusOK, input)
}

func (ts *TaskHandler) DeleteTask(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	err = ts.taskService.DeleteTask(id)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		slog.Error("failed to delete task", "id", id, "err", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": internalErrorMsg})
		return
	}
	slog.Info("task deleted", "id", id)
	c.AbortWithStatus(http.StatusNoContent)
}
