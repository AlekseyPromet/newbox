package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"netbox_go/internal/domain/core/entity"
	"netbox_go/internal/domain/core/repository"
	"netbox_go/pkg/types"
)

// JobHandler is a function that executes a specific job type
type JobHandler func(ctx context.Context, job *entity.Job) error

// Worker handles the execution of background jobs using a DB-backed queue
type Worker struct {
	ctx        context.Context
	cancel     context.CancelFunc
	jobRepo    repository.JobRepository
	handlers   map[string]JobHandler
	pollInterval time.Duration
}

// NewWorker creates a new background worker
func NewWorker(jobRepo repository.JobRepository) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		ctx:          ctx,
		cancel:       cancel,
		jobRepo:      jobRepo,
		handlers:     make(map[string]JobHandler),
		pollInterval: 5 * time.Second,
	}
}

// RegisterHandler registers a handler for a specific job object type
func (w *Worker) RegisterHandler(objectType string, handler JobHandler) {
	w.handlers[objectType] = handler
}

// Start begins polling for and executing jobs
func (w *Worker) Start() {
	log.Println("Starting NetBox Go background worker pool...")
	
	// Start a few worker goroutines
	for i := 0; i < 5; i++ {
		go w.work()
	}
}

func (w *Worker) work() {
	for {
		select {
		case <-w.ctx.Done():
			return
		default:
			w.processNextJob()
			time.Sleep(w.pollInterval)
		}
	}
}

func (w *Worker) processNextJob() {
	// Get a job that is 'scheduled' or 'pending'
	// In a real implementation, we'd use SELECT ... FOR UPDATE SKIP LOCKED
	jobs, err := w.jobRepo.GetScheduled(w.ctx, time.Now(), 1)
	if err != nil || len(jobs) == 0 {
		return
	}

	job := jobs[0]
	
	// Update status to STARTED
	err = w.jobRepo.UpdateStatus(w.ctx, job.ID, types.Status("started"), nil, nil)
	if err != nil {
		log.Printf("Failed to start job %d: %v", job.ID, err)
		return
	}

	log.Printf("Executing job %d (%s) for object %s", job.ID, job.Name, job.ObjectType)
	
	handler, ok := w.handlers[job.ObjectType]
	if !ok {
		errMsg := fmt.Sprintf("no handler registered for object type: %s", job.ObjectType)
		w.jobRepo.UpdateStatus(w.ctx, job.ID, types.Status("errored"), &errMsg, time.Now())
		log.Printf("Job %d errored: %s", job.ID, errMsg)
		return
	}

	if err := handler(w.ctx, job); err != nil {
		errMsg := err.Error()
		w.jobRepo.UpdateStatus(w.ctx, job.ID, types.Status("failed"), &errMsg, time.Now())
		log.Printf("Job %d failed: %v", job.ID, err)
		return
	}

	// Mark as completed
	now := time.Now()
	w.jobRepo.UpdateStatus(w.ctx, job.ID, types.Status("completed"), nil, &now)
	log.Printf("Job %d completed successfully", job.ID)
}

func (w *Worker) Stop() {
	w.cancel()
}

// SystemJob defines a job that runs at regular intervals
type SystemJob struct {
	Name     string
	Interval time.Duration
	Task     func(ctx context.Context) error
}

func (w *Worker) ScheduleSystemJob(job SystemJob) {
	go func() {
		ticker := time.NewTicker(job.Interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				log.Printf("Running system job: %s", job.Name)
				if err := job.Task(w.ctx); err != nil {
					log.Printf("System job %s failed: %v", job.Name, err)
				}
			case <-w.ctx.Done():
				return
			}
		}
	}()
}
