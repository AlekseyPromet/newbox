package graphql

import (
	"context"
	"fmt"
	"netbox_go/internal/domain/core/services"
	"netbox_go/internal/repository"
	"time"
)

// Resolver handles GraphQL queries
type Resolver struct {
	coreService services.CoreService
}

func NewResolver(coreService services.CoreService) *Resolver {
	return &Resolver{
		coreService: coreService,
	}
}

// DataFile resolver
func (r *Resolver) DataFile(ctx context.Context, args struct {
	ID string
}) (*DataFileType, error) {
	df, err := r.coreService.GetDataFile(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if df == nil {
		return nil, nil
	}

	return &DataFileType{
		ID:          df.ID.String(),
		Path:        df.Path,
		Size:        df.Size,
		Hash:        df.Hash,
		Created:     df.Created.Format(time.RFC3339),
		LastUpdated: df.Updated.Format(time.RFC3339),
		Source:      nil,
	}, nil
}

// DataFileList resolver
func (r *Resolver) DataFileList(ctx context.Context, args struct {
	Filter *DataFileFilter
	Limit  *int
	Offset *int
}) ([]*DataFileType, error) {
	limit := 0
	if args.Limit != nil {
		limit = *args.Limit
	}
	offset := 0
	if args.Offset != nil {
		offset = *args.Offset
	}

	var filter repository.DataFileFilter
	if args.Filter != nil {
		if args.Filter.SourceID != nil {
			filter.SourceID = args.Filter.SourceID
		}
		if args.Filter.Path != nil && args.Filter.Path.Exact != nil {
			filter.Path = args.Filter.Path.Exact
		}
	}
	dfs, total, err := r.coreService.ListDataFiles(ctx, filter, limit, offset)
	_ = total
	if err != nil {
		return nil, err
	}

	var result []*DataFileType
	for _, df := range dfs {
		result = append(result, &DataFileType{
			ID:          df.ID.String(),
			Path:        df.Path,
			Size:        df.Size,
			Hash:        df.Hash,
			Created:     df.Created.Format(time.RFC3339),
			LastUpdated: df.Updated.Format(time.RFC3339),
			Source:      nil,
		})
	}
	return result, nil
}

// DataSource resolver
func (r *Resolver) DataSource(ctx context.Context, args struct {
	ID string
}) (*DataSourceType, error) {
	ds, err := r.coreService.GetDataSource(ctx, args.ID)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return nil, nil
	}

	return &DataSourceType{
		ID:          ds.ID.String(),
		Name:        ds.Name,
		Type:        ds.Type,
		SourceURL:   ds.SourceURL,
		Status:      string(ds.Status),
		Enabled:     ds.Enabled,
		IgnoreRules: fmt.Sprintf("%v", ds.IgnoreRules),
		Parameters:  ds.Parameters,
		LastSynced:  "",
		DataFiles:   nil,
	}, nil
}

// DataSourceList resolver
func (r *Resolver) DataSourceList(ctx context.Context, args struct {
	Filter *DataSourceFilter
	Limit  *int
	Offset *int
}) ([]*DataSourceType, error) {
	limit := 0
	if args.Limit != nil {
		limit = *args.Limit
	}
	offset := 0
	if args.Offset != nil {
		offset = *args.Offset
	}

	var filter repository.DataSourceFilter
	if args.Filter != nil {
		if args.Filter.Name != nil && args.Filter.Name.Exact != nil {
			filter.Name = args.Filter.Name.Exact
		}
		if args.Filter.Type != nil && args.Filter.Type.Exact != nil {
			filter.Type = args.Filter.Type.Exact
		}
		if args.Filter.Status != nil {
			filter.Status = args.Filter.Status
		}
		if args.Filter.Enabled != nil {
			filter.Enabled = args.Filter.Enabled
		}
	}
	dss, total, err := r.coreService.ListDataSources(ctx, filter, limit, offset)
	_ = total
	if err != nil {
		return nil, err
	}

	var result []*DataSourceType
	for _, ds := range dss {
		result = append(result, &DataSourceType{
			ID:          ds.ID.String(),
			Name:        ds.Name,
			Type:        ds.Type,
			SourceURL:   ds.SourceURL,
			Status:      string(ds.Status),
			Enabled:     ds.Enabled,
			IgnoreRules: fmt.Sprintf("%v", ds.IgnoreRules),
			Parameters:  ds.Parameters,
			LastSynced:  "",
			DataFiles:   nil,
		})
	}
	return result, nil
}
