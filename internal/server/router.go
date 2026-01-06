package server

import (
	"TaskTracker_/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine      *gin.Engine
	TaskHandler *handlers.TaskHandler
}

func NewRouter(e *gin.Engine, ts *handlers.TaskHandler) *Router {
	if e == nil {
		e = gin.Default()
	}
	if ts == nil {
		ts = handlers.NewTaskHandler(nil)
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
