// internal/api/handlers/verification_handler.go
package handlers

import (
	"net/http"
	"ppdb-backend/internal/core/services"
	"ppdb-backend/utils"

	"github.com/labstack/echo/v4"
)

type VerificationHandler struct {
	verificationService services.VerificationService
}

func NewVerificationHandler(verificationService services.VerificationService) *VerificationHandler {
	return &VerificationHandler{verificationService}
}

func (h *VerificationHandler) VerifyEmail(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Token is required", nil)
	}

	err := h.verificationService.VerifyEmail(token)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Verification failed", err.Error())
	}

	return utils.SuccessResponse(c, "Email verified successfully", nil)
}

func (h *VerificationHandler) ResendVerification(c echo.Context) error {
	var input struct {
		Email string `json:"email" validate:"required,email"`
	}

	if err := c.Bind(&input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(input); err != nil {
		return utils.ValidationErrorResponse(c, "Validation failed", err.Error())
	}

	err := h.verificationService.ResendVerification(input.Email)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Failed to resend verification", err.Error())
	}

	return utils.SuccessResponse(c, "Verification email sent successfully", nil)
}