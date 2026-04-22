// Package postgres содержит PostgreSQL реализации репозиториев домена Core
package postgres

import (
	"context"
	"database/sql"
	"time"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/domain/core/repository"
	coredb "netbox_go/internal/infrastructure/storage/sqlc/core"
	"netbox_go/pkg/types"
)

// ConfigRevisionPostgresRepository реализует ConfigRevisionRepository для PostgreSQL
type ConfigRevisionPostgresRepository struct {
	db *sql.DB
}

// NewConfigRevisionPostgresRepository создаёт новый экземпляр репозитория
func NewConfigRevisionPostgresRepository(db *sql.DB) repository.ConfigRevisionRepository {
	return &ConfigRevisionPostgresRepository{db: db}
}

// GetByID возвращает ревизию по ID
func (r *ConfigRevisionPostgresRepository) GetByID(ctx context.Context, id types.ID) (*entity.ConfigRevision, error) {
	q := coredb.New(r.db)
	row, err := q.GetConfigRevisionByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &entity.ConfigRevision{
		ID:      row.ID,
		Created: row.Created,
		Active:  row.Active,
		Comment: row.Comment.String,
		Data:    row.Data,
	}, nil
}

// GetActive возвращает активную ревизию
func (r *ConfigRevisionPostgresRepository) GetActive(ctx context.Context) (*entity.ConfigRevision, error) {
	q := coredb.New(r.db)
	row, err := q.GetActiveConfigRevision(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &entity.ConfigRevision{
		ID:      row.ID,
		Created: row.Created,
		Active:  row.Active,
		Comment: row.Comment.String,
		Data:    row.Data,
	}, nil
}

// List возвращает список ревизий с пагинацией
func (r *ConfigRevisionPostgresRepository) List(ctx context.Context, limit, offset int) ([]*entity.ConfigRevision, int, error) {
	q := coredb.New(r.db)
	rows, err := q.ListConfigRevisions(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, 0, err
	}

	countRow, err := q.CountConfigRevisions(ctx)
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.ConfigRevision, len(rows))
	for i, row := range rows {
		result[i] = &entity.ConfigRevision{
			ID:      row.ID,
			Created: row.Created,
			Active:  row.Active,
			Comment: row.Comment.String,
			Data:    row.Data,
		}
	}

	return result, int(countRow.Count), nil
}

// Create создаёт новую ревизию
func (r *ConfigRevisionPostgresRepository) Create(ctx context.Context, revision *entity.ConfigRevision) error {
	q := coredb.New(r.db)
	
	data := revision.Data
	if data == nil {
		data = []byte("{}")
	}

	row, err := q.CreateConfigRevision(ctx, coredb.CreateConfigRevisionParams{
		Created: time.Now(),
		Active:  revision.Active,
		Comment: sql.NullString{String: revision.Comment, Valid: revision.Comment != ""},
		Data:    data,
	})
	if err != nil {
		return err
	}

	revision.ID = row.ID
	revision.Created = row.Created
	return nil
}

// Update обновляет ревизию
func (r *ConfigRevisionPostgresRepository) Update(ctx context.Context, revision *entity.ConfigRevision) error {
	q := coredb.New(r.db)

	data := revision.Data
	if data == nil {
		data = []byte("{}")
	}

	_, err := q.UpdateConfigRevision(ctx, coredb.UpdateConfigRevisionParams{
		ID:      revision.ID,
		Active:  revision.Active,
		Comment: sql.NullString{String: revision.Comment, Valid: revision.Comment != ""},
		Data:    data,
	})
	return err
}

// Delete удаляет ревизию
func (r *ConfigRevisionPostgresRepository) Delete(ctx context.Context, id types.ID) error {
	q := coredb.New(r.db)
	_, err := q.DeleteConfigRevision(ctx, id)
	return err
}

// SetActive делает ревизию активной
func (r *ConfigRevisionPostgresRepository) SetActive(ctx context.Context, id types.ID) error {
	q := coredb.New(r.db)
	return q.SetActiveConfigRevision(ctx, id)
}
