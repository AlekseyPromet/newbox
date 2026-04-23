package repository

import (
	"time"

	"netbox_go/pkg/types"
)

// DataSourceFilter defines the filtering criteria for DataSources
type DataSourceFilter struct {
	ID               *types.ID
	Name             *string
	Enabled          *bool
	Description      *string
	SourceURL        *string
	Status           []types.Status
	Type             *string
	SearchQuery      *string
	LastSyncedBefore *time.Time
	LastSyncedAfter  *time.Time
}

// DataFileFilter defines the filtering criteria for DataFiles
type DataFileFilter struct {
	ID                *types.ID
	Path              *string
	SourceID          *types.ID
	SourceName        *string
	LastUpdatedBefore *time.Time
	LastUpdatedAfter  *time.Time
	Size              *int64
	Hash              *string
	SearchQuery       *string
}

// JobFilter defines the filtering criteria for Jobs
type JobFilter struct {
	ID              *types.ID
	ObjectType      *string
	ObjectID        *types.ID
	Name            *string
	Status          []types.Status
	Interval        *string
	User            *string
	QueueName       *string
	CreatedBefore   *time.Time
	CreatedAfter    *time.Time
	ScheduledBefore *time.Time
	ScheduledAfter  *time.Time
	StartedBefore   *time.Time
	StartedAfter    *time.Time
	CompletedBefore *time.Time
	CompletedAfter  *time.Time
	SearchQuery     *string
}

// ObjectTypeFilter defines the filtering criteria for ObjectTypes
type ObjectTypeFilter struct {
	ID          *types.ID
	AppLabel    *string
	Model       *string
	Public      *bool
	Features    *string
	SearchQuery *string
}

// ObjectChangeFilter defines the filtering criteria for ObjectChanges
type ObjectChangeFilter struct {
	ID                *types.ID
	User              *string
	UserID            *types.ID
	RequestID         *string
	Action            *string
	ChangedObjectType *string
	ChangedObjectID   *string
	RelatedObjectType *string
	RelatedObjectID   *string
	TimeBefore        *time.Time
	TimeAfter         *time.Time
	SearchQuery       *string
}

// ConfigRevisionFilter defines the filtering criteria for ConfigRevisions
type ConfigRevisionFilter struct {
	ID            *types.ID
	CreatedBefore *time.Time
	CreatedAfter  *time.Time
	Comment       *string
	SearchQuery   *string
}
