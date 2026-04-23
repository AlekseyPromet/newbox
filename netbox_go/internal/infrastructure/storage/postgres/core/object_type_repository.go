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

	"github.com/sqlc-dev/pqtype"
)

// ObjectTypePostgresRepository реализует ObjectTypeRepository для PostgreSQL
type ObjectTypePostgresRepository struct {
	db *sql.DB
}

// NewObjectTypePostgresRepository создаёт новый экземпляр репозитория
func NewObjectTypePostgresRepository(db *sql.DB) repository.ObjectTypeRepository {
	return &ObjectTypePostgresRepository{db: db}
}

// GetByID возвращает тип объекта по ID
func (r *ObjectTypePostgresRepository) GetByID(ctx context.Context, id types.ID) (*entity.ObjectType, error) {
	q := coredb.Queries{DB: r.db}
	row, err := q.GetObjectTypeByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	var features []string
	if row.Features.Valid && len(row.Features.RawMessage) > 0 {
		if err := types.UnmarshalJSON(row.Features.RawMessage, &features); err != nil {
			features = []string{}
		}
	} else {
		features = []string{}
	}

	return &entity.ObjectType{
		ID:       row.ID,
		AppLabel: row.AppLabel,
		Model:    row.Model,
		Public:   row.Public,
		Features: features,
		Created:  row.Created,
		Updated:  row.Updated,
	}, nil
}

// GetByAppAndModel возвращает тип объекта по app_label и model
func (r *ObjectTypePostgresRepository) GetByAppAndModel(ctx context.Context, appLabel, model string) (*entity.ObjectType, error) {
	q := coredb.Queries{DB: r.db}
	row, err := q.GetObjectTypeByAppAndModel(ctx, coredb.GetObjectTypeByAppAndModelParams{
		AppLabel: appLabel,
		Model:    model,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	var features []string
	if row.Features.Valid && len(row.Features.RawMessage) > 0 {
		if err := types.UnmarshalJSON(row.Features.RawMessage, &features); err != nil {
			features = []string{}
		}
	} else {
		features = []string{}
	}

	return &entity.ObjectType{
		ID:       row.ID,
		AppLabel: row.AppLabel,
		Model:    row.Model,
		Public:   row.Public,
		Features: features,
		Created:  row.Created,
		Updated:  row.Updated,
	}, nil
}

// List возвращает список типов объектов с фильтрацией
func (r *ObjectTypePostgresRepository) List(ctx context.Context, filter repository.ObjectTypeFilter, limit, offset int) ([]*entity.ObjectType, int, error) {
	var query string
	var args []interface{}

	query = `SELECT id, app_label, model, public, features, created, updated FROM django_content_type WHERE 1=1`
	countQuery := `SELECT COUNT(*)::int FROM django_content_type WHERE 1=1`

	if filter.AppLabel != nil && *filter.AppLabel != "" {
		query += ` AND app_label = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.AppLabel)
		countQuery += ` AND app_label = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.Model != nil && *filter.Model != "" {
		query += ` AND model = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.Model)
		countQuery += ` AND model = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.Public != nil {
		query += ` AND public = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.Public)
		countQuery += ` AND public = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.Features != nil && *filter.Features != "" {
		val := "%" + *filter.Features + "%"
		query += ` AND features::text ILIKE $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, val)
		countQuery += ` AND features::text ILIKE $` + fmt.Sprintf("%d", len(args))
	}

	if filter.SearchQuery != nil && *filter.SearchQuery != "" {
		searchVal := "%" + *filter.SearchQuery + "%"
		query += ` AND (app_label ILIKE $` + fmt.Sprintf("%d", len(args)+1) + ` OR model ILIKE $` + fmt.Sprintf("%d", len(args)+1) + `)`
		args = append(args, searchVal)
		countQuery += ` AND (app_label ILIKE $` + fmt.Sprintf("%d", len(args)) + ` OR model ILIKE $` + fmt.Sprintf("%d", len(args)) + `)`
	}

	query += ` ORDER BY app_label, model LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	rows, err := r.db.QueryContext(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*entity.ObjectType
	for rows.Next() {
		var row struct {
			ID       types.ID
			AppLabel string
			Model    string
			Public   bool
			Features []byte
			Created  time.Time
			Updated  time.Time
		}
		if err := rows.Scan(&row.ID, &row.AppLabel, &row.Model, &row.Public, &row.Features, &row.Created, &row.Updated); err != nil {
			return nil, 0, err
		}

		var features []string
		if row.Features != nil {
			types.UnmarshalJSON(row.Features, &features)
		}

		result = append(result, &entity.ObjectType{
			ID:       row.ID,
			AppLabel: row.AppLabel,
			Model:    row.Model,
			Public:   row.Public,
			Features: features,
			Created:  row.Created,
			Updated:  row.Updated,
		})
	}

	var count int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return result, count, nil
}

// Create создаёт новый тип объекта
func (r *ObjectTypePostgresRepository) Create(ctx context.Context, ot *entity.ObjectType) error {
	q := &coredb.Queries{DB: r.db}

	features, err := types.MarshalJSON(ot.Features)
	if err != nil {
		return err
	}

	now := time.Now()
	row, err := q.CreateObjectType(ctx, coredb.CreateObjectTypeParams{
		AppLabel: ot.AppLabel,
		Model:    ot.Model,
		Public:   ot.Public,
		Features: pqtype.NullRawMessage{
			RawMessage: features,
			Valid:      true,
		},
		Created: now,
		Updated: now,
	})
	if err != nil {
		return err
	}

	ot.ID = row.ID
	ot.Created = row.Created
	ot.Updated = row.Updated
	return nil
}

// Update обновляет тип объекта
func (r *ObjectTypePostgresRepository) Update(ctx context.Context, ot *entity.ObjectType) error {
	q := &coredb.Queries{DB: r.db}

	features, err := types.MarshalJSON(ot.Features)
	if err != nil {
		return err
	}

	_, err = q.UpdateObjectType(ctx, coredb.UpdateObjectTypeParams{
		ID:       ot.ID,
		AppLabel: ot.AppLabel,
		Model:    ot.Model,
		Public:   ot.Public,
		Features: pqtype.NullRawMessage{
			RawMessage: features,
			Valid:      true,
		},
		Updated: time.Now(),
	})
	return err
}

// Delete удаляет тип объекта
func (r *ObjectTypePostgresRepository) Delete(ctx context.Context, id types.ID) error {
	q := coredb.Queries{DB: r.db}
	_, err := q.DeleteObjectType(ctx, id)
	return err
}
