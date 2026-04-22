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

// DataSourcePostgresRepository реализует DataSourceRepository для PostgreSQL
type DataSourcePostgresRepository struct {
	db *sql.DB
}

// NewDataSourcePostgresRepository создаёт новый экземпляр репозитория
func NewDataSourcePostgresRepository(db *sql.DB) repository.DataSourceRepository {
	return &DataSourcePostgresRepository{db: db}
}

// GetByID возвращает источник данных по ID
func (r *DataSourcePostgresRepository) GetByID(ctx context.Context, id types.ID) (*entity.DataSource, error) {
	q := coredb.New(r.db)
	row, err := q.GetDataSourceByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	var ignoreRules []string
	if row.IgnoreRules != nil && len(row.IgnoreRules) > 0 {
		if err := types.UnmarshalJSON(row.IgnoreRules, &ignoreRules); err != nil {
			ignoreRules = []string{}
		}
	} else {
		ignoreRules = []string{}
	}

	var lastSynced *time.Time
	if row.LastSynced.Valid {
		lastSynced = &row.LastSynced.Time
	}

	return &entity.DataSource{
		ID:           row.ID,
		Name:         row.Name,
		Type:         row.Type,
		SourceURL:    row.SourceURL,
		Status:       types.Status(row.Status),
		Enabled:      row.Enabled,
		SyncInterval: int(row.SyncInterval),
		IgnoreRules:  ignoreRules,
		Parameters:   row.Parameters,
		LastSynced:   lastSynced,
		Created:      row.Created,
		Updated:      row.Updated,
	}, nil
}

// GetByName возвращает источник данных по имени
func (r *DataSourcePostgresRepository) GetByName(ctx context.Context, name string) (*entity.DataSource, error) {
	q := coredb.New(r.db)
	row, err := q.GetDataSourceByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	var ignoreRules []string
	if row.IgnoreRules != nil && len(row.IgnoreRules) > 0 {
		if err := types.UnmarshalJSON(row.IgnoreRules, &ignoreRules); err != nil {
			ignoreRules = []string{}
		}
	} else {
		ignoreRules = []string{}
	}

	var lastSynced *time.Time
	if row.LastSynced.Valid {
		lastSynced = &row.LastSynced.Time
	}

	return &entity.DataSource{
		ID:           row.ID,
		Name:         row.Name,
		Type:         row.Type,
		SourceURL:    row.SourceURL,
		Status:       types.Status(row.Status),
		Enabled:      row.Enabled,
		SyncInterval: int(row.SyncInterval),
		IgnoreRules:  ignoreRules,
		Parameters:   row.Parameters,
		LastSynced:   lastSynced,
		Created:      row.Created,
		Updated:      row.Updated,
	}, nil
}

// List возвращает список источников данных с фильтрацией
func (r *DataSourcePostgresRepository) List(ctx context.Context, status types.Status, enabled *bool, typeFilter string, limit, offset int) ([]*entity.DataSource, int, error) {
	q := coredb.New(r.db)

	rows, err := q.ListDataSources(ctx, coredb.ListDataSourcesParams{
		Status:   string(status),
		Enabled:  enabled,
		Type_:    typeFilter,
		Limit:    int32(limit),
		Offset:   int32(offset),
	})
	if err != nil {
		return nil, 0, err
	}

	countRow, err := q.CountDataSources(ctx, coredb.CountDataSourcesParams{
		Status:  string(status),
		Enabled: enabled,
		Type_:   typeFilter,
	})
	if err != nil {
		return nil, 0, err
	}

	result := make([]*entity.DataSource, len(rows))
	for i, row := range rows {
		var ignoreRules []string
		if row.IgnoreRules != nil && len(row.IgnoreRules) > 0 {
			if err := types.UnmarshalJSON(row.IgnoreRules, &ignoreRules); err != nil {
				ignoreRules = []string{}
			}
		} else {
			ignoreRules = []string{}
		}

		var lastSynced *time.Time
		if row.LastSynced.Valid {
			lastSynced = &row.LastSynced.Time
		}

		result[i] = &entity.DataSource{
			ID:           row.ID,
			Name:         row.Name,
			Type:         row.Type,
			SourceURL:    row.SourceURL,
			Status:       types.Status(row.Status),
			Enabled:      row.Enabled,
			SyncInterval: int(row.SyncInterval),
			IgnoreRules:  ignoreRules,
			Parameters:   row.Parameters,
			LastSynced:   lastSynced,
			Created:      row.Created,
			Updated:      row.Updated,
		}
	}

	return result, int(countRow.Count), nil
}

