# План миграции модуля Core (NetBox Python → Go)

## Обзор модуля Core

Модуль `core` в NetBox отвечает за системные функции:
- **Data Sources** — источники данных (git, файлы) для синхронизации конфигов
- **Data Files** — файлы, полученные из источников данных
- **Jobs** — фоновые задачи (аналог RQ jobs)
- **Object Changes** — журнал изменений объектов (change logging)
- **Object Types** — типы объектов (аналог Django ContentType)
- **Config Revisions** — ревизии конфигурации NetBox
- **Managed Files** — управляемые файлы (скрипты, отчёты)
- **Plugin Management** — управление плагинами и каталогом плагинов

---

## Текущее состояние Go-проекта

### ✅ Уже реализовано

#### 1. Доменные сущности (`internal/domain/core/entity/core.go`)
- `ConfigRevision` — ревизия конфигурации
- `ObjectType` — тип объекта
- `ObjectChange` — изменение объекта
- `DataSource` — источник данных
- `DataFile` — файл данных
- `Job` — фоновая задача

#### 2. Перечисления (`internal/domain/core/enum/status.go`)
- `DataSourceStatusChoices` — статусы источников данных
- `JobStatusChoices` — статусы задач
- `JobIntervalChoices` — интервалы выполнения
- `ObjectChangeActionChoices` — действия изменений

#### 3. HTTP обработчики (`internal/delivery/http/handlers/core_handler.go`)
- CRUD для Data Sources
- CRUD для Data Files (частично)
- List/Get для Jobs
- List/Get для Object Changes
- List/Get для Object Types
- Заглушки для Background Queues/Workers/Tasks

#### 4. Интерфейсы репозиториев (`internal/repository/interfaces.go`)
Требуется дополнение интерфейсами для Core сущностей.

---

## План миграции

### Этап 1: Анализ исходного кода Python

#### 1.1 Модели данных (`netbox/core/models/`)

| Модель | Файл | Описание | Приоритет |
|--------|------|----------|-----------|
| `ObjectType` | `object_types.py` | Обёртка над ContentType | Высокий |
| `ObjectChange` | `change_logging.py` | Журнал изменений | Высокий |
| `DataSource` | `data.py` | Источники данных | Высокий |
| `DataFile` | `data.py` | Файлы данных | Высокий |
| `AutoSyncRecord` | `data.py` | Авто-синхронизация | Средний |
| `Job` | `jobs.py` | Фоновые задачи | Высокий |
| `ConfigRevision` | `config.py` | Ревизии конфига | Средний |
| `ManagedFile` | `files.py` | Управляемые файлы | Средний |

#### 1.2 API ViewSets (`netbox/core/api/views.py`)

| Endpoint | Методы | Описание |
|----------|--------|----------|
| `/api/core/data-sources/` | GET, POST | Список/создание источников |
| `/api/core/data-sources/:id/` | GET, PUT, DELETE | Операции с источником |
| `/api/core/data-sources/:id/sync/` | POST | Синхронизация источника |
| `/api/core/data-files/` | GET | Список файлов |
| `/api/core/data-files/:id/` | GET | Детали файла |
| `/api/core/jobs/` | GET | Список задач |
| `/api/core/jobs/:id/` | GET | Детали задачи |
| `/api/core/object-changes/` | GET | Список изменений |
| `/api/core/object-changes/:id/` | GET | Детали изменения |
| `/api/core/object-types/` | GET | Список типов объектов |
| `/api/core/object-types/:id/` | GET | Детали типа |
| `/api/core/background-queues/` | GET | Очереди RQ |
| `/api/core/background-workers/` | GET | Воркеры RQ |
| `/api/core/background-tasks/` | GET, POST | Задачи RQ |

#### 1.3 Выборы (choices) (`netbox/core/choices.py`)

```python
# DataSourceStatusChoices
NEW = 'new'
QUEUED = 'queued'
SYNCING = 'syncing'
COMPLETED = 'completed'
FAILED = 'failed'

# JobStatusChoices
STATUS_PENDING = 'pending'
STATUS_SCHEDULED = 'scheduled'
STATUS_RUNNING = 'running'
STATUS_COMPLETED = 'completed'
STATUS_ERRORED = 'errored'
STATUS_FAILED = 'failed'

# JobIntervalChoices
INTERVAL_MINUTELY = 1
INTERVAL_HOURLY = 60
INTERVAL_DAILY = 60 * 24
INTERVAL_WEEKLY = 60 * 24 * 7

# ObjectChangeActionChoices
ACTION_CREATE = 'create'
ACTION_UPDATE = 'update'
ACTION_DELETE = 'delete'

# ManagedFileRootPathChoices
SCRIPTS = 'scripts'
REPORTS = 'reports'
```

#### 1.4 Фильтры (`netbox/core/filtersets.py`)

| Фильтр | Сущность | Поля фильтрации |
|--------|----------|-----------------|
| `DataSourceFilterSet` | DataSource | name, type, status, enabled, sync_interval |
| `DataFileFilterSet` | DataFile | source_id, path, hash |
| `JobFilterSet` | Job | object_type, object_id, status, queue_name |
| `ObjectChangeFilterSet` | ObjectChange | changed_object_type, user_id, action, request_id |
| `ObjectTypeFilterSet` | ObjectType | app_label, model, public, feature |

---

### Этап 2: Реализация интерфейсов репозиториев

#### 2.1 Добавить интерфейсы в `internal/repository/interfaces.go`

```go
// DataSourceRepository определяет интерфейс для работы с источниками данных
type DataSourceRepository interface {
    GetByID(ctx context.Context, id string) (*core_entity.DataSource, error)
    GetByName(ctx context.Context, name string) (*core_entity.DataSource, error)
    List(ctx context.Context, filter DataSourceFilter) ([]*core_entity.DataSource, int64, error)
    Create(ctx context.Context, ds *core_entity.DataSource) error
    Update(ctx context.Context, ds *core_entity.DataSource) error
    Delete(ctx context.Context, id string) error
    Sync(ctx context.Context, id string) error
    Exists(ctx context.Context, id string) (bool, error)
}

// DataSourceFilter представляет фильтры для поиска источников данных
type DataSourceFilter struct {
    Name         *string
    Type         *string
    Status       *string
    Enabled      *bool
    SyncInterval *int
    Limit        int
    Offset       int
}

// DataFileRepository определяет интерфейс для работы с файлами данных
type DataFileRepository interface {
    GetByID(ctx context.Context, id string) (*core_entity.DataFile, error)
    GetBySourceAndPath(ctx context.Context, sourceID string, path string) (*core_entity.DataFile, error)
    List(ctx context.Context, filter DataFileFilter) ([]*core_entity.DataFile, int64, error)
    Create(ctx context.Context, df *core_entity.DataFile) error
    Update(ctx context.Context, df *core_entity.DataFile) error
    Delete(ctx context.Context, id string) error
    BulkCreate(ctx context.Context, files []*core_entity.DataFile) error
    BulkUpdate(ctx context.Context, files []*core_entity.DataFile) error
    BulkDelete(ctx context.Context, ids []string) error
}

// DataFileFilter представляет фильтры для поиска файлов данных
type DataFileFilter struct {
    SourceID *string
    Path     *string
    Hash     *string
    Limit    int
    Offset   int
}

// JobRepository определяет интерфейс для работы с фоновыми задачами
type JobRepository interface {
    GetByID(ctx context.Context, id string) (*core_entity.Job, error)
    GetByJobID(ctx context.Context, jobID string) (*core_entity.Job, error)
    List(ctx context.Context, filter JobFilter) ([]*core_entity.Job, int64, error)
    Create(ctx context.Context, job *core_entity.Job) error
    Update(ctx context.Context, job *core_entity.Job) error
    Delete(ctx context.Context, id string) error
    Start(ctx context.Context, id string) error
    Complete(ctx context.Context, id string, status types.Status, errorText *string) error
    Log(ctx context.Context, id string, entry core_entity.JobLogEntry) error
}

// JobFilter представляет фильтры для поиска задач
type JobFilter struct {
    ObjectType  *string
    ObjectID    *string
    Status      *string
    QueueName   *string
    ScheduledAt *time.Time
    Limit       int
    Offset      int
}

// ObjectChangeRepository определяет интерфейс для журнала изменений
type ObjectChangeRepository interface {
    GetByID(ctx context.Context, id string) (*core_entity.ObjectChange, error)
    List(ctx context.Context, filter ObjectChangeFilter) ([]*core_entity.ObjectChange, int64, error)
    Create(ctx context.Context, oc *core_entity.ObjectChange) error
    LogChange(ctx context.Context, action types.Status, objectType string, objectID string, 
              objectRepr string, preChange, postChange json.RawMessage, userID *types.ID) error
    GetChangesForObject(ctx context.Context, objectType string, objectID string, limit int) ([]*core_entity.ObjectChange, error)
}

// ObjectChangeFilter представляет фильтры для поиска изменений
type ObjectChangeFilter struct {
    ChangedObjectType *string
    ChangedObjectID   *string
    UserID            *string
    Action            *string
    RequestID         *string
    Since             *time.Time
    Until             *time.Time
    Limit             int
    Offset            int
}

// ObjectTypeRepository определяет интерфейс для типов объектов
type ObjectTypeRepository interface {
    GetByID(ctx context.Context, id string) (*core_entity.ObjectType, error)
    GetByAppAndModel(ctx context.Context, appLabel, model string) (*core_entity.ObjectType, error)
    List(ctx context.Context, filter ObjectTypeFilter) ([]*core_entity.ObjectType, int64, error)
    GetForModel(ctx context.Context, modelName string) (*core_entity.ObjectType, error)
    GetForModels(ctx context.Context, modelNames []string) (map[string]*core_entity.ObjectType, error)
    Public(ctx context.Context) ([]*core_entity.ObjectType, error)
    WithFeature(ctx context.Context, feature string) ([]*core_entity.ObjectType, error)
}

// ObjectTypeFilter представляет фильтры для поиска типов объектов
type ObjectTypeFilter struct {
    AppLabel *string
    Model    *string
    Public   *bool
    Feature  *string
    Limit    int
    Offset   int
}

// ConfigRevisionRepository определяет интерфейс для ревизий конфигурации
type ConfigRevisionRepository interface {
    GetByID(ctx context.Context, id string) (*core_entity.ConfigRevision, error)
    GetActive(ctx context.Context) (*core_entity.ConfigRevision, error)
    List(ctx context.Context, limit, offset int) ([]*core_entity.ConfigRevision, int64, error)
    Create(ctx context.Context, cr *core_entity.ConfigRevision) error
    Activate(ctx context.Context, id string) error
    GetLatest(ctx context.Context) (*core_entity.ConfigRevision, error)
}
```

---

### Этап 3: Реализация PostgreSQL репозиториев

#### 3.1 Создать файлы репозиториев

Создать директорию `internal/repository/postgres/core/` со следующими файлами:

1. **`datasource_repo.go`** — реализация `DataSourceRepository`
   - Методы CRUD
   - Метод `Sync()` для запуска синхронизации
   - Поддержка фильтрации

2. **`datafile_repo.go`** — реализация `DataFileRepository`
   - Методы CRUD
   - Методы bulk-операций (Create, Update, Delete)
   - Метод `GetBySourceAndPath()`

3. **`job_repo.go`** — реализация `JobRepository`
   - Методы CRUD
   - Интеграция с EtcdQueue для фоновых задач
   - Логирование выполнения задач

4. **`objectchange_repo.go`** — реализация `ObjectChangeRepository`
   - Методы записи изменений
   - Методы чтения истории изменений
   - Оптимизированные запросы для GIN индексов

5. **`objecttype_repo.go`** — реализация `ObjectTypeRepository`
   - Кэширование результатов (etcd)
   - Метод `GetForModel()` с авто-созданием
   - Фильтрация по features (PostgreSQL JSONB)

6. **`configrevision_repo.go`** — реализация `ConfigRevisionRepository`
   - Метод активации ревизии
   - Получение активной ревизии
   - Валидация уникальности активной ревизии

#### 3.2 SQL запросы через sqlc

Создать файл `internal/repository/postgres/core/queries.sql`:

