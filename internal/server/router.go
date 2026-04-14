package server

import (
	"github.com/gin-gonic/gin"
)

type TaskHandlerInterface interface {
	GetLastTask(c *gin.Context)
	PostTask(c *gin.Context)
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
	taskGroup.POST("", r.TaskHandler.PostTask)
}
