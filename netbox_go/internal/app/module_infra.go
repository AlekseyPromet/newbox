package app

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"go.uber.org/fx"
)

// ModuleInfra provides infrastructure dependencies like DB and Echo
var ModuleInfra = fx.Options(
	fx.Provide(
		NewDatabase,
		NewEcho,
	),
)

func NewDatabase() (*sql.DB, error) {
	cfg := GetConfig()

	var dsn string
	if cfg.Database.URL != "" {
		dsn = cfg.Database.URL
	} else {
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)
	}

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func NewEcho() *echo.Echo {
	e := echo.New()

	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())
	e.Use(echomiddleware.RequestID())

	return e
}
