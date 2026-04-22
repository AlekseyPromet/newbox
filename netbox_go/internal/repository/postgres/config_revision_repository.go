// Package postgres содержит реализацию репозиториев для PostgreSQL с использованием sqlc
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	core_entity "netbox_go/internal/domain/core/entity"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
)

// ConfigRevisionRepositoryPostgres реализует интерфейс ConfigRevisionRepository для PostgreSQL
type ConfigRevisionRepositoryPostgres struct {
	db *sql.DB
}

// NewConfigRevisionRepositoryPostgres создает новый экземпляр репозитория ревизий конфигурации
func NewConfigRevisionRepositoryPostgres(db *sql.DB) *ConfigRevisionRepositoryPostgres {
	return &ConfigRevisionRepositoryPostgres{db: db}
}

// GetByID получает ревизию конфигурации по ID
func (r *ConfigRevisionRepositoryPostgres) GetByID(ctx context.Context, id string) (*core_entity.ConfigRevision, error) {
	query := `
		SELECT id, created, active, comment, data
		FROM core_config_revisions
		WHERE id = $1 AND deleted_at IS NULL
	`

	var revision core_entity.ConfigRevision
	var comment sql.NullString
	var dataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&revision.ID, &revision.Created, &revision.Active,
		&comment, &dataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get config revision by ID: %w", err)
	}

	if comment.Valid {
		revision.Comment = comment.String
	}
	if dataJSON != nil {
		revision.Data = json.RawMessage(dataJSON)
	}

	return &revision, nil
}

// List получает список ревизий конфигурации с фильтрацией
func (r *ConfigRevisionRepositoryPostgres) List(ctx context.Context, filter repository.ConfigRevisionFilter) ([]*core_entity.ConfigRevision, int64, error) {
	query := `
		SELECT id, created, active, comment, data
		FROM core_config_revisions
		WHERE deleted_at IS NULL
	`
	countQuery := `SELECT COUNT(*) FROM core_config_revisions WHERE deleted_at IS NULL`

	args := []interface{}{}
	argIndex := 1

	if filter.Active != nil {
		query += fmt.Sprintf(" AND active = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND active = $%d", argIndex)
		args = append(args, *filter.Active)
		argIndex++
	}

	if filter.CreatedSince != nil {
		query += fmt.Sprintf(" AND created >= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND created >= $%d", argIndex)
		args = append(args, *filter.CreatedSince)
		argIndex++
	}

	if filter.CreatedUntil != nil {
		query += fmt.Sprintf(" AND created <= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND created <= $%d", argIndex)
		args = append(args, *filter.CreatedUntil)
		argIndex++
	}

	query += " ORDER BY created DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Получаем общее количество
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count config revisions: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list config revisions: %w", err)
	}
	defer rows.Close()

	var revisions []*core_entity.ConfigRevision
	for rows.Next() {
		var revision core_entity.ConfigRevision
		var comment sql.NullString
		var dataJSON []byte

		err := rows.Scan(
			&revision.ID, &revision.Created, &revision.Active,
			&comment, &dataJSON,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan config revision: %w", err)
		}

		if comment.Valid {
			revision.Comment = comment.String
		}
		if dataJSON != nil {
			revision.Data = json.RawMessage(dataJSON)
		}

		revisions = append(revisions, &revision)
	}

	return revisions, total, nil
}

// Create создает новую ревизию конфигурации
func (r *ConfigRevisionRepositoryPostgres) Create(ctx context.Context, revision *core_entity.ConfigRevision) error {
	query := `
		INSERT INTO core_config_revisions (id, created, active, comment, data)
		VALUES ($1, $2, $3, $4, $5)
	`

	var dataJSON []byte
	if revision.Data != nil {
		dataJSON = revision.Data
	} else {
		dataJSON = []byte("{}")
	}

	_, err := r.db.ExecContext(ctx, query,
		revision.ID.String(), revision.Created, revision.Active,
		revision.Comment, dataJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to create config revision: %w", err)
	}

	return nil
}

// Activate активирует указанную ревизию конфигурации (деактивируя все остальные)
func (r *ConfigRevisionRepositoryPostgres) Activate(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Деактивируем все ревизии
	_, err = tx.ExecContext(ctx, `UPDATE core_config_revisions SET active = FALSE WHERE active = TRUE`)
	if err != nil {
		return fmt.Errorf("failed to deactivate all revisions: %w", err)
	}

	// Активируем указанную
	_, err = tx.ExecContext(ctx, `UPDATE core_config_revisions SET active = TRUE WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to activate revision: %w", err)
	}

	return tx.Commit()
}

// Delete удаляет ревизию конфигурации (soft delete)
func (r *ConfigRevisionRepositoryPostgres) Delete(ctx context.Context, id string) error {
	query := `UPDATE core_config_revisions SET deleted_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete config revision: %w", err)
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

// GetActive получает активную ревизию конфигурации
func (r *ConfigRevisionRepositoryPostgres) GetActive(ctx context.Context) (*core_entity.ConfigRevision, error) {
	query := `
		SELECT id, created, active, comment, data
		FROM core_config_revisions
		WHERE active = TRUE AND deleted_at IS NULL
		ORDER BY created DESC
		LIMIT 1
	`

	var revision core_entity.ConfigRevision
	var comment sql.NullString
	var dataJSON []byte

	err := r.db.QueryRowContext(ctx, query).Scan(
		&revision.ID, &revision.Created, &revision.Active,
		&comment, &dataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get active config revision: %w", err)
	}

	if comment.Valid {
		revision.Comment = comment.String
	}
	if dataJSON != nil {
		revision.Data = json.RawMessage(dataJSON)
	}

	return &revision, nil
}

// GetLatest получает последнюю созданную ревизию конфигурации
func (r *ConfigRevisionRepositoryPostgres) GetLatest(ctx context.Context) (*core_entity.ConfigRevision, error) {
	query := `
		SELECT id, created, active, comment, data
		FROM core_config_revisions
		WHERE deleted_at IS NULL
		ORDER BY created DESC
		LIMIT 1
	`

	var revision core_entity.ConfigRevision
	var comment sql.NullString
	var dataJSON []byte

	err := r.db.QueryRowContext(ctx, query).Scan(
		&revision.ID, &revision.Created, &revision.Active,
		&comment, &dataJSON,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest config revision: %w", err)
	}

	if comment.Valid {
		revision.Comment = comment.String
	}
	if dataJSON != nil {
		revision.Data = json.RawMessage(dataJSON)
	}

	return &revision, nil
}

// Update обновляет существующую ревизию конфигурации
func (r *ConfigRevisionRepositoryPostgres) Update(ctx context.Context, revision *core_entity.ConfigRevision) error {
	query := `
		UPDATE core_config_revisions
		SET created = $2, active = $3, comment = $4, data = $5
		WHERE id = $1 AND deleted_at IS NULL
	`

	var dataJSON []byte
	if revision.Data != nil {
		dataJSON = revision.Data
	} else {
		dataJSON = []byte("{}")
	}

	result, err := r.db.ExecContext(ctx, query,
		revision.ID.String(), revision.Created, revision.Active,
		revision.Comment, dataJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to update config revision: %w", err)
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
