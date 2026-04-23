package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/domain/core/repository"
	coredb "netbox_go/internal/infrastructure/storage/sqlc/core"
	"netbox_go/pkg/types"

	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
)

// ObjectChangePostgresRepository реализует ObjectChangeRepository для PostgreSQL
type ObjectChangePostgresRepository struct {
	db *sql.DB
}

// NewObjectChangePostgresRepository создаёт новый экземпляр репозитория
func NewObjectChangePostgresRepository(db *sql.DB) repository.ObjectChangeRepository {
	return &ObjectChangePostgresRepository{db: db}
}

// GetByID возвращает запись об изменении по ID
func (r *ObjectChangePostgresRepository) GetByID(ctx context.Context, id types.ID) (*entity.ObjectChange, error) {
	q := coredb.Queries{DB: r.db}
	row, err := q.GetObjectChangeByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &entity.ObjectChange{
		ID:                row.ID,
		Time:              row.Time,
		UserID:            types.ID(row.UserID.UUID),
		RequestID:         &row.RequestID.String,
		Action:            types.Status(row.Action),
		ChangedObjectType: row.ChangedObjectType,
		ChangedObjectID:   row.ChangedObjectID,
		ObjectRepr:        row.ObjectRepr,
		ObjectData:        row.ObjectData.RawMessage,
		RelatedObjectType: &row.RelatedObjectType.String,
		RelatedObjectID:   &row.RelatedObjectID.String,
		RelatedObjectRepr: &row.RelatedObjectRepr.String,
	}, nil
}

