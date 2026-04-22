// Package memory содержит простые in-memory реализации репозиториев account
package memory

import (
	"context"
	"net"
	"sync"
	"time"

	account_entity "netbox_go/internal/domain/account/entity"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
)

// AccountRepositoryMemory реализует UserTokenRepository, BookmarkRepository, NotificationRepository,
// SubscriptionRepository и UserConfigRepository для упрощённого in-memory хранения.
type AccountRepositoryMemory struct {
	mu            sync.RWMutex
	tokens        map[types.ID]*account_entity.UserToken
	bookmarks     map[types.ID]*account_entity.Bookmark
	notifications map[types.ID]*account_entity.Notification
	subscriptions map[types.ID]*account_entity.Subscription
	userConfigs   map[string]*account_entity.UserConfig // ключ = userID.String()
}

// NewAccountRepositoryMemory создаёт новый in-memory репозиторий
func NewAccountRepositoryMemory() *AccountRepositoryMemory {
	return &AccountRepositoryMemory{
		tokens:        make(map[types.ID]*account_entity.UserToken),
		bookmarks:     make(map[types.ID]*account_entity.Bookmark),
		notifications: make(map[types.ID]*account_entity.Notification),
		subscriptions: make(map[types.ID]*account_entity.Subscription),
		userConfigs:   make(map[string]*account_entity.UserConfig),
	}
}

// --- UserTokenRepository ---

func (r *AccountRepositoryMemory) ListByUser(ctx context.Context, userID types.ID) ([]*account_entity.UserToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := []*account_entity.UserToken{}
	for _, t := range r.tokens {
		if t.UserID == userID {
			clone := *t
			res = append(res, &clone)
		}
	}
	return res, nil
}

func (r *AccountRepositoryMemory) Get(ctx context.Context, id types.ID, userID types.ID) (*account_entity.UserToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tokens[id]
	if !ok || t.UserID != userID {
		return nil, repository.ErrNotFound
	}
	clone := *t
	return &clone, nil
}

func (r *AccountRepositoryMemory) Create(ctx context.Context, token *account_entity.UserToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if token.ID.String() == "" {
		token.ID = types.NewID()
	}
	now := time.Now()
	if token.Created.IsZero() {
		token.Created = now
	}
	r.tokens[token.ID] = cloneToken(token)
	return nil
}

func (r *AccountRepositoryMemory) Update(ctx context.Context, token *account_entity.UserToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	existing, ok := r.tokens[token.ID]
	if !ok || existing.UserID != token.UserID {
		return repository.ErrNotFound
	}
	r.tokens[token.ID] = cloneToken(token)
	return nil
}

func (r *AccountRepositoryMemory) Delete(ctx context.Context, id types.ID, userID types.ID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if existing, ok := r.tokens[id]; !ok || existing.UserID != userID {
		return repository.ErrNotFound
	}
	delete(r.tokens, id)
	return nil
}

func cloneToken(t *account_entity.UserToken) *account_entity.UserToken {
	clone := *t
	if t.AllowedIPs != nil {
		clone.AllowedIPs = make([]*net.IPNet, len(t.AllowedIPs))
		copy(clone.AllowedIPs, t.AllowedIPs)
	}
	return &clone
}

// --- BookmarkRepository ---

func (r *AccountRepositoryMemory) ListBookmarksByUser(ctx context.Context, userID types.ID) ([]*account_entity.Bookmark, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := []*account_entity.Bookmark{}
	for _, b := range r.bookmarks {
		if b.UserID == userID {
			clone := *b
			res = append(res, &clone)
		}
	}
	return res, nil
}

// --- NotificationRepository ---

func (r *AccountRepositoryMemory) ListNotificationsByUser(ctx context.Context, userID types.ID) ([]*account_entity.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := []*account_entity.Notification{}
	for _, n := range r.notifications {
		if n.UserID == userID {
			clone := *n
			res = append(res, &clone)
		}
	}
	return res, nil
}

// --- SubscriptionRepository ---

func (r *AccountRepositoryMemory) ListSubscriptionsByUser(ctx context.Context, userID types.ID) ([]*account_entity.Subscription, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	res := []*account_entity.Subscription{}
	for _, s := range r.subscriptions {
		if s.UserID == userID {
			clone := *s
			res = append(res, &clone)
		}
	}
	return res, nil
}

// --- UserConfigRepository ---

func (r *AccountRepositoryMemory) GetByUser(ctx context.Context, userID types.ID) (*account_entity.UserConfig, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	cfg, ok := r.userConfigs[userID.String()]
	if !ok {
		return nil, repository.ErrNotFound
	}
	clone := *cfg
	return &clone, nil
}

func (r *AccountRepositoryMemory) Upsert(ctx context.Context, config *account_entity.UserConfig) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.userConfigs[config.UserID.String()] = &account_entity.UserConfig{
		UserID: config.UserID,
		Data:   append([]byte(nil), config.Data...),
	}
	return nil
}