```sql
-- name: GetDataSourceByID :one
SELECT * FROM core_datasource WHERE id = $1;

-- name: GetDataSources :many
SELECT * FROM core_datasource
WHERE 
    ($1::text IS NULL OR name = $1)
    AND ($2::text IS NULL OR type = $2)
    AND ($3::text IS NULL OR status = $3)
    AND ($4::boolean IS NULL OR enabled = $4)
ORDER BY name
LIMIT $5 OFFSET $6;

-- name: CountDataSources :one
SELECT COUNT(*) FROM core_datasource
WHERE 
    ($1::text IS NULL OR name = $1)
    AND ($2::text IS NULL OR type = $2)
    AND ($3::text IS NULL OR status = $3)
    AND ($4::boolean IS NULL OR enabled = $4);

-- name: InsertDataSource :one
INSERT INTO core_datasource (id, name, type, source_url, status, enabled, sync_interval, ignore_rules, parameters, last_synced, created, updated)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING *;

-- name: UpdateDataSource :exec
UPDATE core_datasource 
SET name = $2, type = $3, source_url = $4, status = $5, enabled = $6, 
    sync_interval = $7, ignore_rules = $8, parameters = $9, last_synced = $10, updated = $11
WHERE id = $1;

-- name: DeleteDataSource :exec
DELETE FROM core_datasource WHERE id = $1;

-- name: GetDataFileByID :one
SELECT * FROM core_datafile WHERE id = $1;

-- name: GetDataFiles :many
SELECT * FROM core_datafile
WHERE 
    ($1::uuid IS NULL OR source_id = $1)
    AND ($2::text IS NULL OR path = $2)
ORDER BY source_id, path
LIMIT $3 OFFSET $4;

-- name: InsertDataFile :one
INSERT INTO core_datafile (id, source_id, path, size, hash, data, created, last_updated)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: BulkInsertDataFiles :copyfrom
INSERT INTO core_datafile (id, source_id, path, size, hash, data, created, last_updated)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetJobByID :one
SELECT * FROM core_job WHERE id = $1;

-- name: GetJobs :many
SELECT * FROM core_job
WHERE 
    ($1::text IS NULL OR object_type = $1)
    AND ($2::uuid IS NULL OR object_id = $2)
    AND ($3::text IS NULL OR status = $3)
    AND ($4::text IS NULL OR queue_name = $4)
ORDER BY created DESC
LIMIT $5 OFFSET $6;

-- name: InsertJob :one
INSERT INTO core_job (id, object_type, object_id, name, status, interval, scheduled_at, started_at, completed_at, user_id, queue_name, job_id, data, error, created, updated)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING *;

-- name: UpdateJobStatus :exec
UPDATE core_job 
SET status = $2, completed_at = $3, error = $4, updated = $5
WHERE id = $1;

-- name: InsertObjectChange :one
INSERT INTO core_objectchange (id, time, user_id, user_name, request_id, action, changed_object_type_id, changed_object_id, related_object_type_id, related_object_id, object_repr, message, prechange_data, postchange_data)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING *;

-- name: GetObjectChanges :many
SELECT * FROM core_objectchange
WHERE 
    ($1::uuid IS NULL OR changed_object_type_id = $1)
    AND ($2::bigint IS NULL OR changed_object_id = $2)
    AND ($3::uuid IS NULL OR user_id = $3)
    AND ($4::text IS NULL OR action = $4)
    AND ($5::timestamptz IS NULL OR time >= $5)
    AND ($6::timestamptz IS NULL OR time <= $6)
ORDER BY time DESC
LIMIT $7 OFFSET $8;

-- name: GetObjectTypeByID :one
SELECT * FROM core_objecttype WHERE id = $1;

-- name: GetObjectTypeByAppAndModel :one
SELECT * FROM core_objecttype WHERE app_label = $1 AND model = $2;

-- name: GetObjectTypes :many
SELECT * FROM core_objecttype
WHERE 
    ($1::text IS NULL OR app_label = $1)
    AND ($2::text IS NULL OR model = $2)
    AND ($3::boolean IS NULL OR public = $3)
    AND ($4::text IS NULL OR $4 = ANY(features))
ORDER BY app_label, model
LIMIT $5 OFFSET $6;

-- name: GetActiveConfigRevision :one
SELECT * FROM core_configrevision WHERE active = true LIMIT 1;

-- name: ActivateConfigRevision :exec
UPDATE core_configrevision SET active = false WHERE active = true;
UPDATE core_configrevision SET active = true WHERE id = $1;
```

---

### Этап 4: Расширение HTTP обработчиков

#### 4.1 Обновить `internal/delivery/http/handlers/core_handler.go`

Добавить следующие методы:

```go
// SyncDataSource обрабатывает POST /api/core/data-sources/:id/sync
func (h *CoreHandlers) SyncDataSource(c echo.Context) error {
    id := c.Param("id")
    if err := h.dataSources.Sync(c.Request().Context(), id); err != nil {
        return handleError(err)
    }
    return c.JSON(http.StatusOK, map[string]string{"status": "sync initiated"})
}

// CreateDataFile обрабатывает POST /api/core/data-files
func (h *CoreHandlers) CreateDataFile(c echo.Context) error {
    // Реализация создания файла данных
}

// UpdateDataFile обрабатывает PUT /api/core/data-files/:id
func (h *CoreHandlers) UpdateDataFile(c echo.Context) error {
    // Реализация обновления файла данных
}

// DeleteDataFile обрабатывает DELETE /api/core/data-files/:id
func (h *CoreHandlers) DeleteDataFile(c echo.Context) error {
    id := c.Param("id")
    if err := h.dataFiles.Delete(c.Request().Context(), id); err != nil {
        return handleError(err)
    }
    return c.NoContent(http.StatusNoContent)
}

// CreateJob обрабатывает POST /api/core/jobs
func (h *CoreHandlers) CreateJob(c echo.Context) error {
    // Создание новой задачи
}

// LogObjectChange обрабатывает POST /api/core/object-changes/log
func (h *CoreHandlers) LogObjectChange(c echo.Context) error {
    // Логирование изменения объекта
}

// GetActiveConfigRevision обрабатывает GET /api/core/config-revisions/active
func (h *CoreHandlers) GetActiveConfigRevision(c echo.Context) error {
    // Получение активной ревизии конфигурации
}

// ActivateConfigRevision обрабатывает POST /api/core/config-revisions/:id/activate
func (h *CoreHandlers) ActivateConfigRevision(c echo.Context) error {
    id := c.Param("id")
    // Активация ревизии
}
```

#### 4.2 Регистрация маршрутов

Обновить роутер в `cmd/api/main.go`:

```go
// Core API routes
coreGroup := api.Group("/core")
{
    // Data Sources
    coreGroup.GET("/data-sources", coreHandlers.ListDataSources)
    coreGroup.GET("/data-sources/:id", coreHandlers.GetDataSource)
    coreGroup.POST("/data-sources", coreHandlers.CreateDataSource)
    coreGroup.PUT("/data-sources/:id", coreHandlers.UpdateDataSource)
    coreGroup.DELETE("/data-sources/:id", coreHandlers.DeleteDataSource)
    coreGroup.POST("/data-sources/:id/sync", coreHandlers.SyncDataSource)
    
    // Data Files
    coreGroup.GET("/data-files", coreHandlers.ListDataFiles)
    coreGroup.GET("/data-files/:id", coreHandlers.GetDataFile)
    coreGroup.POST("/data-files", coreHandlers.CreateDataFile)
    coreGroup.PUT("/data-files/:id", coreHandlers.UpdateDataFile)
    coreGroup.DELETE("/data-files/:id", coreHandlers.DeleteDataFile)
    
    // Jobs
    coreGroup.GET("/jobs", coreHandlers.ListJobs)
    coreGroup.GET("/jobs/:id", coreHandlers.GetJob)
    coreGroup.POST("/jobs", coreHandlers.CreateJob)
    
    // Object Changes
    coreGroup.GET("/object-changes", coreHandlers.ListObjectChanges)
    coreGroup.GET("/object-changes/:id", coreHandlers.GetObjectChange)
    coreGroup.POST("/object-changes/log", coreHandlers.LogObjectChange)
    
    // Object Types
    coreGroup.GET("/object-types", coreHandlers.ListObjectTypes)
    coreGroup.GET("/object-types/:id", coreHandlers.GetObjectType)
    
    // Config Revisions
    coreGroup.GET("/config-revisions/active", coreHandlers.GetActiveConfigRevision)
    coreGroup.POST("/config-revisions/:id/activate", coreHandlers.ActivateConfigRevision)
    
    // Background (RQ) - заглушки
    coreGroup.GET("/background-queues", coreHandlers.ListBackgroundQueues)
    coreGroup.GET("/background-queues/:name", coreHandlers.GetBackgroundQueue)
    coreGroup.GET("/background-workers", coreHandlers.ListBackgroundWorkers)
    coreGroup.GET("/background-workers/:name", coreHandlers.GetBackgroundWorker)
    coreGroup.GET("/background-tasks", coreHandlers.ListBackgroundTasks)
    coreGroup.GET("/background-tasks/:id", coreHandlers.GetBackgroundTask)
    coreGroup.POST("/background-tasks/:id/delete", coreHandlers.DeleteBackgroundTask)
    coreGroup.POST("/background-tasks/:id/requeue", coreHandlers.RequeueBackgroundTask)
    coreGroup.POST("/background-tasks/:id/enqueue", coreHandlers.EnqueueBackgroundTask)
    coreGroup.POST("/background-tasks/:id/stop", coreHandlers.StopBackgroundTask)
}
```

---

### Этап 5: Миграции базы данных

#### 5.1 Создать SQL миграции

Создать файлы в `netbox_go/migrations/`:

**`001_core_initial.up.sql`**:
```sql
-- ObjectType (наследует django_content_type)
CREATE TABLE core_objecttype (
    contenttype_ptr_id INTEGER PRIMARY KEY REFERENCES django_content_type(id),
    public BOOLEAN DEFAULT FALSE,
    features TEXT[] DEFAULT '{}'
);

CREATE INDEX idx_objecttype_features ON core_objecttype USING GIN(features);
CREATE INDEX idx_objecttype_app_model ON core_objecttype(app_label, model);

-- DataSource
CREATE TABLE core_datasource (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    type VARCHAR(50) NOT NULL,
    source_url VARCHAR(200) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'new',
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sync_interval SMALLINT,
    ignore_rules TEXT NOT NULL DEFAULT '',
    parameters JSONB,
    last_synced TIMESTAMPTZ,
    created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_datasource_status ON core_datasource(status);
CREATE INDEX idx_datasource_enabled ON core_datasource(enabled);

-- DataFile
CREATE TABLE core_datafile (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES core_datasource(id) ON DELETE CASCADE,
    path VARCHAR(1000) NOT NULL,
    size BIGINT NOT NULL,
    hash CHAR(64) NOT NULL,
    data BYTEA NOT NULL,
    created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_updated TIMESTAMPTZ NOT NULL,
    CONSTRAINT unique_source_path UNIQUE(source_id, path)
);

CREATE INDEX idx_datafile_source ON core_datafile(source_id);
CREATE INDEX idx_datafile_path ON core_datafile(path);

-- Job
CREATE TABLE core_job (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    object_type_id INTEGER REFERENCES django_content_type(id),
    object_id BIGINT,
    name VARCHAR(200) NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'pending',
    interval INTEGER,
    scheduled_at TIMESTAMPTZ,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    user_id INTEGER REFERENCES auth_user(id),
    queue_name VARCHAR(100),
    job_id UUID NOT NULL UNIQUE,
    data JSONB,
    error TEXT,
    log_entries JSONB[] DEFAULT '{}',
    created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_job_object ON core_job(object_type_id, object_id);
CREATE INDEX idx_job_status ON core_job(status);
CREATE INDEX idx_job_created ON core_job(created DESC);

-- ObjectChange
CREATE TABLE core_objectchange (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id INTEGER REFERENCES auth_user(id),
    user_name VARCHAR(150) NOT NULL,
    request_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    changed_object_type_id INTEGER NOT NULL REFERENCES django_content_type(id),
    changed_object_id BIGINT NOT NULL,
    related_object_type_id INTEGER REFERENCES django_content_type(id),
    related_object_id BIGINT,
    object_repr VARCHAR(200) NOT NULL,
    message VARCHAR(200),
    prechange_data JSONB,
    postchange_data JSONB
);

CREATE INDEX idx_objectchange_time ON core_objectchange(time DESC);
CREATE INDEX idx_objectchange_changed ON core_objectchange(changed_object_type_id, changed_object_id);
CREATE INDEX idx_objectchange_related ON core_objectchange(related_object_type_id, related_object_id);
CREATE INDEX idx_objectchange_request ON core_objectchange(request_id);

-- ConfigRevision
CREATE TABLE core_configrevision (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    active BOOLEAN NOT NULL DEFAULT FALSE,
    created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    comment VARCHAR(200),
    data JSONB,
    CONSTRAINT unique_active_revision UNIQUE(active) WHERE (active = TRUE)
);

CREATE INDEX idx_configrevision_active ON core_configrevision(active);
CREATE INDEX idx_configrevision_created ON core_configrevision(created DESC);

-- AutoSyncRecord
CREATE TABLE core_autosyncrecord (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    datafile_id UUID NOT NULL REFERENCES core_datafile(id) ON DELETE CASCADE,
    object_type_id INTEGER NOT NULL REFERENCES django_content_type(id),
    object_id BIGINT NOT NULL,
    CONSTRAINT unique_autosync_object UNIQUE(object_type_id, object_id)
);

-- ManagedFile
CREATE TABLE core_managedfile (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_updated TIMESTAMPTZ,
    file_root VARCHAR(1000) NOT NULL,
    file_path VARCHAR(1000) NOT NULL,
    datafile_id UUID REFERENCES core_datafile(id),
    CONSTRAINT unique_root_path UNIQUE(file_root, file_path)
);

CREATE INDEX idx_managedfile_root_path ON core_managedfile(file_root, file_path);
```

