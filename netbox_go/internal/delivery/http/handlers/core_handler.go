// Package handlers содержит HTTP обработчики для REST API домена Core
package handlers

import (
    "net/http"
    "strconv"
    "time"

    "github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
    "github.com/AlekseyPromet/netbox_go/internal/repository"
    "github.com/AlekseyPromet/netbox_go/pkg/types"
    "github.com/labstack/echo/v4"
)

// CoreHandlers объединяет обработчики Core API (data sources, data files, jobs, object changes, object types, config revisions).
type CoreHandlers struct {
    dataSources   repository.DataSourceRepository
    dataFiles     repository.DataFileRepository
    jobs          repository.JobRepository
    objectChanges repository.ObjectChangeRepository
    objectTypes   repository.ObjectTypeRepository
    configRevisions repository.ConfigRevisionRepository
}

// NewCoreHandlers конструирует CoreHandlers.
func NewCoreHandlers(
    ds repository.DataSourceRepository,
    df repository.DataFileRepository,
    jobs repository.JobRepository,
    oc repository.ObjectChangeRepository,
    ot repository.ObjectTypeRepository,
    cr repository.ConfigRevisionRepository,
) *CoreHandlers {
    return &CoreHandlers{
        dataSources:     ds,
        dataFiles:       df,
        jobs:            jobs,
        objectChanges:   oc,
        objectTypes:     ot,
        configRevisions: cr,
    }
}

// -------- Data Sources --------

// ListDataSources обрабатывает GET /api/core/data-sources
func (h *CoreHandlers) ListDataSources(c echo.Context) error {
    if h.dataSources == nil {
        return notImplemented(c, "DataSourceRepository")
    }

    filter := repository.DataSourceFilter{}

    if v := c.QueryParam("name"); v != "" {
        filter.Name = &v
    }
    if v := c.QueryParam("type"); v != "" {
        filter.Type = &v
    }
    if v := c.QueryParam("status"); v != "" {
        filter.Status = &v
    }
    if v := c.QueryParam("enabled"); v != "" {
        b, err := strconv.ParseBool(v)
        if err == nil {
            filter.Enabled = &b
        }
    }
    if v := c.QueryParam("sync_interval"); v != "" {
        iv, err := strconv.Atoi(v)
        if err == nil {
            filter.SyncInterval = &iv
        }
    }

    filter.Limit = parseLimit(c.QueryParam("limit"))
    filter.Offset = parseOffset(c.QueryParam("offset"))

    items, total, err := h.dataSources.List(c.Request().Context(), filter)
    if err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "count":    total,
        "next":     getNextURL(c, filter.Limit, filter.Offset+len(items), total),
        "previous": getPreviousURL(c, filter.Offset),
        "results":  items,
    })
}

// GetDataSource обрабатывает GET /api/core/data-sources/:id
func (h *CoreHandlers) GetDataSource(c echo.Context) error {
    if h.dataSources == nil {
        return notImplemented(c, "DataSourceRepository")
    }

    id := c.Param("id")
    item, err := h.dataSources.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }
    return c.JSON(http.StatusOK, item)
}

// CreateDataSource обрабатывает POST /api/core/data-sources
func (h *CoreHandlers) CreateDataSource(c echo.Context) error {
    if h.dataSources == nil {
        return notImplemented(c, "DataSourceRepository")
    }

    var ds entity.DataSource
    if err := c.Bind(&ds); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    if err := ds.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := h.dataSources.Create(c.Request().Context(), &ds); err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusCreated, ds)
}

// UpdateDataSource обрабатывает PUT /api/core/data-sources/:id
func (h *CoreHandlers) UpdateDataSource(c echo.Context) error {
    if h.dataSources == nil {
        return notImplemented(c, "DataSourceRepository")
    }

    id := c.Param("id")
    existing, err := h.dataSources.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }

    var input entity.DataSource
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    input.ID = existing.ID
    if err := input.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := h.dataSources.Update(c.Request().Context(), &input); err != nil {
        return handleError(err)
    }

    updated, err := h.dataSources.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusOK, updated)
}

// DeleteDataSource обрабатывает DELETE /api/core/data-sources/:id
func (h *CoreHandlers) DeleteDataSource(c echo.Context) error {
    if h.dataSources == nil {
        return notImplemented(c, "DataSourceRepository")
    }

    id := c.Param("id")
    if err := h.dataSources.Delete(c.Request().Context(), id); err != nil {
        return handleError(err)
    }
    return c.NoContent(http.StatusNoContent)
}

