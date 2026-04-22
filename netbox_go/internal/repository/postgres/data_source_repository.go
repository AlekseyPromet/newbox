// Package postgres содержит реализацию репозиториев для PostgreSQL с использованием sqlc
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	core_entity "netbox_go/internal/domain/core/entity"
	core_enum "netbox_go/internal/domain/core/enum"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
)

// DataSourceRepositoryPostgres реализует интерфейс DataSourceRepository для PostgreSQL
type DataSourceRepositoryPostgres struct {
	db *sql.DB
}

// NewDataSourceRepositoryPostgres создает новый экземпляр репозитория источников данных
func NewDataSourceRepositoryPostgres(db *sql.DB) *DataSourceRepositoryPostgres {
	return &DataSourceRepositoryPostgres{db: db}
}

// GetByID получает источник данных по ID
func (r *DataSourceRepositoryPostgres) GetByID(ctx context.Context, id string) (*core_entity.DataSource, error) {
	query := `
		SELECT id, name, type, source_url, status, enabled, sync_interval, ignore_rules,
		       parameters, last_synced, created, updated
		FROM core_data_sources
		WHERE id = $1 AND deleted_at IS NULL
	`

	var ds core_entity.DataSource
	var ignoreRulesJSON, parametersJSON []byte
	var lastSynced sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ds.ID, &ds.Name, &ds.Type, &ds.SourceURL, &ds.Status, &ds.Enabled,
		&ds.SyncInterval, &ignoreRulesJSON, &parametersJSON, &lastSynced,
		&ds.Created, &ds.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get data source by ID: %w", err)
	}

	if ignoreRulesJSON != nil {
		if err := json.Unmarshal(ignoreRulesJSON, &ds.IgnoreRules); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ignore rules: %w", err)
		}
	}
	if parametersJSON != nil {
		ds.Parameters = json.RawMessage(parametersJSON)
	}
	if lastSynced.Valid {
		ds.LastSynced = &lastSynced.Time
	}

	return &ds, nil
}

// List получает список источников данных с фильтрацией
func (r *DataSourceRepositoryPostgres) List(ctx context.Context, filter repository.DataSourceFilter) ([]*core_entity.DataSource, int64, error) {
	query := `
		SELECT id, name, type, source_url, status, enabled, sync_interval, ignore_rules,
		       parameters, last_synced, created, updated
		FROM core_data_sources
		WHERE deleted_at IS NULL
	`
	countQuery := `SELECT COUNT(*) FROM core_data_sources WHERE deleted_at IS NULL`

	args := []interface{}{}
	argIndex := 1

	if filter.Name != nil {
		query += fmt.Sprintf(" AND name = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND name = $%d", argIndex)
		args = append(args, *filter.Name)
		argIndex++
	}

	if filter.Type != nil {
		query += fmt.Sprintf(" AND type = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND type = $%d", argIndex)
		args = append(args, *filter.Type)
		argIndex++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.Enabled != nil {
		query += fmt.Sprintf(" AND enabled = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND enabled = $%d", argIndex)
		args = append(args, *filter.Enabled)
		argIndex++
	}

	if filter.SyncInterval != nil {
		query += fmt.Sprintf(" AND sync_interval = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND sync_interval = $%d", argIndex)
		args = append(args, *filter.SyncInterval)
		argIndex++
	}

	query += " ORDER BY name"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Получаем общее количество
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count data sources: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list data sources: %w", err)
	}
	defer rows.Close()

	var dataSources []*core_entity.DataSource
	for rows.Next() {
		var ds core_entity.DataSource
		var ignoreRulesJSON, parametersJSON []byte
		var lastSynced sql.NullTime

		err := rows.Scan(
			&ds.ID, &ds.Name, &ds.Type, &ds.SourceURL, &ds.Status, &ds.Enabled,
			&ds.SyncInterval, &ignoreRulesJSON, &parametersJSON, &lastSynced,
			&ds.Created, &ds.Updated,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan data source: %w", err)
		}

		if ignoreRulesJSON != nil {
			if err := json.Unmarshal(ignoreRulesJSON, &ds.IgnoreRules); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal ignore rules: %w", err)
			}
		}
		if parametersJSON != nil {
			ds.Parameters = json.RawMessage(parametersJSON)
		}
		if lastSynced.Valid {
			ds.LastSynced = &lastSynced.Time
		}

		dataSources = append(dataSources, &ds)
	}

	return dataSources, total, nil
}

