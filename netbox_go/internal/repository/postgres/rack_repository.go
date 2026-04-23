package postgres

import (
	"context"
	"database/sql"
	"fmt"

	dcim_entity "netbox_go/internal/domain/dcim/entity"
	"netbox_go/internal/domain/dcim/enum"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
)

// RackRepositoryPostgres реализует интерфейс RackRepository для PostgreSQL
type RackRepositoryPostgres struct {
	db *sql.DB
}

// NewRackRepositoryPostgres создает новый экземпляр репозитория стоек
func NewRackRepositoryPostgres(db *sql.DB) *RackRepositoryPostgres {
	return &RackRepositoryPostgres{db: db}
}

// GetByID получает стойку по ID
func (r *RackRepositoryPostgres) GetByID(ctx context.Context, id string) (*dcim_entity.Rack, error) {
	query := `
		SELECT id, name, facility_id, site_id, location_id, tenant_id, status, role_id,
		       rack_type_id, form_factor, width, serial, asset_tag, airflow,
		       u_height, starting_unit, desc_units,
		       outer_width, outer_height, outer_depth, outer_unit,
		       mounting_depth, weight, max_weight, weight_unit,
		       description, comments, created, updated
		FROM dcim_racks
		WHERE id = $1 AND deleted_at IS NULL
	`

	var rack dcim_entity.Rack
	var facilityID, locationID, tenantID, roleID, rackTypeID sql.NullString
	var serial, formFactor, airflow, assetTag sql.NullString
	var outerWidth, outerHeight, outerDepth, mountingDepth sql.NullInt64
	var outerUnit, weightUnit sql.NullString
	var weight, maxWeight sql.NullInt64
	var description, comments sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rack.ID, &rack.Name, &facilityID, &rack.SiteID, &locationID, &tenantID,
		&rack.Status, &roleID, &rackTypeID, &formFactor, &rack.Width, &serial,
		&assetTag, &airflow, &rack.UHeight, &rack.StartingUnit, &rack.DescUnits,
		&outerWidth, &outerHeight, &outerDepth, &outerUnit, &mountingDepth,
		&weight, &maxWeight, &weightUnit, &description, &comments,
		&rack.Created, &rack.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get rack by ID: %w", err)
	}

	// Заполнение опциональных полей
	if facilityID.Valid {
		rack.FacilityID = &facilityID.String
	}
	if locationID.Valid {
		locID, _ := types.ParseID(locationID.String)
		rack.LocationID = &locID
	}
	if tenantID.Valid {
		tenID, _ := types.ParseID(tenantID.String)
		rack.TenantID = &tenID
	}
	if roleID.Valid {
		role, _ := types.ParseID(roleID.String)
		rack.RoleID = &role
	}
	if rackTypeID.Valid {
		rType, _ := types.ParseID(rackTypeID.String)
		rack.RackTypeID = &rType
	}
	if formFactor.Valid {
		ff := enum.RackType(formFactor.String)
		rack.FormFactor = &ff
	}
	if airflow.Valid {
		rack.Airflow = &airflow.String
	}
	if assetTag.Valid {
		rack.AssetTag = &assetTag.String
	}
	if outerWidth.Valid {
		v := int16(outerWidth.Int64)
		rack.OuterWidth = &v
	}
	if outerHeight.Valid {
		v := int16(outerHeight.Int64)
		rack.OuterHeight = &v
	}
	if outerDepth.Valid {
		v := int16(outerDepth.Int64)
		rack.OuterDepth = &v
	}
	if outerUnit.Valid {
		ou := enum.RackDimensionUnit(outerUnit.String)
		rack.OuterUnit = &ou
	}
	if mountingDepth.Valid {
		v := int16(mountingDepth.Int64)
		rack.MountingDepth = &v
	}
	if weight.Valid {
		v := int32(weight.Int64)
		rack.Weight = &v
	}
	if maxWeight.Valid {
		v := int32(maxWeight.Int64)
		rack.MaxWeight = &v
	}
	if weightUnit.Valid {
		rack.WeightUnit = &weightUnit.String
	}
	if description.Valid {
		rack.Description = types.Description(description.String)
	}
	if comments.Valid {
		rack.Comments = types.Comments(comments.String)
	}

	return &rack, nil
}