// -------- Data Files --------

// ListDataFiles обрабатывает GET /api/core/data-files
func (h *CoreHandlers) ListDataFiles(c echo.Context) error {
    if h.dataFiles == nil {
        return notImplemented(c, "DataFileRepository")
    }

    filter := repository.DataFileFilter{}
    if v := c.QueryParam("source_id"); v != "" {
        filter.SourceID = &v
    }
    if v := c.QueryParam("path"); v != "" {
        filter.Path = &v
    }
    filter.Limit = parseLimit(c.QueryParam("limit"))
    filter.Offset = parseOffset(c.QueryParam("offset"))

    items, total, err := h.dataFiles.List(c.Request().Context(), filter)
    if err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "count":    total,
        "next":     getNextURL(c, filter.Limit, filter.Offset+len(items), total),
        "previous": getPreviousURL(c, filter.Offset),
        "results":  items,
    })
}

// GetDataFile обрабатывает GET /api/core/data-files/:id
func (h *CoreHandlers) GetDataFile(c echo.Context) error {
    if h.dataFiles == nil {
        return notImplemented(c, "DataFileRepository")
    }

    id := c.Param("id")
    item, err := h.dataFiles.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }
    return c.JSON(http.StatusOK, item)
}

// CreateDataFile обрабатывает POST /api/core/data-files
func (h *CoreHandlers) CreateDataFile(c echo.Context) error {
    if h.dataFiles == nil {
        return notImplemented(c, "DataFileRepository")
    }

    var df entity.DataFile
    if err := c.Bind(&df); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    if err := df.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := h.dataFiles.Create(c.Request().Context(), &df); err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusCreated, df)
}

// UpdateDataFile обрабатывает PUT /api/core/data-files/:id
func (h *CoreHandlers) UpdateDataFile(c echo.Context) error {
    if h.dataFiles == nil {
        return notImplemented(c, "DataFileRepository")
    }

    id := c.Param("id")
    existing, err := h.dataFiles.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }

    var input entity.DataFile
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    input.ID = existing.ID
    if err := input.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := h.dataFiles.Update(c.Request().Context(), &input); err != nil {
        return handleError(err)
    }

    updated, err := h.dataFiles.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusOK, updated)
}

// DeleteDataFile обрабатывает DELETE /api/core/data-files/:id
func (h *CoreHandlers) DeleteDataFile(c echo.Context) error {
    if h.dataFiles == nil {
        return notImplemented(c, "DataFileRepository")
    }

    id := c.Param("id")
    if err := h.dataFiles.Delete(c.Request().Context(), id); err != nil {
        return handleError(err)
    }
    return c.NoContent(http.StatusNoContent)
}

// -------- Jobs --------

// ListJobs обрабатывает GET /api/core/jobs
func (h *CoreHandlers) ListJobs(c echo.Context) error {
    if h.jobs == nil {
        return notImplemented(c, "JobRepository")
    }

    filter := repository.JobFilter{}
    if v := c.QueryParam("object_type"); v != "" {
        filter.ObjectType = &v
    }
    if v := c.QueryParam("object_id"); v != "" {
        filter.ObjectID = &v
    }
    if v := c.QueryParam("status"); v != "" {
        filter.Status = &v
    }
    if v := c.QueryParam("queue_name"); v != "" {
        filter.QueueName = &v
    }
    if v := c.QueryParam("scheduled_at"); v != "" {
        if t, err := time.Parse(time.RFC3339, v); err == nil {
            filter.ScheduledAt = &t
        }
    }
    filter.Limit = parseLimit(c.QueryParam("limit"))
    filter.Offset = parseOffset(c.QueryParam("offset"))

    items, total, err := h.jobs.List(c.Request().Context(), filter)
    if err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "count":    total,
        "next":     getNextURL(c, filter.Limit, filter.Offset+len(items), total),
        "previous": getPreviousURL(c, filter.Offset),
        "results":  items,
    })
}

// GetJob обрабатывает GET /api/core/jobs/:id
func (h *CoreHandlers) GetJob(c echo.Context) error {
    if h.jobs == nil {
        return notImplemented(c, "JobRepository")
    }

    id := c.Param("id")
    item, err := h.jobs.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }
    return c.JSON(http.StatusOK, item)
}

