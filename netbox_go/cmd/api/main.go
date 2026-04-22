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
	"github.com/AlekseyPromet/netbox_go/internal/delivery/http/middleware"
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
	accountRepo := postgres.NewAccountRepositoryPostgres(db)

	// Core репозитории реализованы в postgres
	dataSourceRepo := postgres.NewDataSourceRepositoryPostgres(db)
	dataFileRepo := postgres.NewDataFileRepositoryPostgres(db)
	jobRepo := postgres.NewJobRepositoryPostgres(db)
	objectChangeRepo := postgres.NewObjectChangeRepositoryPostgres(db)
	objectTypeRepo := postgres.NewObjectTypeRepositoryPostgres(db)
	configRevRepo := postgres.NewConfigRevisionRepositoryPostgres(db)

	// Создание обработчиков
	siteHandler := handlers.NewSiteHandler(siteRepo)
	accountHandler := handlers.NewAccountHandler(
		accountRepo, accountRepo, accountRepo, accountRepo, accountRepo,
	)
	coreHandler := handlers.NewCoreHandlers(
		dataSourceRepo, dataFileRepo, jobRepo, objectChangeRepo, objectTypeRepo, configRevRepo,
	)

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

	// API Routes
	api := e.Group("/api")

	// Core
	core := api.Group("/core")
	dataSources := core.Group("/data-sources")
	dataSources.GET("", coreHandler.ListDataSources)
	dataSources.GET("/:id", coreHandler.GetDataSource)
	dataSources.POST("", coreHandler.CreateDataSource)
	dataSources.PUT("/:id", coreHandler.UpdateDataSource)
	dataSources.DELETE("/:id", coreHandler.DeleteDataSource)
	dataSources.POST("/:id/sync", coreHandler.SyncDataSource)

	dataFiles := core.Group("/data-files")
	dataFiles.GET("", coreHandler.ListDataFiles)
	dataFiles.GET("/:id", coreHandler.GetDataFile)
	dataFiles.POST("", coreHandler.CreateDataFile)
	dataFiles.PUT("/:id", coreHandler.UpdateDataFile)
	dataFiles.DELETE("/:id", coreHandler.DeleteDataFile)

	jobs := core.Group("/jobs")
	jobs.GET("", coreHandler.ListJobs)
	jobs.GET("/:id", coreHandler.GetJob)
	jobs.POST("", coreHandler.CreateJob)

	objectChanges := core.Group("/object-changes")
	objectChanges.GET("", coreHandler.ListObjectChanges)
	objectChanges.GET("/:id", coreHandler.GetObjectChange)
	objectChanges.POST("/log", coreHandler.LogObjectChange)

	objectTypes := core.Group("/object-types")
	objectTypes.GET("", coreHandler.ListObjectTypes)
	objectTypes.GET("/:id", coreHandler.GetObjectType)

	configRevisions := core.Group("/config-revisions")
	configRevisions.GET("", coreHandler.ListConfigRevisions)
	configRevisions.GET("/:id", coreHandler.GetConfigRevision)
	configRevisions.POST("", coreHandler.CreateConfigRevision)
	configRevisions.PUT("/:id", coreHandler.UpdateConfigRevision)
	configRevisions.DELETE("/:id", coreHandler.DeleteConfigRevision)
	configRevisions.POST("/:id/activate", coreHandler.ActivateConfigRevision)
	configRevisions.GET("/active", coreHandler.GetActiveConfigRevision)

	bgQueues := core.Group("/background-queues")
	bgQueues.GET("", coreHandler.ListBackgroundQueues)
	bgQueues.GET("/:name", coreHandler.GetBackgroundQueue)

	bgWorkers := core.Group("/background-workers")
	bgWorkers.GET("", coreHandler.ListBackgroundWorkers)
	bgWorkers.GET("/:name", coreHandler.GetBackgroundWorker)

	bgTasks := core.Group("/background-tasks")
	bgTasks.GET("", coreHandler.ListBackgroundTasks)
	bgTasks.GET("/:id", coreHandler.GetBackgroundTask)
	bgTasks.POST("/:id/delete", coreHandler.DeleteBackgroundTask)
	bgTasks.POST("/:id/requeue", coreHandler.RequeueBackgroundTask)
	bgTasks.POST("/:id/enqueue", coreHandler.EnqueueBackgroundTask)
	bgTasks.POST("/:id/stop", coreHandler.StopBackgroundTask)

	// DCIM
	dcim := api.Group("/dcim")

	sites := dcim.Group("/sites")
	sites.GET("", siteHandler.List)
	sites.GET("/:id", siteHandler.GetByID)
	sites.GET("/slug/:slug", siteHandler.GetBySlug)
	sites.POST("", siteHandler.Create)
	sites.PUT("/:id", siteHandler.Update)
	sites.DELETE("/:id", siteHandler.Delete)

	account := api.Group("/account", middleware.KerberosSSOMiddleware())
	account.GET("/profile", accountHandler.Profile)
	account.GET("/bookmarks", accountHandler.ListBookmarks)
	account.GET("/notifications", accountHandler.ListNotifications)
	account.GET("/subscriptions", accountHandler.ListSubscriptions)
	account.GET("/preferences", accountHandler.GetPreferences)
	account.POST("/preferences", accountHandler.UpsertPreferences)

	tokens := account.Group("/api-tokens")
	tokens.GET("", accountHandler.ListTokens)
	tokens.GET("/:id", accountHandler.GetToken)
	tokens.POST("", accountHandler.CreateToken)
	tokens.PUT("/:id", accountHandler.UpdateToken)
	tokens.DELETE("/:id", accountHandler.DeleteToken)

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
