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

// ObjectTypeRepositoryPostgres реализует интерфейс ObjectTypeRepository для PostgreSQL
type ObjectTypeRepositoryPostgres struct {
	db *sql.DB
}

// NewObjectTypeRepositoryPostgres создает новый экземпляр репозитория типов объектов
func NewObjectTypeRepositoryPostgres(db *sql.DB) *ObjectTypeRepositoryPostgres {
	return &ObjectTypeRepositoryPostgres{db: db}
}

// GetByID получает тип объекта по ID
func (r *ObjectTypeRepositoryPostgres) GetByID(ctx context.Context, id string) (*core_entity.ObjectType, error) {
	query := `
		SELECT id, app_label, model, public, features, created, updated
		FROM core_object_types
		WHERE id = $1 AND deleted_at IS NULL
	`

	var ot core_entity.ObjectType
	var featuresJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ot.ID, &ot.AppLabel, &ot.Model, &ot.Public,
		&featuresJSON, &ot.Created, &ot.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get object type by ID: %w", err)
	}

	if featuresJSON != nil {
		if err := json.Unmarshal(featuresJSON, &ot.Features); err != nil {
			return nil, fmt.Errorf("failed to unmarshal features: %w", err)
		}
	}

	return &ot, nil
}

// List получает список типов объектов с фильтрацией
func (r *ObjectTypeRepositoryPostgres) List(ctx context.Context, filter repository.ObjectTypeFilter) ([]*core_entity.ObjectType, int64, error) {
	query := `
		SELECT id, app_label, model, public, features, created, updated
		FROM core_object_types
		WHERE deleted_at IS NULL
	`
	countQuery := `SELECT COUNT(*) FROM core_object_types WHERE deleted_at IS NULL`

	args := []interface{}{}
	argIndex := 1

	if filter.AppLabel != nil {
		query += fmt.Sprintf(" AND app_label = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND app_label = $%d", argIndex)
		args = append(args, *filter.AppLabel)
		argIndex++
	}

	if filter.Model != nil {
		query += fmt.Sprintf(" AND model = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND model = $%d", argIndex)
		args = append(args, *filter.Model)
		argIndex++
	}

	if filter.Public != nil {
		query += fmt.Sprintf(" AND public = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND public = $%d", argIndex)
		args = append(args, *filter.Public)
		argIndex++
	}

	if filter.Feature != nil {
		query += fmt.Sprintf(" AND $%d = ANY(features)", argIndex)
		countQuery += fmt.Sprintf(" AND $%d = ANY(features)", argIndex)
		args = append(args, *filter.Feature)
		argIndex++
	}

	query += " ORDER BY app_label, model"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Получаем общее количество
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count object types: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list object types: %w", err)
	}
	defer rows.Close()

	var objectTypes []*core_entity.ObjectType
	for rows.Next() {
		var ot core_entity.ObjectType
		var featuresJSON []byte

		err := rows.Scan(
			&ot.ID, &ot.AppLabel, &ot.Model, &ot.Public,
			&featuresJSON, &ot.Created, &ot.Updated,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan object type: %w", err)
		}

		if featuresJSON != nil {
			if err := json.Unmarshal(featuresJSON, &ot.Features); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal features: %w", err)
			}
		}

		objectTypes = append(objectTypes, &ot)
	}

	return objectTypes, total, nil
}

// Create создает новый тип объекта
func (r *ObjectTypeRepositoryPostgres) Create(ctx context.Context, ot *core_entity.ObjectType) error {
	query := `
		INSERT INTO core_object_types (id, app_label, model, public, features, created, updated)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`

	featuresJSON, err := json.Marshal(ot.Features)
	if err != nil {
		return fmt.Errorf("failed to marshal features: %w", err)
	}

	_, err = r.db.ExecContext(ctx, query,
		ot.ID.String(), ot.AppLabel, ot.Model, ot.Public, featuresJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to create object type: %w", err)
	}

	return nil
}

// Update обновляет тип объекта
func (r *ObjectTypeRepositoryPostgres) Update(ctx context.Context, ot *core_entity.ObjectType) error {
	query := `
		UPDATE core_object_types
		SET app_label = $1, model = $2, public = $3, features = $4, updated = NOW()
		WHERE id = $5
	`

	featuresJSON, err := json.Marshal(ot.Features)
	if err != nil {
		return fmt.Errorf("failed to marshal features: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query,
		ot.AppLabel, ot.Model, ot.Public, featuresJSON, ot.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update object type: %w", err)
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

// Delete удаляет тип объекта (soft delete)
func (r *ObjectTypeRepositoryPostgres) Delete(ctx context.Context, id string) error {
	query := `UPDATE core_object_types SET deleted_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete object type: %w", err)
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

// GetByAppAndModel получает тип объекта по app_label и model
func (r *ObjectTypeRepositoryPostgres) GetByAppAndModel(ctx context.Context, appLabel string, model string) (*core_entity.ObjectType, error) {
	query := `
		SELECT id, app_label, model, public, features, created, updated
		FROM core_object_types
		WHERE app_label = $1 AND model = $2 AND deleted_at IS NULL
	`

	var ot core_entity.ObjectType
	var featuresJSON []byte

	err := r.db.QueryRowContext(ctx, query, appLabel, model).Scan(
		&ot.ID, &ot.AppLabel, &ot.Model, &ot.Public,
		&featuresJSON, &ot.Created, &ot.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get object type by app and model: %w", err)
	}

	if featuresJSON != nil {
		if err := json.Unmarshal(featuresJSON, &ot.Features); err != nil {
			return nil, fmt.Errorf("failed to unmarshal features: %w", err)
		}
	}

	return &ot, nil
}

// GetForModel получает все типы объектов для данной модели
func (r *ObjectTypeRepositoryPostgres) GetForModel(ctx context.Context, model string) ([]*core_entity.ObjectType, error) {
	query := `
		SELECT id, app_label, model, public, features, created, updated
		FROM core_object_types
		WHERE model = $1 AND deleted_at IS NULL
		ORDER BY app_label
	`

	rows, err := r.db.QueryContext(ctx, query, model)
	if err != nil {
		return nil, fmt.Errorf("failed to list object types for model: %w", err)
	}
	defer rows.Close()

	var objectTypes []*core_entity.ObjectType
	for rows.Next() {
		var ot core_entity.ObjectType
		var featuresJSON []byte

		err := rows.Scan(
			&ot.ID, &ot.AppLabel, &ot.Model, &ot.Public,
			&featuresJSON, &ot.Created, &ot.Updated,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan object type: %w", err)
		}

		if featuresJSON != nil {
			if err := json.Unmarshal(featuresJSON, &ot.Features); err != nil {
				return nil, fmt.Errorf("failed to unmarshal features: %w", err)
			}
		}

		objectTypes = append(objectTypes, &ot)
	}

	return objectTypes, nil
}

// Public возвращает список всех публичных типов объектов (их ID)
func (r *ObjectTypeRepositoryPostgres) Public() []string {
	// В полной реализации здесь был бы запрос к БД
	// Для MVP возвращаем заглушку - в реальности нужно кэшировать или запрашивать
	return []string{}
}

// WithFeature возвращает список ID типов объектов с указанной фичей
func (r *ObjectTypeRepositoryPostgres) WithFeature(feature string) []string {
	// В полной реализации здесь был бы запрос к БД с фильтрацией по features array
	// Для MVP возвращаем заглушку
	return []string{}
}
