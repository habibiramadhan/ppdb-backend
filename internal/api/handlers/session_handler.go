// internal/api/handlers/session_handler.go
package handlers

import (
    "net/http"
    "ppdb-backend/internal/core/services"
    "ppdb-backend/utils"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
)

type SessionHandler struct {
    sessionService services.SessionService
}

func NewSessionHandler(sessionService services.SessionService) *SessionHandler {
    return &SessionHandler{sessionService}
}

func (h *SessionHandler) GetActiveSessions(c echo.Context) error {
    user := c.Get("user").(jwt.MapClaims)
    userID, err := uuid.Parse(user["user_id"].(string))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
    }

    sessions, err := h.sessionService.GetActiveSessions(userID)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get sessions", err.Error())
    }

    return utils.SuccessResponse(c, "Active sessions retrieved successfully", sessions)
}

func (h *SessionHandler) RevokeSession(c echo.Context) error {
    sessionID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid session ID", err.Error())
    }

    if err := h.sessionService.RevokeSession(sessionID); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to revoke session", err.Error())
    }

    return utils.SuccessResponse(c, "Session revoked successfully", nil)
}

func (h *SessionHandler) RevokeAllSessions(c echo.Context) error {
    user := c.Get("user").(jwt.MapClaims)
    userID, err := uuid.Parse(user["user_id"].(string))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
    }

    if err := h.sessionService.RevokeAllSessions(userID); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to revoke all sessions", err.Error())
    }

    return utils.SuccessResponse(c, "All sessions revoked successfully", nil)
}