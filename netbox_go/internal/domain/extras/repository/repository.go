package repository

import (
	"context"

	"github.com/google/uuid"
	"netbox_go/internal/domain/extras/entity"
)

// ExtrasRepository defines the data access layer for the extras module.
type ExtrasRepository interface {
	// EventRule
	CreateEventRule(ctx context.Context, rule *entity.EventRule) error
	GetEventRule(ctx context.Context, id uuid.UUID) (*entity.EventRule, error)
	UpdateEventRule(ctx context.Context, rule *entity.EventRule) error
	DeleteEventRule(ctx context.Context, id uuid.UUID) error
	ListEventRules(ctx context.Context, limit, offset int) ([]*entity.EventRule, error)

	// Webhook
	CreateWebhook(ctx context.Context, webhook *entity.Webhook) error
	GetWebhook(ctx context.Context, id uuid.UUID) (*entity.Webhook, error)
	UpdateWebhook(ctx context.Context, webhook *entity.Webhook) error
	DeleteWebhook(ctx context.Context, id uuid.UUID) error
	ListWebhooks(ctx context.Context, limit, offset int) ([]*entity.Webhook, error)

	// CustomLink
	CreateCustomLink(ctx context.Context, link *entity.CustomLink) error
	GetCustomLink(ctx context.Context, id uuid.UUID) (*entity.CustomLink, error)
	UpdateCustomLink(ctx context.Context, link *entity.CustomLink) error
	DeleteCustomLink(ctx context.Context, id uuid.UUID) error
	ListCustomLinks(ctx context.Context, limit, offset int) ([]*entity.CustomLink, error)

	// ExportTemplate
	CreateExportTemplate(ctx context.Context, template *entity.ExportTemplate) error
	GetExportTemplate(ctx context.Context, id uuid.UUID) (*entity.ExportTemplate, error)
	UpdateExportTemplate(ctx context.Context, template *entity.ExportTemplate) error
	DeleteExportTemplate(ctx context.Context, id uuid.UUID) error
	ListExportTemplates(ctx context.Context, limit, offset int) ([]*entity.ExportTemplate, error)

	// SavedFilter
	CreateSavedFilter(ctx context.Context, filter *entity.SavedFilter) error
	GetSavedFilter(ctx context.Context, id uuid.UUID) (*entity.SavedFilter, error)
	UpdateSavedFilter(ctx context.Context, filter *entity.SavedFilter) error
	DeleteSavedFilter(ctx context.Context, id uuid.UUID) error
	ListSavedFilters(ctx context.Context, limit, offset int) ([]*entity.SavedFilter, error)

	// TableConfig
	CreateTableConfig(ctx context.Context, config *entity.TableConfig) error
	GetTableConfig(ctx context.Context, id uuid.UUID) (*entity.TableConfig, error)
	UpdateTableConfig(ctx context.Context, config *entity.TableConfig) error
	DeleteTableConfig(ctx context.Context, id uuid.UUID) error
	ListTableConfigs(ctx context.Context, limit, offset int) ([]*entity.TableConfig, error)

	// ImageAttachment
	CreateImageAttachment(ctx context.Context, img *entity.ImageAttachment) error
	GetImageAttachment(ctx context.Context, id uuid.UUID) (*entity.ImageAttachment, error)
	UpdateImageAttachment(ctx context.Context, img *entity.ImageAttachment) error
	DeleteImageAttachment(ctx context.Context, id uuid.UUID) error
	ListImageAttachments(ctx context.Context, objectType string, objectID int64) ([]*entity.ImageAttachment, error)

	// JournalEntry
	CreateJournalEntry(ctx context.Context, entry *entity.JournalEntry) error
	GetJournalEntry(ctx context.Context, id uuid.UUID) (*entity.JournalEntry, error)
	ListJournalEntries(ctx context.Context, objectType string, objectID int64) ([]*entity.JournalEntry, error)

	// Bookmark
	CreateBookmark(ctx context.Context, bookmark *entity.Bookmark) error
	DeleteBookmark(ctx context.Context, id uuid.UUID) error
	ListBookmarks(ctx context.Context, userID uuid.UUID) ([]*entity.Bookmark, error)
}
