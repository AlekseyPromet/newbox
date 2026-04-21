-- Migration: 001_create_core_tables.sql
-- Описание: Создание таблиц модуля Core для NetBox на Go
-- Дата: 2024-01-15

-- Таблица: core_configrevision
-- Хранит ревизии конфигурации NetBox
CREATE TABLE IF NOT EXISTS core_configrevision (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    active BOOLEAN NOT NULL DEFAULT FALSE,
    comment TEXT,
    data JSONB NOT NULL DEFAULT '{}'::jsonb
);

CREATE INDEX idx_core_configrevision_active ON core_configrevision(active);
CREATE INDEX idx_core_configrevision_created ON core_configrevision(created DESC);

COMMENT ON TABLE core_configrevision IS 'Ревизии конфигурации NetBox';
COMMENT ON COLUMN core_configrevision.id IS 'Первичный ключ';
COMMENT ON COLUMN core_configrevision.created IS 'Дата создания ревизии';
COMMENT ON COLUMN core_configrevision.active IS 'Флаг активной ревизии';
COMMENT ON COLUMN core_configrevision.comment IS 'Комментарий к ревизии';
COMMENT ON COLUMN core_configrevision.data IS 'Данные конфигурации в формате JSON';

-- Таблица: django_content_type
-- Аналог ContentType из Django для хранения типов объектов
CREATE TABLE IF NOT EXISTS django_content_type (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    app_label VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    public BOOLEAN NOT NULL DEFAULT TRUE,
    features JSONB DEFAULT '[]'::jsonb,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(app_label, model)
);

CREATE INDEX idx_django_content_type_app_label ON django_content_type(app_label);
CREATE INDEX idx_django_content_type_model ON django_content_type(model);
CREATE INDEX idx_django_content_type_public ON django_content_type(public);

COMMENT ON TABLE django_content_type IS 'Типы объектов (аналог Django ContentType)';
COMMENT ON COLUMN django_content_type.app_label IS 'Метка приложения';
COMMENT ON COLUMN django_content_type.model IS 'Имя модели';
COMMENT ON COLUMN django_content_type.public IS 'Публичный тип объекта';
COMMENT ON COLUMN django_content_type.features IS 'Список поддерживаемых функций';

-- Таблица: core_objectchange
-- Журнал изменений объектов (change logging)
CREATE TABLE IF NOT EXISTS core_objectchange (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    time TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    user_id UUID,
    request_id VARCHAR(36),
    action VARCHAR(20) NOT NULL CHECK (action IN ('create', 'update', 'delete')),
    changed_object_type VARCHAR(200) NOT NULL,
    changed_object_id VARCHAR(36) NOT NULL,
    object_repr VARCHAR(500) NOT NULL,
    object_data JSONB,
    related_object_type VARCHAR(200),
    related_object_id VARCHAR(36),
    related_object_repr VARCHAR(500)
);

CREATE INDEX idx_core_objectchange_time ON core_objectchange(time DESC);
CREATE INDEX idx_core_objectchange_user_id ON core_objectchange(user_id);
CREATE INDEX idx_core_objectchange_action ON core_objectchange(action);
CREATE INDEX idx_core_objectchange_changed_object_type ON core_objectchange(changed_object_type);
CREATE INDEX idx_core_objectchange_changed_object_id ON core_objectchange(changed_object_id);
CREATE INDEX idx_core_objectchange_related_object_type ON core_objectchange(related_object_type);
CREATE INDEX idx_core_objectchange_request_id ON core_objectchange(request_id);

-- Композитный индекс для частых запросов фильтрации
CREATE INDEX idx_core_objectchange_filter 
ON core_objectchange(user_id, action, changed_object_type, time DESC);

COMMENT ON TABLE core_objectchange IS 'Журнал изменений объектов';
COMMENT ON COLUMN core_objectchange.time IS 'Время изменения';
COMMENT ON COLUMN core_objectchange.user_id IS 'ID пользователя, внёсшего изменение';
COMMENT ON COLUMN core_objectchange.request_id IS 'ID HTTP запроса';
COMMENT ON COLUMN core_objectchange.action IS 'Тип действия (create/update/delete)';
COMMENT ON COLUMN core_objectchange.changed_object_type IS 'Тип изменённого объекта';
COMMENT ON COLUMN core_objectchange.changed_object_id IS 'ID изменённого объекта';
COMMENT ON COLUMN core_objectchange.object_repr IS 'Строковое представление объекта';
COMMENT ON COLUMN core_objectchange.object_data IS 'Данные объекта до/после изменения';
COMMENT ON COLUMN core_objectchange.related_object_type IS 'Тип связанного объекта';
COMMENT ON COLUMN core_objectchange.related_object_id IS 'ID связанного объекта';
COMMENT ON COLUMN core_objectchange.related_object_repr IS 'Представление связанного объекта';