// List возвращает список изменений с фильтрацией
func (r *ObjectChangePostgresRepository) List(ctx context.Context, filter repository.ObjectChangeFilter, limit, offset int) ([]*entity.ObjectChange, int, error) {
	var query string
	var args []interface{}

	query = `SELECT id, time, user_id, request_id, action, changed_object_type, changed_object_id, object_repr, object_data, related_object_type, related_object_id, related_object_repr FROM core_objectchange WHERE 1=1`
	countQuery := `SELECT COUNT(*)::int FROM core_objectchange WHERE 1=1`

	if filter.UserID != nil {
		query += ` AND user_id = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.UserID)
		countQuery += ` AND user_id = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.User != nil && *filter.User != "" {
		// This requires a join with users table or a subquery, but for simplicity we use the user_name if available in the table
		// In the Python code, it filters by user__username.
		query += ` AND user_id IN (SELECT id FROM users WHERE username ILIKE $` + fmt.Sprintf("%d", len(args)+1) + `)`
		args = append(args, "%"+*filter.User+"%")
		countQuery += ` AND user_id IN (SELECT id FROM users WHERE username ILIKE $` + fmt.Sprintf("%d", len(args)) + `)`
	}

	if filter.Action != nil && *filter.Action != "" {
		query += ` AND action = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.Action)
		countQuery += ` AND action = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.ChangedObjectType != nil && *filter.ChangedObjectType != "" {
		query += ` AND changed_object_type = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.ChangedObjectType)
		countQuery += ` AND changed_object_type = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.ChangedObjectID != nil && *filter.ChangedObjectID != "" {
		query += ` AND changed_object_id = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.ChangedObjectID)
		countQuery += ` AND changed_object_id = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.TimeAfter != nil {
		query += ` AND time >= $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.TimeAfter)
		countQuery += ` AND time >= $` + fmt.Sprintf("%d", len(args))
	}

	if filter.TimeBefore != nil {
		query += ` AND time <= $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.TimeBefore)
		countQuery += ` AND time <= $` + fmt.Sprintf("%d", len(args))
	}

	if filter.SearchQuery != nil && *filter.SearchQuery != "" {
		searchVal := "%" + *filter.SearchQuery + "%"
		query += ` AND (object_repr ILIKE $` + fmt.Sprintf("%d", len(args)+1) + ` OR object_data::text ILIKE $` + fmt.Sprintf("%d", len(args)+1) + `)`
		args = append(args, searchVal)
		countQuery += ` AND (object_repr ILIKE $` + fmt.Sprintf("%d", len(args)) + ` OR object_data::text ILIKE $` + fmt.Sprintf("%d", len(args)) + `)`
	}

	query += ` ORDER BY time DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	rows, err := r.db.QueryContext(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*entity.ObjectChange
	for rows.Next() {
		var row struct {
			ID                types.ID
			Time              time.Time
			UserID            types.ID
			RequestID         sql.NullString
			Action            string
			ChangedObjectType string
			ChangedObjectID   string
			ObjectRepr        sql.NullString
			ObjectData        []byte
			RelatedObjectType sql.NullString
			RelatedObjectID   sql.NullString
			RelatedObjectRepr sql.NullString
		}
		if err := rows.Scan(&row.ID, &row.Time, &row.UserID, &row.RequestID, &row.Action, &row.ChangedObjectType, &row.ChangedObjectID, &row.ObjectRepr, &row.ObjectData, &row.RelatedObjectType, &row.RelatedObjectID, &row.RelatedObjectRepr); err != nil {
			return nil, 0, err
		}
		result = append(result, &entity.ObjectChange{
			ID:                row.ID,
			Time:              row.Time,
			UserID:            row.UserID,
			RequestID:         new(row.RequestID.String),
			Action:            types.Status(row.Action),
			ChangedObjectType: row.ChangedObjectType,
			ChangedObjectID:   row.ChangedObjectID,
			ObjectRepr:        row.ObjectRepr.String,
			ObjectData:        row.ObjectData,
			RelatedObjectType: new(row.RelatedObjectType.String),
			RelatedObjectID:   new(row.RelatedObjectID.String),
			RelatedObjectRepr: new(row.RelatedObjectRepr.String),
		})
	}

	var count int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return result, count, nil
}

func getNullUuid(s types.ID) uuid.NullUUID {
	u, err := uuid.Parse(s.String())
	if err != nil {
		return uuid.NullUUID{}
	}
	return uuid.NullUUID{UUID: u, Valid: true}
}

// Create создаёт запись об изменении
func (r *ObjectChangePostgresRepository) Create(ctx context.Context, change *entity.ObjectChange) error {
	q := coredb.Queries{DB: r.db}

	row, err := q.CreateObjectChange(ctx, coredb.CreateObjectChangeParams{
		Time:              time.Now(),
		UserID:            getNullUuid(change.UserID),
		RequestID:         sql.NullString{String: valOrEmpty(change.RequestID), Valid: change.RequestID != nil},
		Action:            string(change.Action),
		ChangedObjectType: change.ChangedObjectType,
		ChangedObjectID:   change.ChangedObjectID,
		ObjectRepr:        change.ObjectRepr,
		ObjectData:        pqtype.NullRawMessage{RawMessage: change.ObjectData},
		RelatedObjectType: sql.NullString{String: valOrEmpty(change.RelatedObjectType), Valid: change.RelatedObjectType != nil},
		RelatedObjectID:   sql.NullString{String: valOrEmpty(change.RelatedObjectID), Valid: change.RelatedObjectID != nil},
		RelatedObjectRepr: sql.NullString{String: valOrEmpty(change.RelatedObjectRepr), Valid: change.RelatedObjectRepr != nil},
	})
	if err != nil {
		return err
	}

	change.ID = row.ID
	change.Time = row.Time
	return nil
}

// BulkCreate создаёт несколько записей об изменениях
func (r *ObjectChangePostgresRepository) BulkCreate(ctx context.Context, changes []*entity.ObjectChange) error {
	q := coredb.Queries{DB: r.db}

	now := time.Now()

	for i, change := range changes {

		params := coredb.BulkCreateObjectChangesParams{
			Time:              now,
			UserID:            getNullUuid(change.UserID),
			RequestID:         sql.NullString{String: valOrEmpty(change.RequestID), Valid: change.RequestID != nil},
			Action:            string(change.Action),
			ChangedObjectType: change.ChangedObjectType,
			ChangedObjectID:   change.ChangedObjectID,
			ObjectRepr:        change.ObjectRepr,
			ObjectData:        pqtype.NullRawMessage{RawMessage: change.ObjectData},
			RelatedObjectType: sql.NullString{String: valOrEmpty(change.RelatedObjectType), Valid: change.RelatedObjectType != nil},
			RelatedObjectID:   sql.NullString{String: valOrEmpty(change.RelatedObjectID), Valid: change.RelatedObjectID != nil},
			RelatedObjectRepr: sql.NullString{String: valOrEmpty(change.RelatedObjectRepr), Valid: change.RelatedObjectRepr != nil},
		}

		err := q.BulkCreateObjectChanges(ctx, params)
		if err != nil {
			return fmt.Errorf("%d %v", i, err)
		}
	}

	return nil
}

// DeleteOld удаляет старые записи (старше cutoffTime)
func (r *ObjectChangePostgresRepository) DeleteOld(ctx context.Context, cutoffTime time.Time) (int64, error) {
	q := coredb.Queries{DB: r.db}
	result, err := q.DeleteOldObjectChanges(ctx, cutoffTime)
	if err != nil {
		return 0, err
	}
	return result, nil
}
