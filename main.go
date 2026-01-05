package main

import (
	"TaskTracker_/internal/server"
	"TaskTracker_/internal/server/handlers"
	"log"
)

func main() {
	r := server.NewRouter(handlers.NewTaskHandler())
	r.RegisterRoutes()
	err := r.Engine.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
