package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/domain/core/repository"
	coredb "netbox_go/internal/infrastructure/storage/sqlc/core"
	"netbox_go/pkg/types"

	"github.com/lib/pq"
	"github.com/sqlc-dev/pqtype"
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
	q := coredb.Queries{DB: r.db}
	row, err := q.GetDataSourceByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	var ignoreRules []string
	if row.IgnoreRules.Valid && len(row.IgnoreRules.RawMessage) > 0 {
		if err := types.UnmarshalJSON(row.IgnoreRules.RawMessage, &ignoreRules); err != nil {
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
		SourceURL:    row.SourceUrl,
		Status:       types.Status(row.Status),
		Enabled:      row.Enabled,
		SyncInterval: int(row.SyncInterval),
		IgnoreRules:  ignoreRules,
		Parameters:   row.Parameters.RawMessage,
		LastSynced:   lastSynced,
		Created:      row.Created,
		Updated:      row.Updated,
	}, nil
}

// GetByName возвращает источник данных по имени
func (r *DataSourcePostgresRepository) GetByName(ctx context.Context, name string) (*entity.DataSource, error) {
	q := coredb.Queries{DB: r.db}
	row, err := q.GetDataSourceByName(ctx, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	var ignoreRules []string
	if row.IgnoreRules.Valid && len(row.IgnoreRules.RawMessage) > 0 {
		if err := types.UnmarshalJSON(row.IgnoreRules.RawMessage, &ignoreRules); err != nil {
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
		SourceURL:    row.SourceUrl,
		Status:       types.Status(row.Status),
		Enabled:      row.Enabled,
		SyncInterval: int(row.SyncInterval),
		IgnoreRules:  ignoreRules,
		Parameters:   row.Parameters.RawMessage,
		LastSynced:   lastSynced,
		Created:      row.Created,
		Updated:      row.Updated,
	}, nil
}

// BulkCreate создаёт несколько источников данных.
// Использует транзакцию для атомарности.
func (r *DataSourcePostgresRepository) BulkCreate(ctx context.Context, data []entity.DataSource) error {
	if len(data) == 0 {
		return nil
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO core_datasource
			(name, type, source_url, status, enabled, sync_interval,
			 ignore_rules, parameters, last_synced, created, updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, ds := range data {
		var ignoreRules, parameters []byte
		if len(ds.IgnoreRules) > 0 {
			ignoreRules, _ = json.Marshal(ds.IgnoreRules)
		}
		if len(ds.Parameters) > 0 {
			parameters = ds.Parameters
		}

		_, err := stmt.ExecContext(ctx,
			ds.Name, ds.Type, ds.SourceURL, string(ds.Status), ds.Enabled, ds.SyncInterval,
			ignoreRules, parameters, ds.LastSynced, ds.Created, ds.Updated,
		)
		if err != nil {
			return fmt.Errorf("failed to create data source %s: %w", ds.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// BulkUpdate обновляет несколько источников данных с указанными параметрами.
// Использует COALESCE для частичного обновления - nil поля сохраняют существующие значения.
func (r *DataSourcePostgresRepository) BulkUpdate(ctx context.Context, ids []int64, params repository.DataSourceBulkUpdateParams) error {
	if len(ids) == 0 || !params.HasChanges() {
		return nil
	}

	query := `
		UPDATE core_datasource
		SET
			type = COALESCE($1, type),
			enabled = COALESCE($2, enabled),
			description = COALESCE($3, description),
			sync_interval = COALESCE($4, sync_interval),
			parameters = COALESCE($5, parameters),
			ignore_rules = COALESCE($6, ignore_rules),
			comments = COALESCE($7, comments),
			updated = NOW()
		WHERE id = ANY($8)
	`

	var typeVal, descriptionVal, syncIntervalVal, ignoreRulesVal, commentsVal interface{}
	var parametersVal []byte
	var enabledVal *bool

	if params.Type != nil {
		typeVal = *params.Type
	}
	if params.Enabled != nil {
		enabledVal = params.Enabled
	}
	if params.Description != nil {
		descriptionVal = *params.Description
	}
	if params.SyncInterval != nil {
		syncIntervalVal = *params.SyncInterval
	}
	if params.Parameters != nil {
		parametersVal = *params.Parameters
	}
	if params.IgnoreRules != nil {
		ignoreRulesVal = *params.IgnoreRules
	}
	if params.Comments != nil {
		commentsVal = *params.Comments
	}

	_, err := r.db.ExecContext(ctx, query,
		typeVal,
		enabledVal,
		descriptionVal,
		syncIntervalVal,
		parametersVal,
		ignoreRulesVal,
		commentsVal,
		pq.Array(ids),
	)
	if err != nil {
		return fmt.Errorf("failed to bulk update data sources: %w", err)
	}

	return nil
}

// List возвращает список источников данных с фильтрацией
func (r *DataSourcePostgresRepository) List(ctx context.Context, filter repository.DataSourceFilter, limit, offset int) ([]*entity.DataSource, int, error) {
	var query string
	var args []interface{}

	// Базовый запрос
	query = `SELECT id, name, type, source_url, status, enabled, sync_interval, ignore_rules, parameters, last_synced, created, updated FROM core_datasource WHERE 1=1`
	countQuery := `SELECT COUNT(*)::int FROM core_datasource WHERE 1=1`

	// Фильтрация по статусу (MultipleChoiceFilter)
	if len(filter.Status) > 0 {
		statusList := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			statusList[i] = string(s)
		}
		query += ` AND status = ANY($` + appendArg(&args, statusList) + `)`
		countQuery += ` AND status = ANY($` + appendArg(&args, statusList) + `)`
	}

	if filter.Enabled != nil {
		query += ` AND enabled = $` + appendArg(&args, *filter.Enabled) + `)`
		countQuery += ` AND enabled = $` + appendArg(&args, *filter.Enabled) + `)`
	}

	if filter.Type != nil && *filter.Type != "" {
		query += ` AND type = $` + appendArg(&args, *filter.Type) + `)`
		countQuery += ` AND type = $` + appendArg(&args, *filter.Type) + `)`
	}

	// Поиск (Search logic)
	if filter.SearchQuery != nil && *filter.SearchQuery != "" {
		searchVal := "%" + *filter.SearchQuery + "%"
		query += ` AND (name ILIKE $` + appendArg(&args, searchVal) + ` OR description ILIKE $` + appendArg(&args, searchVal) + ` OR comments ILIKE $` + appendArg(&args, searchVal) + `)`
		countQuery += ` AND (name ILIKE $` + appendArg(&args, searchVal) + ` OR description ILIKE $` + appendArg(&args, searchVal) + ` OR comments ILIKE $` + appendArg(&args, searchVal) + `)`
	}

	// Пагинация
	query += ` ORDER BY name LIMIT $` + appendArg(&args, limit) + ` OFFSET $` + appendArg(&args, offset) + `)`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*entity.DataSource
	for rows.Next() {
		var row struct {
			ID           types.ID
			Name         string
			Type         string
			SourceURL    string
			Status       string
			Enabled      bool
			SyncInterval int32
			IgnoreRules  []byte
			Parameters   []byte
			LastSynced   sql.NullTime
			Created      time.Time
			Updated      time.Time
		}
		if err := rows.Scan(&row.ID, &row.Name, &row.Type, &row.SourceURL, &row.Status, &row.Enabled, &row.SyncInterval, &row.IgnoreRules, &row.Parameters, &row.LastSynced, &row.Created, &row.Updated); err != nil {
			return nil, 0, err
		}

		var ignoreRules []string
		if row.IgnoreRules != nil {
			types.UnmarshalJSON(row.IgnoreRules, &ignoreRules)
		}

		var lastSynced *time.Time
		if row.LastSynced.Valid {
			lastSynced = &row.LastSynced.Time
		}

		result = append(result, &entity.DataSource{
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
		})
	}

	var count int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return result, count, nil
}

func appendArg(args *[]interface{}, val interface{}) string {
	*args = append(*args, val)
	return fmt.Sprintf("%d", len(*args))
}

// Create создаёт новый источник данных
func (r *DataSourcePostgresRepository) Create(ctx context.Context, ds *entity.DataSource) error {
	q := coredb.Queries{DB: r.db}

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
		SourceUrl:    ds.SourceURL,
		Status:       string(ds.Status),
		Enabled:      ds.Enabled,
		SyncInterval: int32(ds.SyncInterval),
		IgnoreRules: pqtype.NullRawMessage{
			RawMessage: ignoreRules,
			Valid:      true,
		},
		Parameters: pqtype.NullRawMessage{
			RawMessage: parameters,
			Valid:      true,
		},
		LastSynced: lastSynced,
		Created:    now,
		Updated:    now,
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
	q := coredb.Queries{DB: r.db}

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
		SourceUrl:    ds.SourceURL,
		Status:       string(ds.Status),
		Enabled:      ds.Enabled,
		SyncInterval: int32(ds.SyncInterval),
		IgnoreRules: pqtype.NullRawMessage{
			RawMessage: ignoreRules,
			Valid:      true,
		},
		Parameters: pqtype.NullRawMessage{
			RawMessage: parameters,
			Valid:      true,
		},
		LastSynced: lastSynced,
		Updated:    time.Now(),
	})
	return err
}

// Delete удаляет источник данных
func (r *DataSourcePostgresRepository) Delete(ctx context.Context, id types.ID) error {
	q := coredb.Queries{DB: r.db}
	_, err := q.DeleteDataSource(ctx, id)
	return err
}

// UpdateStatus обновляет статус источника
func (r *DataSourcePostgresRepository) UpdateStatus(ctx context.Context, id types.ID, status types.Status, lastSynced *time.Time) error {
	q := coredb.Queries{DB: r.db}

	var lastSyncedNull sql.NullTime
	if lastSynced != nil {
		lastSyncedNull = sql.NullTime{Time: *lastSynced, Valid: true}
	}

	_, err := q.UpdateDataSourceStatus(ctx, coredb.UpdateDataSourceStatusParams{
		ID:         id,
		Status:     string(status),
		LastSynced: lastSyncedNull,
	})
	return err
}

// GetQueuedForSync возвращает источники, ожидающие синхронизации
func (r *DataSourcePostgresRepository) GetQueuedForSync(ctx context.Context, limit int) ([]*entity.DataSource, error) {
	q := coredb.Queries{DB: r.db}

	rows, err := q.GetQueuedDataSources(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	result := make([]*entity.DataSource, len(rows))
	for i, row := range rows {
		var ignoreRules []string
		if row.IgnoreRules.Valid && len(row.IgnoreRules.RawMessage) > 0 {
			if err := types.UnmarshalJSON(row.IgnoreRules.RawMessage, &ignoreRules); err != nil {
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
			SourceURL:    row.SourceUrl,
			Status:       types.Status(row.Status),
			Enabled:      row.Enabled,
			SyncInterval: int(row.SyncInterval),
			IgnoreRules:  ignoreRules,
			Parameters:   row.Parameters.RawMessage,
			LastSynced:   lastSynced,
			Created:      row.Created,
			Updated:      row.Updated,
		}
	}

	return result, nil
}
