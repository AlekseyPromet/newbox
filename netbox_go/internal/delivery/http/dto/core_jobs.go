package dto

import (
	"time"

	"github.com/google/uuid"
)

// JobRequest is used for creating/scheduling jobs
type JobRequest struct {
	Name       string                 `json:"name" validate:"required,max=200"`
	ObjectType string                 `json:"object_type"`
	ObjectID   *uuid.UUID             `json:"object_id"`
	Interval   *int                   `json:"interval"`
	ScheduledAt *time.Time            `json:"scheduled_at"`
	Data       map[string]interface{} `json:"data"`
	QueueName  string                 `json:"queue_name"`
}

// JobResponse is the full representation of a Job
type JobResponse struct {
	ID          uuid.UUID              `json:"id"`
	ObjectType  string                 `json:"object_type"`
	ObjectID    *uuid.UUID             `json:"object_id"`
	Object      interface{}            `json:"object"` // Dynamic representation of the object
	Name        string                 `json:"name"`
	Status      string                 `json:"status"`
	Created     time.Time              `json:"created"`
	ScheduledAt *time.Time             `json:"scheduled_at"`
	Interval    int                    `json:"interval"`
	StartedAt   *time.Time             `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	User        *UserBriefResponse     `json:"user"`
	Data        map[string]interface{} `json:"data"`
	Error       string                 `json:"error"`
	JobID       string                 `json:"job_id"`
	QueueName   string                 `json:"queue_name"`
}

// JobBriefResponse is a condensed version for lists
type JobBriefResponse struct {
	ID        uuid.UUID  `json:"id"`
	Created   time.Time  `json:"created"`
	Completed *time.Time `json:"completed"`
	User      string     `json:"user"`
	Status    string     `json:"status"`
}

// UserBriefResponse is a simple user representation for nested Job responses
type UserBriefResponse struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	FullName string    `json:"full_name"`
}