// List получает список стоек с фильтрацией и пагинацией
func (r *RackRepositoryPostgres) List(ctx context.Context, filter repository.RackFilter) ([]*dcim_entity.Rack, int64, error) {
	query := `
		SELECT id, name, facility_id, site_id, location_id, tenant_id, status, role_id,
		       rack_type_id, form_factor, width, serial, asset_tag, airflow,
		       u_height, starting_unit, desc_units,
		       outer_width, outer_height, outer_depth, outer_unit,
		       mounting_depth, weight, max_weight, weight_unit,
		       description, comments, created, updated
		FROM dcim_racks
		WHERE deleted_at IS NULL
	`

	args := []interface{}{}
	argIndex := 1

	if filter.SiteID != nil {
		query += fmt.Sprintf(" AND site_id = $%d", argIndex)
		args = append(args, *filter.SiteID)
		argIndex++
	}
	if filter.LocationID != nil {
		query += fmt.Sprintf(" AND location_id = $%d", argIndex)
		args = append(args, *filter.LocationID)
		argIndex++
	}
	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}
	if filter.TypeID != nil {
		query += fmt.Sprintf(" AND rack_type_id = $%d", argIndex)
		args = append(args, *filter.TypeID)
		argIndex++
	}

	// Получение общего количества
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM (%s) AS count_query", query)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count racks: %w", err)
	}

	// Пагинация
	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	offset := filter.Offset

	query += fmt.Sprintf(" ORDER BY name ASC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list racks: %w", err)
	}
	defer rows.Close()

	racks := make([]*dcim_entity.Rack, 0)
	for rows.Next() {
		var rack dcim_entity.Rack
		var facilityID, locationID, tenantID, roleID, rackTypeID sql.NullString
		var formFactor, airflow, assetTag sql.NullString
		var outerWidth, outerHeight, outerDepth, mountingDepth sql.NullInt64
		var outerUnit, weightUnit sql.NullString
		var weight, maxWeight sql.NullInt64
		var description, comments, serial sql.NullString

		err := rows.Scan(
			&rack.ID, &rack.Name, &facilityID, &rack.SiteID, &locationID, &tenantID,
			&rack.Status, &roleID, &rackTypeID, &formFactor, &rack.Width, &serial,
			&assetTag, &airflow, &rack.UHeight, &rack.StartingUnit, &rack.DescUnits,
			&outerWidth, &outerHeight, &outerDepth, &outerUnit, &mountingDepth,
			&weight, &maxWeight, &weightUnit, &description, &comments,
			&rack.Created, &rack.Updated,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan rack: %w", err)
		}

		// Заполнение опциональных полей
		if facilityID.Valid {
			rack.FacilityID = &facilityID.String
		}
		if locationID.Valid {
			locID, _ := types.ParseID(locationID.String)
			rack.LocationID = &locID
		}
		if tenantID.Valid {
			tenID, _ := types.ParseID(tenantID.String)
			rack.TenantID = &tenID
		}
		if roleID.Valid {
			role, _ := types.ParseID(roleID.String)
			rack.RoleID = &role
		}
		if rackTypeID.Valid {
			rType, _ := types.ParseID(rackTypeID.String)
			rack.RackTypeID = &rType
		}
		if formFactor.Valid {
			ff := enum.RackType(formFactor.String)
			rack.FormFactor = &ff
		}
		if airflow.Valid {
			rack.Airflow = &airflow.String
		}
		if assetTag.Valid {
			rack.AssetTag = &assetTag.String
		}
		if outerWidth.Valid {
			v := int16(outerWidth.Int64)
			rack.OuterWidth = &v
		}
		if outerHeight.Valid {
			v := int16(outerHeight.Int64)
			rack.OuterHeight = &v
		}
		if outerDepth.Valid {
			v := int16(outerDepth.Int64)
			rack.OuterDepth = &v
		}
		if outerUnit.Valid {
			ou := enum.RackDimensionUnit(outerUnit.String)
			rack.OuterUnit = &ou
		}
		if mountingDepth.Valid {
			v := int16(mountingDepth.Int64)
			rack.MountingDepth = &v
		}
		if weight.Valid {
			v := int32(weight.Int64)
			rack.Weight = &v
		}
		if maxWeight.Valid {
			v := int32(maxWeight.Int64)
			rack.MaxWeight = &v
		}
		if weightUnit.Valid {
			rack.WeightUnit = &weightUnit.String
		}
		if description.Valid {
			rack.Description = types.Description(description.String)
		}
		if comments.Valid {
			rack.Comments = types.Comments(comments.String)
		}

		racks = append(racks, &rack)
	}

	return racks, total, nil
}

