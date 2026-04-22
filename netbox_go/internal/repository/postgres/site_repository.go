// Package postgres содержит реализацию репозиториев для PostgreSQL с использованием sqlc
package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"netbox_go/internal/domain/dcim/entity"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
)

// SiteRepositoryPostgres реализует интерфейс SiteRepository для PostgreSQL
type SiteRepositoryPostgres struct {
	db *sql.DB
}

// NewSiteRepositoryPostgres создает новый экземпляр репозитория сайтов
func NewSiteRepositoryPostgres(db *sql.DB) *SiteRepositoryPostgres {
	return &SiteRepositoryPostgres{db: db}
}

// GetByID получает сайт по ID
func (r *SiteRepositoryPostgres) GetByID(ctx context.Context, id string) (*entity.Site, error) {
	query := `
		SELECT id, name, slug, status, region_id, group_id, tenant_id, facility, asn_ids, 
		       time_zone, physical_address, shipping_address, latitude, longitude, 
		       description, comments, created, updated
		FROM dcim_sites
		WHERE id = $1 AND deleted_at IS NULL
	`
	
	var site entity.Site
	var regionID, groupID, tenantID sql.NullString
	var facility sql.NullString
	var timeZone sql.NullString
	var physicalAddr, shippingAddr sql.NullString
	var latitude, longitude sql.NullFloat64
	var description, comments sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&site.ID, &site.Name, &site.Slug, &site.Status,
		&regionID, &groupID, &tenantID, &facility,
		&timeZone, &physicalAddr, &shippingAddr,
		&latitude, &longitude,
		&description, &comments,
		&site.Created, &site.Updated,
	)
	
	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get site by ID: %w", err)
	}
	
	// Заполнение опциональных полей
	if regionID.Valid {
		rid, _ := types.ParseID(regionID.String)
		site.RegionID = &rid
	}
	if groupID.Valid {
		gid, _ := types.ParseID(groupID.String)
		site.GroupID = &gid
	}
	if tenantID.Valid {
		tid, _ := types.ParseID(tenantID.String)
		site.TenantID = &tid
	}
	if facility.Valid {
		site.Facility = types.Facility(facility.String)
	}
	if timeZone.Valid {
		site.TimeZone = types.TimeZone(timeZone.String)
	}
	if physicalAddr.Valid {
		site.PhysicalAddress = types.Address(physicalAddr.String)
	}
	if shippingAddr.Valid {
		site.ShippingAddress = types.Address(shippingAddr.String)
	}
	if latitude.Valid {
		site.Latitude = &latitude.Float64
	}
	if longitude.Valid {
		site.Longitude = &longitude.Float64
	}
	if description.Valid {
		site.Description = types.Description(description.String)
	}
	if comments.Valid {
		site.Comments = types.Comments(comments.String)
	}
	
	return &site, nil
}

// GetBySlug получает сайт по slug
func (r *SiteRepositoryPostgres) GetBySlug(ctx context.Context, slug string) (*entity.Site, error) {
	query := `
		SELECT id, name, slug, status, region_id, group_id, tenant_id, facility, asn_ids, 
		       time_zone, physical_address, shipping_address, latitude, longitude, 
		       description, comments, created, updated
		FROM dcim_sites
		WHERE slug = $1 AND deleted_at IS NULL
	`
	
	var site entity.Site
	var regionID, groupID, tenantID sql.NullString
	var facility sql.NullString
	var timeZone sql.NullString
	var physicalAddr, shippingAddr sql.NullString
	var latitude, longitude sql.NullFloat64
	var description, comments sql.NullString
	
	err := r.db.QueryRowContext(ctx, query, slug).Scan(
		&site.ID, &site.Name, &site.Slug, &site.Status,
		&regionID, &groupID, &tenantID, &facility,
		&timeZone, &physicalAddr, &shippingAddr,
		&latitude, &longitude,
		&description, &comments,
		&site.Created, &site.Updated,
	)
	
	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get site by slug: %w", err)
	}
	
	// Заполнение опциональных полей (аналогично GetByID)
	if regionID.Valid {
		rid, _ := types.ParseID(regionID.String)
		site.RegionID = &rid
	}
	if groupID.Valid {
		gid, _ := types.ParseID(groupID.String)
		site.GroupID = &gid
	}
	if tenantID.Valid {
		tid, _ := types.ParseID(tenantID.String)
		site.TenantID = &tid
	}
	if facility.Valid {
		site.Facility = types.Facility(facility.String)
	}
	if timeZone.Valid {
		site.TimeZone = types.TimeZone(timeZone.String)
	}
	if physicalAddr.Valid {
		site.PhysicalAddress = types.Address(physicalAddr.String)
	}
	if shippingAddr.Valid {
		site.ShippingAddress = types.Address(shippingAddr.String)
	}
	if latitude.Valid {
		site.Latitude = &latitude.Float64
	}
	if longitude.Valid {
		site.Longitude = &longitude.Float64
	}
	if description.Valid {
		site.Description = types.Description(description.String)
	}
	if comments.Valid {
		site.Comments = types.Comments(comments.String)
	}
	
	return &site, nil
}

