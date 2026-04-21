package taskqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	TaskQueuePrefix      = "/netbox/tasks/queue/"
	TaskProcessingPrefix = "/netbox/tasks/processing/"
	TaskResultPrefix     = "/netbox/tasks/result/"
	TaskLockPrefix       = "/netbox/tasks/lock/"

	TypeSyncDataSource = "core:sync_datasource"
	TypeProcessJob     = "core:process_job"
)

// TaskStatus статус задачи
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

// Task задача для выполнения
type Task struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Payload    map[string]interface{} `json:"payload"`
	Status     TaskStatus             `json:"status"`
	Priority   int                    `json:"priority"`
	CreatedAt  time.Time              `json:"created_at"`
	StartedAt  *time.Time             `json:"started_at,omitempty"`
	EndedAt    *time.Time             `json:"ended_at,omitempty"`
	Result     interface{}            `json:"result,omitempty"`
	Error      string                 `json:"error,omitempty"`
	RetryCount int                    `json:"retry_count"`
	MaxRetries int                    `json:"max_retries"`
}

// EtcdQueue очередь задач на базе etcd
type EtcdQueue struct {
	client *clientv3.Client
	ctx    context.Context
}

// NewEtcdQueue создаёт новую очередь задач
func NewEtcdQueue(client *clientv3.Client) *EtcdQueue {
	return &EtcdQueue{
		client: client,
		ctx:    context.Background(),
	}
}

// Enqueue добавляет задачу в очередь
func (q *EtcdQueue) Enqueue(ctx context.Context, task *Task) error {
	task.ID = uuid.New().String()
	task.Status = TaskStatusPending
	task.CreatedAt = time.Now()

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// Ключ с приоритетом для сортировки (чем меньше число, тем выше приоритет)
	key := q.buildQueueKey(task.Priority, task.ID)

	_, err = q.client.Put(ctx, key, string(data))
	return err
}

// Dequeue извлекает следующую задачу из очереди
func (q *EtcdQueue) Dequeue(ctx context.Context, workerID string) (*Task, error) {
	// Получаем задачи из очереди, отсортированные по приоритету
	resp, err := q.client.Get(
		ctx,
		TaskQueuePrefix,
		clientv3.WithPrefix(),
		clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend),
		clientv3.WithLimit(1),
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		return nil, nil // Очередь пуста
	}

	kv := resp.Kvs[0]
	var task Task
	if err := json.Unmarshal(kv.Value, &task); err != nil {
		return nil, err
	}

	// Атомарно перемещаем задачу в processing
	task.Status = TaskStatusRunning
	now := time.Now()
	task.StartedAt = &now

	data, err := json.Marshal(task)
	if err != nil {
		return nil, err
	}

	// Используем транзакцию для атомарности
	txnResp, err := q.client.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(string(kv.Key)), "=", kv.CreateRevision)).
		Then(
			clientv3.OpDelete(string(kv.Key)),
			clientv3.OpPut(q.buildProcessingKey(workerID, task.ID), string(data)),
		).
		Commit()

	if err != nil || !txnResp.Succeeded {
		return nil, fmt.Errorf("failed to move task to processing: %w", err)
	}

	return &task, nil
}

// Complete завершает задачу успешно
func (q *EtcdQueue) Complete(ctx context.Context, workerID, taskID string, result interface{}) error {
	processingKey := q.buildProcessingKey(workerID, taskID)

	// Получаем текущую задачу
	resp, err := q.client.Get(ctx, processingKey)
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return fmt.Errorf("task not found in processing")
	}

	var task Task
	if err := json.Unmarshal(resp.Kvs[0].Value, &task); err != nil {
		return err
	}

	task.Status = TaskStatusCompleted
	now := time.Now()
	task.EndedAt = &now
	task.Result = result

	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// Перемещаем в результаты и удаляем из processing
	_, err = q.client.Txn(ctx).
		Then(
			clientv3.OpDelete(processingKey),
			clientv3.OpPut(q.buildResultKey(taskID), string(data)),
		).
		Commit()

	return err
}

