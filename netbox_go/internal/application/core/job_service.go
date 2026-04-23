package core

import (
	"context"
	"fmt"
	"time"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/pkg/taskqueue"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
)

// JobService сервис для управления задачами
type JobService struct {
	jobRepo   repository.JobRepository
	taskQueue *taskqueue.EtcdQueue
}

// NewJobService создаёт новый сервис задач
func NewJobService(jobRepo repository.JobRepository, taskQueue *taskqueue.EtcdQueue) *JobService {
	return &JobService{
		jobRepo:   jobRepo,
		taskQueue: taskQueue,
	}
}

// CreateJobParams параметры создания задачи
type CreateJobParams struct {
	ObjectType  string
	ObjectID    string
	Name        string
	Description string
	Interval    int
	ScheduledAt *time.Time
	Data        map[string]interface{}
	Priority    int
}

// CreateJob создаёт новую задачу
func (s *JobService) CreateJob(ctx context.Context, params CreateJobParams) (*entity.Job, error) {
	now := time.Now()

	job := &entity.Job{
		ID:          types.ID{},
		ObjectType:  params.ObjectType,
		Name:        params.Name,
		Status:      "pending",
		Interval:    params.Interval,
		ScheduledAt: *params.ScheduledAt,
		Created:     now,
		Updated:     now,
	}

	if params.ObjectID != "" {
		objID := types.ID{}
		job.ObjectID = objID
	}

	if err := s.jobRepo.Create(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create job: %w", err)
	}

	// Если задача должна быть выполнена немедленно, добавляем в очередь
	if params.ScheduledAt == nil || params.ScheduledAt.Before(now) {
		task := &taskqueue.Task{
			Type: taskqueue.TypeProcessJob,
			Payload: map[string]interface{}{
				"job_id": job.ID.String(),
				"data":   params.Data,
			},
			Priority:   params.Priority,
			MaxRetries: 3,
		}

		if err := s.taskQueue.Enqueue(ctx, task); err != nil {
			return nil, fmt.Errorf("failed to enqueue job: %w", err)
		}
	}

	return job, nil
}

// ScheduleJob планирует задачу на выполнение
func (s *JobService) ScheduleJob(ctx context.Context, jobID string, scheduledAt time.Time) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	job.ScheduledAt = scheduledAt
	job.Status = "scheduled"
	job.Updated = time.Now()

	return s.jobRepo.Update(ctx, job)
}

// CancelJob отменяет задачу
func (s *JobService) CancelJob(ctx context.Context, jobID string) error {
	job, err := s.jobRepo.GetByID(ctx, jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	job.Status = "cancelled"
	job.Updated = time.Now()

	return s.jobRepo.Update(ctx, job)
}

// GetJob получает задачу по ID
func (s *JobService) GetJob(ctx context.Context, jobID string) (*entity.Job, error) {
	return s.jobRepo.GetByID(ctx, jobID)
}

// ListJobs получает список задач с фильтрацией
func (s *JobService) ListJobs(ctx context.Context, filter repository.JobFilter) ([]*entity.Job, int64, error) {
	return s.jobRepo.List(ctx, filter)
}

// SyncDataSourceParams параметры для синхронизации источника данных
type SyncDataSourceParams struct {
	DataSourceID string
	Priority     int
}

// SyncDataSource запускает синхронизацию источника данных
func (s *JobService) SyncDataSource(ctx context.Context, params SyncDataSourceParams) error {
	task := &taskqueue.Task{
		Type: taskqueue.TypeSyncDataSource,
		Payload: map[string]interface{}{
			"datasource_id": params.DataSourceID,
		},
		Priority:   params.Priority,
		MaxRetries: 3,
	}

	return s.taskQueue.Enqueue(ctx, task)
}
