-- Название: core
-- Описание: SQL запросы для работы с таблицами модуля Core NetBox
-- Версия: 1.0.0

-- name: GetConfigRevisionByID :one
SELECT 
    id, created, active, comment, data
FROM core_configrevision
WHERE id = $1;

-- name: GetActiveConfigRevision :one
SELECT 
    id, created, active, comment, data
FROM core_configrevision
WHERE active = TRUE
ORDER BY created DESC
LIMIT 1;

-- name: ListConfigRevisions :many
SELECT 
    id, created, active, comment, data
FROM core_configrevision
ORDER BY created DESC
LIMIT $1 OFFSET $2;

-- name: CountConfigRevisions :one
SELECT COUNT(*)::int as count
FROM core_configrevision;

-- name: CreateConfigRevision :one
INSERT INTO core_configrevision (created, active, comment, data)
VALUES ($1, $2, $3, $4)
RETURNING id, created, active, comment, data;

-- name: UpdateConfigRevision :execrows
UPDATE core_configrevision
SET active = $2, comment = $3, data = $4
WHERE id = $1;

-- name: DeleteConfigRevision :execrows
DELETE FROM core_configrevision
WHERE id = $1;

-- name: SetActiveConfigRevision :execrows
UPDATE core_configrevision
SET active = FALSE
WHERE active = TRUE;

UPDATE core_configrevision
SET active = TRUE
WHERE id = $1;

-- name: GetObjectTypeByID :one
SELECT 
    id, app_label, model, public, features, created, updated
FROM django_content_type
WHERE id = $1;

-- name: GetObjectTypeByAppAndModel :one
SELECT 
    id, app_label, model, public, features, created, updated
FROM django_content_type
WHERE app_label = $1 AND model = $2;

-- name: ListObjectTypes :many
SELECT 
    id, app_label, model, public, features, created, updated
FROM django_content_type
WHERE ($1 = '' OR app_label = $1)
  AND ($2 = '' OR model = $2)
  AND ($3::boolean IS NULL OR public = $3)
ORDER BY app_label, model
LIMIT $4 OFFSET $5;

-- name: CountObjectTypes :one
SELECT COUNT(*)::int as count
FROM django_content_type
WHERE ($1 = '' OR app_label = $1)
  AND ($2 = '' OR model = $2)
  AND ($3::boolean IS NULL OR public = $3);

-- name: CreateObjectType :one
INSERT INTO django_content_type (app_label, model, public, features, created, updated)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, app_label, model, public, features, created, updated;

-- name: UpdateObjectType :execrows
UPDATE django_content_type
SET app_label = $2, model = $3, public = $4, features = $5, updated = $6
WHERE id = $1;

-- name: DeleteObjectType :execrows
DELETE FROM django_content_type
WHERE id = $1;

-- name: GetObjectChangeByID :one
SELECT 
    id, time, user_id, request_id, action, changed_object_type, 
    changed_object_id, object_repr, object_data, related_object_type,
    related_object_id, related_object_repr
FROM core_objectchange
WHERE id = $1;

-- name: ListObjectChanges :many
SELECT 
    id, time, user_id, request_id, action, changed_object_type, 
    changed_object_id, object_repr, object_data, related_object_type,
    related_object_id, related_object_repr
FROM core_objectchange
WHERE ($1::uuid IS NULL OR user_id = $1)
  AND ($2 = '' OR action = $2)
  AND ($3 = '' OR changed_object_type = $3)
  AND ($4 = '' OR changed_object_id = $4)
  AND ($5::timestamp WITH TIME ZONE IS NULL OR time >= $5)
  AND ($6::timestamp WITH TIME ZONE IS NULL OR time <= $6)
ORDER BY time DESC
LIMIT $7 OFFSET $8;

-- name: CountObjectChanges :one
SELECT COUNT(*)::int as count
FROM core_objectchange
WHERE ($1::uuid IS NULL OR user_id = $1)
  AND ($2 = '' OR action = $2)
  AND ($3 = '' OR changed_object_type = $3)
  AND ($4 = '' OR changed_object_id = $4)
  AND ($5::timestamp WITH TIME ZONE IS NULL OR time >= $5)
  AND ($6::timestamp WITH TIME ZONE IS NULL OR time <= $6);

