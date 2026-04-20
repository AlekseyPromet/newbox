// Package postgres содержит реализации account-репозиториев для PostgreSQL
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"net"
	"time"

	account_entity "github.com/AlekseyPromet/netbox_go/internal/domain/account/entity"
	"github.com/AlekseyPromet/netbox_go/internal/repository"
	"github.com/AlekseyPromet/netbox_go/pkg/types"
	"github.com/lib/pq"
)

// AccountRepositoryPostgres реализует UserTokenRepository, BookmarkRepository,
// NotificationRepository, SubscriptionRepository и UserConfigRepository
// в едином экземпляре поверх PostgreSQL.
type AccountRepositoryPostgres struct {
	db *sql.DB
}

// NewAccountRepositoryPostgres конструктор
func NewAccountRepositoryPostgres(db *sql.DB) *AccountRepositoryPostgres {
	return &AccountRepositoryPostgres{db: db}
}

// --- UserTokenRepository ---

func (r *AccountRepositoryPostgres) ListByUser(ctx context.Context, userID types.ID) ([]*account_entity.UserToken, error) {
	query := `
        SELECT id, version, user_id, description, created, expires, last_used, enabled, write_enabled,
               plaintext, key, pepper_id, hmac_digest, allowed_ips
        FROM users_tokens
        WHERE user_id = $1
        ORDER BY created DESC`

	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("list tokens: %w", err)
	}
	defer rows.Close()

	var result []*account_entity.UserToken
	for rows.Next() {
		var t account_entity.UserToken
		var versionStr string
		var expires, lastUsed sql.NullTime
		var plaintext, key, hmac sql.NullString
		var pepperID sql.NullInt16
		var ips pq.StringArray

		if err := rows.Scan(
			&t.ID, &versionStr, &t.UserID, &t.Description, &t.Created,
			&expires, &lastUsed, &t.Enabled, &t.WriteEnabled,
			&plaintext, &key, &pepperID, &hmac, &ips,
		); err != nil {
			return nil, fmt.Errorf("scan token: %w", err)
		}

		t.Version = parseTokenVersion(versionStr)
		if expires.Valid {
			t.Expires = &expires.Time
		}
		if lastUsed.Valid {
			t.LastUsed = &lastUsed.Time
		}
		if plaintext.Valid {
			val := plaintext.String
			t.Plaintext = &val
		}
		if key.Valid {
			val := key.String
			t.Key = &val
		}
		if pepperID.Valid {
			pid := uint16(pepperID.Int16)
			t.PepperID = &pid
		}
		if hmac.Valid {
			val := hmac.String
			t.HMACDigest = &val
		}
		t.AllowedIPs = parseIPNets(ips)

		result = append(result, &t)
	}

	return result, nil
}

func (r *AccountRepositoryPostgres) Get(ctx context.Context, id types.ID, userID types.ID) (*account_entity.UserToken, error) {
	query := `
        SELECT id, version, user_id, description, created, expires, last_used, enabled, write_enabled,
               plaintext, key, pepper_id, hmac_digest, allowed_ips
        FROM users_tokens
        WHERE id = $1 AND user_id = $2`

	var t account_entity.UserToken
	var versionStr string
	var expires, lastUsed sql.NullTime
	var plaintext, key, hmac sql.NullString
	var pepperID sql.NullInt16
	var ips pq.StringArray

	err := r.db.QueryRowContext(ctx, query, id.String(), userID.String()).Scan(
		&t.ID, &versionStr, &t.UserID, &t.Description, &t.Created,
		&expires, &lastUsed, &t.Enabled, &t.WriteEnabled,
		&plaintext, &key, &pepperID, &hmac, &ips,
	)
	if err == sql.ErrNoRows {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get token: %w", err)
	}

	t.Version = parseTokenVersion(versionStr)
	if expires.Valid {
		t.Expires = &expires.Time
	}
	if lastUsed.Valid {
		t.LastUsed = &lastUsed.Time
	}
	if plaintext.Valid {
		val := plaintext.String
		t.Plaintext = &val
	}
	if key.Valid {
		val := key.String
		t.Key = &val
	}
	if pepperID.Valid {
		pid := uint16(pepperID.Int16)
		t.PepperID = &pid
	}
	if hmac.Valid {
		val := hmac.String
		t.HMACDigest = &val
	}
	t.AllowedIPs = parseIPNets(ips)

	return &t, nil
}

