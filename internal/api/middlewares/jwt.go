package middlewares

import (
    "net/http"
    "ppdb-backend/internal/core/services"
    
    "github.com/labstack/echo/v4"
)

func JWTMiddleware(authService services.AuthService) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            token := c.Request().Header.Get("Authorization")
            if token == "" {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "no token provided",
                })
            }

            // Remove 'Bearer ' prefix if exists
            if len(token) > 7 && token[:7] == "Bearer " {
                token = token[7:]
            }

            // Validate token
            jwtToken, err := authService.ValidateToken(token)
            if err != nil || !jwtToken.Valid {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "invalid token",
                })
            }

            // Set user claims in context
            c.Set("user", jwtToken.Claims)
            return next(c)
        }
    }
}