// CreateJob обрабатывает POST /api/core/jobs
func (h *CoreHandlers) CreateJob(c echo.Context) error {
    if h.jobs == nil {
        return notImplemented(c, "JobRepository")
    }

    var job entity.Job
    if err := c.Bind(&job); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    if err := job.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := h.jobs.Create(c.Request().Context(), &job); err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusCreated, job)
}

// -------- Object Changes --------

// ListObjectChanges обрабатывает GET /api/core/object-changes
func (h *CoreHandlers) ListObjectChanges(c echo.Context) error {
    if h.objectChanges == nil {
        return notImplemented(c, "ObjectChangeRepository")
    }

    filter := repository.ObjectChangeFilter{}
    if v := c.QueryParam("changed_object_type"); v != "" {
        filter.ChangedObjectType = &v
    }
    if v := c.QueryParam("changed_object_id"); v != "" {
        filter.ChangedObjectID = &v
    }
    if v := c.QueryParam("user_id"); v != "" {
        filter.UserID = &v
    }
    if v := c.QueryParam("action"); v != "" {
        filter.Action = &v
    }
    if v := c.QueryParam("request_id"); v != "" {
        filter.RequestID = &v
    }
    if v := c.QueryParam("since"); v != "" {
        if t, err := time.Parse(time.RFC3339, v); err == nil {
            filter.Since = &t
        }
    }
    if v := c.QueryParam("until"); v != "" {
        if t, err := time.Parse(time.RFC3339, v); err == nil {
            filter.Until = &t
        }
    }
    filter.Limit = parseLimit(c.QueryParam("limit"))
    filter.Offset = parseOffset(c.QueryParam("offset"))

    items, total, err := h.objectChanges.List(c.Request().Context(), filter)
    if err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "count":    total,
        "next":     getNextURL(c, filter.Limit, filter.Offset+len(items), total),
        "previous": getPreviousURL(c, filter.Offset),
        "results":  items,
    })
}

// GetObjectChange обрабатывает GET /api/core/object-changes/:id
func (h *CoreHandlers) GetObjectChange(c echo.Context) error {
    if h.objectChanges == nil {
        return notImplemented(c, "ObjectChangeRepository")
    }

    id := c.Param("id")
    item, err := h.objectChanges.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }
    return c.JSON(http.StatusOK, item)
}

// LogObjectChange обрабатывает POST /api/core/object-changes/log
func (h *CoreHandlers) LogObjectChange(c echo.Context) error {
    if h.objectChanges == nil {
        return notImplemented(c, "ObjectChangeRepository")
    }

    var req struct {
        Action       string      `json:"action"`
        ObjectType   string      `json:"object_type"`
        ObjectID     string      `json:"object_id"`
        ObjectRepr   string      `json:"object_repr"`
        ObjectData   interface{} `json:"object_data,omitempty"`
        UserID       *string     `json:"user_id,omitempty"`
        RequestID    *string     `json:"request_id,omitempty"`
    }

    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    if req.Action == "" || req.ObjectType == "" || req.ObjectID == "" || req.ObjectRepr == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "missing required fields"})
    }

    var userID *types.ID
    if req.UserID != nil {
        uid := types.ID(*req.UserID)
        userID = &uid
    }

    action := types.Status(req.Action)
    if err := h.objectChanges.LogChange(c.Request().Context(), action, req.ObjectType, req.ObjectID, req.ObjectRepr, req.ObjectData, userID, req.RequestID); err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusCreated, map[string]string{"status": "logged"})
}

// -------- Object Types --------

// ListObjectTypes обрабатывает GET /api/core/object-types
func (h *CoreHandlers) ListObjectTypes(c echo.Context) error {
    if h.objectTypes == nil {
        return notImplemented(c, "ObjectTypeRepository")
    }

    filter := repository.ObjectTypeFilter{}
    if v := c.QueryParam("app_label"); v != "" {
        filter.AppLabel = &v
    }
    if v := c.QueryParam("model"); v != "" {
        filter.Model = &v
    }
    if v := c.QueryParam("public"); v != "" {
        if b, err := strconv.ParseBool(v); err == nil {
            filter.Public = &b
        }
    }
    if v := c.QueryParam("feature"); v != "" {
        filter.Feature = &v
    }
    filter.Limit = parseLimit(c.QueryParam("limit"))
    filter.Offset = parseOffset(c.QueryParam("offset"))

    items, total, err := h.objectTypes.List(c.Request().Context(), filter)
    if err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "count":    total,
        "next":     getNextURL(c, filter.Limit, filter.Offset+len(items), total),
        "previous": getPreviousURL(c, filter.Offset),
        "results":  items,
    })
}