// Fail отмечает задачу как неудачную
func (q *EtcdQueue) Fail(ctx context.Context, workerID, taskID string, errMsg string) error {
	processingKey := q.buildProcessingKey(workerID, taskID)

	resp, err := q.client.Get(ctx, processingKey)
	if err != nil {
		return err
	}
	if len(resp.Kvs) == 0 {
		return fmt.Errorf("task not found in processing")
	}

	var task Task
	if err := json.Unmarshal(resp.Kvs[0].Value, &task); err != nil {
		return err
	}

	task.RetryCount++

	if task.RetryCount >= task.MaxRetries {
		task.Status = TaskStatusFailed
		now := time.Now()
		task.EndedAt = &now
		task.Error = errMsg

		data, err := json.Marshal(task)
		if err != nil {
			return err
		}

		// Перемещаем в результаты как failed
		_, err = q.client.Txn(ctx).
			Then(
				clientv3.OpDelete(processingKey),
				clientv3.OpPut(q.buildResultKey(taskID), string(data)),
			).
			Commit()
		return err
	} else {
		// Возвращаем в очередь для повторной попытки
		task.Status = TaskStatusPending
		task.StartedAt = nil

		data, err := json.Marshal(task)
		if err != nil {
			return err
		}

		key := q.buildQueueKey(task.Priority, task.ID)
		_, err = q.client.Txn(ctx).
			Then(
				clientv3.OpDelete(processingKey),
				clientv3.OpPut(key, string(data)),
			).
			Commit()
		return err
	}
}

// GetTaskResult получает результат выполненной задачи
func (q *EtcdQueue) GetTaskResult(ctx context.Context, taskID string) (*Task, error) {
	resp, err := q.client.Get(ctx, q.buildResultKey(taskID))
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("task result not found")
	}

	var task Task
	if err := json.Unmarshal(resp.Kvs[0].Value, &task); err != nil {
		return nil, err
	}

	return &task, nil
}

// WatchQueue следит за появлением новых задач
func (q *EtcdQueue) WatchQueue(ctx context.Context) clientv3.WatchChan {
	return q.client.Watch(ctx, TaskQueuePrefix, clientv3.WithPrefix())
}

// buildQueueKey строит ключ для очереди с учётом приоритета
func (q *EtcdQueue) buildQueueKey(priority int, taskID string) string {
	// Инвертируем приоритет для сортировки (меньшее число = выше приоритет)
	invertedPriority := 9999 - priority
	return fmt.Sprintf("%s%04d/%s", TaskQueuePrefix, invertedPriority, taskID)
}

// buildProcessingKey строит ключ для выполняемых задач
func (q *EtcdQueue) buildProcessingKey(workerID, taskID string) string {
	return fmt.Sprintf("%s%s/%s", TaskProcessingPrefix, workerID, taskID)
}

// buildResultKey строит ключ для результатов задач
func (q *EtcdQueue) buildResultKey(taskID string) string {
	return fmt.Sprintf("%s%s", TaskResultPrefix, taskID)
}

// CleanupStaleTasks очищает зависшие задачи (например, после краша воркера)
func (q *EtcdQueue) CleanupStaleTasks(ctx context.Context, timeout time.Duration) error {
	resp, err := q.client.Get(ctx, TaskProcessingPrefix, clientv3.WithPrefix())
	if err != nil {
		return err
	}

	now := time.Now()
	for _, kv := range resp.Kvs {
		var task Task
		if err := json.Unmarshal(kv.Value, &task); err != nil {
			continue
		}

		if task.StartedAt != nil && now.Sub(*task.StartedAt) > timeout {
			// Задача выполняется слишком долго, возвращаем в очередь
			task.Status = TaskStatusPending
			task.StartedAt = nil

			data, err := json.Marshal(task)
			if err != nil {
				continue
			}

			key := q.buildQueueKey(task.Priority, task.ID)
			q.client.Txn(ctx).
				Then(
					clientv3.OpDelete(string(kv.Key)),
					clientv3.OpPut(key, string(data)),
				).
				Commit()
		}
	}

	return nil
}
