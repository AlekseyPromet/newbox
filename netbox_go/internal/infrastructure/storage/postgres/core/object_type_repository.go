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
	q := coredb.New(r.db)
	row, err := q.GetObjectTypeByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	var features []string
	if row.Features != nil && len(row.Features) > 0 {
		if err := types.UnmarshalJSON(row.Features, &features); err != nil {
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
	q := coredb.New(r.db)
	row, err := q.GetObjectTypeByAppAndModel(ctx, appLabel, model)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	var features []string
	if row.Features != nil && len(row.Features) > 0 {
		if err := types.UnmarshalJSON(row.Features, &features); err != nil {
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
func (r *ObjectTypePostgresRepository) List(ctx context.Context, appLabel, model string, public *bool, limit, offset int) ([]*entity.ObjectType, int, error) {
	q := coredb.New(r.db)

	rows, err := q.ListObjectTypes(ctx, coredb.ListObjectTypesParams{
		AppLabel: appLabel,
		Model:    model,
		Public:   public,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	countRow, err := q.CountObjectTypes(ctx, coredb.CountObjectTypesParams{
		AppLabel: appLabel,
		Model:    model,
		Public:   public,
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.ObjectType, len(rows))
	for i, row := range rows {
		var features []string
		if row.Features != nil && len(row.Features) > 0 {
			if err := types.UnmarshalJSON(row.Features, &features); err != nil {
				features = []string{}
			}
		} else {
			features = []string{}
		}

		result[i] = &entity.ObjectType{
			ID:       row.ID,
			AppLabel: row.AppLabel,
			Model:    row.Model,
			Public:   row.Public,
			Features: features,
			Created:  row.Created,
			Updated:  row.Updated,
		}
	}

	return result, int(countRow.Count), nil
}

// Create создаёт новый тип объекта
func (r *ObjectTypePostgresRepository) Create(ctx context.Context, ot *entity.ObjectType) error {
	q := coredb.New(r.db)

	features, err := types.MarshalJSON(ot.Features)
	if err != nil {
		return err
	}

	now := time.Now()
	row, err := q.CreateObjectType(ctx, coredb.CreateObjectTypeParams{
		AppLabel: ot.AppLabel,
		Model:    ot.Model,
		Public:   ot.Public,
		Features: features,
		Created:  now,
		Updated:  now,
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
	q := coredb.New(r.db)

	features, err := types.MarshalJSON(ot.Features)
	if err != nil {
		return err
	}

	_, err = q.UpdateObjectType(ctx, coredb.UpdateObjectTypeParams{
		ID:       ot.ID,
		AppLabel: ot.AppLabel,
		Model:    ot.Model,
		Public:   ot.Public,
		Features: features,
		Updated:  time.Now(),
	})
	return err
}

// Delete удаляет тип объекта
func (r *ObjectTypePostgresRepository) Delete(ctx context.Context, id types.ID) error {
	q := coredb.New(r.db)
	_, err := q.DeleteObjectType(ctx, id)
	return err
}
