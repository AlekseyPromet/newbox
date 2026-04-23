package service

import (
	"context"
	"fmt"
	"path/filepath"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/domain/core/repository"

	ntypes "netbox_go/pkg/types"

	"github.com/google/uuid"
)

// CoreService handles business logic for core components
type CoreService struct {
	dataSourceRepo     repository.DataSourceRepository
	configRevisionRepo repository.ConfigRevisionRepository
	managedFileRepo    repository.DataFileRepository
	storage            repository.FileStorage
}

func NewCoreService(
	dsRepo repository.DataSourceRepository,
	crRepo repository.ConfigRevisionRepository,
	mfRepo repository.DataFileRepository,
	storage repository.FileStorage,
) *CoreService {
	return &CoreService{
		dataSourceRepo:     dsRepo,
		configRevisionRepo: crRepo,
		managedFileRepo:    mfRepo,
		storage:            storage,
	}
}

// --- DataSource Logic ---

func (s *CoreService) BulkUpdateDataSources(ctx context.Context, ids []int64, params repository.DataSourceBulkUpdateParams) error {
	if len(ids) == 0 {
		return nil
	}

	if !params.HasChanges() {
		return nil
	}

	if err := s.dataSourceRepo.BulkUpdate(ctx, ids, params); err != nil {
		return fmt.Errorf("failed to bulk update data sources (ids: %v): %w", ids, err)
	}

	return nil
}

func (s *CoreService) ImportDataSources(ctx context.Context, data []entity.DataSource) error {
	if len(data) == 0 {
		return nil
	}

	// Validate all records before attempting bulk insertion to ensure atomicity
	// and provide clear feedback on which record failed.
	for i := range data {
		if err := data[i].Validate(); err != nil {
			return fmt.Errorf("validation failed for data source at index %d: %w", i, err)
		}
	}

	if err := s.dataSourceRepo.BulkCreate(ctx, data); err != nil {
		return fmt.Errorf("failed to bulk import data sources: %w", err)
	}

	return nil
}

// --- ConfigRevision Logic ---

func (s *CoreService) CreateConfigRevision(ctx context.Context, revision *entity.ConfigRevision) error {
	// Logic from ConfigRevisionForm.save()
	// Ensure data is rendered to JSON
	return s.configRevisionRepo.Create(ctx, revision)
}

// --- ManagedFile Logic ---

type ManagedFileCreateParams struct {
	DataSourceID    ntypes.ID
	DataFile        *string
	AutoSyncEnabled bool
	UploadFile      []byte
	FileName        string
}

func (s *CoreService) CreateManagedFile(ctx context.Context, params ManagedFileCreateParams) error {
	// Logic from ManagedFileForm.clean()
	if params.UploadFile != nil && params.DataFile != nil {
		return fmt.Errorf("cannot upload a file and sync from an existing file")
	}
	if params.UploadFile == nil && params.DataFile == nil {
		return fmt.Errorf("must upload a file or select a data file to sync")
	}

	var finalPath string

	if params.UploadFile != nil {
		// Sanitize filename to prevent path traversal
		safeFileName := sanitizeFileName(params.FileName)

		// Save file to storage
		path, err := s.storage.Save(ctx, safeFileName, params.UploadFile)
		if err != nil {
			return fmt.Errorf("failed to save uploaded file: %w", err)
		}
		finalPath = path
	} else if params.DataFile != nil {
		finalPath = *params.DataFile
	}

	// Persist the record to the database
	err := s.managedFileRepo.Create(ctx, &entity.DataFile{
		SourceID: params.DataSourceID,
		Path:     finalPath,
	})

	if err != nil {
		// Cleanup: If DB fails, remove the uploaded file to prevent orphans
		if params.UploadFile != nil {
			_ = s.storage.Delete(ctx, finalPath)
		}
		return fmt.Errorf("failed to create managed file record: %w", err)
	}

	return nil
}

func sanitizeFileName(name string) string {
	if name == "" {
		return uuid.NewString()
	}

	// filepath.Base returns the last element of path.
	// It handles trailing separators and returns "." if the path is empty.
	name = filepath.Base(name)

	// Prevent directory traversal and hidden files by removing leading dots.
	// We loop because a filename like ".../.hidden" could still result in ".hidden"
	for len(name) > 0 && name[0] == '.' {
		name = name[1:]
	}

	if name == "" {
		return uuid.NewString()
	}

	return name
}
