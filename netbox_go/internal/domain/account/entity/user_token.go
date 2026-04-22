// Package entity содержит сущности домена Account
package entity

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net"
	"time"

	"netbox_go/pkg/types"
)

// TokenVersion представляет версию токена
type TokenVersion uint16

const (
	TokenVersionV1 TokenVersion = 1
	TokenVersionV2 TokenVersion = 2
)

// UserToken представляет прокси-модель для управления пользователями своими API токенами
// Это расширенная версия Token с методами для account модуля
type UserToken struct {
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

// Validate проверяет корректность пользовательского токена
func (ut *UserToken) Validate() error {
	// Проверка версии
	if ut.Version != TokenVersionV1 && ut.Version != TokenVersionV2 {
		return errors.New("invalid token version")
	}

	// Проверка обязательных полей
	if ut.UserID.String() == "" {
		return errors.New("user_id is required")
	}

	// Проверка целостности полей в зависимости от версии
	if ut.Version == TokenVersionV1 {
		if ut.Plaintext == nil || len(*ut.Plaintext) != 40 {
			return errors.New("v1 token requires 40-character plaintext")
		}
		if ut.Key != nil || ut.PepperID != nil || ut.HMACDigest != nil {
			return errors.New("v1 token should not have key, pepper_id, or hmac_digest")
		}
	} else if ut.Version == TokenVersionV2 {
		if ut.Key == nil || len(*ut.Key) != 12 {
			return errors.New("v2 token requires 12-character key")
		}
		if ut.PepperID == nil {
			return errors.New("v2 token requires pepper_id")
		}
		if ut.HMACDigest == nil || len(*ut.HMACDigest) != 64 {
			return errors.New("v2 token requires 64-character hmac_digest")
		}
		if ut.Plaintext != nil {
			return errors.New("v2 token should not have plaintext")
		}
	}

	// Проверка description length
	if len(ut.Description) > 200 {
		return errors.New("description too long (max 200 characters)")
	}

	return nil
}

// IsExpired проверяет, истек ли срок действия токена
func (ut *UserToken) IsExpired() bool {
	if ut.Expires == nil {
		return false
	}
	return time.Now().After(*ut.Expires)
}

// IsActive проверяет, активен ли токен (enabled и не истек)
func (ut *UserToken) IsActive() bool {
	return ut.Enabled && !ut.IsExpired()
}

// ValidateClientIP проверяет IP-адрес клиента against allowed_ips
func (ut *UserToken) ValidateClientIP(clientIP net.IP) bool {
	if ut.AllowedIPs == nil || len(ut.AllowedIPs) == 0 {
		return true // Нет ограничений
	}

	for _, allowedNet := range ut.AllowedIPs {
		if allowedNet.Contains(clientIP) {
			return true
		}
	}
	return false
}

// ComputeHMACDigest вычисляет HMAC digest для v2 токена
func (ut *UserToken) ComputeHMACDigest(tokenValue string, pepper string) string {
	h := hmac.New(sha256.New, []byte(pepper))
	h.Write([]byte(tokenValue))
	return hex.EncodeToString(h.Sum(nil))
}

// GetAuthHeaderPrefix возвращает префикс для HTTP Authorization header
func (ut *UserToken) GetAuthHeaderPrefix() string {
	const TOKEN_PREFIX = "NBX-"
	if ut.Version == TokenVersionV1 {
		return "Token "
	} else if ut.Version == TokenVersionV2 && ut.Key != nil {
		return "Bearer " + TOKEN_PREFIX + *ut.Key + "."
	}
	return ""
}

// IsV1 проверяет, является ли токен версией v1
func (ut *UserToken) IsV1() bool {
	return ut.Version == TokenVersionV1
}

// IsV2 проверяет, является ли токен версией v2
func (ut *UserToken) IsV2() bool {
	return ut.Version == TokenVersionV2
}

// GetPartial возвращает частичное представление токена (для v1)
func (ut *UserToken) GetPartial() string {
	if ut.Plaintext != nil && len(*ut.Plaintext) >= 6 {
		return "**********************************" + (*ut.Plaintext)[len(*ut.Plaintext)-6:]
	}
	return ""
}

// ValidateToken валидирует предоставленный plaintext токен против сохранённого
func (ut *UserToken) ValidateToken(token string, pepper string) bool {
	if ut.IsV1() {
		if ut.Plaintext == nil {
			return false
		}
		return token == *ut.Plaintext
	}
	if ut.IsV2() {
		const TOKEN_PREFIX = "NBX-"
		token = tokenRemovePrefix(token, TOKEN_PREFIX)
		if ut.PepperID == nil || ut.HMACDigest == nil || pepper == "" {
			return false
		}
		digest := ut.ComputeHMACDigest(token, pepper)
		return digest == *ut.HMACDigest
	}
	return false
}

func tokenRemovePrefix(s, prefix string) string {
	if len(s) >= len(prefix) && s[:len(prefix)] == prefix {
		return s[len(prefix):]
	}
	return s
}
