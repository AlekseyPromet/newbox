// Package handlers contains HTTP handlers for account module
package handlers

import (
	"net/http"

	account_entity "netbox_go/internal/domain/account/entity"
	"netbox_go/internal/repository"
	"netbox_go/pkg/types"
	"github.com/labstack/echo/v4"
)

// AccountHandler aggregates handlers for profile, preferences, bookmarks, notifications, subscriptions and API tokens
// NOTE: This is a minimal REST projection of netbox/account views. Authentication/authorization is simplified: user is resolved
// from header X-User-ID (UUID). In real deployment integrate with auth middleware/session.
type AccountHandler struct {
	tokenRepo        repository.UserTokenRepository
	bookmarkRepo     repository.BookmarkRepository
	notificationRepo repository.NotificationRepository
	subscriptionRepo repository.SubscriptionRepository
	userConfigRepo   repository.UserConfigRepository
}

// NewAccountHandler builds AccountHandler
func NewAccountHandler(
	tokenRepo repository.UserTokenRepository,
	bookmarkRepo repository.BookmarkRepository,
	notificationRepo repository.NotificationRepository,
	subscriptionRepo repository.SubscriptionRepository,
	userConfigRepo repository.UserConfigRepository,
) *AccountHandler {
	return &AccountHandler{
		tokenRepo:        tokenRepo,
		bookmarkRepo:     bookmarkRepo,
		notificationRepo: notificationRepo,
		subscriptionRepo: subscriptionRepo,
		userConfigRepo:   userConfigRepo,
	}
}

// Profile returns simple payload for current user
func (h *AccountHandler) Profile(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"user_id":    userID.String(),
		"active_tab": "profile",
	})
}

// ListBookmarks returns bookmarks of current user
func (h *AccountHandler) ListBookmarks(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	items, err := h.bookmarkRepo.ListBookmarksByUser(c.Request().Context(), userID)
	if err != nil {
		return handleRepoError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"active_tab": "bookmarks",
		"results":    items,
	})
}

// ListNotifications returns notifications of current user
func (h *AccountHandler) ListNotifications(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	items, err := h.notificationRepo.ListNotificationsByUser(c.Request().Context(), userID)
	if err != nil {
		return handleRepoError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"active_tab": "notifications",
		"results":    items,
	})
}

// ListSubscriptions returns subscriptions of current user
func (h *AccountHandler) ListSubscriptions(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	items, err := h.subscriptionRepo.ListSubscriptionsByUser(c.Request().Context(), userID)
	if err != nil {
		return handleRepoError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"active_tab": "subscriptions",
		"results":    items,
	})
}

// GetPreferences returns user config
func (h *AccountHandler) GetPreferences(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	cfg, err := h.userConfigRepo.GetByUser(c.Request().Context(), userID)
	if err != nil {
		return handleRepoError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"active_tab": "preferences",
		"config":     cfg,
	})
}

// UpsertPreferences updates user config
func (h *AccountHandler) UpsertPreferences(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	var payload account_entity.UserConfig
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	payload.UserID = userID
	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := h.userConfigRepo.Upsert(c.Request().Context(), &payload); err != nil {
		return handleRepoError(err)
	}
	return c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

// ListTokens returns tokens of current user
func (h *AccountHandler) ListTokens(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	tokens, err := h.tokenRepo.ListByUser(c.Request().Context(), userID)
	if err != nil {
		return handleRepoError(err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"active_tab": "api-tokens",
		"results":    tokens,
	})
}

// GetToken returns token by id for current user
func (h *AccountHandler) GetToken(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	tokenID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	token, err := h.tokenRepo.Get(c.Request().Context(), tokenID, userID)
	if err != nil {
		return handleRepoError(err)
	}
	return c.JSON(http.StatusOK, token)
}

// CreateToken creates new token for current user
func (h *AccountHandler) CreateToken(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	var token account_entity.UserToken
	if err := c.Bind(&token); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	token.UserID = userID
	if err := token.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := h.tokenRepo.Create(c.Request().Context(), &token); err != nil {
		return handleRepoError(err)
	}
	return c.JSON(http.StatusCreated, token)
}

// UpdateToken updates existing token for current user
func (h *AccountHandler) UpdateToken(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	tokenID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	existing, err := h.tokenRepo.Get(c.Request().Context(), tokenID, userID)
	if err != nil {
		return handleRepoError(err)
	}
	if err := c.Bind(existing); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}
	existing.ID = tokenID
	existing.UserID = userID
	if err := existing.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := h.tokenRepo.Update(c.Request().Context(), existing); err != nil {
		return handleRepoError(err)
	}
	updated, err := h.tokenRepo.Get(c.Request().Context(), tokenID, userID)
	if err != nil {
		return handleRepoError(err)
	}
	return c.JSON(http.StatusOK, updated)
}

// DeleteToken removes token for current user
func (h *AccountHandler) DeleteToken(c echo.Context) error {
	userID, err := currentUserID(c)
	if err != nil {
		return err
	}
	tokenID, err := parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.tokenRepo.Delete(c.Request().Context(), tokenID, userID); err != nil {
		return handleRepoError(err)
	}
	return c.NoContent(http.StatusNoContent)
}

// --- helpers ---

func handleRepoError(err error) error {
	switch err {
	case repository.ErrNotFound:
		return echo.NewHTTPError(http.StatusNotFound, "not found")
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
}

func parseIDParam(c echo.Context, name string) (types.ID, error) {
	idStr := c.Param(name)
	id, err := types.ParseID(idStr)
	if err != nil {
		return types.ID{}, echo.NewHTTPError(http.StatusBadRequest, "invalid id")
	}
	return id, nil
}

func currentUserID(c echo.Context) (types.ID, error) {
	// Пробуем получить userID, выставленный SSO middleware
	if v := c.Get("userID"); v != nil {
		switch t := v.(type) {
		case types.ID:
			return t, nil
		case string:
			if id, err := types.ParseID(t); err == nil {
				return id, nil
			}
		}
	}

	// fallback: прямой заголовок (например, reverse-proxy после Kerberos/SPNEGO)
	userIDStr := c.Request().Header.Get("X-User-ID")
	if userIDStr == "" {
		return types.ID{}, echo.NewHTTPError(http.StatusUnauthorized, "missing user identity")
	}
	id, err := types.ParseID(userIDStr)
	if err != nil {
		return types.ID{}, echo.NewHTTPError(http.StatusBadRequest, "invalid user identity")
	}
	return id, nil
}