-- name: CreateObjectChange :one
INSERT INTO core_objectchange (
    time, user_id, request_id, action, changed_object_type, 
    changed_object_id, object_repr, object_data, related_object_type,
    related_object_id, related_object_repr
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING 
    id, time, user_id, request_id, action, changed_object_type, 
    changed_object_id, object_repr, object_data, related_object_type,
    related_object_id, related_object_repr;

-- name: BulkCreateObjectChanges :exec
INSERT INTO core_objectchange (
    time, user_id, request_id, action, changed_object_type,
    changed_object_id, object_repr, object_data, related_object_type,
    related_object_id, related_object_repr
)
VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
);

-- name: DeleteOldObjectChanges :execrows
DELETE FROM core_objectchange
WHERE time < $1;

-- name: GetDataSourceByID :one
SELECT 
    id, name, type, source_url, status, enabled, sync_interval,
    ignore_rules, parameters, last_synced, created, updated
FROM core_datasource
WHERE id = $1;

-- name: GetDataSourceByName :one
SELECT 
    id, name, type, source_url, status, enabled, sync_interval,
    ignore_rules, parameters, last_synced, created, updated
FROM core_datasource
WHERE name = $1;

-- name: ListDataSources :many
SELECT 
    id, name, type, source_url, status, enabled, sync_interval,
    ignore_rules, parameters, last_synced, created, updated
FROM core_datasource
WHERE ($1 = '' OR status = $1)
  AND ($2::boolean IS NULL OR enabled = $2)
  AND ($3 = '' OR type = $3)
ORDER BY name
LIMIT $4 OFFSET $5;

-- name: CountDataSources :one
SELECT COUNT(*)::int as count
FROM core_datasource
WHERE ($1 = '' OR status = $1)
  AND ($2::boolean IS NULL OR enabled = $2)
  AND ($3 = '' OR type = $3);

