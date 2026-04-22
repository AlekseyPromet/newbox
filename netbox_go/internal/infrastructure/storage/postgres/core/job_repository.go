package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/domain/core/repository"
	coredb "netbox_go/internal/infrastructure/storage/sqlc/core"
	"netbox_go/pkg/types"
)

// JobPostgresRepository реализует JobRepository для PostgreSQL
type JobPostgresRepository struct {
	db *sql.DB
}

// NewJobPostgresRepository создаёт новый экземпляр репозитория
func NewJobPostgresRepository(db *sql.DB) repository.JobRepository {
	return &JobPostgresRepository{db: db}
}

// GetByID возвращает задачу по ID
func (r *JobPostgresRepository) GetByID(ctx context.Context, id types.ID) (*entity.Job, error) {
	q := coredb.New(r.db)
	row, err := q.GetJobByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}
		return nil, err
	}

	return &entity.Job{
		ID:           row.ID,
		ObjectType:   row.ObjectType,
		ObjectID:     row.ObjectID,
		Name:         row.Name,
		Status:       types.Status(row.Status),
		Interval:     row.Interval,
		ScheduledAt:  row.ScheduledAt,
		StartedAt:    row.StartedAt,
		CompletedAt:  row.CompletedAt,
		QueueName:    row.QueueName,
		JobID:        row.JobID,
		Data:         row.Data,
		Error:        row.Error,
		Created:      row.Created,
		Updated:      row.Updated,
	}, nil
}

