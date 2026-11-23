package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example.com/go-backend/config"
	"example.com/go-backend/database"
	"example.com/go-backend/handler"
	"example.com/go-backend/middleware"
	"example.com/go-backend/repository"
	"example.com/go-backend/server"
	"example.com/go-backend/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка конфигурации: %v", err)
	}

	startupCtx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	db, err := database.Connect(startupCtx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	taskRepository := repository.NewPostgresTaskRepository(db)
	taskService := service.NewTaskService(taskRepository)
	taskHandler := handler.NewTaskHandler(taskService)
	router := handler.NewRouter(taskHandler)

	corsRouter := middleware.CORSMiddleware(router, cfg.AllowedOrigins)
	recoveredRouter := middleware.RecoveryMiddleware(corsRouter)
	loggedRouter := middleware.LoggingMiddleware(recoveredRouter)

	httpServer := server.New(loggedRouter, ":"+cfg.Port)

	signalCtx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErr := make(chan error, 1)

	go func() {
		serverErr <- httpServer.Start()
	}()

	select {
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("Ошибка HTTP-сервера: %v", err)
		}

	case <-signalCtx.Done():
		log.Println("Получен сигнал остановки")

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := httpServer.Stop(shutdownCtx); err != nil {
			log.Printf("Ошибка остановки сервера: %v", err)
		}

		log.Println("Сервер остановлен")
	}
}