---

### Этап 6: Система фоновых задач на etcd

#### 6.1 Архитектура хранения задач в etcd

Создать `internal/pkg/taskqueue/etcd_queue.go`:

```go
package taskqueue

import (
    "context"
    "encoding/json"
    "time"
    
    "go.etcd.io/etcd/client/v3"
    "github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
    "github.com/google/uuid"
)

const (
    TaskQueuePrefix     = "/netbox/tasks/queue/"
    TaskProcessingPrefix = "/netbox/tasks/processing/"
    TaskResultPrefix    = "/netbox/tasks/result/"
    TaskLockPrefix      = "/netbox/tasks/lock/"
    
    TypeSyncDataSource = "core:sync_datasource"
    TypeProcessJob     = "core:process_job"
)

// TaskStatus статус задачи
type TaskStatus string

const (
    TaskStatusPending   TaskStatus = "pending"
    TaskStatusRunning   TaskStatus = "running"
    TaskStatusCompleted TaskStatus = "completed"
    TaskStatusFailed    TaskStatus = "failed"
)

// Task задача для выполнения
type Task struct {
    ID        string                 `json:"id"`
    Type      string                 `json:"type"`
    Payload   map[string]interface{} `json:"payload"`
    Status    TaskStatus             `json:"status"`
    Priority  int                    `json:"priority"`
    CreatedAt time.Time              `json:"created_at"`
    StartedAt *time.Time             `json:"started_at,omitempty"`
    EndedAt   *time.Time             `json:"ended_at,omitempty"`
    Result    interface{}            `json:"result,omitempty"`
    Error     string                 `json:"error,omitempty"`
    RetryCount int                   `json:"retry_count"`
    MaxRetries int                  `json:"max_retries"`
}

// EtcdQueue очередь задач на базе etcd
type EtcdQueue struct {
    client *clientv3.Client
    ctx    context.Context
}

// NewEtcdQueue создаёт новую очередь задач
func NewEtcdQueue(client *clientv3.Client) *EtcdQueue {
    return &EtcdQueue{
        client: client,
        ctx:    context.Background(),
    }
}

// Enqueue добавляет задачу в очередь
func (q *EtcdQueue) Enqueue(ctx context.Context, task *Task) error {
    task.ID = uuid.New().String()
    task.Status = TaskStatusPending
    task.CreatedAt = time.Now()
    
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }
    
    // Ключ с приоритетом для сортировки (чем меньше число, тем выше приоритет)
    key := q.buildQueueKey(task.Priority, task.ID)
    
    _, err = q.client.Put(ctx, key, string(data))
    return err
}

// Dequeue извлекает следующую задачу из очереди
func (q *EtcdQueue) Dequeue(ctx context.Context, workerID string) (*Task, error) {
    // Получаем задачи из очереди, отсортированные по приоритету
    resp, err := q.client.Get(
        ctx, 
        TaskQueuePrefix,
        clientv3.WithPrefix(),
        clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
        clientv3.WithLimit(1),
    )
    if err != nil {
        return nil, err
    }
    
    if len(resp.Kvs) == 0 {
        return nil, nil // Очередь пуста
    }
    
    kv := resp.Kvs[0]
    var task Task
    if err := json.Unmarshal(kv.Value, &task); err != nil {
        return nil, err
    }
    
    // Атомарно перемещаем задачу в processing
    task.Status = TaskStatusRunning
    now := time.Now()
    task.StartedAt = &now
    
    data, err := json.Marshal(task)
    if err != nil {
        return nil, err
    }
    
    // Используем транзакцию для атомарности
    txnResp, err := q.client.Txn(ctx).
        If(clientv3.Compare(clientv3.CreateRevision(string(kv.Key)), "=", kv.CreateRevision)).
        Then(
            clientv3.OpDelete(string(kv.Key)),
            clientv3.OpPut(q.buildProcessingKey(workerID, task.ID), string(data)),
        ).
        Commit()
    
    if err != nil || !txnResp.Succeeded {
        return nil, fmt.Errorf("failed to move task to processing: %w", err)
    }
    
    return &task, nil
}

// Complete завершает задачу успешно
func (q *EtcdQueue) Complete(ctx context.Context, workerID, taskID string, result interface{}) error {
    processingKey := q.buildProcessingKey(workerID, taskID)
    
    // Получаем текущую задачу
    resp, err := q.client.Get(ctx, processingKey)
    if err != nil {
        return err
    }
    if len(resp.Kvs) == 0 {
        return fmt.Errorf("task not found in processing")
    }
    
    var task Task
    if err := json.Unmarshal(resp.Kvs[0].Value, &task); err != nil {
        return err
    }
    
    task.Status = TaskStatusCompleted
    now := time.Now()
    task.EndedAt = &now
    task.Result = result
    
    data, err := json.Marshal(task)
    if err != nil {
        return err
    }
    
    // Перемещаем в результаты и удаляем из processing
    _, err = q.client.Txn(ctx).
        Then(
            clientv3.OpDelete(processingKey),
            clientv3.OpPut(q.buildResultKey(taskID), string(data)),
        ).
        Commit()
    
    return err
}

// Fail отмечает задачу как неудачную
func (q *EtcdQueue) Fail(ctx context.Context, workerID, taskID string, errMsg string) error {
    processingKey := q.buildProcessingKey(workerID, taskID)
    
    resp, err := q.client.Get(ctx, processingKey)
    if err != nil {
        return err
    }
    if len(resp.Kvs) == 0 {
        return fmt.Errorf("task not found in processing")
    }
    
    var task Task
    if err := json.Unmarshal(resp.Kvs[0].Value, &task); err != nil {
        return err
    }
    
    task.RetryCount++
    
    if task.RetryCount >= task.MaxRetries {
        task.Status = TaskStatusFailed
        now := time.Now()
        task.EndedAt = &now
        task.Error = errMsg
        
        data, err := json.Marshal(task)
        if err != nil {
            return err
        }
        
        // Перемещаем в результаты как failed
        _, err = q.client.Txn(ctx).
            Then(
                clientv3.OpDelete(processingKey),
                clientv3.OpPut(q.buildResultKey(taskID), string(data)),
            ).
            Commit()
        return err
    } else {
        // Возвращаем в очередь для повторной попытки
        task.Status = TaskStatusPending
        task.StartedAt = nil
        
        data, err := json.Marshal(task)
        if err != nil {
            return err
        }
        
        key := q.buildQueueKey(task.Priority, task.ID)
        _, err = q.client.Txn(ctx).
            Then(
                clientv3.OpDelete(processingKey),
                clientv3.OpPut(key, string(data)),
            ).
            Commit()
        return err
    }
}

// GetTaskResult получает результат выполненной задачи
func (q *EtcdQueue) GetTaskResult(ctx context.Context, taskID string) (*Task, error) {
    resp, err := q.client.Get(ctx, q.buildResultKey(taskID))
    if err != nil {
        return nil, err
    }
    if len(resp.Kvs) == 0 {
        return nil, fmt.Errorf("task result not found")
    }
    
    var task Task
    if err := json.Unmarshal(resp.Kvs[0].Value, &task); err != nil {
        return nil, err
    }
    
    return &task, nil
}

// WatchQueue следит за появлением новых задач
func (q *EtcdQueue) WatchQueue(ctx context.Context) clientv3.WatchChan {
    return q.client.Watch(ctx, TaskQueuePrefix, clientv3.WithPrefix())
}

// buildQueueKey строит ключ для очереди с учётом приоритета
func (q *EtcdQueue) buildQueueKey(priority int, taskID string) string {
    // Инвертируем приоритет для сортировки (меньшее число = выше приоритет)
    invertedPriority := 9999 - priority
    return fmt.Sprintf("%s%04d/%s", TaskQueuePrefix, invertedPriority, taskID)
}

// buildProcessingKey строит ключ для выполняемых задач
func (q *EtcdQueue) buildProcessingKey(workerID, taskID string) string {
    return fmt.Sprintf("%s%s/%s", TaskProcessingPrefix, workerID, taskID)
}

// buildResultKey строит ключ для результатов задач
func (q *EtcdQueue) buildResultKey(taskID string) string {
    return fmt.Sprintf("%s%s", TaskResultPrefix, taskID)
}

// CleanupStaleTasks очищает зависшие задачи (например, после краша воркера)
func (q *EtcdQueue) CleanupStaleTasks(ctx context.Context, timeout time.Duration) error {
    resp, err := q.client.Get(ctx, TaskProcessingPrefix, clientv3.WithPrefix())
    if err != nil {
        return err
    }
    
    now := time.Now()
    for _, kv := range resp.Kvs {
        var task Task
        if err := json.Unmarshal(kv.Value, &task); err != nil {
            continue
        }
        
        if task.StartedAt != nil && now.Sub(*task.StartedAt) > timeout {
            // Задача выполняется слишком долго, возвращаем в очередь
            task.Status = TaskStatusPending
            task.StartedAt = nil
            
            data, err := json.Marshal(task)
            if err != nil {
                continue
            }
            
            key := q.buildQueueKey(task.Priority, task.ID)
            q.client.Txn(ctx).
                Then(
                    clientv3.OpDelete(string(kv.Key)),
                    clientv3.OpPut(key, string(data)),
                ).
                Commit()
        }
    }
    
    return nil
}
```

#### 6.2 Worker пул для обработки задач

Создать `internal/pkg/taskqueue/worker_pool.go`:

```go
package taskqueue

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
    
    "go.etcd.io/etcd/client/v3"
    "github.com/google/uuid"
)

// Worker обработчик задач
type Worker struct {
    ID       string
    queue    *EtcdQueue
    handlers map[string]TaskHandlerFunc
    ctx      context.Context
    cancel   context.CancelFunc
}

// TaskHandlerFunc функция обработки задачи
type TaskHandlerFunc func(ctx context.Context, payload map[string]interface{}) error

// WorkerPool пул воркеров
type WorkerPool struct {
    workers   []*Worker
    queue     *EtcdQueue
    wg        sync.WaitGroup
    ctx       context.Context
    cancel    context.CancelFunc
}

// NewWorker создаёт нового воркера
func NewWorker(queue *EtcdQueue, handlers map[string]TaskHandlerFunc) *Worker {
    ctx, cancel := context.WithCancel(context.Background())
    return &Worker{
        ID:       uuid.New().String(),
        queue:    queue,
        handlers: handlers,
        ctx:      ctx,
        cancel:   cancel,
    }
}

// Start запускает воркера
func (w *Worker) Start() {
    log.Printf("Worker %s started", w.ID)
    
    for {
        select {
        case <-w.ctx.Done():
            log.Printf("Worker %s stopped", w.ID)
            return
        default:
            task, err := w.queue.Dequeue(w.ctx, w.ID)
            if err != nil {
                log.Printf("Worker %s dequeue error: %v", w.ID, err)
                time.Sleep(100 * time.Millisecond)
                continue
            }
            
            if task == nil {
                // Очередь пуста, ждём
                time.Sleep(500 * time.Millisecond)
                continue
            }
            
            // Обрабатываем задачу
            handler, ok := w.handlers[task.Type]
            if !ok {
                log.Printf("Unknown task type: %s", task.Type)
                w.queue.Fail(w.ctx, w.ID, task.ID, fmt.Sprintf("unknown task type: %s", task.Type))
                continue
            }
            
            err = handler(w.ctx, task.Payload)
            if err != nil {
                log.Printf("Task %s failed: %v", task.ID, err)
                w.queue.Fail(w.ctx, w.ID, task.ID, err.Error())
            } else {
                log.Printf("Task %s completed", task.ID)
                w.queue.Complete(w.ctx, w.ID, task.ID, nil)
            }
        }
    }
}

// Stop останавливает воркера
func (w *Worker) Stop() {
    w.cancel()
}

// NewWorkerPool создаёт пул воркеров
func NewWorkerPool(queue *EtcdQueue, handlers map[string]TaskHandlerFunc, poolSize int) *WorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    
    pool := &WorkerPool{
        workers: make([]*Worker, poolSize),
        queue:   queue,
        ctx:     ctx,
        cancel:  cancel,
    }
    
    for i := 0; i < poolSize; i++ {
        pool.workers[i] = NewWorker(queue, handlers)
    }
    
    return pool
}

// Start запускает все воркеры в пуле
func (p *WorkerPool) Start() {
    for _, worker := range p.workers {
        p.wg.Add(1)
        go func(w *Worker) {
            defer p.wg.Done()
            w.Start()
        }(worker)
    }
    log.Printf("Started %d workers", len(p.workers))
}

// Stop останавливает все воркеры
func (p *WorkerPool) Stop() {
    p.cancel()
    for _, worker := range p.workers {
        worker.Stop()
    }
    p.wg.Wait()
    log.Println("All workers stopped")
}

// RegisterHandler регистрирует обработчик для типа задач
func (p *WorkerPool) RegisterHandler(taskType string, handler TaskHandlerFunc) {
    for _, worker := range p.workers {
        worker.handlers[taskType] = handler
    }
}
```

