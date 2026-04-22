package service

import (
	"context"
	"encoding/json"
	"fmt"

	"netbox_go/internal/domain/core/model"
	"netbox_go/internal/domain/core/repository"
)

// CoreService handles business logic for core components
type CoreService struct {
	dataSourceRepo    repository.DataSourceRepository
	configRevisionRepo repository.ConfigRevisionRepository
	managedFileRepo    repository.ManagedFileRepository
}

func NewCoreService(
	dsRepo repository.DataSourceRepository,
	crRepo repository.ConfigRevisionRepository,
	mfRepo repository.ManagedFileRepository,
) *CoreService {
	return &CoreService{
		dataSourceRepo:    dsRepo,
		configRevisionRepo: crRepo,
		managedFileRepo:    mfRepo,
	}
}

// --- DataSource Logic ---

type DataSourceBulkUpdateParams struct {
	Type          *string
	Enabled       *bool
	Description   *string
	SyncInterval  *string
	Parameters    *json.RawMessage
	IgnoreRules   *string
	Comments      *string
}

func (s *CoreService) BulkUpdateDataSources(ctx context.Context, ids []int64, params DataSourceBulkUpdateParams) error {
	// In a real implementation, this would use a dynamic query builder
	// based on the provided params.
	// For now, we'll implement the logic as a series of updates or a single dynamic SQL call.
	return s.dataSourceRepo.BulkUpdate(ctx, ids, params)
}

func (s *CoreService) ImportDataSources(ctx context.Context, data []model.DataSource) error {
	// Bulk import logic from netbox/core/forms/bulk_import.py
	// This would typically involve validation of each record and a bulk insert.
	return s.dataSourceRepo.BulkCreate(ctx, data)
}

// --- ConfigRevision Logic ---

func (s *CoreService) CreateConfigRevision(ctx context.Context, revision *model.ConfigRevision) error {
	// Logic from ConfigRevisionForm.save()
	// Ensure data is rendered to JSON
	return s.configRevisionRepo.Create(ctx, revision)
}

// --- ManagedFile Logic ---

type ManagedFileCreateParams struct {
	DataSourceID   int64
	DataFile       *string
	AutoSyncEnabled bool
	UploadFile     []byte
	FileName       string
}

func (s *CoreService) CreateManagedFile(ctx context.Context, params ManagedFileCreateParams) error {
	// Logic from ManagedFileForm.clean()
	if params.UploadFile != nil && params.DataFile != nil {
		return fmt.Errorf("cannot upload a file and sync from an existing file")
	}
	if params.UploadFile == nil && params.DataFile == nil {
		return fmt.Errorf("must upload a file or select a data file to sync")
	}

	managedFile := &model.ManagedFile{
		DataSourceID:    params.DataSourceID,
		DataFile:        params.DataFile,
		AutoSyncEnabled: params.AutoSyncEnabled,
	}

	// Handle file upload logic from ManagedFileForm.save()
	if params.UploadFile != nil {
		// In a real implementation, we would save the file to the filesystem
		// and set the file_path on the model.
		// For now, we assume the repository handles the storage or the path is pre-calculated.
		managedFile.FilePath = params.FileName
	}

	return s.managedFileRepo.Create(ctx, managedFile)
}
