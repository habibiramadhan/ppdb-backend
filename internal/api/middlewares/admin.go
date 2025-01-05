// internal/api/middlewares/admin.go
package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func AdminMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(jwt.MapClaims)
			role := user["role"].(string)

			if role != "admin" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "Access denied: admin role required",
				})
			}

			return next(c)
		}
	}
}
