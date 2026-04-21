package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
	"github.com/AlekseyPromet/netbox_go/internal/domain/core/repository"
	coredb "github.com/AlekseyPromet/netbox_go/internal/infrastructure/storage/sqlc/core"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
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
	q := coredb.New(r.db)
	row, err := q.GetObjectChangeByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &entity.ObjectChange{
		ID:                 row.ID,
		Time:               row.Time,
		UserID:             row.UserID,
		RequestID:          row.RequestID.String,
		Action:             types.Status(row.Action),
		ChangedObjectType:  row.ChangedObjectType,
		ChangedObjectID:    row.ChangedObjectID,
		ObjectRepr:         row.ObjectRepr,
		ObjectData:         row.ObjectData,
		RelatedObjectType:  row.RelatedObjectType.String,
		RelatedObjectID:    row.RelatedObjectID.String,
		RelatedObjectRepr:  row.RelatedObjectRepr.String,
	}, nil
}

// List возвращает список изменений с фильтрацией
func (r *ObjectChangePostgresRepository) List(ctx context.Context, userID *types.ID, action types.Status, objectType, objectID string, timeFrom, timeTo time.Time, limit, offset int) ([]*entity.ObjectChange, int, error) {
	q := coredb.New(r.db)

	var timeFromPtr, timeToPtr sql.NullTime
	if !timeFrom.IsZero() {
		timeFromPtr = sql.NullTime{Time: timeFrom, Valid: true}
	}
	if !timeTo.IsZero() {
		timeToPtr = sql.NullTime{Time: timeTo, Valid: true}
	}

	rows, err := q.ListObjectChanges(ctx, coredb.ListObjectChangesParams{
		UserID:              userID,
		Action:              string(action),
		ChangedObjectType:   objectType,
		ChangedObjectID:     objectID,
		TimeFrom:            timeFromPtr,
		TimeTo:              timeToPtr,
		Limit:               int32(limit),
		Offset:              int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	countRow, err := q.CountObjectChanges(ctx, coredb.CountObjectChangesParams{
		UserID:            userID,
		Action:            string(action),
		ChangedObjectType: objectType,
		ChangedObjectID:   objectID,
		TimeFrom:          timeFromPtr,
		TimeTo:            timeToPtr,
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.ObjectChange, len(rows))
	for i, row := range rows {
		result[i] = &entity.ObjectChange{
			ID:                 row.ID,
			Time:               row.Time,
			UserID:             row.UserID,
			RequestID:          row.RequestID.String,
			Action:             types.Status(row.Action),
			ChangedObjectType:  row.ChangedObjectType,
			ChangedObjectID:    row.ChangedObjectID,
			ObjectRepr:         row.ObjectRepr,
			ObjectData:         row.ObjectData,
			RelatedObjectType:  row.RelatedObjectType.String,
			RelatedObjectID:    row.RelatedObjectID.String,
			RelatedObjectRepr:  row.RelatedObjectRepr.String,
		}
	}

	return result, int(countRow.Count), nil
}

// Create создаёт запись об изменении
func (r *ObjectChangePostgresRepository) Create(ctx context.Context, change *entity.ObjectChange) error {
	q := coredb.New(r.db)

	data := change.ObjectData
	if data == nil {
		data = []byte("{}")
	}

	row, err := q.CreateObjectChange(ctx, coredb.CreateObjectChangeParams{
		Time:               time.Now(),
		UserID:             change.UserID,
		RequestID:          sql.NullString{String: valOrEmpty(change.RequestID), Valid: change.RequestID != nil},
		Action:             string(change.Action),
		ChangedObjectType:  change.ChangedObjectType,
		ChangedObjectID:    change.ChangedObjectID,
		ObjectRepr:         change.ObjectRepr,
		ObjectData:         data,
		RelatedObjectType:  sql.NullString{String: valOrEmpty(change.RelatedObjectType), Valid: change.RelatedObjectType != nil},
		RelatedObjectID:    sql.NullString{String: valOrEmpty(change.RelatedObjectID), Valid: change.RelatedObjectID != nil},
		RelatedObjectRepr:  sql.NullString{String: valOrEmpty(change.RelatedObjectRepr), Valid: change.RelatedObjectRepr != nil},
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
	q := coredb.New(r.db)

	params := make([]coredb.BulkCreateObjectChangesParams, len(changes))
	now := time.Now()

	for i, change := range changes {
		data := change.ObjectData
		if data == nil {
			data = []byte("{}")
		}

		params[i] = coredb.BulkCreateObjectChangesParams{
			Time:               now,
			UserID:             change.UserID,
			RequestID:          sql.NullString{String: valOrEmpty(change.RequestID), Valid: change.RequestID != nil},
			Action:             string(change.Action),
			ChangedObjectType:  change.ChangedObjectType,
			ChangedObjectID:    change.ChangedObjectID,
			ObjectRepr:         change.ObjectRepr,
			ObjectData:         data,
			RelatedObjectType:  sql.NullString{String: valOrEmpty(change.RelatedObjectType), Valid: change.RelatedObjectType != nil},
			RelatedObjectID:    sql.NullString{String: valOrEmpty(change.RelatedObjectID), Valid: change.RelatedObjectID != nil},
			RelatedObjectRepr:  sql.NullString{String: valOrEmpty(change.RelatedObjectRepr), Valid: change.RelatedObjectRepr != nil},
		}
	}

	return q.BulkCreateObjectChanges(ctx, params)
}

// DeleteOld удаляет старые записи (старше cutoffTime)
func (r *ObjectChangePostgresRepository) DeleteOld(ctx context.Context, cutoffTime time.Time) (int64, error) {
	q := coredb.New(r.db)
	result, err := q.DeleteOldObjectChanges(ctx, cutoffTime)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected(), nil
}

func valOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
