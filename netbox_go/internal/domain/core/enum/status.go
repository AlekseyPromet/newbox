// Package enum содержит перечисления для домена Core
package enum

import "github.com/AlekseyPromet/netbox_go/pkg/types"

// DataSourceStatus определяет статусы источника данных
// См. netbox/core/choices.py DataSourceStatusChoices
const (
	DataSourceStatusNew      types.Status = "new"
	DataSourceStatusQueued   types.Status = "queued"
	DataSourceStatusSyncing  types.Status = "syncing"
	DataSourceStatusCompleted types.Status = "completed"
	DataSourceStatusFailed   types.Status = "failed"
)

// ValidDataSourceStatuses список допустимых статусов
var ValidDataSourceStatuses = map[types.Status]struct{}{
	DataSourceStatusNew:       {},
	DataSourceStatusQueued:    {},
	DataSourceStatusSyncing:   {},
	DataSourceStatusCompleted: {},
	DataSourceStatusFailed:    {},
}

// ValidateDataSourceStatus проверяет корректность статуса источника данных
func ValidateDataSourceStatus(status types.Status) error {
	if _, ok := ValidDataSourceStatuses[status]; !ok {
		return types.ErrInvalidStatus
	}
	return nil
}

// JobStatus определяет статусы фоновой задачи
// См. netbox/core/choices.py JobStatusChoices
const (
	JobStatusPending   types.Status = "pending"
	JobStatusScheduled types.Status = "scheduled"
	JobStatusRunning   types.Status = "running"
	JobStatusCompleted types.Status = "completed"
	JobStatusErrored   types.Status = "errored"
	JobStatusFailed    types.Status = "failed"
)

// ValidJobStatuses список допустимых статусов задач
var ValidJobStatuses = map[types.Status]struct{}{
	JobStatusPending:   {},
	JobStatusScheduled: {},
	JobStatusRunning:   {},
	JobStatusCompleted: {},
	JobStatusErrored:   {},
	JobStatusFailed:    {},
}

// ValidateJobStatus проверяет корректность статуса задачи
func ValidateJobStatus(status types.Status) error {
	if _, ok := ValidJobStatuses[status]; !ok {
		return types.ErrInvalidStatus
	}
	return nil
}

// JobInterval задает интервал выполнения задач (минуты)
// См. netbox/core/choices.py JobIntervalChoices
const (
	JobIntervalMinutely int = 1
	JobIntervalHourly   int = 60
	JobIntervalDaily    int = 60 * 24
	JobIntervalWeekly   int = 60 * 24 * 7
	JobInterval12Hours  int = 60 * 12
	JobInterval30Days   int = 60 * 24 * 30
)

// ValidJobIntervals список допустимых интервалов (минуты)
var ValidJobIntervals = map[int]struct{}{
	JobIntervalMinutely: {},
	JobIntervalHourly:   {},
	JobInterval12Hours:  {},
	JobIntervalDaily:    {},
	JobIntervalWeekly:   {},
	JobInterval30Days:   {},
}

// ValidateJobInterval проверяет корректность интервала
func ValidateJobInterval(interval int) error {
	if interval == 0 {
		return nil
	}
	if _, ok := ValidJobIntervals[interval]; !ok {
		return types.ErrInvalidStatus
	}
	return nil
}

// ObjectChangeAction определяет действие в журнале изменений
const (
	ObjectChangeActionCreate types.Status = "create"
	ObjectChangeActionUpdate types.Status = "update"
	ObjectChangeActionDelete types.Status = "delete"
)

// ValidObjectChangeActions список допустимых действий
var ValidObjectChangeActions = map[types.Status]struct{}{
	ObjectChangeActionCreate: {},
	ObjectChangeActionUpdate: {},
	ObjectChangeActionDelete: {},
}

// ValidateObjectChangeAction проверяет корректность действия
func ValidateObjectChangeAction(action types.Status) error {
	if _, ok := ValidObjectChangeActions[action]; !ok {
		return types.ErrInvalidStatus
	}
	return nil
}
