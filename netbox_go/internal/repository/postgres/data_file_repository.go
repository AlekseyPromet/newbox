// Package postgres содержит реализацию репозиториев для PostgreSQL с использованием sqlc
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	core_entity "github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
	"github.com/AlekseyPromet/netbox_go/internal/repository"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// DataFileRepositoryPostgres реализует интерфейс DataFileRepository для PostgreSQL
type DataFileRepositoryPostgres struct {
	db *sql.DB
}

// NewDataFileRepositoryPostgres создает новый экземпляр репозитория файлов данных
func NewDataFileRepositoryPostgres(db *sql.DB) *DataFileRepositoryPostgres {
	return &DataFileRepositoryPostgres{db: db}
}

// GetByID получает файл данных по ID
func (r *DataFileRepositoryPostgres) GetByID(ctx context.Context, id string) (*core_entity.DataFile, error) {
	query := `
		SELECT id, source_id, path, size, hash, data, created, updated
		FROM core_data_files
		WHERE id = $1 AND deleted_at IS NULL
	`

	var df core_entity.DataFile
	var dataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&df.ID, &df.SourceID, &df.Path, &df.Size, &df.Hash,
		&dataJSON, &df.Created, &df.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get data file by ID: %w", err)
	}

	if dataJSON != nil {
		df.Data = json.RawMessage(dataJSON)
	}

	return &df, nil
}

// List получает список файлов данных с фильтрацией
func (r *DataFileRepositoryPostgres) List(ctx context.Context, filter repository.DataFileFilter) ([]*core_entity.DataFile, int64, error) {
	query := `
		SELECT id, source_id, path, size, hash, data, created, updated
		FROM core_data_files
		WHERE deleted_at IS NULL
	`
	countQuery := `SELECT COUNT(*) FROM core_data_files WHERE deleted_at IS NULL`

	args := []interface{}{}
	argIndex := 1

	if filter.SourceID != nil {
		query += fmt.Sprintf(" AND source_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND source_id = $%d", argIndex)
		args = append(args, *filter.SourceID)
		argIndex++
	}

	if filter.Path != nil {
		query += fmt.Sprintf(" AND path = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND path = $%d", argIndex)
		args = append(args, *filter.Path)
		argIndex++
	}

	query += " ORDER BY path"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Получаем общее количество
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count data files: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list data files: %w", err)
	}
	defer rows.Close()

	var dataFiles []*core_entity.DataFile
	for rows.Next() {
		var df core_entity.DataFile
		var dataJSON []byte

		err := rows.Scan(
			&df.ID, &df.SourceID, &df.Path, &df.Size, &df.Hash,
			&dataJSON, &df.Created, &df.Updated,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan data file: %w", err)
		}

		if dataJSON != nil {
			df.Data = json.RawMessage(dataJSON)
		}

		dataFiles = append(dataFiles, &df)
	}

	return dataFiles, total, nil
}

// Create создает новый файл данных
func (r *DataFileRepositoryPostgres) Create(ctx context.Context, df *core_entity.DataFile) error {
	query := `
		INSERT INTO core_data_files (id, source_id, path, size, hash, data, created, updated)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`

	var dataJSON []byte
	if df.Data != nil {
		dataJSON = df.Data
	}

	_, err := r.db.ExecContext(ctx, query,
		df.ID.String(), df.SourceID.String(), df.Path, df.Size, df.Hash, dataJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to create data file: %w", err)
	}

	return nil
}

// Update обновляет файл данных
func (r *DataFileRepositoryPostgres) Update(ctx context.Context, df *core_entity.DataFile) error {
	query := `
		UPDATE core_data_files
		SET source_id = $1, path = $2, size = $3, hash = $4, data = $5, updated = NOW()
		WHERE id = $6
	`

	var dataJSON []byte
	if df.Data != nil {
		dataJSON = df.Data
	}

	result, err := r.db.ExecContext(ctx, query,
		df.SourceID.String(), df.Path, df.Size, df.Hash, dataJSON, df.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update data file: %w", err)
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

// Delete удаляет файл данных (soft delete)
func (r *DataFileRepositoryPostgres) Delete(ctx context.Context, id string) error {
	query := `UPDATE core_data_files SET deleted_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete data file: %w", err)
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

// GetBySourceAndPath получает файл данных по ID источника и пути
func (r *DataFileRepositoryPostgres) GetBySourceAndPath(ctx context.Context, sourceID string, path string) (*core_entity.DataFile, error) {
	query := `
		SELECT id, source_id, path, size, hash, data, created, updated
		FROM core_data_files
		WHERE source_id = $1 AND path = $2 AND deleted_at IS NULL
	`

	var df core_entity.DataFile
	var dataJSON []byte

	err := r.db.QueryRowContext(ctx, query, sourceID, path).Scan(
		&df.ID, &df.SourceID, &df.Path, &df.Size, &df.Hash,
		&dataJSON, &df.Created, &df.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get data file by source and path: %w", err)
	}

	if dataJSON != nil {
		df.Data = json.RawMessage(dataJSON)
	}

	return &df, nil
}

// DeleteBySourceID удаляет все файлы данных указанного источника (soft delete)
func (r *DataFileRepositoryPostgres) DeleteBySourceID(ctx context.Context, sourceID string) error {
	query := `UPDATE core_data_files SET deleted_at = NOW() WHERE source_id = $1`

	_, err := r.db.ExecContext(ctx, query, sourceID)
	if err != nil {
		return fmt.Errorf("failed to delete data files by source ID: %w", err)
	}

	return nil
}
