package worker

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Job represents a background task to be executed.
type Job struct {
	ID       string
	Name     string
	Payload  interface{}
	Queue    string
	Priority int
	CreatedAt time.Time
}

// Worker handles the execution of background jobs.
type Worker struct {
	ctx    context.Context
	cancel context.CancelFunc
	// In a real implementation, this would interface with Redis/Postgres
	queue chan Job
}

func NewWorker() *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		ctx:    ctx,
		cancel: cancel,
		queue:  make(chan Job, 100),
	}
}

func (w *Worker) Start() {
	log.Println("Starting NetBox Go background worker...")
	go func() {
		for {
			select {
			case job := <-w.queue:
				w.executeJob(job)
			case <-w.ctx.Done():
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.cancel()
}

func (w *Worker) Enqueue(job Job) {
	w.queue <- job
}

func (w *Worker) executeJob(job Job) {
	log.Printf("Executing job %s (%s)", job.ID, job.Name)
	// Job execution logic would go here
	// For now, we just simulate work
	time.Sleep(1 * time.Second)
	log.Printf("Job %s completed", job.ID)
}

// SystemJob defines a job that runs at regular intervals.
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