-- name: CreateDataSource :one
INSERT INTO core_datasource (
    name, type, source_url, status, enabled, sync_interval,
    ignore_rules, parameters, last_synced, created, updated
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING 
    id, name, type, source_url, status, enabled, sync_interval,
    ignore_rules, parameters, last_synced, created, updated;

-- name: UpdateDataSource :execrows
UPDATE core_datasource
SET name = $2, type = $3, source_url = $4, status = $5, enabled = $6,
    sync_interval = $7, ignore_rules = $8, parameters = $9, 
    last_synced = $10, updated = $11
WHERE id = $1;

-- name: DeleteDataSource :execrows
DELETE FROM core_datasource
WHERE id = $1;

-- name: UpdateDataSourceStatus :execrows
UPDATE core_datasource
SET status = $2, last_synced = $3, updated = NOW()
WHERE id = $1;

-- name: GetQueuedDataSources :many
SELECT 
    id, name, type, source_url, status, enabled, sync_interval,
    ignore_rules, parameters, last_synced, created, updated
FROM core_datasource
WHERE status = 'queued'
ORDER BY updated ASC
LIMIT $1;

-- name: GetDataFileByID :one
SELECT 
    id, source_id, path, size, hash, data, created, updated
FROM core_datafile
WHERE id = $1;

-- name: GetDataFileBySourceAndPath :one
SELECT 
    id, source_id, path, size, hash, data, created, updated
FROM core_datafile
WHERE source_id = $1 AND path = $2;

-- name: ListDataFilesBySource :many
SELECT 
    id, source_id, path, size, hash, data, created, updated
FROM core_datafile
WHERE source_id = $1
ORDER BY path
LIMIT $2 OFFSET $3;

-- name: CountDataFilesBySource :one
SELECT COUNT(*)::int as count
FROM core_datafile
WHERE source_id = $1;

-- name: CreateDataFile :one
INSERT INTO core_datafile (source_id, path, size, hash, data, created, updated)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, source_id, path, size, hash, data, created, updated;

-- name: UpdateDataFile :execrows
UPDATE core_datafile
SET path = $2, size = $3, hash = $4, data = $5, updated = $6
WHERE id = $1;

-- name: DeleteDataFile :execrows
DELETE FROM core_datafile
WHERE id = $1;

-- name: BulkDeleteDataFilesBySource :execrows
DELETE FROM core_datafile
WHERE source_id = $1;

-- name: GetJobByID :one
SELECT 
    id, object_type, object_id, name, status, interval, scheduled_at,
    started_at, completed_at, queue_name, job_id, data, error, created, updated
FROM core_job
WHERE id = $1;

-- name: ListJobs :many
SELECT 
    id, object_type, object_id, name, status, interval, scheduled_at,
    started_at, completed_at, queue_name, job_id, data, error, created, updated
FROM core_job
WHERE ($1 = '' OR status = $1)
  AND ($2 = '' OR object_type = $2)
  AND ($3::uuid IS NULL OR object_id = $3)
  AND ($4 = '' OR queue_name = $4)
ORDER BY created DESC
LIMIT $5 OFFSET $6;

-- name: CountJobs :one
SELECT COUNT(*)::int as count
FROM core_job
WHERE ($1 = '' OR status = $1)
  AND ($2 = '' OR object_type = $2)
  AND ($3::uuid IS NULL OR object_id = $3)
  AND ($4 = '' OR queue_name = $4);

-- name: CreateJob :one
INSERT INTO core_job (
    object_type, object_id, name, status, interval, scheduled_at,
    started_at, completed_at, queue_name, job_id, data, error, created, updated
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING 
    id, object_type, object_id, name, status, interval, scheduled_at,
    started_at, completed_at, queue_name, job_id, data, error, created, updated;

-- name: UpdateJob :execrows
UPDATE core_job
SET object_type = $2, object_id = $3, name = $4, status = $5, interval = $6,
    scheduled_at = $7, started_at = $8, completed_at = $9, queue_name = $10,
    job_id = $11, data = $12, error = $13, updated = $14
WHERE id = $1;

-- name: DeleteJob :execrows
DELETE FROM core_job
WHERE id = $1;

-- name: UpdateJobStatus :execrows
UPDATE core_job
SET status = $2, error = $3, completed_at = $4, updated = NOW()
WHERE id = $1;

-- name: GetScheduledJobs :many
SELECT 
    id, object_type, object_id, name, status, interval, scheduled_at,
    started_at, completed_at, queue_name, job_id, data, error, created, updated
FROM core_job
WHERE status = 'scheduled' AND scheduled_at <= $1
ORDER BY scheduled_at ASC
LIMIT $2;

-- name: CleanupOldJobs :execrows
DELETE FROM core_job
WHERE status IN ('completed', 'failed', 'errored') 
  AND completed_at < $1;

-- ============================================================================
-- JSON Helper Functions
-- ============================================================================

-- name: GetJSONField :one
SELECT data->>$2 as value
FROM core_configrevision
WHERE id = $1;

-- name: UpdateJSONField :execrows
UPDATE core_configrevision
SET data = jsonb_set(data::jsonb, $2, to_jsonb($3))
WHERE id = $1;

-- name: GetJSONPath :many
SELECT jsonb_path_query(data, $2::jsonpath)
FROM core_configrevision
WHERE id = $1;

-- name: MergeJSONData :execrows
UPDATE core_configrevision
SET data = COALESCE(data, '{}'::jsonb)::jsonb || $2::jsonb
WHERE id = $1;

-- ============================================================================
-- Additional ObjectChange Queries
-- ============================================================================

-- name: GetObjectChangesByObjectType :many
SELECT 
    id, time, user_id, request_id, action, changed_object_type, 
    changed_object_id, object_repr, object_data, related_object_type,
    related_object_id, related_object_repr
FROM core_objectchange
WHERE changed_object_type = $1
ORDER BY time DESC
LIMIT $2 OFFSET $3;

-- name: GetRecentObjectChanges :many
SELECT 
    id, time, user_id, request_id, action, changed_object_type, 
    changed_object_id, object_repr, object_data, related_object_type,
    related_object_id, related_object_repr
FROM core_objectchange
WHERE changed_object_id = $1 AND changed_object_type = $2
ORDER BY time DESC
LIMIT $3;

-- name: GetObjectChangeCountByAction :one
SELECT action, COUNT(*)::int as count
FROM core_objectchange
WHERE ($1::timestamp WITH TIME ZONE IS NULL OR time >= $1)
  AND ($2::timestamp WITH TIME ZONE IS NULL OR time <= $2)
GROUP BY action;

-- ============================================================================
-- Additional DataSource Queries
-- ============================================================================

-- name: GetEnabledDataSources :many
SELECT 
    id, name, type, source_url, status, enabled, sync_interval,
    ignore_rules, parameters, last_synced, created, updated
FROM core_datasource
WHERE enabled = TRUE
ORDER BY name;

-- name: GetDataSourcesDueForSync :many
SELECT 
    id, name, type, source_url, status, enabled, sync_interval,
    ignore_rules, parameters, last_synced, created, updated
FROM core_datasource
WHERE enabled = TRUE 
  AND (last_synced IS NULL OR NOW() - last_synced >= sync_interval * INTERVAL '1 minute')
ORDER BY last_synced ASC NULLS FIRST
LIMIT $1;

-- name: UpdateDataSourceLastSynced :execrows
UPDATE core_datasource
SET last_synced = $2, updated = NOW()
WHERE id = $1;

-- name: GetDataSourcesByType :many
SELECT 
    id, name, type, source_url, status, enabled, sync_interval,
    ignore_rules, parameters, last_synced, created, updated
FROM core_datasource
WHERE type = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- ============================================================================
-- Additional DataFile Queries
-- ============================================================================

-- name: SearchDataFilesByPathPattern :many
SELECT 
    id, source_id, path, size, hash, data, created, updated
FROM core_datafile
WHERE source_id = $1 AND path LIKE $2
ORDER BY path
LIMIT $3 OFFSET $4;

-- name: GetDataFilesWithHash :many
SELECT 
    id, source_id, path, hash, created
FROM core_datafile
WHERE source_id = $1 AND hash = $2;

-- name: CountDataFilesTotal :one
SELECT COUNT(*)::int as count
FROM core_datafile;

-- name: GetDataFilesSummary :one
SELECT 
    COUNT(*)::int as total_files,
    COALESCE(SUM(size), 0)::bigint as total_size,
    COUNT(DISTINCT source_id)::int as total_sources
FROM core_datafile;

-- ============================================================================
-- Additional Job Queries
-- ============================================================================

-- name: GetJobsByObjectTypeAndID :many
SELECT 
    id, object_type, object_id, name, status, interval, scheduled_at,
    started_at, completed_at, queue_name, job_id, data, error, created, updated
FROM core_job
WHERE object_type = $1 AND object_id = $2
ORDER BY created DESC
LIMIT $3;

-- name: GetRunningJobs :many
SELECT 
    id, object_type, object_id, name, status, interval, scheduled_at,
    started_at, completed_at, queue_name, job_id, data, error, created, updated
FROM core_job
WHERE status IN ('running', 'started')
ORDER BY started_at ASC;

-- name: GetJobStatistics :one
SELECT 
    COUNT(*) FILTER (WHERE status = 'pending')::int as pending_count,
    COUNT(*) FILTER (WHERE status = 'running')::int as running_count,
    COUNT(*) FILTER (WHERE status = 'completed')::int as completed_count,
    COUNT(*) FILTER (WHERE status = 'failed')::int as failed_count,
    COUNT(*) FILTER (WHERE status = 'scheduled')::int as scheduled_count
FROM core_job;

-- name: CancelScheduledJob :execrows
UPDATE core_job
SET status = 'cancelled', completed_at = NOW(), updated = NOW()
WHERE id = $1 AND status = 'scheduled';

-- name: RetryFailedJob :execrows
UPDATE core_job
SET status = 'pending', error = NULL, completed_at = NULL, updated = NOW()
WHERE id = $1 AND status IN ('failed', 'errored');

-- ============================================================================
-- Bulk Operations
-- ============================================================================

-- name: BulkUpdateObjectTypes :exec
INSERT INTO django_content_type (app_label, model, public, features, created, updated)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (app_label, model) DO UPDATE SET
    public = EXCLUDED.public,
    features = EXCLUDED.features,
    updated = EXCLUDED.updated;

-- name: BulkDeleteOldObjectChanges :execrows
DELETE FROM core_objectchange
WHERE id IN (
    SELECT id FROM core_objectchange
    WHERE ($1::timestamp WITH TIME ZONE IS NULL OR time < $1)
    ORDER BY time DESC
    OFFSET $2
);

-- name: ArchiveCompletedJobs :execrows
UPDATE core_job
SET data = NULL
WHERE status IN ('completed', 'cancelled')
  AND completed_at < $1
  AND data IS NOT NULL;