// List возвращает список задач с фильтрацией и пагинацией
func (r *JobPostgresRepository) List(ctx context.Context, filter repository.JobFilter, limit, offset int) ([]*entity.Job, int, error) {
	var query string
	var args []interface{}

	query = `SELECT id, object_type, object_id, name, status, interval, scheduled_at, started_at, completed_at, queue_name, job_id, data, error, created, updated FROM core_job WHERE 1=1`
	countQuery := `SELECT COUNT(*)::int FROM core_job WHERE 1=1`

	if len(filter.Status) > 0 {
		statusList := make([]string, len(filter.Status))
		for i, s := range filter.Status {
			statusList[i] = string(s)
		}
		query += ` AND status = ANY($` + fmt.Sprintf("%d", len(args)+1) + `)`
		args = append(args, statusList)
		countQuery += ` AND status = ANY($` + fmt.Sprintf("%d", len(args)) + `)`
	}

	if filter.ObjectType != nil && *filter.ObjectType != "" {
		query += ` AND object_type = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.ObjectType)
		countQuery += ` AND object_type = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.ObjectID != nil {
		query += ` AND object_id = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.ObjectID)
		countQuery += ` AND object_id = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.QueueName != nil && *filter.QueueName != "" {
		query += ` AND queue_name = $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.QueueName)
		countQuery += ` AND queue_name = $` + fmt.Sprintf("%d", len(args))
	}

	if filter.CreatedAfter != nil {
		query += ` AND created >= $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.CreatedAfter)
		countQuery += ` AND created >= $` + fmt.Sprintf("%d", len(args))
	}

	if filter.CreatedBefore != nil {
		query += ` AND created <= $` + fmt.Sprintf("%d", len(args)+1)
		args = append(args, *filter.CreatedBefore)
		countQuery += ` AND created <= $` + fmt.Sprintf("%d", len(args))
	}

	if filter.SearchQuery != nil && *filter.SearchQuery != "" {
		searchVal := "%" + *filter.SearchQuery + "%"
		query += ` AND (name ILIKE $` + fmt.Sprintf("%d", len(args)+1) + ` OR data::text ILIKE $` + fmt.Sprintf("%d", len(args)+1) + `)`
		args = append(args, searchVal)
		countQuery += ` AND (name ILIKE $` + fmt.Sprintf("%d", len(args)) + ` OR data::text ILIKE $` + fmt.Sprintf("%d", len(args)) + `)`
	}

	query += ` ORDER BY created DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)

	rows, err := r.db.QueryContext(ctx, query, append(args, limit, offset)...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var result []*entity.Job
	for rows.Next() {
		var row struct {
			ID           types.ID
			ObjectType   string
			ObjectID     types.ID
			Name         string
			Status       string
			Interval     string
			ScheduledAt  time.Time
			StartedAt    sql.NullTime
			CompletedAt  sql.NullTime
			QueueName    string
			JobID        sql.NullString
			Data         []byte
			Error        sql.NullString
			Created      time.Time
			Updated      time.Time
		}
		if err := rows.Scan(&row.ID, &row.ObjectType, &row.ObjectID, &row.Name, &row.Status, &row.Interval, &row.ScheduledAt, &row.StartedAt, &row.CompletedAt, &row.QueueName, &row.JobID, &row.Data, &row.Error, &row.Created, &row.Updated); err != nil {
			return nil, 0, err
		}

		var startedAt, completedAt *time.Time
		if row.StartedAt.Valid {
			startedAt = &row.StartedAt.Time
		}
		if row.CompletedAt.Valid {
			completedAt = &row.CompletedAt.Time
		}

		result = append(result, &entity.Job{
			ID:           row.ID,
			ObjectType:   row.ObjectType,
			ObjectID:     row.ObjectID,
			Name:         row.Name,
			Status:       types.Status(row.Status),
			Interval:     row.Interval,
			ScheduledAt:  row.ScheduledAt,
			StartedAt:    startedAt,
			CompletedAt:  completedAt,
			QueueName:    row.QueueName,
			JobID:        row.JobID.String,
			Data:         row.Data,
			Error:        row.Error.String,
			Created:      row.Created,
			Updated:      row.Updated,
		})
	}

	var count int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&count)
	if err != nil {
		return nil, 0, err
	}

	return result, count, nil
}

// Create создаёт новую задачу
func (r *JobPostgresRepository) Create(ctx context.Context, job *entity.Job) error {
	q := coredb.New(r.db)
	
	data := job.Data
	if data == nil {
		data = []byte("{}")
	}

	var scheduledAt sql.NullTime
	if job.ScheduledAt.IsZero() {
		scheduledAt = sql.NullTime{Valid: false}
	} else {
		scheduledAt = sql.NullTime{Time: job.ScheduledAt, Valid: true}
	}

	var startedAt, completedAt sql.NullTime
	if job.StartedAt != nil {
		startedAt = sql.NullTime{Time: *job.StartedAt, Valid: true}
	}
	if job.CompletedAt != nil {
		completedAt = sql.NullTime{Time: *job.CompletedAt, Valid: true}
	}

	row, err := q.CreateJob(ctx, coredb.CreateJobParams{
		ObjectType:   job.ObjectType,
		ObjectID:     job.ObjectID,
		Name:         job.Name,
		Status:       string(job.Status),
		Interval:     job.Interval,
		ScheduledAt:  scheduledAt,
		StartedAt:    startedAt,
		CompletedAt:  completedAt,
		QueueName:    job.QueueName,
		JobID:        sql.NullString{String: job.JobID, Valid: job.JobID != ""},
		Data:         data,
		Error:        sql.NullString{String: job.Error, Valid: job.Error != ""},
		Created:      time.Now(),
		Updated:      time.Now(),
	})
	if err != nil {
		return err
	}

	job.ID = row.ID
	job.Created = row.Created
	job.Updated = row.Updated
	return nil
}

// Update обновляет задачу
func (r *JobPostgresRepository) Update(ctx context.Context, job *entity.Job) error {
	q := coredb.New(r.db)

	data := job.Data
	if data == nil {
		data = []byte("{}")
	}

	var scheduledAt, startedAt, completedAt sql.NullTime
	if !job.ScheduledAt.IsZero() {
		scheduledAt = sql.NullTime{Time: job.ScheduledAt, Valid: true}
	}
	if job.StartedAt != nil {
		startedAt = sql.NullTime{Time: *job.StartedAt, Valid: true}
	}
	if job.CompletedAt != nil {
		completedAt = sql.NullTime{Time: *job.CompletedAt, Valid: true}
	}

	_, err := q.UpdateJob(ctx, coredb.UpdateJobParams{
		ID:           job.ID,
		ObjectType:   job.ObjectType,
		ObjectID:     job.ObjectID,
		Name:         job.Name,
		Status:       string(job.Status),
		Interval:     job.Interval,
		ScheduledAt:  scheduledAt,
		StartedAt:    startedAt,
		CompletedAt:  completedAt,
		QueueName:    job.QueueName,
		JobID:        sql.NullString{String: job.JobID, Valid: job.JobID != ""},
		Data:         data,
		Error:        sql.NullString{String: job.Error, Valid: job.Error != ""},
		Updated:      time.Now(),
	})
	return err
}

// Delete удаляет задачу
func (r *JobPostgresRepository) Delete(ctx context.Context, id types.ID) error {
	q := coredb.New(r.db)
	_, err := q.DeleteJob(ctx, id)
	return err
}

// UpdateStatus обновляет статус задачи
func (r *JobPostgresRepository) UpdateStatus(ctx context.Context, id types.ID, status types.Status, jobError *string, completedAt *time.Time) error {
	q := coredb.New(r.db)

	var completedAtNull sql.NullTime
	if completedAt != nil {
		completedAtNull = sql.NullTime{Time: *completedAt, Valid: true}
	}

	return q.UpdateJobStatus(ctx, coredb.UpdateJobStatusParams{
		ID:          id,
		Status:      string(status),
		Error:       sql.NullString{String: valOrEmpty(jobError), Valid: jobError != nil},
		CompletedAt: completedAtNull,
	})
}

// GetScheduled возвращает запланированные задачи
func (r *JobPostgresRepository) GetScheduled(ctx context.Context, before time.Time, limit int) ([]*entity.Job, error) {
	q := coredb.New(r.db)
	rows, err := q.GetScheduledJobs(ctx, before, int32(limit))
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Job, len(rows))
	for i, row := range rows {
		var startedAt, completedAt *time.Time
		if row.StartedAt.Valid {
			startedAt = &row.StartedAt.Time
		}
		if row.CompletedAt.Valid {
			completedAt = &row.CompletedAt.Time
		}

		result[i] = &entity.Job{
			ID:           row.ID,
			ObjectType:   row.ObjectType,
			ObjectID:     row.ObjectID,
			Name:         row.Name,
			Status:       types.Status(row.Status),
			Interval:     row.Interval,
			ScheduledAt:  row.ScheduledAt,
			StartedAt:    startedAt,
			CompletedAt:  completedAt,
			QueueName:    row.QueueName,
			JobID:        row.JobID.String,
			Data:         row.Data,
			Error:        row.Error.String,
			Created:      row.Created,
			Updated:      row.Updated,
		}
	}

	return result, nil
}

// CleanupOld удаляет старые завершённые задачи
func (r *JobPostgresRepository) CleanupOld(ctx context.Context, olderThan time.Time) (int64, error) {
	q := coredb.New(r.db)
	res, err := q.CleanupOldJobs(ctx, olderThan)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected(), nil
}

func valOrEmpty(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
