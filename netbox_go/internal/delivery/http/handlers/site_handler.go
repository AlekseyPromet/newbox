// Package handlers содержит HTTP обработчики для REST API
package handlers

import (
	"net/http"
	"strconv"

	"github.com/AlekseyPromet/netbox_go/internal/domain/dcim/entity"
	"github.com/AlekseyPromet/netbox_go/internal/repository"
	"github.com/labstack/echo/v4"
)

// SiteHandler обрабатывает HTTP запросы для сайтов
type SiteHandler struct {
	repo repository.SiteRepository
}

// NewSiteHandler создает новый экземпляр обработчика сайтов
func NewSiteHandler(repo repository.SiteRepository) *SiteHandler {
	return &SiteHandler{repo: repo}
}

// GetByID получает сайт по ID
// GET /api/dcim/sites/:id
func (h *SiteHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	
	site, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return handleError(err)
	}
	
	return c.JSON(http.StatusOK, site)
}

// GetBySlug получает сайт по slug
// GET /api/dcim/sites/slug/:slug
func (h *SiteHandler) GetBySlug(c echo.Context) error {
	slug := c.Param("slug")
	
	site, err := h.repo.GetBySlug(c.Request().Context(), slug)
	if err != nil {
		return handleError(err)
	}
	
	return c.JSON(http.StatusOK, site)
}

// List получает список сайтов с фильтрацией и пагинацией
// GET /api/dcim/sites?status=active&region_id=xxx&limit=50&offset=0
func (h *SiteHandler) List(c echo.Context) error {
	filter := repository.SiteFilter{}
	
	if status := c.QueryParam("status"); status != "" {
		filter.Status = &status
	}
	if regionID := c.QueryParam("region_id"); regionID != "" {
		filter.RegionID = &regionID
	}
	if groupID := c.QueryParam("group_id"); groupID != "" {
		filter.GroupID = &groupID
	}
	if tenantID := c.QueryParam("tenant_id"); tenantID != "" {
		filter.TenantID = &tenantID
	}
	
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil && limit > 0 {
			filter.Limit = limit
		}
	}
	
	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil && offset >= 0 {
			filter.Offset = offset
		}
	}
	
	sites, total, err := h.repo.List(c.Request().Context(), filter)
	if err != nil {
		return handleError(err)
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"count":  total,
		"next":   getNextURL(c, filter.Limit, filter.Offset+len(sites), total),
		"previous": getPreviousURL(c, filter.Offset),
		"results": sites,
	})
}

// Create создает новый сайт
// POST /api/dcim/sites
func (h *SiteHandler) Create(c echo.Context) error {
	var site entity.Site
	if err := c.Bind(&site); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	
	if err := site.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	
	if err := h.repo.Create(c.Request().Context(), &site); err != nil {
		return handleError(err)
	}
	
	return c.JSON(http.StatusCreated, site)
}

// Update обновляет существующий сайт
// PUT /api/dcim/sites/:id
func (h *SiteHandler) Update(c echo.Context) error {
	id := c.Param("id")
	
	existing, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return handleError(err)
	}
	
	var updateData entity.Site
	if err := c.Bind(&updateData); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	
	// Сохраняем ID
	updateData.ID = existing.ID
	
	if err := updateData.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	
	if err := h.repo.Update(c.Request().Context(), &updateData); err != nil {
		return handleError(err)
	}
	
	updated, err := h.repo.GetByID(c.Request().Context(), id)
	if err != nil {
		return handleError(err)
	}
	
	return c.JSON(http.StatusOK, updated)
}

// Delete удаляет сайт
// DELETE /api/dcim/sites/:id
func (h *SiteHandler) Delete(c echo.Context) error {
	id := c.Param("id")
	
	if err := h.repo.Delete(c.Request().Context(), id); err != nil {
		return handleError(err)
	}
	
	return c.NoContent(http.StatusNoContent)
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
