// Package main содержит точку входа REST API сервера
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"

	"github.com/AlekseyPromet/netbox_go/internal/delivery/http/handlers"
	"github.com/AlekseyPromet/netbox_go/internal/repository/postgres"
)

func main() {
	// Инициализация базы данных
	db, err := initDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Создание репозиториев
	siteRepo := postgres.NewSiteRepositoryPostgres(db)

	// Создание обработчиков
	siteHandler := handlers.NewSiteHandler(siteRepo)

	// Инициализация Echo
	e := echo.New()

	// Middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(echomiddleware.RequestID())

	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API Routes - DCIM
	api := e.Group("/api")
	dcim := api.Group("/dcim")
	
	sites := dcim.Group("/sites")
	sites.GET("", siteHandler.List)
	sites.GET("/:id", siteHandler.GetByID)
	sites.GET("/slug/:slug", siteHandler.GetBySlug)
	sites.POST("", siteHandler.Create)
	sites.PUT("/:id", siteHandler.Update)
	sites.DELETE("/:id", siteHandler.Delete)

	// Запуск сервера
	go func() {
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		
		log.Printf("Starting HTTP server on port %s", port)
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}

func initDatabase() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "postgres://netbox:netbox@localhost:5432/netbox?sslmode=disable"
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Настройка пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