func (r *AccountRepositoryPostgres) Create(ctx context.Context, token *account_entity.UserToken) error {
	if err := token.Validate(); err != nil {
		return err
	}
	query := `
        INSERT INTO users_tokens (
            id, version, user_id, description, created, expires, last_used, enabled, write_enabled,
            plaintext, key, pepper_id, hmac_digest, allowed_ips
        ) VALUES ($1, $2, $3, $4, COALESCE($5, NOW()), $6, $7, $8, $9,
                  $10, $11, $12, $13, $14)
        RETURNING id`

	if token.ID.String() == "" {
		token.ID = types.NewID()
	}
	var expires, lastUsed interface{}
	if token.Expires != nil {
		expires = *token.Expires
	}
	if token.LastUsed != nil {
		lastUsed = *token.LastUsed
	}

	var plaintext, key, hmac interface{}
	if token.Plaintext != nil {
		plaintext = *token.Plaintext
	}
	if token.Key != nil {
		key = *token.Key
	}
	if token.HMACDigest != nil {
		hmac = *token.HMACDigest
	}

	var pepper interface{}
	if token.PepperID != nil {
		pepper = *token.PepperID
	}

	ips := stringifyIPNets(token.AllowedIPs)

	err := r.db.QueryRowContext(ctx, query,
		token.ID,
		formatTokenVersion(token.Version),
		token.UserID,
		token.Description,
		nullableTime(token.Created),
		expires,
		lastUsed,
		token.Enabled,
		token.WriteEnabled,
		plaintext,
		key,
		pepper,
		hmac,
		pq.Array(ips),
	).Scan(&token.ID)
	if err != nil {
		return fmt.Errorf("create token: %w", err)
	}
	return nil
}