// List получает список сайтов с фильтрацией и пагинацией
func (r *SiteRepositoryPostgres) List(ctx context.Context, filter repository.SiteFilter) ([]*entity.Site, int64, error) {
	query := `
		SELECT id, name, slug, status, region_id, group_id, tenant_id, facility, asn_ids, 
		       time_zone, physical_address, shipping_address, latitude, longitude, 
		       description, comments, created, updated
		FROM dcim_sites
		WHERE deleted_at IS NULL
	`
	
	args := []interface{}{}
	argIndex := 1
	
	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}
	if filter.RegionID != nil {
		query += fmt.Sprintf(" AND region_id = $%d", argIndex)
		args = append(args, *filter.RegionID)
		argIndex++
	}
	if filter.GroupID != nil {
		query += fmt.Sprintf(" AND group_id = $%d", argIndex)
		args = append(args, *filter.GroupID)
		argIndex++
	}
	if filter.TenantID != nil {
		query += fmt.Sprintf(" AND tenant_id = $%d", argIndex)
		args = append(args, *filter.TenantID)
		argIndex++
	}
	
	// Получение общего количества
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", query)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count sites: %w", err)
	}
	
	// Пагинация
	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	offset := filter.Offset
	
	query += fmt.Sprintf(" ORDER BY created DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list sites: %w", err)
	}
	defer rows.Close()
	
	sites := make([]*entity.Site, 0)
	for rows.Next() {
		var site entity.Site
		var regionID, groupID, tenantID sql.NullString
		var facility sql.NullString
		var timeZone sql.NullString
		var physicalAddr, shippingAddr sql.NullString
		var latitude, longitude sql.NullFloat64
		var description, comments sql.NullString
		
		err := rows.Scan(
			&site.ID, &site.Name, &site.Slug, &site.Status,
			&regionID, &groupID, &tenantID, &facility,
			&timeZone, &physicalAddr, &shippingAddr,
			&latitude, &longitude,
			&description, &comments,
			&site.Created, &site.Updated,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan site: %w", err)
		}
		
		// Заполнение опциональных полей
		if regionID.Valid {
			rid, _ := types.ParseID(regionID.String)
			site.RegionID = &rid
		}
		if groupID.Valid {
			gid, _ := types.ParseID(groupID.String)
			site.GroupID = &gid
		}
		if tenantID.Valid {
			tid, _ := types.ParseID(tenantID.String)
			site.TenantID = &tid
		}
		if facility.Valid {
			site.Facility = types.Facility(facility.String)
		}
		if timeZone.Valid {
			site.TimeZone = types.TimeZone(timeZone.String)
		}
		if physicalAddr.Valid {
			site.PhysicalAddress = types.Address(physicalAddr.String)
		}
		if shippingAddr.Valid {
			site.ShippingAddress = types.Address(shippingAddr.String)
		}
		if latitude.Valid {
			site.Latitude = &latitude.Float64
		}
		if longitude.Valid {
			site.Longitude = &longitude.Float64
		}
		if description.Valid {
			site.Description = types.Description(description.String)
		}
		if comments.Valid {
			site.Comments = types.Comments(comments.String)
		}
		
		sites = append(sites, &site)
	}
	
	return sites, total, nil
}

