package handlers

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"netbox_go/internal/domain/telemetry/entity"
	"netbox_go/internal/domain/telemetry/service"
)

// TelemetryHandler handles HTTP requests for telemetry operations
type TelemetryHandler struct {
	service *service.TelemetryService
	logger  *zap.Logger
}

// NewTelemetryHandler creates a new telemetry handler
func NewTelemetryHandler(svc *service.TelemetryService, logger *zap.Logger) *TelemetryHandler {
	return &TelemetryHandler{
		service: svc,
		logger:  logger,
	}
}

// RegisterRoutes registers telemetry routes
func (h *TelemetryHandler) RegisterRoutes(e *echo.Echo) {
	telemetry := e.Group("/api/telemetry")

	// Device routes
	telemetry.GET("/devices", h.ListDevices)
	telemetry.GET("/devices/:id", h.GetDevice)
	telemetry.POST("/devices", h.CreateDevice)
	telemetry.PUT("/devices/:id", h.UpdateDevice)
	telemetry.DELETE("/devices/:id", h.DeleteDevice)

	// Collection routes
	telemetry.GET("/collections", h.ListCollections)
	telemetry.GET("/collections/:id", h.GetCollection)
	telemetry.POST("/collections", h.CreateCollection)
	telemetry.PUT("/collections/:id", h.UpdateCollection)
	telemetry.DELETE("/collections/:id", h.DeleteCollection)

	// Job routes
	telemetry.GET("/jobs", h.ListJobs)
	telemetry.GET("/jobs/:id", h.GetJob)

	// Ping target routes
	telemetry.GET("/ping-targets", h.ListPingTargets)
	telemetry.GET("/ping-targets/:id", h.GetPingTarget)
	telemetry.POST("/ping-targets", h.CreatePingTarget)
	telemetry.PUT("/ping-targets/:id", h.UpdatePingTarget)
	telemetry.DELETE("/ping-targets/:id", h.DeletePingTarget)

	// DNS query routes
	telemetry.GET("/dns-queries", h.ListDNSQueries)
	telemetry.GET("/dns-queries/:id", h.GetDNSQuery)
	telemetry.POST("/dns-queries", h.CreateDNSQuery)
	telemetry.PUT("/dns-queries/:id", h.UpdateDNSQuery)
	telemetry.DELETE("/dns-queries/:id", h.DeleteDNSQuery)
}

// Device handlers

// ListDevices returns all telemetry devices
func (h *TelemetryHandler) ListDevices(c echo.Context) error {
	ctx := c.Request().Context()

	devices, err := h.service.ListDevices(ctx)
	if err != nil {
		h.logger.Error("failed to list devices", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, devices)
}

// GetDevice returns a telemetry device by ID
func (h *TelemetryHandler) GetDevice(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid device ID"})
	}

	device, err := h.service.GetDevice(ctx, id)
	if err != nil {
		h.logger.Error("failed to get device", zap.Error(err), zap.String("id", idStr))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "device not found"})
	}

	return c.JSON(http.StatusOK, device)
}

// CreateDeviceRequest represents a device creation request
type CreateDeviceRequest struct {
	DeviceID        uuid.UUID `json:"device_id"`
	CollectionType  string    `json:"collection_type"`
	GNMIAddress     string    `json:"gnmi_address"`
	GNmiPort        int       `json:"gnmi_port"`
	VaultSecretPath string    `json:"vault_secret_path"`
	Enabled         bool      `json:"enabled"`
}

// CreateDevice creates a new telemetry device
func (h *TelemetryHandler) CreateDevice(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateDeviceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	device := &entity.TelemetryDevice{
		ID:              uuid.New(),
		DeviceID:        req.DeviceID,
		CollectionType:  entity.CollectionType(req.CollectionType),
		GNMIAddress:     req.GNMIAddress,
		GNmiPort:        req.GNmiPort,
		VaultSecretPath: req.VaultSecretPath,
		Enabled:         req.Enabled,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.service.CreateDevice(ctx, device); err != nil {
		h.logger.Error("failed to create device", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, device)
}

// UpdateDevice updates a telemetry device
func (h *TelemetryHandler) UpdateDevice(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid device ID"})
	}

	var req CreateDeviceRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	device, err := h.service.GetDevice(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "device not found"})
	}

	device.CollectionType = entity.CollectionType(req.CollectionType)
	device.GNMIAddress = req.GNMIAddress
	device.GNmiPort = req.GNmiPort
	device.VaultSecretPath = req.VaultSecretPath
	device.Enabled = req.Enabled
	device.UpdatedAt = time.Now()

	if err := h.service.UpdateDevice(ctx, device); err != nil {
		h.logger.Error("failed to update device", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, device)
}

