// Package entity содержит сущности домена Users/Account
package entity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net"
	"time"

	"github.com/AlekseyPromet/netbox_go/pkg/types"
)

// TokenVersion представляет версию токена
type TokenVersion string

const (
	TokenVersionV1 TokenVersion = "v1"
	TokenVersionV2 TokenVersion = "v2"
)

// Token представляет API токен пользователя для аутентификации
type Token struct {
	ID           types.ID     `json:"id"`
	Version      TokenVersion `json:"version"`
	UserID       types.ID     `json:"user_id"`
	Description  string       `json:"description,omitempty"`
	Created      time.Time    `json:"created"`
	Expires      *time.Time   `json:"expires,omitempty"`
	LastUsed     *time.Time   `json:"last_used,omitempty"`
	Enabled      bool         `json:"enabled"`
	WriteEnabled bool         `json:"write_enabled"`

	// Поля для v1 токенов (plaintext)
	Plaintext *string `json:"plaintext,omitempty"`

	// Поля для v2 токенов
	Key        *string `json:"key,omitempty"`
	PepperID   *uint16 `json:"pepper_id,omitempty"`
	HMACDigest *string `json:"hmac_digest,omitempty"`

	// Allowed IP addresses (CIDR notation)
	AllowedIPs []*net.IPNet `json:"allowed_ips,omitempty"`
}

// Validate проверяет корректность токена
func (t *Token) Validate() error {
	// Проверка версии
	if t.Version != TokenVersionV1 && t.Version != TokenVersionV2 {
		return errors.New("invalid token version")
	}

	// Проверка обязательных полей
	if t.UserID.String() == "" {
		return errors.New("user_id is required")
	}

	// Проверка целостности полей в зависимости от версии
	if t.Version == TokenVersionV1 {
		if t.Plaintext == nil || len(*t.Plaintext) != 40 {
			return errors.New("v1 token requires 40-character plaintext")
		}
		if t.Key != nil || t.PepperID != nil || t.HMACDigest != nil {
			return errors.New("v1 token should not have key, pepper_id, or hmac_digest")
		}
	} else if t.Version == TokenVersionV2 {
		if t.Key == nil || len(*t.Key) != 12 {
			return errors.New("v2 token requires 12-character key")
		}
		if t.PepperID == nil {
			return errors.New("v2 token requires pepper_id")
		}
		if t.HMACDigest == nil || len(*t.HMACDigest) != 64 {
			return errors.New("v2 token requires 64-character hmac_digest")
		}
		if t.Plaintext != nil {
			return errors.New("v2 token should not have plaintext")
		}
	}

	// Проверка description length
	if len(t.Description) > 200 {
		return errors.New("description too long (max 200 characters)")
	}

	return nil
}

// IsExpired проверяет, истек ли срок действия токена
func (t *Token) IsExpired() bool {
	if t.Expires == nil {
		return false
	}
	return time.Now().After(*t.Expires)
}

// IsActive проверяет, активен ли токен (enabled и не истек)
func (t *Token) IsActive() bool {
	return t.Enabled && !t.IsExpired()
}

// ValidateClientIP проверяет IP-адрес клиента against allowed_ips
func (t *Token) ValidateClientIP(clientIP net.IP) bool {
	if t.AllowedIPs == nil || len(t.AllowedIPs) == 0 {
		return true // Нет ограничений
	}

	for _, allowedNet := range t.AllowedIPs {
		if allowedNet.Contains(clientIP) {
			return true
		}
	}
	return false
}

// ComputeHMACDigest вычисляет HMAC digest для v2 токена
func (t *Token) ComputeHMACDigest(tokenValue string, pepper string) string {
	h := hmac.New(sha256.New, []byte(pepper))
	h.Write([]byte(tokenValue))
	return hex.EncodeToString(h.Sum(nil))
}

// GetAuthHeaderPrefix возвращает префикс для HTTP Authorization header
func (t *Token) GetAuthHeaderPrefix() string {
	const TOKEN_PREFIX = "NBX-"
	if t.Version == TokenVersionV1 {
		return "Token "
	} else if t.Version == TokenVersionV2 {
		return "Bearer " + TOKEN_PREFIX + *t.Key + "."
	}
	return ""
}

// OwnerGroup представляет группу владельцев объектов
type OwnerGroup struct {
	ID          types.ID        `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Created     time.Time       `json:"created"`
	Updated     time.Time       `json:"updated"`
}

// Validate проверяет корректность группы владельцев
func (og *OwnerGroup) Validate() error {
	if len(og.Name) == 0 || len(og.Name) > 100 {
		return errors.New("name is required and must be <= 100 characters")
	}
	if len(og.Description) > 200 {
		return errors.New("description too long (max 200 characters)")
	}
	return nil
}

// Owner представляет владельца объекта (пользователь или группа)
type Owner struct {
	ID          types.ID   `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	GroupID     *types.ID  `json:"group_id,omitempty"`
	UserIDs     []types.ID `json:"user_ids,omitempty"`
	GroupIDs    []types.ID `json:"group_ids,omitempty"`
	Created     time.Time  `json:"created"`
	Updated     time.Time  `json:"updated"`
}

// Validate проверяет корректность владельца
func (o *Owner) Validate() error {
	if len(o.Name) == 0 || len(o.Name) > 100 {
		return errors.New("name is required and must be <= 100 characters")
	}
	if len(o.Description) > 200 {
		return errors.New("description too long (max 200 characters)")
	}
	return nil
}

// IsGroupOwner проверяет, является ли владелец групповым
func (o *Owner) IsGroupOwner() bool {
	return o.GroupID != nil
}