// Create создает новую стойку
func (r *RackRepositoryPostgres) Create(ctx context.Context, rack *dcim_entity.Rack) error {
	query := `
		INSERT INTO dcim_racks (
			id, name, facility_id, site_id, location_id, tenant_id, status, role_id,
			rack_type_id, form_factor, width, serial, asset_tag, airflow,
			u_height, starting_unit, desc_units,
			outer_width, outer_height, outer_depth, outer_unit,
			mounting_depth, weight, max_weight, weight_unit,
			description, comments, created, updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, NOW(), NOW())
	`

	var facilityID, locationID, tenantID, roleID, rackTypeID interface{}
	if rack.FacilityID != nil {
		facilityID = *rack.FacilityID
	}
	if rack.LocationID != nil {
		locationID = rack.LocationID.String()
	}
	if rack.TenantID != nil {
		tenantID = rack.TenantID.String()
	}
	if rack.RoleID != nil {
		roleID = rack.RoleID.String()
	}
	if rack.RackTypeID != nil {
		rackTypeID = rack.RackTypeID.String()
	}

	var formFactor, airflow, assetTag interface{}
	if rack.FormFactor != nil {
		formFactor = string(*rack.FormFactor)
	}
	if rack.Airflow != nil {
		airflow = *rack.Airflow
	}
	if rack.AssetTag != nil {
		assetTag = *rack.AssetTag
	}

	var outerWidth, outerHeight, outerDepth, mountingDepth interface{}
	if rack.OuterWidth != nil {
		outerWidth = *rack.OuterWidth
	}
	if rack.OuterHeight != nil {
		outerHeight = *rack.OuterHeight
	}
	if rack.OuterDepth != nil {
		outerDepth = *rack.OuterDepth
	}
	if rack.MountingDepth != nil {
		mountingDepth = *rack.MountingDepth
	}

	var outerUnit, weightUnit interface{}
	if rack.OuterUnit != nil {
		outerUnit = string(*rack.OuterUnit)
	}
	if rack.WeightUnit != nil {
		weightUnit = *rack.WeightUnit
	}

	var weight, maxWeight interface{}
	if rack.Weight != nil {
		weight = *rack.Weight
	}
	if rack.MaxWeight != nil {
		maxWeight = *rack.MaxWeight
	}

	var description, comments interface{}
	if rack.Description != "" {
		description = string(rack.Description)
	}
	if rack.Comments != "" {
		comments = string(rack.Comments)
	}

	_, err := r.db.ExecContext(ctx, query,
		rack.ID.String(), rack.Name, facilityID, rack.SiteID.String(), locationID, tenantID,
		rack.Status, roleID, rackTypeID, formFactor, rack.Width, rack.Serial, assetTag, airflow,
		rack.UHeight, rack.StartingUnit, rack.DescUnits,
		outerWidth, outerHeight, outerDepth, outerUnit, mountingDepth,
		weight, maxWeight, weightUnit, description, comments,
	)

	if err != nil {
		return fmt.Errorf("failed to create rack: %w", err)
	}

	return nil
}

