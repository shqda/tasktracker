package main

import (
	"log"
	"tasktracker/internal/server"
)

func main() {
	db, err := sqlx.Connect("postgres", "host=localhost port=5432 user=postgres password=postgres dbname=tasktracker sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	repo := &postgres.PostgresDB{DB: db}
	svc := service.NewTaskService(repo)
	hndlr := handler.NewTaskHandler(svc)

	r := server.NewRouter(nil, hndlr)
	r.RegisterRoutes()
	err := r.Engine.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
