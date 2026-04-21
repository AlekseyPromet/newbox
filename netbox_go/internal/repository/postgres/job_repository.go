// Package postgres содержит реализацию репозиториев для PostgreSQL с использованием sqlc
package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	core_entity "github.com/AlekseyPromet/netbox_go/internal/domain/core/entity"
	"github.com/AlekseyPromet/netbox_go/internal/repository"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// JobRepositoryPostgres реализует интерфейс JobRepository для PostgreSQL
type JobRepositoryPostgres struct {
	db *sql.DB
}

// NewJobRepositoryPostgres создает новый экземпляр репозитория задач
func NewJobRepositoryPostgres(db *sql.DB) *JobRepositoryPostgres {
	return &JobRepositoryPostgres{db: db}
}

// GetByID получает задачу по ID
func (r *JobRepositoryPostgres) GetByID(ctx context.Context, id string) (*core_entity.Job, error) {
	query := `
		SELECT id, object_type, object_id, name, status, interval, scheduled_at,
		       started_at, completed_at, queue_name, job_id, data, error, created, updated
		FROM core_jobs
		WHERE id = $1 AND deleted_at IS NULL
	`

	var job core_entity.Job
	var objectType sql.NullString
	var objectID sql.NullString
	var scheduledAt, startedAt, completedAt sql.NullTime
	var queueName sql.NullString
	var jobID sql.NullString
	var dataJSON, errorJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&job.ID, &objectType, &objectID, &job.Name, &job.Status, &job.Interval,
		&scheduledAt, &startedAt, &completedAt, &queueName, &jobID,
		&dataJSON, &errorJSON, &job.Created, &job.Updated,
	)

	if err == sql.ErrNoRows {
		return nil, types.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get job by ID: %w", err)
	}

	// Заполнение опциональных полей
	if objectType.Valid {
		job.ObjectType = &objectType.String
	}
	if objectID.Valid {
		oid, _ := types.ParseID(objectID.String)
		job.ObjectID = &oid
	}
	if scheduledAt.Valid {
		job.ScheduledAt = &scheduledAt.Time
	}
	if startedAt.Valid {
		job.StartedAt = &startedAt.Time
	}
	if completedAt.Valid {
		job.CompletedAt = &completedAt.Time
	}
	if queueName.Valid {
		job.QueueName = queueName.String
	}
	if jobID.Valid {
		job.JobID = &jobID.String
	}
	if dataJSON != nil {
		job.Data = json.RawMessage(dataJSON)
	}
	if errorJSON != nil {
		errStr := string(errorJSON)
		job.Error = &errStr
	}

	return &job, nil
}

// List получает список задач с фильтрацией
func (r *JobRepositoryPostgres) List(ctx context.Context, filter repository.JobFilter) ([]*core_entity.Job, int64, error) {
	query := `
		SELECT id, object_type, object_id, name, status, interval, scheduled_at,
		       started_at, completed_at, queue_name, job_id, data, error, created, updated
		FROM core_jobs
		WHERE deleted_at IS NULL
	`
	countQuery := `SELECT COUNT(*) FROM core_jobs WHERE deleted_at IS NULL`

	args := []interface{}{}
	argIndex := 1

	if filter.ObjectType != nil {
		query += fmt.Sprintf(" AND object_type = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND object_type = $%d", argIndex)
		args = append(args, *filter.ObjectType)
		argIndex++
	}

	if filter.ObjectID != nil {
		query += fmt.Sprintf(" AND object_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND object_id = $%d", argIndex)
		args = append(args, *filter.ObjectID)
		argIndex++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, *filter.Status)
		argIndex++
	}

	if filter.QueueName != nil {
		query += fmt.Sprintf(" AND queue_name = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND queue_name = $%d", argIndex)
		args = append(args, *filter.QueueName)
		argIndex++
	}

	if filter.ScheduledAt != nil {
		query += fmt.Sprintf(" AND scheduled_at = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND scheduled_at = $%d", argIndex)
		args = append(args, *filter.ScheduledAt)
		argIndex++
	}

	query += " ORDER BY created DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	// Получаем общее количество
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count jobs: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list jobs: %w", err)
	}
	defer rows.Close()

	var jobs []*core_entity.Job
	for rows.Next() {
		var job core_entity.Job
		var objectType sql.NullString
		var objectID sql.NullString
		var scheduledAt, startedAt, completedAt sql.NullTime
		var queueName sql.NullString
		var jobID sql.NullString
		var dataJSON, errorJSON []byte

		err := rows.Scan(
			&job.ID, &objectType, &objectID, &job.Name, &job.Status, &job.Interval,
			&scheduledAt, &startedAt, &completedAt, &queueName, &jobID,
			&dataJSON, &errorJSON, &job.Created, &job.Updated,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan job: %w", err)
		}

		if objectType.Valid {
			job.ObjectType = &objectType.String
		}
		if objectID.Valid {
			oid, _ := types.ParseID(objectID.String)
			job.ObjectID = &oid
		}
		if scheduledAt.Valid {
			job.ScheduledAt = &scheduledAt.Time
		}
		if startedAt.Valid {
			job.StartedAt = &startedAt.Time
		}
		if completedAt.Valid {
			job.CompletedAt = &completedAt.Time
		}
		if queueName.Valid {
			job.QueueName = queueName.String
		}
		if jobID.Valid {
			job.JobID = &jobID.String
		}
		if dataJSON != nil {
			job.Data = json.RawMessage(dataJSON)
		}
		if errorJSON != nil {
			errStr := string(errorJSON)
			job.Error = &errStr
		}

		jobs = append(jobs, &job)
	}

	return jobs, total, nil
}