// Update обновляет существующую стойку
func (r *RackRepositoryPostgres) Update(ctx context.Context, rack *dcim_entity.Rack) error {
	query := `
		UPDATE dcim_racks
		SET name = $2, facility_id = $3, site_id = $4, location_id = $5, tenant_id = $6,
		    status = $7, role_id = $8, rack_type_id = $9, form_factor = $10, width = $11,
		    serial = $12, asset_tag = $13, airflow = $14, u_height = $15, starting_unit = $16,
		    desc_units = $17, outer_width = $18, outer_height = $19, outer_depth = $20,
		    outer_unit = $21, mounting_depth = $22, weight = $23, max_weight = $24,
		    weight_unit = $25, description = $26, comments = $27, updated = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	var facilityID, locationID, tenantID, roleID, rackTypeID interface{}
	if rack.FacilityID != nil {
		facilityID = *rack.FacilityID
	}
	if rack.LocationID != nil {
		locationID = rack.LocationID.String()
	}
	if rack.TenantID != nil {
		tenantID = rack.TenantID.String()
	}
	if rack.RoleID != nil {
		roleID = rack.RoleID.String()
	}
	if rack.RackTypeID != nil {
		rackTypeID = rack.RackTypeID.String()
	}

	var formFactor, airflow, assetTag interface{}
	if rack.FormFactor != nil {
		formFactor = string(*rack.FormFactor)
	}
	if rack.Airflow != nil {
		airflow = *rack.Airflow
	}
	if rack.AssetTag != nil {
		assetTag = *rack.AssetTag
	}

	var outerWidth, outerHeight, outerDepth, mountingDepth interface{}
	if rack.OuterWidth != nil {
		outerWidth = *rack.OuterWidth
	}
	if rack.OuterHeight != nil {
		outerHeight = *rack.OuterHeight
	}
	if rack.OuterDepth != nil {
		outerDepth = *rack.OuterDepth
	}
	if rack.MountingDepth != nil {
		mountingDepth = *rack.MountingDepth
	}

	var outerUnit, weightUnit interface{}
	if rack.OuterUnit != nil {
		outerUnit = string(*rack.OuterUnit)
	}
	if rack.WeightUnit != nil {
		weightUnit = *rack.WeightUnit
	}

	var weight, maxWeight interface{}
	if rack.Weight != nil {
		weight = *rack.Weight
	}
	if rack.MaxWeight != nil {
		maxWeight = *rack.MaxWeight
	}

	var description, comments interface{}
	if rack.Description != "" {
		description = string(rack.Description)
	}
	if rack.Comments != "" {
		comments = string(rack.Comments)
	}

	result, err := r.db.ExecContext(ctx, query,
		rack.ID.String(), rack.Name, facilityID, rack.SiteID.String(), locationID, tenantID,
		rack.Status, roleID, rackTypeID, formFactor, rack.Width, rack.Serial, assetTag, airflow,
		rack.UHeight, rack.StartingUnit, rack.DescUnits,
		outerWidth, outerHeight, outerDepth, outerUnit, mountingDepth,
		weight, maxWeight, weightUnit, description, comments,
	)

	if err != nil {
		return fmt.Errorf("failed to update rack: %w", err)
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

// Delete удаляет стойку (мягкое удаление)
func (r *RackRepositoryPostgres) Delete(ctx context.Context, id string) error {
	query := `
		UPDATE dcim_racks
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete rack: %w", err)
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

// Exists проверяет существование стойки
func (r *RackRepositoryPostgres) Exists(ctx context.Context, id string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM dcim_racks
			WHERE id = $1 AND deleted_at IS NULL
		)
	`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check rack existence: %w", err)
	}

	return exists, nil
}