#### 6.3 Интеграция с Job моделью

Создать `internal/application/core/job_service.go`:

```go
package core

import (
    "context"
    "time"
    
    "github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
    "github.com/AlekseyPromet/netbox_go/internal/repository"
    "github.com/AlekseyPromet/netbox_go/internal/pkg/taskqueue"
)

// JobService сервис для управления задачами
type JobService struct {
    jobRepo   repository.JobRepository
    taskQueue *taskqueue.EtcdQueue
}

func NewJobService(jobRepo repository.JobRepository, taskQueue *taskqueue.EtcdQueue) *JobService {
    return &JobService{
        jobRepo:   jobRepo,
        taskQueue: taskQueue,
    }
}

// CreateJobParams параметры создания задачи
type CreateJobParams struct {
    ObjectType  string
    ObjectID    string
    Name        string
    Description string
    Interval    string
    ScheduledAt *time.Time
}

// CreateJob создаёт новую задачу
func (s *JobService) CreateJob(ctx context.Context, params CreateJobParams) (*entity.Job, error) {
    job := &entity.Job{
        ObjectType:  params.ObjectType,
        ObjectID:    params.ObjectID,
        Name:        params.Name,
        Description: params.Description,
        Interval:    params.Interval,
        ScheduledAt: params.ScheduledAt,
        Status:      "pending",
    }
    
    if err := s.jobRepo.Create(ctx, job); err != nil {
        return nil, err
    }
    
    // Если задача должна быть выполнена немедленно, добавляем в очередь
    if params.ScheduledAt == nil || params.ScheduledAt.Before(time.Now()) {
        task := &taskqueue.Task{
            Type: taskqueue.TypeProcessJob,
            Payload: map[string]interface{}{
                "job_id": job.ID,
            },
            Priority:   5,
            MaxRetries: 3,
        }
        
        if err := s.taskQueue.Enqueue(ctx, task); err != nil {
            return nil, err
        }
    }
    
    return job, nil
}

// ScheduleJob планирует задачу на выполнение
func (s *JobService) ScheduleJob(ctx context.Context, jobID string, scheduledAt time.Time) error {
    job, err := s.jobRepo.GetByID(ctx, jobID)
    if err != nil {
        return err
    }
    
    job.ScheduledAt = &scheduledAt
    job.Status = "scheduled"
    
    return s.jobRepo.Update(ctx, job)
}

// CancelJob отменяет задачу
func (s *JobService) CancelJob(ctx context.Context, jobID string) error {
    job, err := s.jobRepo.GetByID(ctx, jobID)
    if err != nil {
        return err
    }
    
    job.Status = "cancelled"
    return s.jobRepo.Update(ctx, job)
}
```

---

### Этап 7: Change Logging система

#### 7.1 Сервис логирования изменений

Создать `internal/application/core/changelog_service.go`:

```go
package core

import (
    "context"
    "encoding/json"
    "time"
    
    "github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
    "github.com/AlekseyPromet/netbox_go/internal/repository"
    "github.com/AlekseyPromet/netbox_go/pkg/types"
    "github.com/google/uuid"
)

// ChangeLogService сервис для логирования изменений объектов
type ChangeLogService struct {
    objectChangeRepo repository.ObjectChangeRepository
    objectTypeRepo   repository.ObjectTypeRepository
}

func NewChangeLogService(
    ocRepo repository.ObjectChangeRepository,
    otRepo repository.ObjectTypeRepository,
) *ChangeLogService {
    return &ChangeLogService{
        objectChangeRepo: ocRepo,
        objectTypeRepo:   otRepo,
    }
}

// LogChangeParams параметры для логирования изменения
type LogChangeParams struct {
    Action          types.Status
    ObjectType      string
    ObjectID        string
    ObjectRepr      string
    PreChangeData   interface{}
    PostChangeData  interface{}
    UserID          *types.ID
    UserName        string
    RequestID       *string
    RelatedObject   *RelatedObjectInfo
    Message         string
}

// RelatedObjectInfo информация о связанном объекте
type RelatedObjectInfo struct {
    ObjectType string
    ObjectID   string
    ObjectRepr string
}

// LogChange логирует изменение объекта
func (s *ChangeLogService) LogChange(ctx context.Context, params LogChangeParams) error {
    // Сериализация данных изменений
    var preChangeJSON, postChangeJSON json.RawMessage
    var err error
    
    if params.PreChangeData != nil {
        preChangeJSON, err = json.Marshal(params.PreChangeData)
        if err != nil {
            return err
        }
    }
    
    if params.PostChangeData != nil {
        postChangeJSON, err = json.Marshal(params.PostChangeData)
        if err != nil {
            return err
        }
    }
    
    // Генерация request_id если не предоставлен
    requestID := params.RequestID
    if requestID == nil {
        id := uuid.New().String()
        requestID = &id
    }
    
    oc := &entity.ObjectChange{
        Time:              time.Now(),
        UserID:            params.UserID,
        RequestID:         requestID,
        Action:            params.Action,
        ChangedObjectType: params.ObjectType,
        ChangedObjectID:   params.ObjectID,
        ObjectRepr:        params.ObjectRepr,
        RelatedObjectType: nil,
        RelatedObjectID:   nil,
        RelatedObjectRepr: nil,
    }
    
    if params.RelatedObject != nil {
        oc.RelatedObjectType = &params.RelatedObject.ObjectType
        oc.RelatedObjectID = &params.RelatedObject.ObjectID
        oc.RelatedObjectRepr = &params.RelatedObject.ObjectRepr
    }
    
    return s.objectChangeRepo.Create(ctx, oc)
}

// GetObjectHistory возвращает историю изменений объекта
func (s *ChangeLogService) GetObjectHistory(
    ctx context.Context, 
    objectType string, 
    objectID string,
    limit int,
) ([]*entity.ObjectChange, error) {
    filter := repository.ObjectChangeFilter{
        ChangedObjectType: &objectType,
        ChangedObjectID:   &objectID,
        Limit:             limit,
    }
    
    changes, _, err := s.objectChangeRepo.List(ctx, filter)
    return changes, err
}

// GetRecentChanges возвращает последние изменения
func (s *ChangeLogService) GetRecentChanges(
    ctx context.Context,
    limit int,
    since *time.Time,
) ([]*entity.ObjectChange, error) {
    filter := repository.ObjectChangeFilter{
        Since: since,
        Limit: limit,
    }
    
    changes, _, err := s.objectChangeRepo.List(ctx, filter)
    return changes, err
}
```

---

### Этап 8: Тестирование

#### 8.1 Unit тесты

Создать `internal/repository/postgres/core/core_test.go`:

```go
package core

import (
    "context"
    "testing"
    "time"
    
    "github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
    "github.com/AlekseyPromet/netbox_go/pkg/types"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestDataSourceRepository_CreateAndGet(t *testing.T) {
    // Setup
    repo := NewDataSourceRepository(testDB)
    ctx := context.Background()
    
    ds := &entity.DataSource{
        ID:         types.NewID(),
        Name:       "test-source",
        Type:       "local",
        SourceURL:  "file:///tmp/test",
        Status:     "new",
        Enabled:    true,
        Created:    time.Now(),
        Updated:    time.Now(),
    }
    
    // Create
    err := repo.Create(ctx, ds)
    require.NoError(t, err)
    
    // Get
    retrieved, err := repo.GetByID(ctx, ds.ID.String())
    require.NoError(t, err)
    assert.Equal(t, ds.Name, retrieved.Name)
    assert.Equal(t, ds.Type, retrieved.Type)
}

func TestObjectChangeRepository_LogChange(t *testing.T) {
    // Setup
    repo := NewObjectChangeRepository(testDB)
    ctx := context.Background()
    
    // Log change
    err := repo.LogChange(ctx, "dcim.device", "123", "device1", 
        "create", nil, nil, nil)
    require.NoError(t, err)
    
    // Verify
    filter := repository.ObjectChangeFilter{
        ChangedObjectType: ptr("dcim.device"),
        ChangedObjectID:   ptr("123"),
    }
    changes, count, err := repo.List(ctx, filter)
    require.NoError(t, err)
    assert.Equal(t, int64(1), count)
    assert.Len(t, changes, 1)
}
```

#### 8.2 Integration тесты

Создать `tests/integration/core_integration_test.go`:

```go
package integration

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    
    "github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
    "github.com/labstack/echo/v4"
    "github.com/stretchr/testify/assert"
)

func TestDataSourceCRUD(t *testing.T) {
    e := echo.New()
    server := httptest.NewServer(e)
    defer server.Close()
    
    // Create DataSource
    ds := entity.DataSource{
        Name:      "test-git",
        Type:      "git",
        SourceURL: "https://github.com/example/configs.git",
        Status:    "new",
        Enabled:   true,
    }
    
    body, _ := json.Marshal(ds)
    resp, err := http.Post(server.URL+"/api/core/data-sources", 
        "application/json", bytes.NewReader(body))
    assert.NoError(t, err)
    assert.Equal(t, http.StatusCreated, resp.StatusCode)
    
    // Get DataSource
    resp, err = http.Get(server.URL + "/api/core/data-sources/" + ds.ID.String())
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // List DataSources
    resp, err = http.Get(server.URL + "/api/core/data-sources")
    assert.NoError(t, err)
    assert.Equal(t, http.StatusOK, resp.StatusCode)
    
    // Delete DataSource
    req, _ := http.NewRequest(http.MethodDelete, 
        server.URL+"/api/core/data-sources/"+ds.ID.String(), nil)
    client := &http.Client{}
    resp, err = client.Do(req)
    assert.NoError(t, err)
    assert.Equal(t, http.StatusNoContent, resp.StatusCode)
}
```

---

### Этап 9: Документация

#### 9.1 Обновить документацию API

Добавить секцию в `netbox_go/docs/api.md`:

```markdown
## Core API

### Data Sources

#### GET /api/core/data-sources/
Получить список источников данных

**Query Parameters:**
- `name` (string) — фильтр по имени
- `type` (string) — фильтр по типу
- `status` (string) — фильтр по статусу
- `enabled` (boolean) — фильтр по состоянию
- `limit` (int) — количество записей (default: 100, max: 1000)
- `offset` (int) — смещение

**Response:**
```json
{
  "count": 10,
  "next": "/api/core/data-sources/?limit=10&offset=10",
  "previous": null,
  "results": [
    {
      "id": "uuid",
      "name": "example-git",
      "type": "git",
      "source_url": "https://...",
      "status": "completed",
      "enabled": true,
      "sync_interval": 60,
      "last_synced": "2024-01-01T00:00:00Z"
    }
  ]
}
```

#### POST /api/core/data-sources/
Создать новый источник данных

**Request Body:**
```json
{
  "name": "my-configs",
  "type": "git",
  "source_url": "https://github.com/org/configs.git",
  "enabled": true,
  "sync_interval": 60
}
```

#### POST /api/core/data-sources/:id/sync/
Запустить синхронизацию источника данных

**Response:**
```json
{
  "status": "sync initiated"
}
```

### Jobs

#### GET /api/core/jobs/
Получить список фоновых задач

**Query Parameters:**
- `status` (string) — фильтр по статусу
- `queue_name` (string) — фильтр по очереди
- `scheduled_at` (datetime) — фильтр по времени планирования

### Object Changes

#### GET /api/core/object-changes/
Получить журнал изменений объектов

**Query Parameters:**
- `changed_object_type` (string) — тип изменённого объекта
- `changed_object_id` (string) — ID изменённого объекта
- `user_id` (string) — фильтр по пользователю
- `action` (string) — тип действия (create/update/delete)
- `since` (datetime) — изменения после даты
- `until` (datetime) — изменения до даты
```

---

## Чеклист завершения миграции

### Обязательные компоненты

- [ ] Интерфейсы репозиториев определены в `internal/repository/interfaces.go`
- [ ] PostgreSQL репозитории реализованы в `internal/repository/postgres/core/`
- [ ] SQL миграции созданы в `netbox_go/migrations/`
- [ ] HTTP обработчики обновлены с полным CRUD
- [ ] Маршруты API зарегистрированы в роутере
- [ ] Интеграция с EtcdQueue для фоновых задач
- [ ] Worker Pool настроен и запущен
- [ ] Change Logging сервис реализован
- [ ] Unit тесты написаны (покрытие > 80%)
- [ ] Integration тесты проходят успешно