// DeleteDevice deletes a telemetry device
func (h *TelemetryHandler) DeleteDevice(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid device ID"})
	}

	if err := h.service.DeleteDevice(ctx, id); err != nil {
		h.logger.Error("failed to delete device", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}

// Collection handlers

// ListCollections returns all telemetry collections
func (h *TelemetryHandler) ListCollections(c echo.Context) error {
	ctx := c.Request().Context()

	collections, err := h.service.ListCollections(ctx)
	if err != nil {
		h.logger.Error("failed to list collections", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, collections)
}

// GetCollection returns a telemetry collection by ID
func (h *TelemetryHandler) GetCollection(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid collection ID"})
	}

	collection, err := h.service.GetCollection(ctx, id)
	if err != nil {
		h.logger.Error("failed to get collection", zap.Error(err), zap.String("id", idStr))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "collection not found"})
	}

	return c.JSON(http.StatusOK, collection)
}

// CreateCollectionRequest represents a collection creation request
type CreateCollectionRequest struct {
	Name            string `json:"name"`
	Description     string `json:"description"`
	CollectionType  string `json:"collection_type"`
	TelemetryType   string `json:"telemetry_type"`
	TargetPath      string `json:"target_path"`
	IntervalSeconds int    `json:"interval_seconds"`
	Enabled         bool   `json:"enabled"`
	Filters         string `json:"filters"`
}

// CreateCollection creates a new telemetry collection
func (h *TelemetryHandler) CreateCollection(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateCollectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	collection := &entity.TelemetryCollection{
		ID:              uuid.New(),
		Name:            req.Name,
		Description:     req.Description,
		CollectionType:  entity.CollectionType(req.CollectionType),
		TelemetryType:   entity.TelemetryType(req.TelemetryType),
		TargetPath:      req.TargetPath,
		IntervalSeconds: req.IntervalSeconds,
		Enabled:         req.Enabled,
		Filters:         req.Filters,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.service.CreateCollection(ctx, collection); err != nil {
		h.logger.Error("failed to create collection", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, collection)
}

// UpdateCollection updates a telemetry collection
func (h *TelemetryHandler) UpdateCollection(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid collection ID"})
	}

	var req CreateCollectionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	collection, err := h.service.GetCollection(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "collection not found"})
	}

	collection.Name = req.Name
	collection.Description = req.Description
	collection.CollectionType = entity.CollectionType(req.CollectionType)
	collection.TelemetryType = entity.TelemetryType(req.TelemetryType)
	collection.TargetPath = req.TargetPath
	collection.IntervalSeconds = req.IntervalSeconds
	collection.Enabled = req.Enabled
	collection.Filters = req.Filters
	collection.UpdatedAt = time.Now()

	if err := h.service.UpdateCollection(ctx, collection); err != nil {
		h.logger.Error("failed to update collection", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, collection)
}

// DeleteCollection deletes a telemetry collection
func (h *TelemetryHandler) DeleteCollection(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid collection ID"})
	}

	if err := h.service.DeleteCollection(ctx, id); err != nil {
		h.logger.Error("failed to delete collection", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}

// Job handlers

// ListJobs returns recent collection jobs
func (h *TelemetryHandler) ListJobs(c echo.Context) error {
	ctx := c.Request().Context()

	limit := 100
	jobs, err := h.service.ListJobs(ctx, limit)
	if err != nil {
		h.logger.Error("failed to list jobs", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, jobs)
}

// GetJob returns a collection job by ID
func (h *TelemetryHandler) GetJob(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid job ID"})
	}

	job, err := h.service.GetJob(ctx, id)
	if err != nil {
		h.logger.Error("failed to get job", zap.Error(err), zap.String("id", idStr))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "job not found"})
	}

	return c.JSON(http.StatusOK, job)
}

// Ping target handlers

// ListPingTargets returns all ping targets
func (h *TelemetryHandler) ListPingTargets(c echo.Context) error {
	ctx := c.Request().Context()

	targets, err := h.service.ListPingTargets(ctx)
	if err != nil {
		h.logger.Error("failed to list ping targets", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, targets)
}

// GetPingTarget returns a ping target by ID
func (h *TelemetryHandler) GetPingTarget(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ping target ID"})
	}

	target, err := h.service.GetPingTarget(ctx, id)
	if err != nil {
		h.logger.Error("failed to get ping target", zap.Error(err), zap.String("id", idStr))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "ping target not found"})
	}

	return c.JSON(http.StatusOK, target)
}

// CreatePingTargetRequest represents a ping target creation request
type CreatePingTargetRequest struct {
	DeviceID        uuid.UUID `json:"device_id"`
	TargetAddress   string    `json:"target_address"`
	TargetType      string    `json:"target_type"`
	IntervalSeconds int       `json:"interval_seconds"`
	PacketCount     int       `json:"packet_count"`
	PacketSize      int       `json:"packet_size"`
	TimeoutSeconds  int       `json:"timeout_seconds"`
	Enabled         bool      `json:"enabled"`
}

