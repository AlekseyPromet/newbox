package services

import (
	"context"
	"fmt"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
)

// CoreService предоставляет бизнес-логику для домена Core.
type CoreService struct {
	dataSources repository.DataSourceRepository
}

// NewCoreService создает новый экземпляр CoreService.
func NewCoreService(ds repository.DataSourceRepository) *CoreService {
	return &CoreService{
		dataSources: ds,
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
			existing.SyncInterval = &val
		}
		if val, ok := updates["parameters"].(string); ok {
			existing.Parameters = val
		}
		if val, ok := updates["ignore_rules"].(string); ok {
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