### Дополнительные компоненты

- [ ] GraphQL резолверы для Core сущностей
- [ ] Кэширование ObjectType в etcd
- [ ] Bulk операции для DataFile
- [ ] Автоматическая синхронизация по расписанию
- [ ] Мониторинг и метрики для задач
- [ ] Документация API обновлена
- [ ] Примеры использования в `examples/`

---

## Риски и зависимости

### Риски

1. **Сложность реализации очереди на etcd**
   - Необходимо обеспечить атомарность операций
   - Требуется обработка зависших задач
   - Мониторинг состояния очереди
   
2. **Производительность ObjectChange**
   - Большой объём записей изменений
   - Необходима оптимизация запросов и индексов

3. **Синхронизация Data Sources**
   - Различные бэкенды (git, local, S3)
   - Обработка ошибок сети и файловых систем

### Зависимости

1. **Базовая инфраструктура**
   - PostgreSQL настроен и доступен
   - etcd кластер развёрнут (минимум 3 ноды для HA)

2. **Смежные модули**
   - Модуль `account` для работы с пользователями
   - Модуль `extras` для интеграции со скриптами
   - Система плагинов

---

## Поддержка вендоров АСУ ТП и сетевых протоколов

Для расширения функциональности NetBox в области промышленной автоматизации и сетевого мониторинга добавлена поддержка вендоров и протоколов:

### Список поддерживаемых вендоров

| № | Вендор | Система / Продукт | Протокол / Интерфейс | Приоритет | Статус |
|---|--------|-------------------|----------------------|-----------|--------|
| 1 | **ABB** | 800xA | OPC DA/UA, Modbus TCP | Высокий | Планируется |
| 2 | **ABB** | Symphony/Harmony | INFI-90 Protocol | Высокий | Планируется |
| 3 | **ABB** | Infi90 | INFI-90 Loop Interface | Средний | Планируется |
| 4 | **ABB** | Network Manager | DNP3, IEC 61850 | Высокий | Планируется |
| 5 | **ABB** | FACTS | Proprietary ABB | Низкий | Планируется |
| 6 | **ABB** | SYS600 | IEC 61850, LON | Средний | Планируется |
| 7 | **ABB** | MicroSCADA | DNP3, Modbus, IEC 60870-5-104 | Высокий | Планируется |
| 8 | **Automsoft** | RAPID Historian | ODBC, API REST | Средний | Планируется |
| 9 | **Emerson** | DeltaV | OPC DA/UA, Modbus TCP | Высокий | Планируется |
| 10 | **Emerson** | Ovation | OPC, Modbus, Serial | Высокий | Планируется |
| 11 | **Emerson/Westinghouse** | WDPF | Vnet/IP, Modbus | Средний | Планируется |
| 12 | **GE** | XA/21 | GEnet, SRTP | Средний | Планируется |
| 13 | **GE** | PowerOn Fusion | Proficy Historian API | Средний | Планируется |
| 14 | **Foxboro (Schneider)** | I/A Series | Foxboro Protocol, Modbus | Средний | Планируется |
| 15 | **Honeywell** | Experion | OPC UA, EtherNet/IP | Высокий | Планируется |
| 16 | **Itron** | OpenWay System | ANSI C12.18/C12.19, DLMS | Низкий | Планируется |
| 17 | **Rockwell** | RSView (FactoryTalk) | OPC DA/UA, Allen-Bradley DF1 | Высокий | Планируется |
| 18 | **Schneider/Telvent** | Oasys | Modbus, OPC | Низкий | Планируется |
| 19 | **Schneider** | Citect | OPC, Modbus, Ethernet | Средний | Планируется |
| 20 | **Schneider** | Momentum | Modbus TCP/RTU | Средний | Планируется |
| 21 | **Schneider** | Quantum | Modbus, Unity Protocol | Высокий | Планируется |
| 22 | **Siemens** | PCS7 | S7 Comm, OPC UA, Profibus | Высокий | Планируется |
| 23 | **Yokogawa** | Centrum CS 3000 | Vnet/IP, Modbus, OPC | Высокий | Планируется |

### Список поддерживаемых протоколов

#### Сетевые протоколы управления и мониторинга

| № | Протокол | Версии / Варианты | Назначение | Приоритет | Статус |
|---|----------|-------------------|------------|-----------|--------|
| 1 | **SNMP** | v1, v2c, v3 | Управление сетевыми устройствами | Высокий | Планируется |
| 2 | **ICMP** | IPv4, IPv6 | Ping, трассировка маршрута | Высокий | Планируется |
| 3 | **Syslog** | RFC 3164, RFC 5424 | Сбор логов | Высокий | Планируется |
| 4 | **NetFlow** | v5, v9, IPFIX | Анализ трафика | Высокий | Планируется |
| 5 | **BACnet** | BACnet/IP, MS/TP | Автоматизация зданий | Средний | Планируется |
| 6 | **Modbus** | TCP, RTU, ASCII | Промышленная автоматизация | Высокий | Планируется |
| 7 | **OPC** | DA, UA, HDA | Промышленная интеграция | Высокий | Планируется |
| 8 | **WMI** | Windows Management Instrumentation | Мониторинг Windows | Средний | Планируется |
| 9 | **DHCP** | v4, v6 | Мониторинг аренды адресов | Средний | Планируется |
| 10 | **DNS** | A, AAAA, MX, TXT записи | Мониторинг DNS | Высокий | Планируется |

#### Протоколы прикладного уровня

| № | Протокол | Версии / Варианты | Назначение | Приоритет | Статус |
|---|----------|-------------------|------------|-----------|--------|
| 11 | **HTTP/HTTPS** | 1.1, 2, 3 | Web-мониторинг, API | Высокий | Планируется |
| 12 | **FTP** | FTP, FTPS, SFTP | Передача файлов | Средний | Планируется |
| 13 | **SSH** | SSH-2 | Удалённое управление | Высокий | Планируется |
| 14 | **SMTP** | SMTP, ESMTP | Отправка почты | Средний | Планируется |
| 15 | **IMAP** | IMAP4 | Проверка почты | Низкий | Планируется |
| 16 | **POP3** | POP3, POP3S | Проверка почты | Низкий | Планируется |
| 17 | **LDAP** | LDAPv3, LDAPS | Каталог пользователей | Высокий | Планируется |
| 18 | **Radius** | RADIUS, RadSec | Аутентификация | Высокий | Планируется |
| 19 | **JMX** | JMX RMI, JMXMP | Мониторинг Java | Средний | Планируется |
| 20 | **TCP** | Raw TCP | Универсальный мониторинг | Высокий | Планируется |

#### Протоколы баз данных и интеграции

| № | Протокол | Версии / Варианты | Назначение | Приоритет | Статус |
|---|----------|-------------------|------------|-----------|--------|
| 21 | **JDBC** | Все драйверы | Подключение к БД | Средний | Планируется |
| 22 | **SQL** | Native protocols | Запросы к БД | Высокий | Планируется |
| 23 | **CORBA** | IIOP | Legacy интеграция | Низкий | Планируется |
| 24 | **WBEM** | CIM-XML, WS-Man | Управление предприятием | Низкий | Планируется |

#### Специализированные протоколы

| № | Протокол | Версии / Варианты | Назначение | Приоритет | Статус |
|---|----------|-------------------|------------|-----------|--------|
| 25 | **Keytroller** | Proprietary | Контроль доступа | Низкий | Планируется |
| 26 | **NMEA 0183** | GPS/GLONASS | Навигационные данные | Низкий | Планируется |
| 27 | **BMP** | BMP Monitoring | Мониторинг BMP | Низкий | Планируется |
| 28 | **BGP** | BGP-4 | Маршрутизация | Высокий | Планируется |

### Требования к реализации поддержки протоколов

#### 1. Расширение модели данных для протоколов

```go
// internal/domain/core/entity/protocol.go

package entity

type Protocol struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`         // SNMP, Modbus, OPC UA
    Slug        string            `json:"slug"`         // snmp, modbus, opc-ua
    Description string            `json:"description"`
    Category    string            `json:"category"`     // Network, Industrial, Database
    Versions    []string          `json:"versions"`     // v1, v2c, v3
    DefaultPort int               `json:"default_port"` // 161, 502, 4840
    Transport   string            `json:"transport"`    // TCP, UDP, Both
    Encrypted   bool              `json:"encrypted"`    // TLS/SSL support
    AuthMethods []string          `json:"auth_methods"` // None, Password, Certificate
    Properties  map[string]string `json:"properties"`   // Специфические свойства
    Created     time.Time         `json:"created"`
    Updated     time.Time         `json:"updated"`
}

type ProtocolConfig struct {
    ID              string            `json:"id"`
    ProtocolID      string            `json:"protocol_id"`
    DeviceID        string            `json:"device_id"`
    Endpoint        string            `json:"endpoint"`        // IP:Port
    Version         string            `json:"version"`         // v2c, v3, TCP
    Community       string            `json:"community"`       // SNMP community
    Username        string            `json:"username"`
    Password        string            `json:"password"`
    AuthProtocol    string            `json:"auth_protocol"`   // MD5, SHA
    PrivProtocol    string            `json:"priv_protocol"`   // DES, AES
    Certificates    []TLSCertificate  `json:"certificates"`
    Timeout         int               `json:"timeout"`         // секунды
    RetryCount      int               `json:"retry_count"`
    PollInterval    int               `json:"poll_interval"`   // секунды
    OIDs            []string          `json:"oids"`            // SNMP OIDs
    Registers       []ModbusRegister  `json:"registers"`       // Modbus registers
    Tags            []string          `json:"tags"`
    Status          string            `json:"status"`          // active, inactive, error
    LastPollTime    *time.Time        `json:"last_poll_time"`
    LastPollStatus  string            `json:"last_poll_status"`
    Metadata        map[string]string `json:"metadata"`
}

type ModbusRegister struct {
    Address    uint16 `json:"address"`
    Type       string `json:"type"`       // Coil, DiscreteInput, HoldingRegister, InputRegister
    DataType   string `json:"data_type"`  // BOOL, INT16, UINT16, INT32, FLOAT32
    Name       string `json:"name"`
    Unit       string `json:"unit"`
    Multiplier float64 `json:"multiplier"`
}

type TLSCertificate struct {
    ID         string    `json:"id"`
    Name       string    `json:"name"`
    Type       string    `json:"type"`        // CA, Client, Server
    Content    string    `json:"content"`     // PEM encoded
    ExpiryDate time.Time `json:"expiry_date"`
    Fingerprint string   `json:"fingerprint"`
}
```

#### 2. Интерфейс универсального поллера протоколов

```go
// internal/pkg/pollers/poller.go

package pollers

type Poller interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    Read(ctx context.Context, query Query) ([]Result, error)
    Write(ctx context.Context, command Command) error
    Subscribe(ctx context.Context, subscription Subscription, callback func(Result)) error
    HealthCheck(ctx context.Context) error
    GetMetrics() PollerMetrics
}

type Query struct {
    Protocol    string                 `json:"protocol"`
    Targets     []string               `json:"targets"`
    Parameters  map[string]interface{} `json:"parameters"`
    Timeout     time.Duration          `json:"timeout"`
    RetryCount  int                    `json:"retry_count"`
}

type Result struct {
    Target    string                 `json:"target"`
    Timestamp time.Time              `json:"timestamp"`
    Success   bool                   `json:"success"`
    Data      interface{}            `json:"data"`
    Error     string                 `json:"error,omitempty"`
    Latency   time.Duration          `json:"latency"`
}

type Command struct {
    Protocol    string                 `json:"protocol"`
    Target      string                 `json:"target"`
    Action      string                 `json:"action"`
    Parameters  map[string]interface{} `json:"parameters"`
}

type Subscription struct {
    Protocol   string                 `json:"protocol"`
    Target     string                 `json:"target"`
    Interval   time.Duration          `json:"interval"`
    Parameters map[string]interface{} `json:"parameters"`
}

type PollerMetrics struct {
    TotalRequests   int64         `json:"total_requests"`
    SuccessfulReqs  int64         `json:"successful_requests"`
    FailedReqs      int64         `json:"failed_requests"`
    AvgLatency      time.Duration `json:"avg_latency"`
    LastPollTime    time.Time     `json:"last_poll_time"`
    ConnectionState string        `json:"connection_state"`
}
```

#### 3. Реестр поллеров протоколов

```go
// internal/pkg/pollers/registry.go

package pollers

type Registry struct {
    pollers map[string]PollerFactory
    mu      sync.RWMutex
}

type PollerFactory func(config ProtocolConfig) (Poller, error)

func (r *Registry) Register(protocol string, factory PollerFactory) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.pollers[protocol] = factory
}

func (r *Registry) Get(protocol string, config ProtocolConfig) (Poller, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    factory, ok := r.pollers[protocol]
    if !ok {
        return nil, ErrPollerNotFound
    }
    return factory(config)
}

