package server

import (
	"TaskTracker_/internal/server/handlers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine      *gin.Engine
	TaskHandler *handlers.TaskHandler
}

func NewRouter(ts *handlers.TaskHandler) *Router {
	return &Router{
		Engine:      gin.Default(),
		TaskHandler: ts,
	}
}

func (r *Router) RegisterRoutes() {
	taskGroup := r.Engine.Group("/tasks")

	taskGroup.GET("/last", r.TaskHandler.GetLastTask)
	taskGroup.POST("", r.TaskHandler.PostTask)
}