// Create создает новый источник данных
func (r *DataSourceRepositoryPostgres) Create(ctx context.Context, ds *core_entity.DataSource) error {
	query := `
		INSERT INTO core_data_sources (id, name, type, source_url, status, enabled, sync_interval,
		                               ignore_rules, parameters, last_synced, created, updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
	`

	ignoreRulesJSON, err := json.Marshal(ds.IgnoreRules)
	if err != nil {
		return fmt.Errorf("failed to marshal ignore rules: %w", err)
	}

	var parametersJSON []byte
	if ds.Parameters != nil {
		parametersJSON = ds.Parameters
	}

	var lastSynced *time.Time
	if ds.LastSynced != nil {
		lastSynced = ds.LastSynced
	}

	_, err = r.db.ExecContext(ctx, query,
		ds.ID.String(), ds.Name, ds.Type, ds.SourceURL, ds.Status, ds.Enabled,
		ds.SyncInterval, ignoreRulesJSON, parametersJSON, lastSynced,
	)
	if err != nil {
		return fmt.Errorf("failed to create data source: %w", err)
	}

	return nil
}

// Update обновляет источник данных
func (r *DataSourceRepositoryPostgres) Update(ctx context.Context, ds *core_entity.DataSource) error {
	query := `
		UPDATE core_data_sources
		SET name = $1, type = $2, source_url = $3, status = $4, enabled = $5,
		    sync_interval = $6, ignore_rules = $7, parameters = $8, last_synced = $9, updated = NOW()
		WHERE id = $10
	`

	ignoreRulesJSON, err := json.Marshal(ds.IgnoreRules)
	if err != nil {
		return fmt.Errorf("failed to marshal ignore rules: %w", err)
	}

	var parametersJSON []byte
	if ds.Parameters != nil {
		parametersJSON = ds.Parameters
	}

	_, err = r.db.ExecContext(ctx, query,
		ds.Name, ds.Type, ds.SourceURL, ds.Status, ds.Enabled,
		ds.SyncInterval, ignoreRulesJSON, parametersJSON, ds.LastSynced, ds.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update data source: %w", err)
	}

	return nil
}

// Delete удаляет источник данных (soft delete)
func (r *DataSourceRepositoryPostgres) Delete(ctx context.Context, id string) error {
	query := `UPDATE core_data_sources SET deleted_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete data source: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return types.ErrNotFound
	}

	return nil
}

// UpdateStatus обновляет статус и время последней синхронизации источника данных
func (r *DataSourceRepositoryPostgres) UpdateStatus(ctx context.Context, id string, status string, lastSynced *time.Time) error {
	query := `
		UPDATE core_data_sources
		SET status = $1, last_synced = COALESCE($2, last_synced), updated = NOW()
		WHERE id = $3
	`

	result, err := r.db.ExecContext(ctx, query, status, lastSynced, id)
	if err != nil {
		return fmt.Errorf("failed to update data source status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return types.ErrNotFound
	}

	return nil
}

// Sync помечает источник данных как синхронизируемый
func (r *DataSourceRepositoryPostgres) Sync(ctx context.Context, id string) error {
	query := `
		UPDATE core_data_sources
		SET status = $1, updated = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, core_enum.DataSourceStatusSyncing, id)
	if err != nil {
		return fmt.Errorf("failed to sync data source: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return types.ErrNotFound
	}

	return nil
}

// Exists проверяет существование источника данных по имени
func (r *DataSourceRepositoryPostgres) Exists(ctx context.Context, name string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM core_data_sources WHERE name = $1 AND deleted_at IS NULL)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, name).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check data source existence: %w", err)
	}

	return exists, nil
}

// GetByName получает источник данных по имени
func (r *DataSourceRepositoryPostgres) GetByName(ctx context.Context, name string) (*core_entity.DataSource, error) {
	query := `
		SELECT id, name, type, source_url, status, enabled, sync_interval, ignore_rules,
		       parameters, last_synced, created, updated
		FROM core_data_sources
		WHERE name = $1 AND deleted_at IS NULL
	`

	var ds core_entity.DataSource
	var ignoreRulesJSON, parametersJSON []byte
	var lastSynced sql.NullTime

	err := r.db.QueryRowContext(ctx, query, name).Scan(
		&ds.ID, &ds.Name, &ds.Type, &ds.SourceURL, &ds.Status, &ds.Enabled,
		&ds.SyncInterval, &ignoreRulesJSON, &parametersJSON, &lastSynced,
		&ds.Created, &ds.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get data source by name: %w", err)
	}

	if ignoreRulesJSON != nil {
		if err := json.Unmarshal(ignoreRulesJSON, &ds.IgnoreRules); err != nil {
			return nil, fmt.Errorf("failed to unmarshal ignore rules: %w", err)
		}
	}
	if parametersJSON != nil {
		ds.Parameters = json.RawMessage(parametersJSON)
	}
	if lastSynced.Valid {
		ds.LastSynced = &lastSynced.Time
	}

	return &ds, nil
}
