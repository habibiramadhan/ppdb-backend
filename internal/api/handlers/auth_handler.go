package handlers

import (
    "net/http"
    "ppdb-backend/internal/core/services"
    
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
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    if err := h.authService.Register(input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusCreated, map[string]string{
        "message": "Registration successful",
    })
}

func (h *AuthHandler) Login(c echo.Context) error {
    var input services.LoginInput
    if err := c.Bind(&input); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{
            "error": err.Error(),
        })
    }

    token, err := h.authService.Login(input)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{
            "error": err.Error(),
        })
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "status": "success",
        "data": map[string]string{
            "access_token": token,
            "token_type": "Bearer",
        },
        "message": "Login successful",
    })
}