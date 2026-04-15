package server

import (
	"github.com/gin-gonic/gin"
)

type TaskHandlerInterface interface {
	GetTaskByID(c *gin.Context)
	GetLastTask(c *gin.Context)
	GetAllTasks(c *gin.Context)
	PostTask(c *gin.Context)
	DeleteTask(c *gin.Context)
	RenameTask(c *gin.Context)
}

type Router struct {
	Engine      *gin.Engine
	TaskHandler TaskHandlerInterface
}

func NewRouter(e *gin.Engine, ts TaskHandlerInterface) *Router {
	if e == nil {
		e = gin.Default()
	}
	return &Router{
		Engine:      e,
		TaskHandler: ts,
	}
}

func (r *Router) RegisterRoutes() {
	taskGroup := r.Engine.Group("/tasks")

	taskGroup.GET("/last", r.TaskHandler.GetLastTask)
	taskGroup.GET("/:id", r.TaskHandler.GetTaskByID)
	taskGroup.GET("/", r.TaskHandler.GetAllTasks)
	taskGroup.POST("", r.TaskHandler.PostTask)
	taskGroup.DELETE("/:id", r.TaskHandler.DeleteTask)
	taskGroup.PUT("/:id", r.TaskHandler.RenameTask)
}
