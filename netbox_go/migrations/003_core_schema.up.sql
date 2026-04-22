-- Миграция 003: Создание таблиц CORE (ConfigRevisions, ObjectTypes, ObjectChanges, DataSources, DataFiles, Jobs)
-- Версия: 1.0.0
-- Дата: 2024-01-15

-- ============================================
-- Таблицы для ConfigRevisions
-- ============================================

CREATE TABLE core_config_revisions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    active BOOLEAN DEFAULT FALSE,
    comment TEXT,
    data JSONB,
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_core_config_revisions_active ON core_config_revisions(active);
CREATE INDEX idx_core_config_revisions_created ON core_config_revisions(created);

-- ============================================
-- Таблицы для ObjectTypes (аналог Django ContentType)
-- ============================================

CREATE TABLE core_object_types (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_label VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    public BOOLEAN DEFAULT TRUE,
    features TEXT[], -- массив фич: "export", "import", "graphql", etc.
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(app_label, model)
);

CREATE INDEX idx_core_object_types_app_label ON core_object_types(app_label);
CREATE INDEX idx_core_object_types_model ON core_object_types(model);
CREATE INDEX idx_core_object_types_public ON core_object_types(public);

-- ============================================
-- Таблицы для ObjectChanges (Change Logging)
-- ============================================

CREATE TYPE core_objectchange_action AS ENUM ('create', 'update', 'delete');

CREATE TABLE core_object_changes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    time TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    user_id UUID REFERENCES users_users(id) ON DELETE SET NULL,
    request_id VARCHAR(128),
    action core_objectchange_action NOT NULL,
    changed_object_type VARCHAR(200) NOT NULL, -- format: "app_label.model"
    changed_object_id VARCHAR(128) NOT NULL,
    object_repr VARCHAR(500) NOT NULL,
    object_data JSONB,
    related_object_type VARCHAR(200),
    related_object_id VARCHAR(128),
    related_object_repr VARCHAR(500),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_core_object_changes_time ON core_object_changes(time);
CREATE INDEX idx_core_object_changes_user_id ON core_object_changes(user_id);
CREATE INDEX idx_core_object_changes_action ON core_object_changes(action);
CREATE INDEX idx_core_object_changes_changed_object_type ON core_object_changes(changed_object_type);
CREATE INDEX idx_core_object_changes_changed_object_id ON core_object_changes(changed_object_id);
CREATE INDEX idx_core_object_changes_request_id ON core_object_changes(request_id);
CREATE INDEX idx_core_object_changes_related_object_type ON core_object_changes(related_object_type);

-- ============================================
-- Таблицы для DataSources
-- ============================================

CREATE TYPE core_datasource_status AS ENUM ('new', 'queued', 'syncing', 'completed', 'failed');

CREATE TABLE core_data_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    type VARCHAR(100) NOT NULL, -- "local", "git", "s3", etc.
    source_url VARCHAR(500) NOT NULL,
    status core_datasource_status DEFAULT 'new',
    enabled BOOLEAN DEFAULT TRUE,
    sync_interval INTEGER DEFAULT 0, -- минуты между синхронизациями (0 = disabled)
    ignore_rules TEXT[], -- правила игнорирования файлов
    parameters JSONB, -- дополнительные параметры подключения
    last_synced TIMESTAMP WITH TIME ZONE,
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name)
);

CREATE INDEX idx_core_data_sources_name ON core_data_sources(name);
CREATE INDEX idx_core_data_sources_type ON core_data_sources(type);
CREATE INDEX idx_core_data_sources_status ON core_data_sources(status);
CREATE INDEX idx_core_data_sources_enabled ON core_data_sources(enabled);

-- ============================================
-- Таблицы для DataFiles
-- ============================================

CREATE TABLE core_data_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES core_data_sources(id) ON DELETE CASCADE,
    path VARCHAR(500) NOT NULL, -- относительный путь файла
    size BIGINT NOT NULL DEFAULT 0, -- размер в байтах
    hash VARCHAR(64) NOT NULL, -- SHA256 хэш содержимого
    data JSONB, -- кэшированное содержимое (для небольших файлов)
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(source_id, path)
);

