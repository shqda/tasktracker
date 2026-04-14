package main

import (
	"log"
	"tasktracker/internal/server"
	"tasktracker/internal/server/handler"
	"tasktracker/internal/service"
	"tasktracker/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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

	if err := r.Engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
