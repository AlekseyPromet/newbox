// Package etcd содержит реализацию кэширования и распределенных блокировок через Etcd
package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AlekseyPromet/netbox_go/pkg/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// CacheClient представляет клиент для работы с кэшем в Etcd
type CacheClient struct {
	client *clientv3.Client
	prefix string
	ttl    time.Duration
}

// NewCacheClient создает новый экземпляр клиента кэша
func NewCacheClient(endpoints []string, prefix string, ttl time.Duration) (*CacheClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &CacheClient{
		client: client,
		prefix: prefix,
		ttl:    ttl,
	}, nil
}

// Close закрывает соединение с Etcd
func (c *CacheClient) Close() error {
	return c.client.Close()
}

// Get получает значение из кэша по ключу
func (c *CacheClient) Get(ctx context.Context, key string, result interface{}) error {
	fullKey := fmt.Sprintf("%s/%s", c.prefix, key)

	resp, err := c.client.Get(ctx, fullKey)
	if err != nil {
		return fmt.Errorf("failed to get from cache: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return types.ErrNotFound
	}

	if err := json.Unmarshal(resp.Kvs[0].Value, result); err != nil {
		return fmt.Errorf("failed to unmarshal cached value: %w", err)
	}

	return nil
}

// Set устанавливает значение в кэш с TTL
func (c *CacheClient) Set(ctx context.Context, key string, value interface{}) error {
	fullKey := fmt.Sprintf("%s/%s", c.prefix, key)

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	leaseResp, err := c.client.Grant(ctx, int64(c.ttl.Seconds()))
	if err != nil {
		return fmt.Errorf("failed to grant lease: %w", err)
	}

	_, err = c.client.Put(ctx, fullKey, string(data), clientv3.WithLease(leaseResp.ID))
	if err != nil {
		return fmt.Errorf("failed to put to cache: %w", err)
	}

	return nil
}

// Delete удаляет значение из кэша
func (c *CacheClient) Delete(ctx context.Context, key string) error {
	fullKey := fmt.Sprintf("%s/%s", c.prefix, key)

	_, err := c.client.Delete(ctx, fullKey)
	if err != nil {
		return fmt.Errorf("failed to delete from cache: %w", err)
	}

	return nil
}

// InvalidateByPattern инвалидирует кэш по паттерну
func (c *CacheClient) InvalidateByPattern(ctx context.Context, pattern string) error {
	fullPattern := fmt.Sprintf("%s/%s", c.prefix, pattern)

	_, err := c.client.Delete(ctx, fullPattern, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("failed to invalidate cache by pattern: %w", err)
	}

	return nil
}

// LockClient представляет клиент для распределенных блокировок
type LockClient struct {
	client *clientv3.Client
	prefix string
}

// NewLockClient создает новый экземпляр клиента блокировок
func NewLockClient(endpoints []string, prefix string) (*LockClient, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return &LockClient{
		client: client,
		prefix: prefix,
	}, nil
}

// Close закрывает соединение с Etcd
func (l *LockClient) Close() error {
	return l.client.Close()
}

// DistributedLock представляет распределенную блокировку
type DistributedLock struct {
	key      string
	leaseID  clientv3.LeaseID
	cancelFn context.CancelFunc
}

// Acquire захватывает распределенную блокировку
func (l *LockClient) Acquire(ctx context.Context, resource string, ttl time.Duration) (*DistributedLock, error) {
	lockKey := fmt.Sprintf("%s/locks/%s", l.prefix, resource)

	// Создаем lease
	leaseResp, err := l.client.Grant(ctx, int64(ttl.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("failed to grant lease: %w", err)
	}

	// Пытаемся захватить блокировку
	resp, err := l.client.Txn(ctx).
		If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseResp.ID))).
		Commit()

	if err != nil {
		l.client.Revoke(ctx, leaseResp.ID)
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	if !resp.Succeeded {
		l.client.Revoke(ctx, leaseResp.ID)
		return nil, fmt.Errorf("lock is already held by another client")
	}

	// Запускаем keepalive для продления lease
	keepAliveCtx, cancel := context.WithCancel(ctx)
	keepAliveChan, err := l.client.KeepAlive(keepAliveCtx, leaseResp.ID)
	if err != nil {
		cancel()
		l.client.Revoke(ctx, leaseResp.ID)
		return nil, fmt.Errorf("failed to start keepalive: %w", err)
	}

	// Читаем первый ответ keepalive чтобы убедиться что все работает
	<-keepAliveChan

	return &DistributedLock{
		key:      lockKey,
		leaseID:  leaseResp.ID,
		cancelFn: cancel,
	}, nil
}

// Release освобождает распределенную блокировку
func (dl *DistributedLock) Release(ctx context.Context, client *clientv3.Client) error {
	if dl.cancelFn != nil {
		dl.cancelFn()
	}

	_, err := client.Delete(ctx, dl.key)
	if err != nil {
		return fmt.Errorf("failed to delete lock key: %w", err)
	}

	_, err = client.Revoke(ctx, dl.leaseID)
	if err != nil {
		return fmt.Errorf("failed to revoke lease: %w", err)
	}

	return nil
}

// WithLock выполняет функцию с захваченной блокировкой
func (l *LockClient) WithLock(ctx context.Context, resource string, ttl time.Duration, fn func() error) error {
	lock, err := l.Acquire(ctx, resource, ttl)
	if err != nil {
		return err
	}
	defer lock.Release(ctx, l.client)

	return fn()
}
