package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"tasktracker/config"
	"tasktracker/internal/server"
	"tasktracker/internal/server/handler"
	"tasktracker/internal/service"
	"tasktracker/internal/storage/postgres"
	"time"

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
	defer func() {
		err = db.Close()
		if err != nil {
			log.Printf("error closing db: %v", err)
		}
	}()

	repo := &postgres.PostgresDB{DB: db}
	svc := service.NewTaskService(repo)
	hndlr := handler.NewTaskHandler(svc)

	r := server.NewRouter(nil, hndlr)
	r.RegisterRoutes()

	addr := fmt.Sprintf(":%s", cfg.Serv.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r.Engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Stopped listening: %v\n", err)
			log.Fatal(err)
		}
	}()

	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()
	<-shutdown.Done()

	fmt.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Shutdown with error: %v", err)
	}
}