func (r *Registry) ListProtocols() []string {
    r.mu.RLock()
    defer r.mu.RUnlock()
    protocols := make([]string, 0, len(r.pollers))
    for proto := range r.pollers {
        protocols = append(protocols, proto)
    }
    return protocols
}
```

### Примеры реализаций поллеров

#### SNMP Poller (v1, v2c, v3)

```go
// internal/pkg/pollers/snmp_poller.go

package pollers

import (
    "context"
    "fmt"
    "time"
    
    "github.com/gosnmp/gosnmp"
)

type SNMPPoller struct {
    client *gosnmp.GoSNMP
    config ProtocolConfig
    metrics PollerMetrics
}

func NewSNMPPoller(config ProtocolConfig) (Poller, error) {
    var version gosnmp.SnmpVersion
    
    switch config.Version {
    case "v1":
        version = gosnmp.Version1
    case "v2c":
        version = gosnmp.Version2c
    case "v3":
        version = gosnmp.Version3
    default:
        version = gosnmp.Version2c
    }
    
    client := &gosnmp.GoSNMP{
        Target:    config.Endpoint,
        Port:      uint16(config.DefaultPort),
        Community: config.Community,
        Version:   version,
        Timeout:   time.Duration(config.Timeout) * time.Second,
        Retries:   config.RetryCount,
    }
    
    if version == gosnmp.Version3 {
        client.SecurityModel = gosnmp.UserSecurityModel
        client.SecurityParameters = &gosnmp.UsmSecurityParameters{
            UserName:                 config.Username,
            AuthenticationProtocol:   getAuthProtocol(config.AuthProtocol),
            AuthenticationPassphrase: config.Password,
            PrivacyProtocol:          getPrivProtocol(config.PrivProtocol),
            PrivacyPassphrase:        config.Password,
        }
    }
    
    return &SNMPPoller{
        client: client,
        config: config,
        metrics: PollerMetrics{ConnectionState: "disconnected"},
    }, nil
}

func (p *SNMPPoller) Connect(ctx context.Context) error {
    return p.client.Connect()
}

func (p *SNMPPoller) Disconnect(ctx context.Context) error {
    return p.client.Conn.Close()
}

func (p *SNMPPoller) Read(ctx context.Context, query Query) ([]Result, error) {
    start := time.Now()
    p.metrics.TotalRequests++
    
    oids, ok := query.Parameters["oids"].([]string)
    if !ok {
        return nil, fmt.Errorf("invalid OIDs parameter")
    }
    
    result, err := p.client.Get(oids)
    latency := time.Since(start)
    
    if err != nil {
        p.metrics.FailedReqs++
        return []Result{{
            Target:    p.config.Endpoint,
            Timestamp: time.Now(),
            Success:   false,
            Error:     err.Error(),
            Latency:   latency,
        }}, err
    }
    
    p.metrics.SuccessfulReqs++
    p.metrics.AvgLatency = latency
    
    data := make(map[string]interface{})
    for _, variable := range result.Variables {
        data[variable.Name] = gosnmp.ToPrettyFormat(variable.Value)
    }
    
    return []Result{{
        Target:    p.config.Endpoint,
        Timestamp: time.Now(),
        Success:   true,
        Data:      data,
        Latency:   latency,
    }}, nil
}

func (p *SNMPPoller) HealthCheck(ctx context.Context) error {
    return p.client.Connect()
}

func (p *SNMPPoller) GetMetrics() PollerMetrics {
    return p.metrics
}

// Заглушки для методов Write и Subscribe
func (p *SNMPPoller) Write(ctx context.Context, command Command) error {
    return fmt.Errorf("SNMP write not implemented")
}

func (p *SNMPPoller) Subscribe(ctx context.Context, subscription Subscription, callback func(Result)) error {
    return fmt.Errorf("SNMP subscribe not implemented")
}

func getAuthProtocol(proto string) gosnmp.SnmpV3AuthProtocol {
    switch proto {
    case "MD5":
        return gosnmp.MD5
    case "SHA":
        return gosnmp.SHA
    case "SHA224":
        return gosnmp.SHA224
    case "SHA256":
        return gosnmp.SHA256
    case "SHA384":
        return gosnmp.SHA384
    case "SHA512":
        return gosnmp.SHA512
    default:
        return gosnmp.NoAuth
    }
}

func getPrivProtocol(proto string) gosnmp.SnmpV3PrivProtocol {
    switch proto {
    case "DES":
        return gosnmp.DES
    case "AES":
        return gosnmp.AES
    case "AES192":
        return gosnmp.AES192
    case "AES256":
        return gosnmp.AES256
    case "AES192C":
        return gosnmp.AES192C
    case "AES256C":
        return gosnmp.AES256C
    default:
        return gosnmp.NoPriv
    }
}
```

#### Modbus TCP Poller

```go
// internal/pkg/pollers/modbus_poller.go

package pollers

import (
    "context"
    "fmt"
    "time"
    
    "github.com/grid-x/modbus"
)

type ModbusPoller struct {
    client modbus.Client
    config ProtocolConfig
    metrics PollerMetrics
}

func NewModbusPoller(config ProtocolConfig) (Poller, error) {
    client := modbus.NewClient(config.Endpoint)
    
    return &ModbusPoller{
        client: client,
        config: config,
        metrics: PollerMetrics{ConnectionState: "disconnected"},
    }, nil
}

func (p *ModbusPoller) Connect(ctx context.Context) error {
    p.metrics.ConnectionState = "connected"
    return nil
}

func (p *ModbusPoller) Disconnect(ctx context.Context) error {
    p.metrics.ConnectionState = "disconnected"
    return nil
}

func (p *ModbusPoller) Read(ctx context.Context, query Query) ([]Result, error) {
    start := time.Now()
    p.metrics.TotalRequests++
    
    registers, ok := query.Parameters["registers"].([]ModbusRegister)
    if !ok {
        return nil, fmt.Errorf("invalid registers parameter")
    }
    
    data := make(map[string]interface{})
    
    for _, reg := range registers {
        var result []byte
        var err error
        
        switch reg.Type {
        case "HoldingRegister":
            result, err = p.client.ReadHoldingRegisters(reg.Address, 1)
        case "InputRegister":
            result, err = p.client.ReadInputRegisters(reg.Address, 1)
        case "Coil":
            result, err = p.client.ReadCoils(reg.Address, 1)
        case "DiscreteInput":
            result, err = p.client.ReadDiscreteInputs(reg.Address, 1)
        default:
            err = fmt.Errorf("unknown register type: %s", reg.Type)
        }
        
        if err != nil {
            p.metrics.FailedReqs++
            return []Result{{
                Target:    p.config.Endpoint,
                Timestamp: time.Now(),
                Success:   false,
                Error:     err.Error(),
                Latency:   time.Since(start),
            }}, err
        }
        
        value := parseModbusValue(result, reg.DataType)
        data[reg.Name] = value
    }
    
    latency := time.Since(start)
    p.metrics.SuccessfulReqs++
    p.metrics.AvgLatency = latency
    
    return []Result{{
        Target:    p.config.Endpoint,
        Timestamp: time.Now(),
        Success:   true,
        Data:      data,
        Latency:   latency,
    }}, nil
}

func (p *ModbusPoller) HealthCheck(ctx context.Context) error {
    _, err := p.client.ReadHoldingRegisters(0, 1)
    return err
}

func (p *ModbusPoller) GetMetrics() PollerMetrics {
    return p.metrics
}

// Заглушки
func (p *ModbusPoller) Write(ctx context.Context, command Command) error {
    return fmt.Errorf("Modbus write not implemented")
}

func (p *ModbusPoller) Subscribe(ctx context.Context, subscription Subscription, callback func(Result)) error {
    return fmt.Errorf("Modbus subscribe not implemented")
}

func parseModbusValue(data []byte, dataType string) interface{} {
    switch dataType {
    case "BOOL":
        return data[0] != 0
    case "INT16":
        return int16(data[0])<<8 | int16(data[1])
    case "UINT16":
        return uint16(data[0])<<8 | uint16(data[1])
    case "INT32":
        return int32(data[0])<<24 | int32(data[1])<<16 | int32(data[2])<<8 | int32(data[3])
    case "FLOAT32":
        // IEEE 754 Float32 conversion needed
        return float32(0)
    default:
        return data
    }
}
```

#### ICMP Poller (Ping)

```go
// internal/pkg/pollers/icmp_poller.go

package pollers

import (
    "context"
    "fmt"
    "net"
    "time"
    
    "github.com/go-ping/ping"
)

type ICMPPoller struct {
    config ProtocolConfig
    metrics PollerMetrics
}

func NewICMPPoller(config ProtocolConfig) (Poller, error) {
    return &ICMPPoller{
        config: config,
        metrics: PollerMetrics{ConnectionState: "ready"},
    }, nil
}

func (p *ICMPPoller) Connect(ctx context.Context) error {
    return nil
}

func (p *ICMPPoller) Disconnect(ctx context.Context) error {
    return nil
}

func (p *ICMPPoller) Read(ctx context.Context, query Query) ([]Result, error) {
    start := time.Now()
    p.metrics.TotalRequests++
    
    target := p.config.Endpoint
    timeout := time.Duration(p.config.Timeout) * time.Second
    
    pinger, err := ping.NewPinger(target)
    if err != nil {
        p.metrics.FailedReqs++
        return []Result{{
            Target:    target,
            Timestamp: time.Now(),
            Success:   false,
            Error:     err.Error(),
            Latency:   0,
        }}, err
    }
    
    pinger.SetPrivileged(true)
    pinger.Count = 1
    pinger.Timeout = timeout
    
    err = pinger.Run()
    if err != nil {
        p.metrics.FailedReqs++
        return []Result{{
            Target:    target,
            Timestamp: time.Now(),
            Success:   false,
            Error:     err.Error(),
            Latency:   0,
        }}, err
    }
    
    stats := pinger.Statistics()
    latency := time.Since(start)
    p.metrics.SuccessfulReqs++
    p.metrics.AvgLatency = stats.AvgRtt
    
    data := map[string]interface{}{
        "packets_sent":     stats.PacketsSent,
        "packets_recv":     stats.PacketsRecv,
        "packet_loss":      stats.PacketLoss,
        "min_rtt":          stats.MinRtt.String(),
        "avg_rtt":          stats.AvgRtt.String(),
        "max_rtt":          stats.MaxRtt.String(),
        "stddev_rtt":       stats.StdDevRtt.String(),
        "ip_addr":          stats.IPAddr.String(),
    }
    
    return []Result{{
        Target:    target,
        Timestamp: time.Now(),
        Success:   stats.PacketsRecv > 0,
        Data:      data,
        Latency:   latency,
    }}, nil
}

func (p *ICMPPoller) HealthCheck(ctx context.Context) error {
    _, err := net.LookupHost(p.config.Endpoint)
    return err
}

func (p *ICMPPoller) GetMetrics() PollerMetrics {
    return p.metrics
}

// Заглушки
func (p *ICMPPoller) Write(ctx context.Context, command Command) error {
    return fmt.Errorf("ICMP write not supported")
}

func (p *ICMPPoller) Subscribe(ctx context.Context, subscription Subscription, callback func(Result)) error {
    return fmt.Errorf("ICMP subscribe not implemented")
}
```

### Интеграция с системой фоновых задач etcd

Поллеры протоколов интегрируются с EtcdQueue для выполнения периодических опросов:

```go
// internal/application/core/polling_service.go

package core

import (
    "context"
    "encoding/json"
    "fmt"
    "time"
    
    "netbox/internal/pkg/taskqueue"
    "netbox/internal/pkg/pollers"
)

type PollingService struct {
    queue      *taskqueue.EtcdQueue
    registry   *pollers.Registry
    etcdClient *clientv3.Client
}

type PollingTask struct {
    ProtocolConfigID string            `json:"protocol_config_id"`
    Protocol         string            `json:"protocol"`
    Query            pollers.Query     `json:"query"`
    Schedule         string            `json:"schedule"` // cron expression
    Priority         int               `json:"priority"`
    Metadata         map[string]string `json:"metadata"`
}

func NewPollingService(queue *taskqueue.EtcdQueue, registry *pollers.Registry, etcdClient *clientv3.Client) *PollingService {
    return &PollingService{
        queue:      queue,
        registry:   registry,
        etcdClient: etcdClient,
    }
}