// Create создает новый сайт
func (r *SiteRepositoryPostgres) Create(ctx context.Context, site *entity.Site) error {
	query := `
		INSERT INTO dcim_sites (
			id, name, slug, status, region_id, group_id, tenant_id, facility, 
			time_zone, physical_address, shipping_address, latitude, longitude,
			description, comments, created, updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, NOW(), NOW())
	`
	
	var regionID, groupID, tenantID interface{}
	if site.RegionID != nil {
		regionID = site.RegionID.String()
	} else {
		regionID = nil
	}
	if site.GroupID != nil {
		groupID = site.GroupID.String()
	} else {
		groupID = nil
	}
	if site.TenantID != nil {
		tenantID = site.TenantID.String()
	} else {
		tenantID = nil
	}
	
	var facility, timeZone, physicalAddr, shippingAddr interface{}
	if site.Facility != "" {
		facility = string(site.Facility)
	}
	if site.TimeZone != "" {
		timeZone = string(site.TimeZone)
	}
	if site.PhysicalAddress != "" {
		physicalAddr = string(site.PhysicalAddress)
	}
	if site.ShippingAddress != "" {
		shippingAddr = string(site.ShippingAddress)
	}
	
	var latitude, longitude interface{}
	if site.Latitude != nil {
		latitude = *site.Latitude
	}
	if site.Longitude != nil {
		longitude = *site.Longitude
	}
	
	var description, comments interface{}
	if site.Description != "" {
		description = string(site.Description)
	}
	if site.Comments != "" {
		comments = string(site.Comments)
	}
	
	_, err := r.db.ExecContext(ctx, query,
		site.ID.String(), site.Name, site.Slug, site.Status,
		regionID, groupID, tenantID, facility,
		timeZone, physicalAddr, shippingAddr,
		latitude, longitude,
		description, comments,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create site: %w", err)
	}
	
	return nil
}

// Update обновляет существующий сайт
func (r *SiteRepositoryPostgres) Update(ctx context.Context, site *entity.Site) error {
	query := `
		UPDATE dcim_sites
		SET name = $2, slug = $3, status = $4, region_id = $5, group_id = $6, 
		    tenant_id = $7, facility = $8, time_zone = $9, physical_address = $10,
		    shipping_address = $11, latitude = $12, longitude = $13,
		    description = $14, comments = $15, updated = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	
	// Аналогично Create, обработка опциональных полей
	var regionID, groupID, tenantID interface{}
	if site.RegionID != nil {
		regionID = site.RegionID.String()
	}
	if site.GroupID != nil {
		groupID = site.GroupID.String()
	}
	if site.TenantID != nil {
		tenantID = site.TenantID.String()
	}
	
	var facility, timeZone, physicalAddr, shippingAddr interface{}
	if site.Facility != "" {
		facility = string(site.Facility)
	}
	if site.TimeZone != "" {
		timeZone = string(site.TimeZone)
	}
	if site.PhysicalAddress != "" {
		physicalAddr = string(site.PhysicalAddress)
	}
	if site.ShippingAddress != "" {
		shippingAddr = string(site.ShippingAddress)
	}
	
	var latitude, longitude interface{}
	if site.Latitude != nil {
		latitude = *site.Latitude
	}
	if site.Longitude != nil {
		longitude = *site.Longitude
	}
	
	var description, comments interface{}
	if site.Description != "" {
		description = string(site.Description)
	}
	if site.Comments != "" {
		comments = string(site.Comments)
	}
	
	result, err := r.db.ExecContext(ctx, query,
		site.ID.String(), site.Name, site.Slug, site.Status,
		regionID, groupID, tenantID, facility,
		timeZone, physicalAddr, shippingAddr,
		latitude, longitude,
		description, comments,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update site: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return types.ErrNotFound
	}
	
	return nil
}

// Delete удаляет сайт (мягкое удаление)
func (r *SiteRepositoryPostgres) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE dcim_sites
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return types.ErrNotFound
	}
	
	return nil
}

// Exists проверяет существование сайта
func (r *SiteRepositoryPostgres) Exists(ctx context.Context, id string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM dcim_sites
			WHERE id = $1 AND deleted_at IS NULL
		)
	`
	
	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check site existence: %w", err)
	}
	
	return exists, nil
}
