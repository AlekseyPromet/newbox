// Package middleware содержит middleware для аутентификации (Kerberos/SSO stub)
package middleware

import (
	"net/http"
	"strings"

	"github.com/AlekseyPromet/netbox_go/pkg/types"
	"github.com/labstack/echo/v4"
)

// KerberosSSOMiddleware заглушка для интеграции с AD/Kerberos (SPNEGO).
// В бою следует заменить на полноценный SPNEGO (например, gokrb5 + reverse-proxy).
// Поведение:
// 1) Пытается прочитать заголовок X-User-ID (UUID) — как от reverse-proxy после успешного Kerberos.
// 2) Пытается разобрать user principal из заголовка X-Remote-User (DOMAIN\\user или user@realm) и конвертировать в UUID невозможно — поэтому возвращает 401 без маппинга.
// 3) Кладёт types.ID в контекст (c.Set("userID", id)).
func KerberosSSOMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Приоритет: UUID из прокси
			if userID := strings.TrimSpace(c.Request().Header.Get("X-User-ID")); userID != "" {
				if id, err := types.ParseID(userID); err == nil {
					c.Set("userID", id)
					return next(c)
				}
				return echo.NewHTTPError(http.StatusUnauthorized, "invalid X-User-ID")
			}

			// Kerberos principal passthrough (если прокси кладёт X-Remote-User)
			if principal := strings.TrimSpace(c.Request().Header.Get("X-Remote-User")); principal != "" {
				// Здесь должен быть маппинг principal -> UUID пользователя (LDAP/AD lookup)
				// Пока возвращаем 401, чтобы явно указать на отсутствие маппинга
				return echo.NewHTTPError(http.StatusUnauthorized, "principal mapping required")
			}

			return echo.NewHTTPError(http.StatusUnauthorized, "missing Kerberos identity")
		}
	}
}
