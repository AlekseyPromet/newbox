// Package postgres содержит PostgreSQL реализации репозиториев домена Core
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
	q := coredb.Queries{DB: r.db}
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
	q := coredb.Queries{DB: r.db}
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

// List возвращает список ревизий с фильтрацией и пагинацией
func (r *ConfigRevisionPostgresRepository) List(ctx context.Context, filter repository.ConfigRevisionFilter, limit, offset int) ([]*entity.ConfigRevision, int, error) {
	var query string
	var args []interface{}

	query = `SELECT id, created, active, comment, data FROM core_configrevision WHERE 1=1`
	countQuery := `SELECT COUNT(*)::int FROM core_configrevision WHERE 1=1`

	if filter.Comment != nil && *filter.Comment != "" {
		val := "%" + *filter.Comment + "%"
		query += ` AND comment ILIKE $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, val)
		countQuery += ` AND comment ILIKE $` + fmt.Sprintf("%d", len(args))
	}

	if filter.CreatedAfter != nil {
		query += ` AND created >= $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.CreatedAfter)
		countQuery += ` AND created >= $` + fmt.Sprintf("%d", len(args))
	}

	if filter.CreatedBefore != nil {
		query += ` AND created <= $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.CreatedBefore)
		countQuery += ` AND created <= $` + fmt.Sprintf("%d", len(args))
	}

	if filter.SearchQuery != nil && *filter.SearchQuery != "" {
		searchVal := "%" + *filter.SearchQuery + "%"
		query += ` AND comment ILIKE $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, searchVal)
		countQuery += ` AND comment ILIKE $` + fmt.Sprintf("%d", len(args))
	}

	query += ` ORDER BY created DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	rows, err := r.db.QueryContext(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*entity.ConfigRevision
	for rows.Next() {
		var row struct {
			ID      types.ID
			Created time.Time
			Active  bool
			Comment sql.NullString
			Data    []byte
		}
		if err := rows.Scan(&row.ID, &row.Created, &row.Active, &row.Comment, &row.Data); err != nil {
			return nil, 0, err
		}
		result = append(result, &entity.ConfigRevision{
			ID:      row.ID,
			Created: row.Created,
			Active:  row.Active,
			Comment: row.Comment.String,
			Data:    row.Data,
		})
	}

	var count int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return result, count, nil
}

// Create создаёт новую ревизию
func (r *ConfigRevisionPostgresRepository) Create(ctx context.Context, revision *entity.ConfigRevision) error {
	q := coredb.Queries{DB: r.db}

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
	q := coredb.Queries{DB: r.db}

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
	q := coredb.Queries{DB: r.db}
	_, err := q.DeleteConfigRevision(ctx, id)
	return err
}

// SetActive делает ревизию активной
func (r *ConfigRevisionPostgresRepository) SetActive(ctx context.Context, id types.ID) error {
	q := coredb.Queries{DB: r.db}
	_, err := q.SetActiveConfigRevision(ctx, id)
	return err
}
