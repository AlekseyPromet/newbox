// Package entity содержит сущности домена Core
package entity

import (
	"encoding/json"
	"time"

	"netbox_go/internal/domain/core/enum"
	"netbox_go/pkg/types"
)

// ConfigRevision представляет сохранённую ревизию конфигурации NetBox.
type ConfigRevision struct {
	ID          types.ID        `json:"id"`
	Created     time.Time       `json:"created"`
	Active      bool            `json:"active"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Comment     string          `json:"comment,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
}

// Validate проверяет корректность ревизии конфигурации.
func (cr *ConfigRevision) Validate() error {
	if cr.Data == nil {
		// пустой JSON допустим, но nil означает отсутствие данных
		cr.Data = json.RawMessage([]byte("{}"))
	}
	return nil
}

// ObjectType описывает тип объекта (аналог ContentType в Django).
type ObjectType struct {
	ID       types.ID  `json:"id"`
	AppLabel string    `json:"app_label"`
	Model    string    `json:"model"`
	Public   bool      `json:"public"`
	Features []string  `json:"features,omitempty"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
}

// Validate проверяет корректность ObjectType.
func (ot *ObjectType) Validate() error {
	if ot.AppLabel == "" || ot.Model == "" {
		return types.ErrValidationFailed
	}
	return nil
}

// ObjectChange фиксирует изменение объекта (change logging).
type ObjectChange struct {
	ID                types.ID        `json:"id"`
	Time              time.Time       `json:"time"`
	UserID            types.ID        `json:"user_id,omitempty"`
	RequestID         *string         `json:"request_id,omitempty"`
	Action            types.Status    `json:"action"`
	ChangedObjectType string          `json:"changed_object_type"`
	ChangedObjectID   string          `json:"changed_object_id"`
	ObjectRepr        string          `json:"object_repr"`
	ObjectData        json.RawMessage `json:"object_data,omitempty"`
	RelatedObjectType *string         `json:"related_object_type,omitempty"`
	RelatedObjectID   *string         `json:"related_object_id,omitempty"`
	RelatedObjectRepr *string         `json:"related_object_repr,omitempty"`
}

// Validate проверяет корректность ObjectChange.
func (oc *ObjectChange) Validate() error {
	if err := enum.ValidateObjectChangeAction(oc.Action); err != nil {
		return err
	}
	if oc.ChangedObjectType == "" || oc.ChangedObjectID == "" {
		return types.ErrValidationFailed
	}
	if oc.ObjectRepr == "" {
		return types.ErrValidationFailed
	}
	return nil
}

// DataSource описывает внешний источник данных.
type DataSource struct {
	ID           types.ID        `json:"id"`
	Name         string          `json:"name"`
	Type         string          `json:"type"`
	SourceURL    string          `json:"source_url"`
	Status       types.Status    `json:"status"`
	Enabled      bool            `json:"enabled"`
	SyncInterval int             `json:"sync_interval"` // минуты
	Description  string          `json:"description,omitempty"`
	IgnoreRules  []string        `json:"ignore_rules,omitempty"`
	Parameters   json.RawMessage `json:"parameters,omitempty"`
	LastSynced   *time.Time      `json:"last_synced,omitempty"`
	Created      time.Time       `json:"created"`
	Updated      time.Time       `json:"updated"`
}

// Validate проверяет корректность DataSource.
func (ds *DataSource) Validate() error {
	if ds.Name == "" || ds.Type == "" {
		return types.ErrValidationFailed
	}
	if err := enum.ValidateDataSourceStatus(ds.Status); err != nil {
		return err
	}
	if ds.SyncInterval < 0 {
		return types.ErrValidationFailed
	}
	return nil
}

// DataFile представляет файл, полученный из источника данных.
type DataFile struct {
	ID       types.ID        `json:"id"`
	SourceID types.ID        `json:"source_id"`
	Path     string          `json:"path"`
	FileType string          `json:"file_type"` // csv, yaml, json
	Size     int64           `json:"size"`
	Hash     string          `json:"hash"`
	Data     json.RawMessage `json:"data,omitempty"`
	Created  time.Time       `json:"created"`
	Updated  time.Time       `json:"updated"`
}

// Validate проверяет корректность DataFile.
func (df *DataFile) Validate() error {
	if df.SourceID.String() == "" || df.Path == "" {
		return types.ErrValidationFailed
	}
	if df.FileType == "" {
		return types.ErrValidationFailed
	}
	if df.FileType != "csv" && df.FileType != "yaml" && df.FileType != "json" {
		return types.ErrValidationFailed
	}
	return nil
}

// Job описывает фоновую задачу (jobs/tasks).
type Job struct {
	ID          types.ID        `json:"id"`
	ObjectType  string          `json:"object_type,omitempty"`
	ObjectID    types.ID        `json:"object_id,omitempty"`
	UserID      types.ID        `json:"user_id,omitempty"`
	Object      interface{}     `json:"object,omitempty"`
	Name        string          `json:"name"`
	Status      types.Status    `json:"status"`
	Interval    int             `json:"interval,omitempty"` // минуты
	ScheduledAt time.Time       `json:"scheduled_at,omitempty"`
	StartedAt   time.Time       `json:"started_at,omitempty"`
	CompletedAt time.Time       `json:"completed_at,omitempty"`
	QueueName   string          `json:"queue_name,omitempty"`
	JobID       string          `json:"job_id,omitempty"`
	Data        json.RawMessage `json:"data,omitempty"`
	Error       string          `json:"error,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность Job.
func (j *Job) Validate() error {
	if j.Name == "" {
		return types.ErrValidationFailed
	}
	if err := enum.ValidateJobStatus(j.Status); err != nil {
		return err
	}
	if err := enum.ValidateJobInterval(j.Interval); err != nil {
		return err
	}
	return nil
}

// AutoSyncRecord представляет запись об автоматической синхронизации объекта из файла данных.
type AutoSyncRecord struct {
	ID         types.ID  `json:"id"`
	DataFileID types.ID  `json:"datafile_id"`
	ObjectType string    `json:"object_type"` // формат: "app_label.model"
	ObjectID   string    `json:"object_id"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}

// Validate проверяет корректность AutoSyncRecord.
func (asr *AutoSyncRecord) Validate() error {
	if asr.DataFileID.String() == "" || asr.ObjectType == "" || asr.ObjectID == "" {
		return types.ErrValidationFailed
	}
	return nil
}

// ManagedFile представляет управляемый файл (скрипт, отчёт, фильтр).
type ManagedFile struct {
	ID          types.ID   `json:"id"`
	Created     time.Time  `json:"created"`
	LastUpdated *time.Time `json:"last_updated,omitempty"`
	FileRoot    string     `json:"file_root"` // scripts, reports, filters
	FilePath    string     `json:"file_path"`
	DataFileID  *types.ID  `json:"datafile_id,omitempty"`
}

// Validate проверяет корректность ManagedFile.
func (mf *ManagedFile) Validate() error {
	if mf.FileRoot == "" || mf.FilePath == "" {
		return types.ErrValidationFailed
	}
	// Проверка допустимых значений file_root
	validRoots := map[string]bool{
		"scripts": true,
		"reports": true,
		"filters": true,
	}
	if !validRoots[mf.FileRoot] {
		return types.ErrValidationFailed
	}
	return nil
}