func (r *AccountRepositoryPostgres) Update(ctx context.Context, token *account_entity.UserToken) error {
	if err := token.Validate(); err != nil {
		return err
	}
	query := `
        UPDATE users_tokens
        SET version = $1, description = $2, expires = $3, last_used = $4,
            enabled = $5, write_enabled = $6, plaintext = $7, key = $8,
            pepper_id = $9, hmac_digest = $10, allowed_ips = $11
        WHERE id = $12 AND user_id = $13`

	var expires, lastUsed interface{}
	if token.Expires != nil {
		expires = *token.Expires
	}
	if token.LastUsed != nil {
		lastUsed = *token.LastUsed
	}
	var plaintext, key, hmac interface{}
	if token.Plaintext != nil {
		plaintext = *token.Plaintext
	}
	if token.Key != nil {
		key = *token.Key
	}
	if token.HMACDigest != nil {
		hmac = *token.HMACDigest
	}
	var pepper interface{}
	if token.PepperID != nil {
		pepper = *token.PepperID
	}
	ips := stringifyIPNets(token.AllowedIPs)

	res, err := r.db.ExecContext(ctx, query,
		formatTokenVersion(token.Version), token.Description, expires, lastUsed,
		token.Enabled, token.WriteEnabled, plaintext, key, pepper, hmac, pq.Array(ips),
		token.ID, token.UserID,
	)
	if err != nil {
		return fmt.Errorf("update token: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func (r *AccountRepositoryPostgres) Delete(ctx context.Context, id types.ID, userID types.ID) error {
	query := `DELETE FROM users_tokens WHERE id = $1 AND user_id = $2`
	res, err := r.db.ExecContext(ctx, query, id.String(), userID.String())
	if err != nil {
		return fmt.Errorf("delete token: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// --- BookmarkRepository ---

func (r *AccountRepositoryPostgres) ListBookmarksByUser(ctx context.Context, userID types.ID) ([]*account_entity.Bookmark, error) {
	query := `SELECT id, user_id, title, url, created FROM account_bookmarks WHERE user_id = $1 ORDER BY created DESC`
	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("list bookmarks: %w", err)
	}
	defer rows.Close()
	var res []*account_entity.Bookmark
	for rows.Next() {
		var b account_entity.Bookmark
		if err := rows.Scan(&b.ID, &b.UserID, &b.Title, &b.URL, &b.Created); err != nil {
			return nil, fmt.Errorf("scan bookmark: %w", err)
		}
		res = append(res, &b)
	}
	return res, nil
}

// --- NotificationRepository ---

func (r *AccountRepositoryPostgres) ListNotificationsByUser(ctx context.Context, userID types.ID) ([]*account_entity.Notification, error) {
	query := `SELECT id, user_id, title, message, level, created, read FROM account_notifications WHERE user_id = $1 ORDER BY created DESC`
	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("list notifications: %w", err)
	}
	defer rows.Close()
	var res []*account_entity.Notification
	for rows.Next() {
		var n account_entity.Notification
		if err := rows.Scan(&n.ID, &n.UserID, &n.Title, &n.Message, &n.Level, &n.Created, &n.Read); err != nil {
			return nil, fmt.Errorf("scan notification: %w", err)
		}
		res = append(res, &n)
	}
	return res, nil
}

// --- SubscriptionRepository ---

func (r *AccountRepositoryPostgres) ListSubscriptionsByUser(ctx context.Context, userID types.ID) ([]*account_entity.Subscription, error) {
	query := `SELECT id, user_id, object_id, object_type, created FROM account_subscriptions WHERE user_id = $1 ORDER BY created DESC`
	rows, err := r.db.QueryContext(ctx, query, userID.String())
	if err != nil {
		return nil, fmt.Errorf("list subscriptions: %w", err)
	}
	defer rows.Close()
	var res []*account_entity.Subscription
	for rows.Next() {
		var s account_entity.Subscription
		if err := rows.Scan(&s.ID, &s.UserID, &s.ObjectID, &s.ObjectType, &s.Created); err != nil {
			return nil, fmt.Errorf("scan subscription: %w", err)
		}
		res = append(res, &s)
	}
	return res, nil
}

// --- UserConfigRepository ---

func (r *AccountRepositoryPostgres) GetByUser(ctx context.Context, userID types.ID) (*account_entity.UserConfig, error) {
	query := `SELECT user_id, data FROM account_user_configs WHERE user_id = $1`
	var cfg account_entity.UserConfig
	if err := r.db.QueryRowContext(ctx, query, userID.String()).Scan(&cfg.UserID, &cfg.Data); err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("get user config: %w", err)
	}
	return &cfg, nil
}

func (r *AccountRepositoryPostgres) Upsert(ctx context.Context, config *account_entity.UserConfig) error {
	if err := config.Validate(); err != nil {
		return err
	}
	query := `
        INSERT INTO account_user_configs (user_id, data, updated)
        VALUES ($1, $2, NOW())
        ON CONFLICT (user_id) DO UPDATE SET data = EXCLUDED.data, updated = NOW()`
	_, err := r.db.ExecContext(ctx, query, config.UserID.String(), config.Data)
	if err != nil {
		return fmt.Errorf("upsert user config: %w", err)
	}
	return nil
}

// --- helpers ---

func parseTokenVersion(s string) account_entity.TokenVersion {
	switch s {
	case "v1":
		return account_entity.TokenVersionV1
	case "v2":
		return account_entity.TokenVersionV2
	default:
		return account_entity.TokenVersionV2
	}
}

func formatTokenVersion(v account_entity.TokenVersion) string {
	if v == account_entity.TokenVersionV1 {
		return "v1"
	}
	return "v2"
}

func parseIPNets(arr []string) []*net.IPNet {
	if len(arr) == 0 {
		return nil
	}
	var res []*net.IPNet
	for _, s := range arr {
		if _, ipnet, err := net.ParseCIDR(s); err == nil {
			res = append(res, ipnet)
		}
	}
	return res
}

func stringifyIPNets(nets []*net.IPNet) []string {
	if nets == nil {
		return nil
	}
	res := make([]string, 0, len(nets))
	for _, n := range nets {
		res = append(res, n.String())
	}
	return res
}

func nullableTime(t time.Time) interface{} {
	if t.IsZero() {
		return nil
	}
	return t
}
