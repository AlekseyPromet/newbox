package app

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

// ModuleServer provides server lifecycle management
var ModuleServer = fx.Options(
	fx.Invoke(StartServer),
)

func StartServer(lc fx.Lifecycle, e *echo.Echo) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			port := os.Getenv("PORT")
			if port == "" {
				port = "8080"
			}

			log.Printf("Starting HTTP server on port %s", port)
			go func() {
				if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
					log.Fatalf("HTTP server failed: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down server...")
			if err := e.Shutdown(ctx); err != nil {
				log.Printf("Server forced to shutdown: %v", err)
				return err
			}
			log.Println("Server exited gracefully")
			return nil
		},
	})
}
