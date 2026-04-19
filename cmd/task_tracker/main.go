package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
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
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		slog.Error("failed to load config", "err", err)
		os.Exit(1)
	}
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
		"disable",
	)
	slog.Info("connecting to database")
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Error("failed connect to database", "err", err)
		os.Exit(1)
	}
	defer func() {
		err = db.Close()
		if err != nil {
			slog.Error("error closing database", "err", err)
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

	listenErrChan := make(chan error, 1)
	go func() {
		slog.Info("starting server", "addr", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			listenErrChan <- err
		}
	}()
	shutdown, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	select {
	case <-shutdown.Done():
	case err := <-listenErrChan:
		slog.Error("stopped listening", "err", err)
	}

	slog.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("shutdown with error", "err", err)
	}
}
