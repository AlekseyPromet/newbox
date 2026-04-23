package services

import (
	"context"
	"encoding/json"
	"fmt"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/repository"
)

// CoreService предоставляет бизнес-логику для домена Core.
type CoreService struct {
	dataSources repository.DataSourceRepository
	dataFiles   repository.DataFileRepository
}

// NewCoreService создает новый экземпляр CoreService.
func NewCoreService(ds repository.DataSourceRepository, df repository.DataFileRepository) *CoreService {
	return &CoreService{
		dataSources: ds,
		dataFiles:   df,
	}
}

// BulkEditDataSources обновляет несколько источников данных одновременно.
// В NetBox это реализуется через DataSourceBulkEditForm.
func (s *CoreService) BulkEditDataSources(ctx context.Context, ids []string, updates map[string]interface{}) error {
	for _, id := range ids {
		existing, err := s.dataSources.GetByID(ctx, id)
		if err != nil {
			return fmt.Errorf("failed to get data source %s: %w", id, err)
		}

		// Применяем обновления на основе предоставленных полей
		if val, ok := updates["type"].(string); ok {
			existing.Type = val
		}
		if val, ok := updates["enabled"].(*bool); ok {
			existing.Enabled = *val
		}
		if val, ok := updates["sync_interval"].(int); ok {
			existing.SyncInterval = val
		}
		if val, ok := updates["parameters"].(string); ok {
			existing.Parameters = json.RawMessage(val)
		}
		if val, ok := updates["ignore_rules"].([]string); ok {
			existing.IgnoreRules = val
		}
		if val, ok := updates["description"].(string); ok {
			existing.Description = val
		}

		if err := existing.Validate(); err != nil {
			return fmt.Errorf("validation failed for data source %s: %w", id, err)
		}

		if err := s.dataSources.Update(ctx, existing); err != nil {
			return fmt.Errorf("failed to update data source %s: %w", id, err)
		}

	}
	return nil
}

func (s *CoreService) GetDataSource(ctx context.Context, id string) (*entity.DataSource, error) {
	return s.dataSources.GetByID(ctx, id)
}

func (s *CoreService) ListDataSources(ctx context.Context, filter repository.DataSourceFilter, limit, offset int) ([]*entity.DataSource, int64, error) {
	filter.Limit = limit
	filter.Offset = offset
	return s.dataSources.List(ctx, filter)
}

func (s *CoreService) GetDataFile(ctx context.Context, id string) (*entity.DataFile, error) {
	return s.dataFiles.GetByID(ctx, id)
}

func (s *CoreService) ListDataFiles(ctx context.Context, filter repository.DataFileFilter, limit, offset int) ([]*entity.DataFile, int64, error) {
	filter.Limit = limit
	filter.Offset = offset
	return s.dataFiles.List(ctx, filter)
}

// BulkImportDataSources обрабатывает импорт источников данных.
// В NetBox это реализуется через DataSourceImportForm.
func (s *CoreService) BulkImportDataSources(ctx context.Context, data []entity.DataSource) error {
	for _, ds := range data {
		if err := ds.Validate(); err != nil {
			return fmt.Errorf("validation failed for imported data source %s: %w", ds.Name, err)
		}

		if err := s.dataSources.Create(ctx, &ds); err != nil {
			return fmt.Errorf("failed to import data source %s: %w", ds.Name, err)
		}
	}
	return nil
}
