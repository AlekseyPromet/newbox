// Package repository содержит интерфейсы репозиториев домена Core
package repository

import (
	"context"
	"time"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/pkg/types"
)

// ConfigRevisionRepository интерфейс для работы с ревизиями конфигурации
type ConfigRevisionRepository interface {
	// GetByID возвращает ревизию по ID
	GetByID(ctx context.Context, id types.ID) (*entity.ConfigRevision, error)
	// GetActive возвращает активную ревизию
	GetActive(ctx context.Context) (*entity.ConfigRevision, error)
	// List возвращает список ревизий с фильтрацией и пагинацией
	List(ctx context.Context, filter ConfigRevisionFilter, limit, offset int) ([]*entity.ConfigRevision, int, error)
	// Create создаёт новую ревизию
	Create(ctx context.Context, revision *entity.ConfigRevision) error
	// Update обновляет ревизию
	Update(ctx context.Context, revision *entity.ConfigRevision) error
	// Delete удаляет ревизию
	Delete(ctx context.Context, id types.ID) error
	// SetActive делает ревизию активной
	SetActive(ctx context.Context, id types.ID) error
}

// ObjectTypeRepository интерфейс для работы с типами объектов
type ObjectTypeRepository interface {
	// GetByID возвращает тип объекта по ID
	GetByID(ctx context.Context, id types.ID) (*entity.ObjectType, error)
	// GetByAppAndModel возвращает тип объекта по app_label и model
	GetByAppAndModel(ctx context.Context, appLabel, model string) (*entity.ObjectType, error)
	// List возвращает список типов объектов с фильтрацией и пагинацией
	List(ctx context.Context, filter ObjectTypeFilter, limit, offset int) ([]*entity.ObjectType, int, error)
	// Create создаёт новый тип объекта
	Create(ctx context.Context, ot *entity.ObjectType) error
	// Update обновляет тип объекта
	Update(ctx context.Context, ot *entity.ObjectType) error
	// Delete удаляет тип объекта
	Delete(ctx context.Context, id types.ID) error
}

// ObjectChangeRepository интерфейс для журнала изменений
type ObjectChangeRepository interface {
	// GetByID возвращает запись об изменении по ID
	GetByID(ctx context.Context, id types.ID) (*entity.ObjectChange, error)
	// List возвращает список изменений с фильтрацией и пагинацией
	List(ctx context.Context, filter ObjectChangeFilter, limit, offset int) ([]*entity.ObjectChange, int, error)
	// Create создаёт запись об изменении
	Create(ctx context.Context, change *entity.ObjectChange) error
	// BulkCreate создаёт несколько записей об изменениях
	BulkCreate(ctx context.Context, changes []*entity.ObjectChange) error
	// DeleteOld удаляет старые записи (старше cutoffTime)
	DeleteOld(ctx context.Context, cutoffTime time.Time) (int64, error)
}

// DataSourceRepository интерфейс для источников данных
type DataSourceRepository interface {
	// GetByID возвращает источник данных по ID
	GetByID(ctx context.Context, id types.ID) (*entity.DataSource, error)
	// GetByName возвращает источник данных по имени
	GetByName(ctx context.Context, name string) (*entity.DataSource, error)
	// List возвращает список источников данных с фильтрацией и пагинацией
	List(ctx context.Context, filter DataSourceFilter, limit, offset int) ([]*entity.DataSource, int, error)
	// Create создаёт новый источник данных
	Create(ctx context.Context, ds *entity.DataSource) error
	// Update обновляет источник данных
	Update(ctx context.Context, ds *entity.DataSource) error
	// Delete удаляет источник данных
	Delete(ctx context.Context, id types.ID) error
	// UpdateStatus обновляет статус источника
	UpdateStatus(ctx context.Context, id types.ID, status types.Status, lastSynced *time.Time) error
	// GetQueuedForSync возвращает источники, ожидающие синхронизации
	GetQueuedForSync(ctx context.Context, limit int) ([]*entity.DataSource, error)
}

// DataFileRepository интерфейс для файлов данных
type DataFileRepository interface {
	// GetByID возвращает файл данных по ID
	GetByID(ctx context.Context, id types.ID) (*entity.DataFile, error)
	// GetBySourceAndPath возвращает файл по источнику и пути
	GetBySourceAndPath(ctx context.Context, sourceID types.ID, path string) (*entity.DataFile, error)
	// List возвращает список файлов данных с фильтрацией и пагинацией
	List(ctx context.Context, filter DataFileFilter, limit, offset int) ([]*entity.DataFile, int, error)
	// Create создаёт новый файл данных
	Create(ctx context.Context, df *entity.DataFile) error
	// Update обновляет файл данных
	Update(ctx context.Context, df *entity.DataFile) error
	// Delete удаляет файл данных
	Delete(ctx context.Context, id types.ID) error
	// BulkDeleteBySource удаляет все файлы источника
	BulkDeleteBySource(ctx context.Context, sourceID types.ID) (int64, error)
}

// JobRepository интерфейс для фоновых задач
type JobRepository interface {
	// GetByID возвращает задачу по ID
	GetByID(ctx context.Context, id types.ID) (*entity.Job, error)
	// List возвращает список задач с фильтрацией и пагинацией
	List(ctx context.Context, filter JobFilter, limit, offset int) ([]*entity.Job, int, error)
	// Create создаёт новую задачу
	Create(ctx context.Context, job *entity.Job) error
	// Update обновляет задачу
	Update(ctx context.Context, job *entity.Job) error
	// Delete удаляет задачу
	Delete(ctx context.Context, id types.ID) error
	// UpdateStatus обновляет статус задачи
	UpdateStatus(ctx context.Context, id types.ID, status types.Status, error *string, completedAt *time.Time) error
	// GetScheduled возвращает запланированные задачи
	GetScheduled(ctx context.Context, before time.Time, limit int) ([]*entity.Job, error)
	// CleanupOld удаляет старые завершённые задачи
	CleanupOld(ctx context.Context, olderThan time.Time) (int64, error)
}
