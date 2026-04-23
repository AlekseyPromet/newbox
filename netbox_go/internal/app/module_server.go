package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

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
			cfg := GetConfig()
			port := cfg.Server.Port
			if port == 0 {
				port = 8080
			}

			log.Printf("Starting HTTP server on port %d", port)
			go func() {
				if err := e.Start(fmt.Sprintf(":%d", port)); err != nil && err != http.ErrServerClosed {
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