// GetObjectType обрабатывает GET /api/core/object-types/:id
func (h *CoreHandlers) GetObjectType(c echo.Context) error {
    if h.objectTypes == nil {
        return notImplemented(c, "ObjectTypeRepository")
    }

    id := c.Param("id")
    item, err := h.objectTypes.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }
    return c.JSON(http.StatusOK, item)
}

// -------- Config Revisions --------

// ListConfigRevisions обрабатывает GET /api/core/config-revisions
func (h *CoreHandlers) ListConfigRevisions(c echo.Context) error {
    if h.configRevisions == nil {
        return notImplemented(c, "ConfigRevisionRepository")
    }

    filter := repository.ConfigRevisionFilter{}
    if v := c.QueryParam("active"); v != "" {
        b, err := strconv.ParseBool(v)
        if err == nil {
            filter.Active = &b
        }
    }
    if v := c.QueryParam("created_since"); v != "" {
        if t, err := time.Parse(time.RFC3339, v); err == nil {
            filter.CreatedSince = &t
        }
    }
    if v := c.QueryParam("created_until"); v != "" {
        if t, err := time.Parse(time.RFC3339, v); err == nil {
            filter.CreatedUntil = &t
        }
    }
    filter.Limit = parseLimit(c.QueryParam("limit"))
    filter.Offset = parseOffset(c.QueryParam("offset"))

    items, total, err := h.configRevisions.List(c.Request().Context(), filter)
    if err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "count":    total,
        "next":     getNextURL(c, filter.Limit, filter.Offset+len(items), total),
        "previous": getPreviousURL(c, filter.Offset),
        "results":  items,
    })
}

// GetConfigRevision обрабатывает GET /api/core/config-revisions/:id
func (h *CoreHandlers) GetConfigRevision(c echo.Context) error {
    if h.configRevisions == nil {
        return notImplemented(c, "ConfigRevisionRepository")
    }

    id := c.Param("id")
    item, err := h.configRevisions.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }
    return c.JSON(http.StatusOK, item)
}

// CreateConfigRevision обрабатывает POST /api/core/config-revisions
func (h *CoreHandlers) CreateConfigRevision(c echo.Context) error {
    if h.configRevisions == nil {
        return notImplemented(c, "ConfigRevisionRepository")
    }

    var cr entity.ConfigRevision
    if err := c.Bind(&cr); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    if err := cr.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := h.configRevisions.Create(c.Request().Context(), &cr); err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusCreated, cr)
}

// UpdateConfigRevision обрабатывает PUT /api/core/config-revisions/:id
func (h *CoreHandlers) UpdateConfigRevision(c echo.Context) error {
    if h.configRevisions == nil {
        return notImplemented(c, "ConfigRevisionRepository")
    }

    id := c.Param("id")
    existing, err := h.configRevisions.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }

    var input entity.ConfigRevision
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
    }

    input.ID = existing.ID
    if err := input.Validate(); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
    }

    if err := h.configRevisions.Update(c.Request().Context(), &input); err != nil {
        return handleError(err)
    }

    updated, err := h.configRevisions.GetByID(c.Request().Context(), id)
    if err != nil {
        return handleError(err)
    }

    return c.JSON(http.StatusOK, updated)
}

// DeleteConfigRevision обрабатывает DELETE /api/core/config-revisions/:id
func (h *CoreHandlers) DeleteConfigRevision(c echo.Context) error {
    if h.configRevisions == nil {
        return notImplemented(c, "ConfigRevisionRepository")
    }

    id := c.Param("id")
    if err := h.configRevisions.Delete(c.Request().Context(), id); err != nil {
        return handleError(err)
    }
    return c.NoContent(http.StatusNoContent)
}

// ActivateConfigRevision обрабатывает POST /api/core/config-revisions/:id/activate
func (h *CoreHandlers) ActivateConfigRevision(c echo.Context) error {
    if h.configRevisions == nil {
        return notImplemented(c, "ConfigRevisionRepository")
    }

    id := c.Param("id")
    if err := h.configRevisions.Activate(c.Request().Context(), id); err != nil {
        return handleError(err)
    }
    return c.JSON(http.StatusOK, map[string]string{"status": "activated"})
}

// GetActiveConfigRevision обрабатывает GET /api/core/config-revisions/active
func (h *CoreHandlers) GetActiveConfigRevision(c echo.Context) error {
    if h.configRevisions == nil {
        return notImplemented(c, "ConfigRevisionRepository")
    }

    item, err := h.configRevisions.GetActive(c.Request().Context())
    if err != nil {
        return handleError(err)
    }
    return c.JSON(http.StatusOK, item)
}

