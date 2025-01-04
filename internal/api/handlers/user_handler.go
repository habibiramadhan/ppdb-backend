// internal/api/handlers/user_handler.go
package handlers

import (
    "net/http"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v4"
)

func GetProfile(c echo.Context) error {
    user := c.Get("user").(jwt.MapClaims)
    
    return c.JSON(http.StatusOK, map[string]interface{}{
        "user_id": user["user_id"],
        "email":   user["email"],
        "role":    user["role"],
    })
}