func (s *PollingService) SchedulePoll(ctx context.Context, config ProtocolConfig) error {
    task := PollingTask{
        ProtocolConfigID: config.ID,
        Protocol:         config.ProtocolID,
        Query: pollers.Query{
            Protocol:   config.ProtocolID,
            Targets:    []string{config.Endpoint},
            Parameters: buildQueryParams(config),
            Timeout:    time.Duration(config.Timeout) * time.Second,
            RetryCount: config.RetryCount,
        },
        Schedule: fmt.Sprintf("*/%d * * * *", config.PollInterval/60),
        Priority: 5,
        Metadata: map[string]string{
            "device_id": config.DeviceID,
            "source":    "polling_service",
        },
    }
    
    taskData, err := json.Marshal(task)
    if err != nil {
        return err
    }
    
    // Сохраняем расписание в etcd
    scheduleKey := fmt.Sprintf("/netbox/schedules/polling/%s", config.ID)
    _, err = s.etcdClient.Put(ctx, scheduleKey, string(taskData))
    if err != nil {
        return err
    }
    
    // Добавляем первую задачу в очередь
    return s.queue.Enqueue(ctx, "protocol_poll", taskData, task.Priority)
}

func (s *PollingService) ExecutePoll(ctx context.Context, taskData []byte) (*pollers.Result, error) {
    var task PollingTask
    if err := json.Unmarshal(taskData, &task); err != nil {
        return nil, err
    }
    
    // Получаем конфигурацию из БД
    config, err := s.getProtocolConfig(ctx, task.ProtocolConfigID)
    if err != nil {
        return nil, err
    }
    
    // Создаём поллер
    poller, err := s.registry.Get(task.Protocol, *config)
    if err != nil {
        return nil, err
    }
    
    // Подключаемся и выполняем опрос
    if err := poller.Connect(ctx); err != nil {
        return nil, err
    }
    defer poller.Disconnect(ctx)
    
    results, err := poller.Read(ctx, task.Query)
    if err != nil {
        return nil, err
    }
    
    if len(results) > 0 {
        // Сохраняем результаты в БД
        err = s.savePollResult(ctx, task.ProtocolConfigID, &results[0])
        if err != nil {
            return nil, err
        }
        
        // Генерируем событие изменения если статус изменился
        if results[0].Success != (config.Status == "active") {
            newStatus := "active"
            if !results[0].Success {
                newStatus = "error"
            }
            err = s.updateDeviceStatus(ctx, config.DeviceID, newStatus)
            if err != nil {
                return nil, err
            }
        }
        
        return &results[0], nil
    }
    
    return nil, fmt.Errorf("no results from poller")
}

func buildQueryParams(config ProtocolConfig) map[string]interface{} {
    params := make(map[string]interface{})
    
    switch config.ProtocolID {
    case "snmp":
        params["oids"] = config.OIDs
    case "modbus":
        params["registers"] = config.Registers
    }
    
    return params
}

// Заглушки для методов БД
func (s *PollingService) getProtocolConfig(ctx context.Context, id string) (*ProtocolConfig, error) {
    // TODO: Implement database query
    return &ProtocolConfig{}, nil
}

func (s *PollingService) savePollResult(ctx context.Context, configID string, result *pollers.Result) error {
    // TODO: Implement database insert
    return nil
}

func (s *PollingService) updateDeviceStatus(ctx context.Context, deviceID, status string) error {
    // TODO: Implement device status update with change logging
    return nil
}
```

### Обновлённые требования к реализации поддержки вендоров

#### 1. Расширение модели данных

Добавить новые сущности и поля в модуль Core:

```go
// internal/domain/core/entity/vendor.go

package entity

type Vendor struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`         // ABB, Emerson, Siemens, etc.
    Slug        string            `json:"slug"`         // abb, emerson, siemens
    Description string            `json:"description"`
    Systems     []VendorSystem    `json:"systems"`      // 800xA, DeltaV, PCS7, etc.
    Protocols   []string          `json:"protocols"`    // OPC UA, Modbus, etc.
    Created     time.Time         `json:"created"`
    Updated     time.Time         `json:"updated"`
}

type VendorSystem struct {
    ID           string   `json:"id"`
    VendorID     string   `json:"vendor_id"`
    Name         string   `json:"name"`         // 800xA, Symphony, DeltaV
    Version      string   `json:"version"`
    ProtocolType string   `json:"protocol_type"` // OPC UA, Modbus, INFI-90
    Properties   []string `json:"properties"`    // Специфические свойства
    ConfigSchema string   `json:"config_schema"` // JSON Schema для конфигурации
}

type Platform struct {
    ID              string            `json:"id"`
    Name            string            `json:"name"`
    VendorID        string            `json:"vendor_id"`
    SystemID        string            `json:"system_id"`
    ConnectionType  string            `json:"connection_type"` // TCP, Serial, OPC
    Endpoint        string            `json:"endpoint"`        // IP:Port или путь
    AuthMethod      string            `json:"auth_method"`     // None, Certificate, Password
    Certificates    []TLSCertificate  `json:"certificates"`
    Timeout         int               `json:"timeout"`         // секунды
    RetryCount      int               `json:"retry_count"`
    PollInterval    int               `json:"poll_interval"`   // секунды
    Status          string            `json:"status"`          // active, inactive, error
    LastPollTime    *time.Time        `json:"last_poll_time"`
    Metadata        map[string]string `json:"metadata"`
}
```

#### 2. Адаптеры сбора данных

Создать интерфейс и реализации адаптеров:

```go
// internal/pkg/adapters/adapter.go

package adapters

type Adapter interface {
    Connect(ctx context.Context) error
    Disconnect(ctx context.Context) error
    Read(ctx context.Context, points []DataPoint) ([]DataValue, error)
    Write(ctx context.Context, points []DataPoint) error
    Subscribe(ctx context.Context, points []DataPoint, callback func(DataValue)) error
    HealthCheck(ctx context.Context) error
}

type DataPoint struct {
    Address    string `json:"address"`     // Тег/адрес точки
    DataType   string `json:"data_type"`   // BOOL, INT32, FLOAT, STRING
    PollRate   int    `json:"poll_rate"`   // мс
    Quality    string `json:"quality"`     // Good, Bad, Uncertain
}

type DataValue struct {
    PointID   string      `json:"point_id"`
    Value     interface{} `json:"value"`
    Timestamp time.Time   `json:"timestamp"`
    Quality   string      `json:"quality"`
}
```

#### 3. Реестр адаптеров

```go
// internal/pkg/adapters/registry.go

package adapters

type Registry struct {
    adapters map[string]AdapterFactory
    mu       sync.RWMutex
}

type AdapterFactory func(config PlatformConfig) (Adapter, error)

func (r *Registry) Register(vendor, system string, factory AdapterFactory) {
    r.mu.Lock()
    defer r.mu.Unlock()
    key := fmt.Sprintf("%s:%s", vendor, system)
    r.adapters[key] = factory
}

func (r *Registry) Get(vendor, system string, config PlatformConfig) (Adapter, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    key := fmt.Sprintf("%s:%s", vendor, system)
    factory, ok := r.adapters[key]
    if !ok {
        return nil, ErrAdapterNotFound
    }
    return factory(config)
}
```

#### 4. Примеры реализаций адаптеров

##### ABB 800xA (OPC UA)

```go
// internal/pkg/adapters/abb_800xa.go

package adapters

import "github.com/gopcua/opcua"

type ABB800xAAdapter struct {
    client *opcua.Client
    config PlatformConfig
}

func NewABB800xAAdapter(config PlatformConfig) (Adapter, error) {
    // Инициализация OPC UA клиента
    endpoint := config.Endpoint
    opts := []opcua.Option{
        opcua.SecurityPolicy(opcua.SecurityPolicyBasic256Sha256),
        opcua.MessageSecurityMode(opcua.MessageSecurityModeSignAndEncrypt),
        opcua.AuthAnonymous(),
        opcua.RequestTimeout(time.Duration(config.Timeout) * time.Second),
    }
    
    client, err := opcua.NewClient(endpoint, opts...)
    if err != nil {
        return nil, err
    }
    
    return &ABB800xAAdapter{
        client: client,
        config: config,
    }, nil
}

func (a *ABB800xAAdapter) Connect(ctx context.Context) error {
    return a.client.Connect(ctx)
}

func (a *ABB800xAAdapter) Read(ctx context.Context, points []DataPoint) ([]DataValue, error) {
    nodeIDs := make([]string, len(points))
    for i, p := range points {
        nodeIDs[i] = p.Address
    }
    
    req := &opcuapb.ReadRequest{
        Nodes: nodeIDs,
    }
    
    resp, err := a.client.Read(ctx, req)
    if err != nil {
        return nil, err
    }
    
    values := make([]DataValue, len(resp.Results))
    for i, result := range resp.Results {
        values[i] = DataValue{
            PointID:   points[i].Address,
            Value:     result.Value().Value(),
            Timestamp: result.SourceTimestamp(),
            Quality:   mapQuality(result.Status),
        }
    }
    
    return values, nil
}
```

##### Siemens PCS7 (S7 Comm)

```go
// internal/pkg/adapters/siemens_pcs7.go

package adapters

import "github.com/robinson/gos7"

type SiemensPCS7Adapter struct {
    client gos7.Client
    config PlatformConfig
}

func NewSiemensPCS7Adapter(config PlatformConfig) (Adapter, error) {
    handler := gos7.NewTCPClientHandler(config.Endpoint, 102, 1)
    handler.Timeout = time.Duration(config.Timeout) * time.Second
    handler.IdleTimeout = time.Duration(config.PollInterval) * time.Second
    
    if err := handler.Connect(); err != nil {
        return nil, err
    }
    
    return &SiemensPCS7Adapter{
        client: gos7.NewClient(handler),
        config: config,
    }, nil
}

func (a *SiemensPCS7Adapter) Read(ctx context.Context, points []DataPoint) ([]DataValue, error) {
    values := make([]DataValue, len(points))
    
    for i, point := range points {
        // Парсинг адреса S7 (DB, M, I, Q)
        addr := parseS7Address(point.Address)
        
        switch point.DataType {
        case "BOOL":
            var buf []byte
            err := a.client.ReadArea(addr.Area, addr.DB, addr.Start, 1, buf)
            if err != nil {
                return nil, err
            }
            values[i] = DataValue{
                PointID:   point.Address,
                Value:     buf[0]&addr.Bit != 0,
                Timestamp: time.Now(),
                Quality:   "Good",
            }
        case "INT16":
            var buf [2]byte
            err := a.client.ReadArea(addr.Area, addr.DB, addr.Start, 2, buf[:])
            if err != nil {
                return nil, err
            }
            values[i] = DataValue{
                PointID:   point.Address,
                Value:     int16(binary.BigEndian.Uint16(buf[:])),
                Timestamp: time.Now(),
                Quality:   "Good",
            }
        }
    }
    
    return values, nil
}
```

##### Emerson DeltaV (Modbus TCP)

```go
// internal/pkg/adapters/emerson_deltav.go

package adapters

import "github.com/grid-x/modbus"

type EmersonDeltaVAdapter struct {
    client modbus.Client
    config PlatformConfig
}

func NewEmersonDeltaVAdapter(config PlatformConfig) (Adapter, error) {
    handler := modbus.NewTCPClientHandler(config.Endpoint)
    handler.Timeout = time.Duration(config.Timeout) * time.Second
    handler.SlaveId = 1
    
    if err := handler.Connect(); err != nil {
        return nil, err
    }
    
    return &EmersonDeltaVAdapter{
        client: modbus.NewClient(handler),
        config: config,
    }, nil
}

func (a *EmersonDeltaVAdapter) Read(ctx context.Context, points []DataPoint) ([]DataValue, error) {
    values := make([]DataValue, len(points))
    
    for i, point := range points {
        addr := parseModbusAddress(point.Address)
        
        switch point.DataType {
        case "BOOL":
            results, err := a.client.ReadCoils(addr.Address, 1)
            if err != nil {
                return nil, err
            }
            values[i] = DataValue{
                PointID:   point.Address,
                Value:     results[0] == 1,
                Timestamp: time.Now(),
                Quality:   "Good",
            }
        case "INT16":
            results, err := a.client.ReadHoldingRegisters(addr.Address, 1)
            if err != nil {
                return nil, err
            }
            values[i] = DataValue{
                PointID:   point.Address,
                Value:     int16(binary.BigEndian.Uint16(results)),
                Timestamp: time.Now(),
                Quality:   "Good",
            }
        }
    }
    
    return values, nil
}
```

#### 5. Регистрация адаптеров при старте

```go
// cmd/netbox/main.go

func main() {
    // ... инициализация
    
    // Регистрация адаптеров вендоров
    adapters.Register("ABB", "800xA", adapters.NewABB800xAAdapter)
    adapters.Register("ABB", "MicroSCADA", adapters.NewMicroSCADAAdapter)
    adapters.Register("Emerson", "DeltaV", adapters.NewEmersonDeltaVAdapter)
    adapters.Register("Emerson", "Ovation", adapters.NewEmersonOvationAdapter)
    adapters.Register("Siemens", "PCS7", adapters.NewSiemensPCS7Adapter)
    adapters.Register("Rockwell", "RSView", adapters.NewRockwellRSViewAdapter)
    adapters.Register("Honeywell", "Experion", adapters.NewHoneywellExperionAdapter)
    adapters.Register("Yokogawa", "CentumCS3000", adapters.NewYokogawaCentumAdapter)
    adapters.Register("Schneider", "Quantum", adapters.NewSchneiderQuantumAdapter)
    adapters.Register("GE", "XA21", adapters.NewGEXA21Adapter)
    
    // ... запуск сервера
}
```

