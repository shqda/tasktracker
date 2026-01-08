package main

import (
	"TaskTracker_/internal/server"
	"log"
)

func main() {
	r := server.NewRouter(nil, nil)
	r.RegisterRoutes()
	err := r.Engine.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
