package main

import (
	"TaskTracker_/internal/server"
	"TaskTracker_/internal/server/handlers"
	"TaskTracker_/internal/services"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := server.NewRouter(gin.Default(), handlers.NewTaskHandler(services.NewTaskService()))
	r.RegisterRoutes()
	err := r.Engine.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