// CreatePingTarget creates a new ping target
func (h *TelemetryHandler) CreatePingTarget(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreatePingTargetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	target := &entity.PingTarget{
		ID:              uuid.New(),
		DeviceID:        req.DeviceID,
		TargetAddress:   req.TargetAddress,
		TargetType:      req.TargetType,
		IntervalSeconds: req.IntervalSeconds,
		PacketCount:     req.PacketCount,
		PacketSize:      req.PacketSize,
		TimeoutSeconds:  req.TimeoutSeconds,
		Enabled:         req.Enabled,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.service.CreatePingTarget(ctx, target); err != nil {
		h.logger.Error("failed to create ping target", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, target)
}

// UpdatePingTarget updates a ping target
func (h *TelemetryHandler) UpdatePingTarget(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ping target ID"})
	}

	var req CreatePingTargetRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	target, err := h.service.GetPingTarget(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "ping target not found"})
	}

	target.DeviceID = req.DeviceID
	target.TargetAddress = req.TargetAddress
	target.TargetType = req.TargetType
	target.IntervalSeconds = req.IntervalSeconds
	target.PacketCount = req.PacketCount
	target.PacketSize = req.PacketSize
	target.TimeoutSeconds = req.TimeoutSeconds
	target.Enabled = req.Enabled
	target.UpdatedAt = time.Now()

	if err := h.service.UpdatePingTarget(ctx, target); err != nil {
		h.logger.Error("failed to update ping target", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, target)
}

// DeletePingTarget deletes a ping target
func (h *TelemetryHandler) DeletePingTarget(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid ping target ID"})
	}

	if err := h.service.DeletePingTarget(ctx, id); err != nil {
		h.logger.Error("failed to delete ping target", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}

// DNS query handlers

// ListDNSQueries returns all DNS queries
func (h *TelemetryHandler) ListDNSQueries(c echo.Context) error {
	ctx := c.Request().Context()

	queries, err := h.service.ListDNSQueries(ctx)
	if err != nil {
		h.logger.Error("failed to list DNS queries", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, queries)
}

// GetDNSQuery returns a DNS query by ID
func (h *TelemetryHandler) GetDNSQuery(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid DNS query ID"})
	}

	query, err := h.service.GetDNSQuery(ctx, id)
	if err != nil {
		h.logger.Error("failed to get DNS query", zap.Error(err), zap.String("id", idStr))
		return c.JSON(http.StatusNotFound, map[string]string{"error": "DNS query not found"})
	}

	return c.JSON(http.StatusOK, query)
}

// CreateDNSQueryRequest represents a DNS query creation request
type CreateDNSQueryRequest struct {
	DeviceID        uuid.UUID `json:"device_id"`
	QueryName       string    `json:"query_name"`
	QueryType       string    `json:"query_type"`
	DNSServer       string    `json:"dns_server"`
	IntervalSeconds int       `json:"interval_seconds"`
	Enabled         bool      `json:"enabled"`
}

// CreateDNSQuery creates a new DNS query
func (h *TelemetryHandler) CreateDNSQuery(c echo.Context) error {
	ctx := c.Request().Context()

	var req CreateDNSQueryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	query := &entity.DNSQuery{
		ID:              uuid.New(),
		DeviceID:        req.DeviceID,
		QueryName:       req.QueryName,
		QueryType:       req.QueryType,
		DNSServer:       req.DNSServer,
		IntervalSeconds: req.IntervalSeconds,
		Enabled:         req.Enabled,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.service.CreateDNSQuery(ctx, query); err != nil {
		h.logger.Error("failed to create DNS query", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, query)
}

// UpdateDNSQuery updates a DNS query
func (h *TelemetryHandler) UpdateDNSQuery(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid DNS query ID"})
	}

	var req CreateDNSQueryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	query, err := h.service.GetDNSQuery(ctx, id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "DNS query not found"})
	}

	query.DeviceID = req.DeviceID
	query.QueryName = req.QueryName
	query.QueryType = req.QueryType
	query.DNSServer = req.DNSServer
	query.IntervalSeconds = req.IntervalSeconds
	query.Enabled = req.Enabled
	query.UpdatedAt = time.Now()

	if err := h.service.UpdateDNSQuery(ctx, query); err != nil {
		h.logger.Error("failed to update DNS query", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, query)
}

// DeleteDNSQuery deletes a DNS query
func (h *TelemetryHandler) DeleteDNSQuery(c echo.Context) error {
	ctx := c.Request().Context()
	idStr := c.Param("id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid DNS query ID"})
	}

	if err := h.service.DeleteDNSQuery(ctx, id); err != nil {
		h.logger.Error("failed to delete DNS query", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusNoContent, nil)
}
