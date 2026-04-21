package taskqueue

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Worker обработчик задач
type Worker struct {
	ID       string
	queue    *EtcdQueue
	handlers map[string]TaskHandlerFunc
	ctx      context.Context
	cancel   context.CancelFunc
}

// TaskHandlerFunc функция обработки задачи
type TaskHandlerFunc func(ctx context.Context, payload map[string]interface{}) error

// WorkerPool пул воркеров
type WorkerPool struct {
	workers []*Worker
	queue   *EtcdQueue
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
}

// NewWorker создаёт нового воркера
func NewWorker(queue *EtcdQueue, handlers map[string]TaskHandlerFunc) *Worker {
	ctx, cancel := context.WithCancel(context.Background())
	return &Worker{
		ID:       uuid.New().String(),
		queue:    queue,
		handlers: handlers,
		ctx:      ctx,
		cancel:   cancel,
	}
}

// Start запускает воркера
func (w *Worker) Start() {
	log.Printf("Worker %s started", w.ID)

	for {
		select {
		case <-w.ctx.Done():
			log.Printf("Worker %s stopped", w.ID)
			return
		default:
			task, err := w.queue.Dequeue(w.ctx, w.ID)
			if err != nil {
				log.Printf("Worker %s dequeue error: %v", w.ID, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			if task == nil {
				// Очередь пуста, ждём
				time.Sleep(500 * time.Millisecond)
				continue
			}

			// Обрабатываем задачу
			handler, ok := w.handlers[task.Type]
			if !ok {
				log.Printf("Unknown task type: %s", task.Type)
				w.queue.Fail(w.ctx, w.ID, task.ID, fmt.Sprintf("unknown task type: %s", task.Type))
				continue
			}

			err = handler(w.ctx, task.Payload)
			if err != nil {
				log.Printf("Task %s failed: %v", task.ID, err)
				w.queue.Fail(w.ctx, w.ID, task.ID, err.Error())
			} else {
				log.Printf("Task %s completed", task.ID)
				w.queue.Complete(w.ctx, w.ID, task.ID, nil)
			}
		}
	}
}

// Stop останавливает воркера
func (w *Worker) Stop() {
	w.cancel()
}

// NewWorkerPool создаёт пул воркеров
func NewWorkerPool(queue *EtcdQueue, handlers map[string]TaskHandlerFunc, poolSize int) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	pool := &WorkerPool{
		workers: make([]*Worker, poolSize),
		queue:   queue,
		ctx:     ctx,
		cancel:  cancel,
	}

	for i := 0; i < poolSize; i++ {
		pool.workers[i] = NewWorker(queue, handlers)
	}

	return pool
}

// Start запускает все воркеры в пуле
func (p *WorkerPool) Start() {
	for _, worker := range p.workers {
		p.wg.Add(1)
		go func(w *Worker) {
			defer p.wg.Done()
			w.Start()
		}(worker)
	}
	log.Printf("Started %d workers", len(p.workers))
}

// Stop останавливает все воркеры
func (p *WorkerPool) Stop() {
	p.cancel()
	for _, worker := range p.workers {
		worker.Stop()
	}
	p.wg.Wait()
	log.Println("All workers stopped")
}

// RegisterHandler регистрирует обработчик для типа задач
func (p *WorkerPool) RegisterHandler(taskType string, handler TaskHandlerFunc) {
	for _, worker := range p.workers {
		worker.handlers[taskType] = handler
	}
}