-- Таблица: core_datasource
-- Источники внешних данных
CREATE TABLE IF NOT EXISTS core_datasource (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL UNIQUE,
    type VARCHAR(100) NOT NULL,
    source_url VARCHAR(500) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'new' 
        CHECK (status IN ('new', 'queued', 'syncing', 'completed', 'failed')),
    enabled BOOLEAN NOT NULL DEFAULT TRUE,
    sync_interval INTEGER NOT NULL DEFAULT 0,
    ignore_rules JSONB DEFAULT '[]'::jsonb,
    parameters JSONB DEFAULT '{}'::jsonb,
    last_synced TIMESTAMP WITH TIME ZONE,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_core_datasource_status ON core_datasource(status);
CREATE INDEX idx_core_datasource_enabled ON core_datasource(enabled);
CREATE INDEX idx_core_datasource_type ON core_datasource(type);
CREATE INDEX idx_core_datasource_name ON core_datasource(name);

COMMENT ON TABLE core_datasource IS 'Источники внешних данных';
COMMENT ON COLUMN core_datasource.name IS 'Уникальное имя источника';
COMMENT ON COLUMN core_datasource.type IS 'Тип источника (git, http, local, etc.)';
COMMENT ON COLUMN core_datasource.source_url IS 'URL или путь к источнику';
COMMENT ON COLUMN core_datasource.status IS 'Статус источника';
COMMENT ON COLUMN core_datasource.enabled IS 'Флаг включения';
COMMENT ON COLUMN core_datasource.sync_interval IS 'Интервал синхронизации в минутах';
COMMENT ON COLUMN core_datasource.ignore_rules IS 'Правила игнорирования файлов';
COMMENT ON COLUMN core_datasource.parameters IS 'Параметры подключения';
COMMENT ON COLUMN core_datasource.last_synced IS 'Время последней синхронизации';

-- Таблица: core_datafile
-- Файлы данных из источников
CREATE TABLE IF NOT EXISTS core_datafile (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id UUID NOT NULL REFERENCES core_datasource(id) ON DELETE CASCADE,
    path VARCHAR(500) NOT NULL,
    size BIGINT NOT NULL DEFAULT 0,
    hash VARCHAR(64) NOT NULL,
    data JSONB,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    UNIQUE(source_id, path)
);

CREATE INDEX idx_core_datafile_source_id ON core_datafile(source_id);
CREATE INDEX idx_core_datafile_path ON core_datafile(path);
CREATE INDEX idx_core_datafile_hash ON core_datafile(hash);

COMMENT ON TABLE core_datafile IS 'Файлы данных из источников';
COMMENT ON COLUMN core_datafile.source_id IS 'Ссылка на источник данных';
COMMENT ON COLUMN core_datafile.path IS 'Путь к файлу относительно источника';
COMMENT ON COLUMN core_datafile.size IS 'Размер файла в байтах';
COMMENT ON COLUMN core_datafile.hash IS 'Хэш содержимого файла (SHA256)';
COMMENT ON COLUMN core_datafile.data IS 'Парсенные данные файла (если применимо)';

-- Таблица: core_job
-- Фоновые задачи
CREATE TABLE IF NOT EXISTS core_job (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    object_type VARCHAR(200),
    object_id UUID,
    name VARCHAR(200) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'scheduled', 'running', 'completed', 'errored', 'failed')),
    interval INTEGER DEFAULT 0,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    queue_name VARCHAR(100),
    job_id VARCHAR(100),
    data JSONB DEFAULT '{}'::jsonb,
    error TEXT,
    created TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_core_job_status ON core_job(status);
CREATE INDEX idx_core_job_scheduled_at ON core_job(scheduled_at);
CREATE INDEX idx_core_job_object_type ON core_job(object_type);
CREATE INDEX idx_core_job_object_id ON core_job(object_id);
CREATE INDEX idx_core_job_queue_name ON core_job(queue_name);
CREATE INDEX idx_core_job_created ON core_job(created DESC);

-- Композитный индекс для фильтрации задач
CREATE INDEX idx_core_job_filter 
ON core_job(status, object_type, object_id, created DESC);

COMMENT ON TABLE core_job IS 'Фоновые задачи';
COMMENT ON COLUMN core_job.object_type IS 'Тип связанного объекта';
COMMENT ON COLUMN core_job.object_id IS 'ID связанного объекта';
COMMENT ON COLUMN core_job.name IS 'Имя задачи';
COMMENT ON COLUMN core_job.status IS 'Статус выполнения';
COMMENT ON COLUMN core_job.interval IS 'Интервал повторения в минутах';
COMMENT ON COLUMN core_job.scheduled_at IS 'Время запланированного запуска';
COMMENT ON COLUMN core_job.started_at IS 'Время начала выполнения';
COMMENT ON COLUMN core_job.completed_at IS 'Время завершения';
COMMENT ON COLUMN core_job.queue_name IS 'Имя очереди etcd';
COMMENT ON COLUMN core_job.job_id IS 'Внешний ID задачи в системе очередей';
COMMENT ON COLUMN core_job.data IS 'Параметры задачи';
COMMENT ON COLUMN core_job.error IS 'Сообщение об ошибке';

-- Триггер для обновления updated при изменении записи
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_django_content_type_updated 
    BEFORE UPDATE ON django_content_type 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_core_datasource_updated 
    BEFORE UPDATE ON core_datasource 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_core_datafile_updated 
    BEFORE UPDATE ON core_datafile 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_core_job_updated 
    BEFORE UPDATE ON core_job 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Начальные данные: базовые типы объектов
INSERT INTO django_content_type (app_label, model, public, features) VALUES
    ('core', 'configrevision', TRUE, '["export", "import"]'::jsonb),
    ('core', 'objectchange', TRUE, '["filtering", "export"]'::jsonb),
    ('core', 'datasource', TRUE, '["sync", "export", "import"]'::jsonb),
    ('core', 'datafile', TRUE, '["export"]'::jsonb),
    ('core', 'job', TRUE, '["filtering", "export"]'::jsonb)
ON CONFLICT (app_label, model) DO NOTHING;
