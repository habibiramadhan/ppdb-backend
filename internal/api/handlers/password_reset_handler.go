// internal/api/handlers/password_reset_handler.go
package handlers

import (
    "net/http"
    "ppdb-backend/internal/core/services"
    "ppdb-backend/utils"

    "github.com/labstack/echo/v4"
)

type PasswordResetHandler struct {
    passwordResetService services.PasswordResetService
}

func NewPasswordResetHandler(passwordResetService services.PasswordResetService) *PasswordResetHandler {
    return &PasswordResetHandler{passwordResetService}
}

func (h *PasswordResetHandler) RequestReset(c echo.Context) error {
    var input struct {
        Email string `json:"email" validate:"required,email"`
    }

    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
    }

    if err := c.Validate(input); err != nil {
        return utils.ValidationErrorResponse(c, "Validation failed", err.Error())
    }

    if err := h.passwordResetService.RequestReset(input.Email); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to request password reset", err.Error())
    }

    return utils.SuccessResponse(c, "Password reset email sent successfully", nil)
}

func (h *PasswordResetHandler) ValidateToken(c echo.Context) error {
    token := c.QueryParam("token")
    if token == "" {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Token is required", nil)
    }

    if err := h.passwordResetService.ValidateToken(token); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid token", err.Error())
    }

    return utils.SuccessResponse(c, "Token is valid", nil)
}

func (h *PasswordResetHandler) ResetPassword(c echo.Context) error {
    var input struct {
        Token       string `json:"token" validate:"required"`
        NewPassword string `json:"new_password" validate:"required,min=6"`
    }

    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
    }

    if err := c.Validate(input); err != nil {
        return utils.ValidationErrorResponse(c, "Validation failed", err.Error())
    }

    if err := h.passwordResetService.ResetPassword(input.Token, input.NewPassword); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to reset password", err.Error())
    }

    return utils.SuccessResponse(c, "Password reset successful", nil)
}