CREATE INDEX idx_core_data_files_source_id ON core_data_files(source_id);
CREATE INDEX idx_core_data_files_path ON core_data_files(path);
CREATE INDEX idx_core_data_files_hash ON core_data_files(hash);

-- ============================================
-- Таблицы для Jobs (фоновые задачи)
-- ============================================

CREATE TYPE core_job_status AS ENUM ('pending', 'scheduled', 'running', 'completed', 'errored', 'failed');

CREATE TABLE core_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    object_type VARCHAR(200), -- формат: "app_label.model"
    object_id UUID,
    name VARCHAR(200) NOT NULL,
    status core_job_status DEFAULT 'pending',
    interval INTEGER DEFAULT 0, -- интервал повторения в минутах (0 = однократная)
    scheduled_at TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    queue_name VARCHAR(100), -- имя очереди Asynq/RQ
    job_id VARCHAR(128), -- внешний ID задачи (например, Asynq task ID)
    data JSONB, -- параметры задачи
    error TEXT, -- сообщение об ошибке
    created TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_core_jobs_object_type ON core_jobs(object_type);
CREATE INDEX idx_core_jobs_object_id ON core_jobs(object_id);
CREATE INDEX idx_core_jobs_status ON core_jobs(status);
CREATE INDEX idx_core_jobs_scheduled_at ON core_jobs(scheduled_at);
CREATE INDEX idx_core_jobs_queue_name ON core_jobs(queue_name);
CREATE INDEX idx_core_jobs_job_id ON core_jobs(job_id);

-- ============================================
-- Начальные данные
-- ============================================

-- Добавляем базовые типы объектов для основных моделей NetBox
INSERT INTO core_object_types (app_label, model, public, features, created, updated) VALUES
    ('dcim', 'site', TRUE, ARRAY['export', 'import', 'graphql', 'filter', 'clone']),
    ('dcim', 'rack', TRUE, ARRAY['export', 'import', 'graphql', 'filter', 'clone']),
    ('dcim', 'device', TRUE, ARRAY['export', 'import', 'graphql', 'filter', 'clone']),
    ('dcim', 'cable', TRUE, ARRAY['export', 'import', 'graphql', 'filter']),
    ('circuits', 'circuit', TRUE, ARRAY['export', 'import', 'graphql', 'filter', 'clone']),
    ('ipam', 'ipaddress', TRUE, ARRAY['export', 'import', 'graphql', 'filter', 'clone']),
    ('ipam', 'prefix', TRUE, ARRAY['export', 'import', 'graphql', 'filter', 'clone']),
    ('tenancy', 'tenant', TRUE, ARRAY['export', 'import', 'graphql', 'filter', 'clone']),
    ('users', 'user', TRUE, ARRAY['export', 'import', 'graphql', 'filter']),
    ('extras', 'customfield', TRUE, ARRAY['export', 'import', 'graphql', 'filter']),
    ('core', 'datasource', TRUE, ARRAY['export', 'import', 'graphql', 'filter']),
    ('core', 'datafile', TRUE, ARRAY['export', 'import', 'graphql', 'filter']),
    ('core', 'job', TRUE, ARRAY['export', 'import', 'graphql', 'filter']),
    ('core', 'objectchange', TRUE, ARRAY['export', 'import', 'graphql', 'filter']),
    ('core', 'configrevision', TRUE, ARRAY['export', 'import', 'graphql']);

COMMENT ON TABLE core_config_revisions IS 'Хранит ревизии конфигурации NetBox';
COMMENT ON TABLE core_object_types IS 'Аналог Django ContentType - типы объектов NetBox';
COMMENT ON TABLE core_object_changes IS 'Журнал изменений объектов (change logging)';
COMMENT ON TABLE core_data_sources IS 'Внешние источники данных (git, s3, local)';
COMMENT ON TABLE core_data_files IS 'Файлы, полученные из источников данных';
COMMENT ON TABLE core_jobs IS 'Фоновые задачи и периодические джобы';