// Create создаёт новый источник данных
func (r *DataSourcePostgresRepository) Create(ctx context.Context, ds *entity.DataSource) error {
	q := coredb.New(r.db)

	ignoreRules, err := types.MarshalJSON(ds.IgnoreRules)
	if err != nil {
		return err
	}

	parameters := ds.Parameters
	if parameters == nil {
		parameters = []byte("{}")
	}

	now := time.Now()
	var lastSynced sql.NullTime
	if ds.LastSynced != nil {
		lastSynced = sql.NullTime{Time: *ds.LastSynced, Valid: true}
	}

	row, err := q.CreateDataSource(ctx, coredb.CreateDataSourceParams{
		Name:         ds.Name,
		Type:         ds.Type,
		SourceURL:    ds.SourceURL,
		Status:       string(ds.Status),
		Enabled:      ds.Enabled,
		SyncInterval: int32(ds.SyncInterval),
		IgnoreRules:  ignoreRules,
		Parameters:   parameters,
		LastSynced:   lastSynced,
		Created:      now,
		Updated:      now,
	})
	if err != nil {
		return err
	}

	ds.ID = row.ID
	ds.Created = row.Created
	ds.Updated = row.Updated
	return nil
}

// Update обновляет источник данных
func (r *DataSourcePostgresRepository) Update(ctx context.Context, ds *entity.DataSource) error {
	q := coredb.New(r.db)

	ignoreRules, err := types.MarshalJSON(ds.IgnoreRules)
	if err != nil {
		return err
	}

	parameters := ds.Parameters
	if parameters == nil {
		parameters = []byte("{}")
	}

	var lastSynced sql.NullTime
	if ds.LastSynced != nil {
		lastSynced = sql.NullTime{Time: *ds.LastSynced, Valid: true}
	}

	_, err = q.UpdateDataSource(ctx, coredb.UpdateDataSourceParams{
		ID:           ds.ID,
		Name:         ds.Name,
		Type:         ds.Type,
		SourceURL:    ds.SourceURL,
		Status:       string(ds.Status),
		Enabled:      ds.Enabled,
		SyncInterval: int32(ds.SyncInterval),
		IgnoreRules:  ignoreRules,
		Parameters:   parameters,
		LastSynced:   lastSynced,
		Updated:      time.Now(),
	})
	return err
}

// Delete удаляет источник данных
func (r *DataSourcePostgresRepository) Delete(ctx context.Context, id types.ID) error {
	q := coredb.New(r.db)
	_, err := q.DeleteDataSource(ctx, id)
	return err
}

// UpdateStatus обновляет статус источника
func (r *DataSourcePostgresRepository) UpdateStatus(ctx context.Context, id types.ID, status types.Status, lastSynced *time.Time) error {
	q := coredb.New(r.db)

	var lastSyncedNull sql.NullTime
	if lastSynced != nil {
		lastSyncedNull = sql.NullTime{Time: *lastSynced, Valid: true}
	}

	return q.UpdateDataSourceStatus(ctx, coredb.UpdateDataSourceStatusParams{
		ID:         id,
		Status:     string(status),
		LastSynced: lastSyncedNull,
	})
}

// GetQueuedForSync возвращает источники, ожидающие синхронизации
func (r *DataSourcePostgresRepository) GetQueuedForSync(ctx context.Context, limit int) ([]*entity.DataSource, error) {
	q := coredb.New(r.db)

	rows, err := q.GetQueuedDataSources(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	result := make([]*entity.DataSource, len(rows))
	for i, row := range rows {
		var ignoreRules []string
		if row.IgnoreRules != nil && len(row.IgnoreRules) > 0 {
			if err := types.UnmarshalJSON(row.IgnoreRules, &ignoreRules); err != nil {
				ignoreRules = []string{}
			}
		} else {
			ignoreRules = []string{}
		}

		var lastSynced *time.Time
		if row.LastSynced.Valid {
			lastSynced = &row.LastSynced.Time
		}

		result[i] = &entity.DataSource{
			ID:           row.ID,
			Name:         row.Name,
			Type:         row.Type,
			SourceURL:    row.SourceURL,
			Status:       types.Status(row.Status),
			Enabled:      row.Enabled,
			SyncInterval: int(row.SyncInterval),
			IgnoreRules:  ignoreRules,
			Parameters:   row.Parameters,
			LastSynced:   lastSynced,
			Created:      row.Created,
			Updated:      row.Updated,
		}
	}

	return result, nil
}