// Create создает новую задачу
func (r *JobRepositoryPostgres) Create(ctx context.Context, job *core_entity.Job) error {
	query := `
		INSERT INTO core_jobs (id, object_type, object_id, name, status, interval,
		                        scheduled_at, started_at, completed_at, queue_name, job_id,
		                        data, error, created, updated)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, NOW(), NOW())
	`

	var objectType, objectID *string
	if job.ObjectType != nil {
		objectType = job.ObjectType
	}
	if job.ObjectID != nil {
		oid := job.ObjectID.String()
		objectID = &oid
	}

	var scheduledAt, startedAt, completedAt *time.Time
	if job.ScheduledAt != nil {
		scheduledAt = job.ScheduledAt
	}
	if job.StartedAt != nil {
		startedAt = job.StartedAt
	}
	if job.CompletedAt != nil {
		completedAt = job.CompletedAt
	}

	var queueName, jobID *string
	if job.QueueName != "" {
		queueName = &job.QueueName
	}
	if job.JobID != nil {
		jobID = job.JobID
	}

	var dataJSON []byte
	if job.Data != nil {
		dataJSON = job.Data
	}

	var errorStr *string
	if job.Error != nil {
		errorStr = job.Error
	}

	_, err := r.db.ExecContext(ctx, query,
		job.ID.String(), objectType, objectID, job.Name, job.Status, job.Interval,
		scheduledAt, startedAt, completedAt, queueName, jobID, dataJSON, errorStr,
	)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil
}

// Update обновляет задачу
func (r *JobRepositoryPostgres) Update(ctx context.Context, job *core_entity.Job) error {
	query := `
		UPDATE core_jobs
		SET object_type = $1, object_id = $2, name = $3, status = $4, interval = $5,
		    scheduled_at = $6, started_at = $7, completed_at = $8, queue_name = $9,
		    job_id = $10, data = $11, error = $12, updated = NOW()
		WHERE id = $13
	`

	var objectType, objectID *string
	if job.ObjectType != nil {
		objectType = job.ObjectType
	}
	if job.ObjectID != nil {
		oid := job.ObjectID.String()
		objectID = &oid
	}

	var scheduledAt, startedAt, completedAt *time.Time
	if job.ScheduledAt != nil {
		scheduledAt = job.ScheduledAt
	}
	if job.StartedAt != nil {
		startedAt = job.StartedAt
	}
	if job.CompletedAt != nil {
		completedAt = job.CompletedAt
	}

	var queueName, jobID *string
	if job.QueueName != "" {
		queueName = &job.QueueName
	}
	if job.JobID != nil {
		jobID = job.JobID
	}

	var dataJSON []byte
	if job.Data != nil {
		dataJSON = job.Data
	}

	result, err := r.db.ExecContext(ctx, query,
		objectType, objectID, job.Name, job.Status, job.Interval,
		scheduledAt, startedAt, completedAt, queueName, jobID, dataJSON, job.Error, job.ID.String(),
	)
	if err != nil {
		return fmt.Errorf("failed to update job: %w", err)
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

// Delete удаляет задачу (soft delete)
func (r *JobRepositoryPostgres) Delete(ctx context.Context, id string) error {
	query := `UPDATE core_jobs SET deleted_at = NOW() WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
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

// UpdateStatus обновляет статус задачи
func (r *JobRepositoryPostgres) UpdateStatus(ctx context.Context, id string, status string) error {
	query := `
		UPDATE core_jobs
		SET status = $1, updated = NOW()
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
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

// MarkCompleted помечает задачу как завершенную
func (r *JobRepositoryPostgres) MarkCompleted(ctx context.Context, id string, hasError bool, errorMsg *string) error {
	now := time.Now()
	query := `
		UPDATE core_jobs
		SET status = $1, completed_at = $2, error = COALESCE($3, error), updated = NOW()
		WHERE id = $4
	`

	status := core_entity.JobStatusCompleted
	if hasError && errorMsg != nil {
		status = core_entity.JobStatusErrored
	}

	result, err := r.db.ExecContext(ctx, query, status, now, errorMsg, id)
	if err != nil {
		return fmt.Errorf("failed to mark job as completed: %w", err)
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
