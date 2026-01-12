package main

import (
	"log"
	"tasktracker/internal/server"
)

func main() {
	r := server.NewRouter(nil, nil)
	r.RegisterRoutes()
	err := r.Engine.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