#### 6. Конфигурация подключения

Пример YAML конфигурации для платформы:

```yaml
# config/platforms/abb_800xa.yaml
platform:
  name: "ABB 800xA Production"
  vendor: "ABB"
  system: "800xA"
  version: "6.1"
  
  connection:
    type: "OPC_UA"
    endpoint: "opc.tcp://192.168.1.100:4840"
    timeout: 30
    retry_count: 3
    poll_interval: 1000
  
  security:
    auth_method: "certificate"
    cert_file: "/etc/netbox/certs/client.crt"
    key_file: "/etc/netbox/certs/client.key"
    ca_file: "/etc/netbox/certs/ca.crt"
    policy: "Basic256Sha256"
    mode: "SignAndEncrypt"
  
  tags:
    environment: "production"
    site: "factory-1"
    area: "assembly-line"
```

#### 7. Миграции базы данных

```sql
-- migrations/0000XX_add_vendors_systems.sql

CREATE TABLE core_vendor (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    protocols TEXT[],
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE core_vendorsystem (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id UUID NOT NULL REFERENCES core_vendor(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    version VARCHAR(50),
    protocol_type VARCHAR(50) NOT NULL,
    properties TEXT[],
    config_schema JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(vendor_id, name)
);

CREATE TABLE core_platform (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    vendor_id UUID NOT NULL REFERENCES core_vendor(id),
    system_id UUID REFERENCES core_vendorsystem(id),
    connection_type VARCHAR(20) NOT NULL,
    endpoint VARCHAR(500),
    auth_method VARCHAR(30),
    certificates JSONB,
    timeout INTEGER DEFAULT 30,
    retry_count INTEGER DEFAULT 3,
    poll_interval INTEGER DEFAULT 1000,
    status VARCHAR(20) DEFAULT 'inactive',
    last_poll_time TIMESTAMPTZ,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_core_platform_vendor ON core_platform(vendor_id);
CREATE INDEX idx_core_platform_status ON core_platform(status);
CREATE INDEX idx_core_platform_metadata ON core_platform USING GIN(metadata);

-- Начальные данные для вендоров
INSERT INTO core_vendor (name, slug, description, protocols) VALUES
('ABB', 'abb', 'ABB Industrial Automation', ARRAY['OPC UA', 'Modbus', 'DNP3', 'IEC 61850']),
('Emerson', 'emerson', 'Emerson Automation Solutions', ARRAY['OPC UA', 'Modbus', 'Vnet/IP']),
('Siemens', 'siemens', 'Siemens Process Industries', ARRAY['S7 Comm', 'OPC UA', 'Profibus']),
('Rockwell', 'rockwell', 'Rockwell Automation', ARRAY['OPC UA', 'EtherNet/IP', 'DF1']),
('Honeywell', 'honeywell', 'Honeywell Process Solutions', ARRAY['OPC UA', 'EtherNet/IP']),
('Schneider', 'schneider', 'Schneider Electric', ARRAY['Modbus', 'OPC', 'Unity']),
('Yokogawa', 'yokogawa', 'Yokogawa Electric', ARRAY['Vnet/IP', 'Modbus', 'OPC']),
('GE', 'ge', 'General Electric Power', ARRAY['GEnet', 'SRTP', 'Modbus']);
```

### Этапы внедрения поддержки вендоров

#### Фаза 1: Базовая инфраструктура (Недели 1-4)
- [ ] Создание моделей данных Vendor, VendorSystem, Platform
- [ ] Реализация интерфейса Adapter
- [ ] Создание реестра адаптеров
- [ ] Базовые адаптеры OPC UA и Modbus TCP

#### Фаза 2: Ключевые вендоры (Недели 5-10)
- [ ] ABB 800xA (OPC UA)
- [ ] Emerson DeltaV (OPC UA, Modbus)
- [ ] Siemens PCS7 (S7 Comm)
- [ ] Rockwell RSView (OPC UA, EtherNet/IP)
- [ ] Honeywell Experion (OPC UA)

#### Фаза 3: Дополнительные вендоры (Недели 11-16)
- [ ] ABB MicroSCADA (DNP3, IEC 60870-5-104)
- [ ] Yokogawa Centum CS 3000 (Vnet/IP)
- [ ] Schneider Quantum (Modbus, Unity)
- [ ] GE XA/21 (GEnet, SRTP)
- [ ] Foxboro I/A Series

#### Фаза 4: Специализированные протоколы (Недели 17-22)
- [ ] ABB Symphony/Infi90 (INFI-90 Protocol)
- [ ] Emerson WDPF (Vnet/IP)
- [ ] Интеграция с Historians (Automsoft RAPID, GE Proficy)

#### Фаза 5: Мониторинг и алертинг (Недели 23-26)
- [ ] Сбор метрик состояния соединений
- [ ] Генерация событий при потере связи
- [ ] Дашборды мониторинга
- [ ] Интеграция с системами уведомлений

### Обновленная оценка трудозатрат

| Компонент | Оценка (часы) | Приоритет |
|-----------|---------------|-----------|
| **Базовая миграция Core** | | |
| Интерфейсы репозиториев | 4 | Высокий |
| PostgreSQL репозитории | 16 | Высокий |
| SQL миграции | 8 | Высокий |
| HTTP обработчики | 12 | Высокий |
| EtcdQueue интеграция | 12 | Высокий |
| Worker Pool реализация | 8 | Высокий |
| Change Logging сервис | 8 | Высокий |
| Тесты (unit + integration) | 16 | Высокий |
| Документация | 4 | Средний |
| **Поддержка вендоров** | | |
| Модели данных Vendor/Platform | 8 | Высокий |
| Интерфейс Adapter + Registry | 12 | Высокий |
| Адаптер OPC UA (базовый) | 16 | Высокий |
| Адаптер Modbus TCP (базовый) | 12 | Высокий |
| Адаптер ABB 800xA | 12 | Высокий |
| Адаптер Emerson DeltaV | 12 | Высокий |
| Адаптер Siemens PCS7 | 16 | Высокий |
| Адаптер Rockwell RSView | 12 | Высокий |
| Адаптер Honeywell Experion | 12 | Высокий |
| Адаптер Yokogawa Centum | 12 | Средний |
| Адаптер Schneider Quantum | 12 | Средний |
| Адаптер GE XA/21 | 12 | Средний |
| Адаптер ABB MicroSCADA | 16 | Средний |
| Адаптер INFI-90 Protocol | 20 | Низкий |
| Интеграция Historians | 16 | Низкий |
| Мониторинг и алертинг | 16 | Средний |
| Документация вендоров | 8 | Средний |
| **Поддержка протоколов** | | |
| Поллер SNMP (v1, v2c, v3) | 20 | Высокий |
| Поллер ICMP (Ping) | 8 | Высокий |
| Поллер Modbus TCP/RTU | 16 | Высокий |
| Поллер HTTP/HTTPS | 12 | Высокий |
| Поллер DNS | 12 | Высокий |
| Поллер SSH | 16 | Высокий |
| Поллер LDAP/LDAPS | 12 | Высокий |
| Поллер Radius | 12 | Высокий |
| Поллер BGP (BGP-4) | 20 | Высокий |
| Поллер Syslog | 12 | Высокий |
| Поллер NetFlow | 16 | Высокий |
| Поллер SMTP/IMAP/POP3 | 12 | Средний |
| Поллер FTP/SFTP | 10 | Средний |
| Поллер DHCP | 10 | Средний |
| Поллер JMX | 12 | Средний |
| Поллер WMI | 16 | Средний |
| Поллер BACnet | 16 | Средний |
| Поллер OPC UA | 20 | Высокий |
| Поллер NMEA 0183 | 10 | Низкий |
| Поллер CORBA/WBEM | 20 | Низкий |
| Интеграция PollingService с etcd | 16 | Высокий |
| Документация протоколов | 8 | Средний |
| **Итого базовая миграция** | **88 часов** | |
| **Итого поддержка вендоров** | **232 часа** | |
| **Итого поддержка протоколов** | **308 часов** | |
| **ОБЩАЯ ОЦЕНКА** | **628 часов** | |

### Рекомендуемые Go библиотеки для реализации

| Библиотека | Назначение | Ссылка |
|------------|------------|--------|
| `gopcua/opcua` | OPC UA клиент | github.com/gopcua/opcua |
| `grid-x/modbus` | Modbus TCP/RTU | github.com/grid-x/modbus |
| `robinson/gos7` | Siemens S7 Comm | github.com/robinson/gos7 |
| `eclipse/paho.mqtt.golang` | MQTT для телеметрии | github.com/eclipse/paho.mqtt.golang |
| `apache/thrift` | Thrift RPC (некоторые вендоры) | github.com/apache/thrift |
| `grpc/grpc-go` | gRPC для внутренних API | google.golang.org/grpc |
| `jackc/pgx` | PostgreSQL драйвер | github.com/jackc/pgx/v5 |
| `etcd-io/etcd/client/v3` | etcd клиент | go.etcd.io/etcd/client/v3 |

---

## Оценка трудозатрат

| Компонент | Оценка (часы) | Приоритет |
|-----------|---------------|-----------|
| Интерфейсы репозиториев | 4 | Высокий |
| PostgreSQL репозитории | 16 | Высокий |
| SQL миграции | 8 | Высокий |
| HTTP обработчики | 12 | Высокий |
| EtcdQueue интеграция | 12 | Высокий |
| Worker Pool реализация | 8 | Высокий |
| Change Logging сервис | 8 | Высокий |
| Тесты (unit + integration) | 16 | Высокий |
| Документация | 4 | Средний |
| **Итого базовая миграция** | **88 часов** | |
| **Поддержка вендоров АСУ ТП** | **232 часа** | |
| **Поддержка сетевых протоколов** | **308 часов** | |
| **ОБЩАЯ ОЦЕНКА ПРОЕКТА** | **628 часов (~16 недель)** | |

---

## Следующие шаги

1. **Неделя 1**: Реализация интерфейсов и PostgreSQL репозиториев
2. **Неделя 2**: HTTP обработчики и маршруты API
3. **Неделя 3**: Интеграция с EtcdQueue и Worker Pool
4. **Неделя 4**: Change Logging и тестирование
5. **Неделя 5**: Документация и финальная отладка
6. **Недели 6-13**: Реализация поддержки вендоров (Фазы 1-3)
7. **Недели 14-22**: Специализированные протоколы и Historians (Фазы 4-5)
8. **Недели 23-30**: Реализация поллеров сетевых протоколов (SNMP, ICMP, DNS, SSH, LDAP, BGP, Syslog, NetFlow)
9. **Недели 31-36**: Поллеры прикладных протоколов (HTTP, FTP, SMTP, JMX, WMI, BACnet)
10. **Недели 37-40**: Интеграция, мониторинг и финальное тестирование

### Рекомендуемые Go библиотеки для реализации

| Библиотека | Назначение | Ссылка |
|------------|------------|--------|
| `gopcua/opcua` | OPC UA клиент | github.com/gopcua/opcua |
| `grid-x/modbus` | Modbus TCP/RTU | github.com/grid-x/modbus |
| `robinson/gos7` | Siemens S7 Comm | github.com/robinson/gos7 |
| `gosnmp/gosnmp` | SNMP v1/v2c/v3 | github.com/gosnmp/gosnmp |
| `go-ping/ping` | ICMP Ping | github.com/go-ping/ping |
| `miekg/dns` | DNS клиент/сервер | github.com/miekg/dns |
| `gliderlabs/ssh` | SSH сервер/клиент | github.com/gliderlabs/ssh |
| `go-ldap/ldap` | LDAP клиент | github.com/go-ldap/ldap |
| `eclipse/paho.mqtt.golang` | MQTT для телеметрии | github.com/eclipse/paho.mqtt.golang |
| `apache/thrift` | Thrift RPC (некоторые вендоры) | github.com/apache/thrift |
| `grpc/grpc-go` | gRPC для внутренних API | google.golang.org/grpc |
| `jackc/pgx` | PostgreSQL драйвер | github.com/jackc/pgx/v5 |
| `etcd-io/etcd/client/v3` | etcd клиент | go.etcd.io/etcd/client/v3 |
| `nwaples/radius` | RADIUS клиент | github.com/nwaples/radius |
| `sirikothe/gotextfsm` | Парсинг вывода CLI (для BGP/SSH) | github.com/sirikothe/gotextfsm |
| `google/gopacket` | Анализ пакетов (NetFlow, PCAP) | github.com/google/gopacket |
