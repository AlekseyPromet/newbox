// Package postgres содержит реализацию репозиториев для PostgreSQL
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"netbox_go/internal/domain/extras/entity"
	"netbox_go/pkg/types"
)

// DashboardRepositoryPostgres реализует интерфейс DashboardRepository для PostgreSQL
type DashboardRepositoryPostgres struct {
	db *sql.DB
}

// NewDashboardRepositoryPostgres создает новый экземпляр репозитория панелей управления
func NewDashboardRepositoryPostgres(db *sql.DB) *DashboardRepositoryPostgres {
	return &DashboardRepositoryPostgres{db: db}
}

// GetByUserID получает панель управления пользователя по ID пользователя
func (r *DashboardRepositoryPostgres) GetByUserID(ctx context.Context, userID int64) (*entity.Dashboard, error) {
	query := `
		SELECT id, user_id, layout, config, created, updated
		FROM extras_dashboard
		WHERE user_id = $1
	`

	var dashboard entity.Dashboard
	var layoutJSON, configJSON []byte

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&dashboard.ID,
		&dashboard.UserID,
		&layoutJSON,
		&configJSON,
		&dashboard.Created,
		&dashboard.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrDashboardNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard by user ID: %w", err)
	}

	// Парсинг JSON layout
	if err := json.Unmarshal(layoutJSON, &dashboard.Layout); err != nil {
		return nil, fmt.Errorf("failed to parse layout JSON: %w", err)
	}

	// Парсинг JSON config
	if err := json.Unmarshal(configJSON, &dashboard.Config); err != nil {
		return nil, fmt.Errorf("failed to parse config JSON: %w", err)
	}

	return &dashboard, nil
}

// Create создает новую панель управления
func (r *DashboardRepositoryPostgres) Create(ctx context.Context, dashboard *entity.Dashboard) error {
	query := `
		INSERT INTO extras_dashboard (user_id, layout, config, created, updated)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`

	layoutJSON, err := json.Marshal(dashboard.Layout)
	if err != nil {
		return fmt.Errorf("failed to marshal layout: %w", err)
	}

	configJSON, err := json.Marshal(dashboard.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	var userID interface{}
	if dashboard.UserID != nil {
		userID = *dashboard.UserID
	} else {
		userID = nil
	}

	err = r.db.QueryRowContext(ctx, query, userID, layoutJSON, configJSON).Scan(&dashboard.ID)
	if err != nil {
		return fmt.Errorf("failed to create dashboard: %w", err)
	}

	dashboard.Created = time.Now()
	dashboard.Updated = time.Now()

	return nil
}

// Update обновляет существующую панель управления
func (r *DashboardRepositoryPostgres) Update(ctx context.Context, dashboard *entity.Dashboard) error {
	query := `
		UPDATE extras_dashboard
		SET layout = $2, config = $3, updated = NOW()
		WHERE id = $1
	`

	layoutJSON, err := json.Marshal(dashboard.Layout)
	if err != nil {
		return fmt.Errorf("failed to marshal layout: %w", err)
	}

	configJSON, err := json.Marshal(dashboard.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, dashboard.ID, layoutJSON, configJSON)
	if err != nil {
		return fmt.Errorf("failed to update dashboard: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return types.ErrDashboardNotFound
	}

	dashboard.Updated = time.Now()

	return nil
}

// Delete удаляет панель управления
func (r *DashboardRepositoryPostgres) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM extras_dashboard
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete dashboard: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return types.ErrDashboardNotFound
	}

	return nil
}
