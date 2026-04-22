// Package postgres содержит реализацию репозиториев для PostgreSQL с использованием sqlc
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	core_entity "netbox_go/internal/domain/core/entity"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
)

// ObjectChangeRepositoryPostgres реализует интерфейс ObjectChangeRepository для PostgreSQL
type ObjectChangeRepositoryPostgres struct {
	db *sql.DB
}

// NewObjectChangeRepositoryPostgres создает новый экземпляр репозитория изменений объектов
func NewObjectChangeRepositoryPostgres(db *sql.DB) *ObjectChangeRepositoryPostgres {
	return &ObjectChangeRepositoryPostgres{db: db}
}

// GetByID получает запись об изменении объекта по ID
func (r *ObjectChangeRepositoryPostgres) GetByID(ctx context.Context, id string) (*core_entity.ObjectChange, error) {
	query := `
		SELECT id, time, user_id, request_id, action, changed_object_type, changed_object_id,
		       object_repr, object_data, related_object_type, related_object_id, related_object_repr
		FROM core_object_changes
		WHERE id = $1 AND deleted_at IS NULL
	`

	var oc core_entity.ObjectChange
	var userID sql.NullString
	var requestID sql.NullString
	var relatedObjectType, relatedObjectID, relatedObjectRepr sql.NullString
	var objectDataJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&oc.ID, &oc.Time, &userID, &requestID, &oc.Action,
		&oc.ChangedObjectType, &oc.ChangedObjectID, &oc.ObjectRepr,
		&objectDataJSON, &relatedObjectType, &relatedObjectID, &relatedObjectRepr,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get object change by ID: %w", err)
	}

	// Заполнение опциональных полей
	if userID.Valid {
		uid, _ := types.ParseID(userID.String)
		oc.UserID = &uid
	}
	if requestID.Valid {
		oc.RequestID = &requestID.String
	}
	if objectDataJSON != nil {
		oc.ObjectData = json.RawMessage(objectDataJSON)
	}
	if relatedObjectType.Valid {
		oc.RelatedObjectType = &relatedObjectType.String
	}
	if relatedObjectID.Valid {
		oc.RelatedObjectID = &relatedObjectID.String
	}
	if relatedObjectRepr.Valid {
		oc.RelatedObjectRepr = &relatedObjectRepr.String
	}

	return &oc, nil
}

// List получает список записей об изменениях с фильтрацией
func (r *ObjectChangeRepositoryPostgres) List(ctx context.Context, filter repository.ObjectChangeFilter) ([]*core_entity.ObjectChange, int64, error) {
	query := `
		SELECT id, time, user_id, request_id, action, changed_object_type, changed_object_id,
		       object_repr, object_data, related_object_type, related_object_id, related_object_repr
		FROM core_object_changes
		WHERE deleted_at IS NULL
	`
	countQuery := `SELECT COUNT(*) FROM core_object_changes WHERE deleted_at IS NULL`

	args := []interface{}{}
	argIndex := 1

	if filter.ChangedObjectType != nil {
		query += fmt.Sprintf(" AND changed_object_type = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND changed_object_type = $%d", argIndex)
		args = append(args, *filter.ChangedObjectType)
		argIndex++
	}

	if filter.ChangedObjectID != nil {
		query += fmt.Sprintf(" AND changed_object_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND changed_object_id = $%d", argIndex)
		args = append(args, *filter.ChangedObjectID)
		argIndex++
	}

	if filter.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *filter.UserID)
		argIndex++
	}

	if filter.Action != nil {
		query += fmt.Sprintf(" AND action = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND action = $%d", argIndex)
		args = append(args, *filter.Action)
		argIndex++
	}

	if filter.RequestID != nil {
		query += fmt.Sprintf(" AND request_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND request_id = $%d", argIndex)
		args = append(args, *filter.RequestID)
		argIndex++
	}

	if filter.Since != nil {
		query += fmt.Sprintf(" AND time >= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND time >= $%d", argIndex)
		args = append(args, *filter.Since)
		argIndex++
	}

	if filter.Until != nil {
		query += fmt.Sprintf(" AND time <= $%d", argIndex)
		countQuery += fmt.Sprintf(" AND time <= $%d", argIndex)
		args = append(args, *filter.Until)
		argIndex++
	}

	query += " ORDER BY time DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Получаем общее количество
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count object changes: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list object changes: %w", err)
	}
	defer rows.Close()

	var changes []*core_entity.ObjectChange
	for rows.Next() {
		var oc core_entity.ObjectChange
		var userID sql.NullString
		var requestID sql.NullString
		var relatedObjectType, relatedObjectID, relatedObjectRepr sql.NullString
		var objectDataJSON []byte

		err := rows.Scan(
			&oc.ID, &oc.Time, &userID, &requestID, &oc.Action,
			&oc.ChangedObjectType, &oc.ChangedObjectID, &oc.ObjectRepr,
			&objectDataJSON, &relatedObjectType, &relatedObjectID, &relatedObjectRepr,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan object change: %w", err)
		}

		if userID.Valid {
			uid, _ := types.ParseID(userID.String)
			oc.UserID = &uid
		}
		if requestID.Valid {
			oc.RequestID = &requestID.String
		}
		if objectDataJSON != nil {
			oc.ObjectData = json.RawMessage(objectDataJSON)
		}
		if relatedObjectType.Valid {
			oc.RelatedObjectType = &relatedObjectType.String
		}
		if relatedObjectID.Valid {
			oc.RelatedObjectID = &relatedObjectID.String
		}
		if relatedObjectRepr.Valid {
			oc.RelatedObjectRepr = &relatedObjectRepr.String
		}

		changes = append(changes, &oc)
	}

	return changes, total, nil
}

