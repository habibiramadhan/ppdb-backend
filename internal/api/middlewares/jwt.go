package middlewares

import (
    "net/http"
    "strings"
    "ppdb-backend/internal/core/services"
    
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(authService services.AuthService) echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            authHeader := c.Request().Header.Get("Authorization")
            
            // Check if Authorization header exists
            if authHeader == "" {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "Authorization header required",
                })
            }

            // Check Bearer scheme
            if !strings.HasPrefix(authHeader, "Bearer ") {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "Invalid authorization format. Format is Authorization: Bearer [token]",
                })
            }

            // Extract token from Bearer prefix
            tokenString := strings.TrimPrefix(authHeader, "Bearer ")
            if tokenString == "" {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "Token not provided",
                })
            }

            // Validate token
            token, err := authService.ValidateToken(tokenString)
            if err != nil {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "Invalid or expired token",
                })
            }

            // Get claims from token
            claims, ok := token.Claims.(jwt.MapClaims)
            if !ok || !token.Valid {
                return c.JSON(http.StatusUnauthorized, map[string]string{
                    "error": "Invalid token claims",
                })
            }

            // Set claims in context
            c.Set("user", claims)
            return next(c)
        }
    }
}