// -------- Background (RQ) --------

// ListBackgroundQueues обрабатывает GET /api/core/background-queues
func (h *CoreHandlers) ListBackgroundQueues(c echo.Context) error {
    return notImplemented(c, "Background queues are not yet supported in Go port")
}

// GetBackgroundQueue обрабатывает GET /api/core/background-queues/:name
func (h *CoreHandlers) GetBackgroundQueue(c echo.Context) error {
    return notImplemented(c, "Background queues are not yet supported in Go port")
}

// ListBackgroundWorkers обрабатывает GET /api/core/background-workers
func (h *CoreHandlers) ListBackgroundWorkers(c echo.Context) error {
    return notImplemented(c, "Background workers are not yet supported in Go port")
}

// GetBackgroundWorker обрабатывает GET /api/core/background-workers/:name
func (h *CoreHandlers) GetBackgroundWorker(c echo.Context) error {
    return notImplemented(c, "Background workers are not yet supported in Go port")
}

// ListBackgroundTasks обрабатывает GET /api/core/background-tasks
func (h *CoreHandlers) ListBackgroundTasks(c echo.Context) error {
    return notImplemented(c, "Background tasks are not yet supported in Go port")
}

// GetBackgroundTask обрабатывает GET /api/core/background-tasks/:id
func (h *CoreHandlers) GetBackgroundTask(c echo.Context) error {
    return notImplemented(c, "Background tasks are not yet supported in Go port")
}

// DeleteBackgroundTask обрабатывает POST /api/core/background-tasks/:id/delete
func (h *CoreHandlers) DeleteBackgroundTask(c echo.Context) error {
    return notImplemented(c, "Background task delete is not yet supported in Go port")
}

// RequeueBackgroundTask обрабатывает POST /api/core/background-tasks/:id/requeue
func (h *CoreHandlers) RequeueBackgroundTask(c echo.Context) error {
    return notImplemented(c, "Background task requeue is not yet supported in Go port")
}

// EnqueueBackgroundTask обрабатывает POST /api/core/background-tasks/:id/enqueue
func (h *CoreHandlers) EnqueueBackgroundTask(c echo.Context) error {
    return notImplemented(c, "Background task enqueue is not yet supported in Go port")
}

// StopBackgroundTask обрабатывает POST /api/core/background-tasks/:id/stop
func (h *CoreHandlers) StopBackgroundTask(c echo.Context) error {
    return notImplemented(c, "Background task stop is not yet supported in Go port")
}

// -------- Вспомогательные функции --------

// notImplemented возвращает 501 Not Implemented с сообщением.
func notImplemented(c echo.Context, msg string) error {
    return c.JSON(http.StatusNotImplemented, map[string]string{"error": msg})
}

// parseLimit извлекает limit (по умолчанию 100, максимум 1000).
func parseLimit(raw string) int {
    if raw == "" {
        return 100
    }
    v, err := strconv.Atoi(raw)
    if err != nil || v <= 0 {
        return 100
    }
    if v > 1000 {
        return 1000
    }
    return v
}

// parseOffset извлекает offset (по умолчанию 0).
func parseOffset(raw string) int {
    if raw == "" {
        return 0
    }
    v, err := strconv.Atoi(raw)
    if err != nil || v < 0 {
        return 0
    }
    return v
}

// handleError преобразует ошибки в HTTP ответы
func handleError(err error) error {
    switch err {
    case repository.ErrNotFound:
        return echo.NewHTTPError(http.StatusNotFound, "entity not found")
    default:
        return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
    }
}

// getNextURL генерирует URL для следующей страницы
func getNextURL(c echo.Context, limit, offset int, total int64) interface{} {
    if offset+limit >= int(total) {
        return nil
    }
    
    url := c.Request().URL.Path + "?limit=" + strconv.Itoa(limit) + "&offset=" + strconv.Itoa(offset)
    return url
}

// getPreviousURL генерирует URL для предыдущей страницы
func getPreviousURL(c echo.Context, offset int) interface{} {
    if offset <= 0 {
        return nil
    }
    
    prevOffset := offset - 100
    if prevOffset < 0 {
        prevOffset = 0
    }
    
    url := c.Request().URL.Path + "?offset=" + strconv.Itoa(prevOffset)
    return url
}