// Create создает новую запись об изменении объекта
func (r *ObjectChangeRepositoryPostgres) Create(ctx context.Context, change *core_entity.ObjectChange) error {
	query := `
		INSERT INTO core_object_changes (id, time, user_id, request_id, action, changed_object_type, 
		                                 changed_object_id, object_repr, object_data, 
		                                 related_object_type, related_object_id, related_object_repr)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	var userID *string
	if change.UserID != nil {
		uid := change.UserID.String()
		userID = &uid
	}

	var dataJSON []byte
	if change.ObjectData != nil {
		dataJSON = change.ObjectData
	}

	var relatedObjectType, relatedObjectID, relatedObjectRepr *string
	if change.RelatedObjectType != nil {
		relatedObjectType = change.RelatedObjectType
	}
	if change.RelatedObjectID != nil {
		relatedObjectID = change.RelatedObjectID
	}
	if change.RelatedObjectRepr != nil {
		relatedObjectRepr = change.RelatedObjectRepr
	}

	_, err := r.db.ExecContext(ctx, query,
		change.ID.String(), change.Time, userID, change.RequestID, change.Action,
		change.ChangedObjectType, change.ChangedObjectID, change.ObjectRepr, dataJSON,
		relatedObjectType, relatedObjectID, relatedObjectRepr,
	)
	if err != nil {
		return fmt.Errorf("failed to create object change: %w", err)
	}

	return nil
}

// LogChange удобный метод для логирования изменения объекта
func (r *ObjectChangeRepositoryPostgres) LogChange(ctx context.Context, action types.Status, objectType string, objectID string, objectRepr string, objectData interface{}, userID *types.ID, requestID *string) error {
	var dataJSON []byte
	var err error

	if objectData != nil {
		dataJSON, err = json.Marshal(objectData)
		if err != nil {
			return fmt.Errorf("failed to marshal object data: %w", err)
		}
	}

	change := &core_entity.ObjectChange{
		ID:                types.NewID(),
		Time:              time.Now(),
		UserID:            userID,
		RequestID:         requestID,
		Action:            action,
		ChangedObjectType: objectType,
		ChangedObjectID:   objectID,
		ObjectRepr:        objectRepr,
		ObjectData:        json.RawMessage(dataJSON),
	}

	return r.Create(ctx, change)
}

// GetChangesForObject получает историю изменений конкретного объекта
func (r *ObjectChangeRepositoryPostgres) GetChangesForObject(ctx context.Context, objectType string, objectID string, limit int, offset int) ([]*core_entity.ObjectChange, int64, error) {
	query := `
		SELECT id, time, user_id, request_id, action, changed_object_type, changed_object_id,
		       object_repr, object_data, related_object_type, related_object_id, related_object_repr
		FROM core_object_changes
		WHERE changed_object_type = $1 AND changed_object_id = $2 AND deleted_at IS NULL
	`
	countQuery := `SELECT COUNT(*) FROM core_object_changes WHERE changed_object_type = $1 AND changed_object_id = $2 AND deleted_at IS NULL`

	// Получаем общее количество
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, objectType, objectID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count object changes: %w", err)
	}

	query += " ORDER BY time DESC"
	query += fmt.Sprintf(" LIMIT $3 OFFSET $4")

	rows, err := r.db.QueryContext(ctx, query, objectType, objectID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list object changes: %w", err)
	}
	defer rows.Close()

	var changes []*core_entity.ObjectChange
	for rows.Next() {
		var oc core_entity.ObjectChange
		var userID sql.NullString
		var requestID sql.NullString
		var relatedObjectType, relatedObjectID, relatedObjectRepr sql.NullString
		var objectDataJSON []byte

		err := rows.Scan(
			&oc.ID, &oc.Time, &userID, &requestID, &oc.Action,
			&oc.ChangedObjectType, &oc.ChangedObjectID, &oc.ObjectRepr,
			&objectDataJSON, &relatedObjectType, &relatedObjectID, &relatedObjectRepr,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan object change: %w", err)
		}

		if userID.Valid {
			uid, _ := types.ParseID(userID.String)
			oc.UserID = &uid
		}
		if requestID.Valid {
			oc.RequestID = &requestID.String
		}
		if objectDataJSON != nil {
			oc.ObjectData = json.RawMessage(objectDataJSON)
		}
		if relatedObjectType.Valid {
			oc.RelatedObjectType = &relatedObjectType.String
		}
		if relatedObjectID.Valid {
			oc.RelatedObjectID = &relatedObjectID.String
		}
		if relatedObjectRepr.Valid {
			oc.RelatedObjectRepr = &relatedObjectRepr.String
		}

		changes = append(changes, &oc)
	}

	return changes, total, nil
}
