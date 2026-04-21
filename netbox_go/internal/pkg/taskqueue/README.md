# Task Queue на базе etcd

Этот пакет реализует систему фоновых задач для NetBox Go, используя etcd в качестве хранилища.

## Архитектура

### Компоненты

1. **EtcdQueue** - очередь задач на базе etcd
2. **WorkerPool** - пул воркеров для обработки задач
3. **Task** - структура задачи

### Структура ключей в etcd

```
/netbox/tasks/queue/<priority>/<task_id>      # Очередь задач
/netbox/tasks/processing/<worker_id>/<task_id> # Выполняемые задачи
/netbox/tasks/result/<task_id>                 # Результаты задач
/netbox/tasks/lock/<resource>                  # Блокировки
```

## Использование

### Создание очереди

```go
import (
    "github.com/AlekseyPromet/netbox_go/internal/pkg/taskqueue"
    clientv3 "go.etcd.io/etcd/client/v3"
)

// Подключение к etcd
client, err := clientv3.New(clientv3.Config{
    Endpoints:   []string{"localhost:2379"},
    DialTimeout: 5 * time.Second,
})

// Создание очереди
queue := taskqueue.NewEtcdQueue(client)
```

### Добавление задачи

```go
task := &taskqueue.Task{
    Type: taskqueue.TypeProcessJob,
    Payload: map[string]interface{}{
        "job_id": "123",
        "data":   map[string]interface{}{"key": "value"},
    },
    Priority:   5,       // 0-9999 (меньше = выше приоритет)
    MaxRetries: 3,
}

err := queue.Enqueue(ctx, task)
```

### Создание пула воркеров

```go
// Обработчики задач
handlers := map[string]taskqueue.TaskHandlerFunc{
    taskqueue.TypeProcessJob: func(ctx context.Context, payload map[string]interface{}) error {
        jobID := payload["job_id"].(string)
        // Логика обработки задачи
        return nil
    },
    taskqueue.TypeSyncDataSource: func(ctx context.Context, payload map[string]interface{}) error {
        dsID := payload["datasource_id"].(string)
        // Логика синхронизации
        return nil
    },
}

// Создание и запуск пула
pool := taskqueue.NewWorkerPool(queue, handlers, 4) // 4 воркера
pool.Start()

// Остановка (например, по сигналу)
defer pool.Stop()
```

### Получение результата задачи

```go
result, err := queue.GetTaskResult(ctx, taskID)
if err != nil {
    // Обработка ошибки
}

if result.Status == taskqueue.TaskStatusCompleted {
    // Задача выполнена успешно
} else if result.Status == taskqueue.TaskStatusFailed {
    // Задача не выполнена
    log.Printf("Error: %s", result.Error)
}
```

## Типы задач

| Тип задачи | Описание | Payload |
|------------|----------|---------|
| `core:sync_datasource` | Синхронизация источника данных | `datasource_id` |
| `core:process_job` | Обработка фоновой задачи | `job_id`, `data` |

## Приоритеты

- 0-999: Критические задачи
- 1000-4999: Высокий приоритет
- 5000-7999: Нормальный приоритет
- 8000-9999: Низкий приоритет

## Надёжность

### Повторные попытки

Задачи с ошибками автоматически повторяются до `MaxRetries` раз.

### Очистка зависших задач

```go
// Очистка задач, выполняющихся более 30 минут
err := queue.CleanupStaleTasks(ctx, 30*time.Minute)
```

### Мониторинг

Используйте WatchQueue для отслеживания новых задач:

```go
watchChan := queue.WatchQueue(ctx)
for watchResp := range watchChan {
    for _, event := range watchResp.Events {
        log.Printf("New task event: %s", event.Type)
    }
}
```

## Отличия от Asynq

| Характеристика | Asynq | EtcdQueue |
|----------------|-------|-----------|
| Хранилище | Redis | etcd |
| Приоритеты | Поддерживаются | Поддерживаются |
| Повторные попытки | Поддерживаются | Поддерживаются |
| Распределённые блокировки | Нет | Да (через etcd) |
| Наблюдаемость | Dashboard | etcdctl / API |
| HA | Redis Sentinel/Cluster | etcd Raft |

## Требования

- etcd кластер (минимум 3 ноды для HA)
- Go 1.19+

