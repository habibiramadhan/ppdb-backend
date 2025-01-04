package handlers

import (
	"net/http"
	"ppdb-backend/internal/core/services"
	"ppdb-backend/utils"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{authService}
}

func (h *AuthHandler) Register(c echo.Context) error {
	var input services.RegisterInput
	if err := c.Bind(&input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(input); err != nil {
		return utils.ValidationErrorResponse(c, "Validation failed", err.Error())
	}

	if err := h.authService.Register(input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Registration failed", err.Error())
	}

	return utils.CreatedResponse(c, "Registration successful", nil)
}

func (h *AuthHandler) Login(c echo.Context) error {
	var input services.LoginInput
	if err := c.Bind(&input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(input); err != nil {
		return utils.ValidationErrorResponse(c, "Validation failed", err.Error())
	}

	token, err := h.authService.Login(input)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusUnauthorized, "Login failed", err.Error())
	}

	tokenResponse := map[string]string{
		"access_token": token,
		"token_type":   "Bearer",
	}

	return utils.SuccessResponse(c, "Login successful", tokenResponse)
}

func (h *AuthHandler) GetProfile(c echo.Context) error {
	user := c.Get("user").(jwt.MapClaims)

	userData := map[string]interface{}{
		"user_id": user["user_id"],
		"email":   user["email"],
		"role":    user["role"],
		"name":    user["name"],
	}

	return utils.SuccessResponse(c, "Profile retrieved successfully", userData)
}
