package main

import (
	"fmt"
	"log"
	"tasktracker/config"
	"tasktracker/internal/server"
	"tasktracker/internal/server/handler"
	"tasktracker/internal/service"
	"tasktracker/internal/storage/postgres"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		"disable",
	)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	repo := &postgres.PostgresDB{DB: db}
	svc := service.NewTaskService(repo)
	hndlr := handler.NewTaskHandler(svc)

	r := server.NewRouter(nil, hndlr)
	r.RegisterRoutes()

	addr := fmt.Sprintf(":%s", cfg.Serv.Port)
	if err := r.Engine.Run(addr); err != nil {
		log.Fatal(err)
	}
}
