package routes

import (
    "ppdb-backend/config"
    
    "github.com/labstack/echo/v4"
)

func Setup(e *echo.Echo, cfg *config.Config) {
    // Health Check
    e.GET("/health", func(c echo.Context) error {
        return c.JSON(200, map[string]string{
            "status": "OK",
        })
    })

    // API v1 group
    v1 := e.Group("/api/v1")

    // Auth routes will be added here
    setupAuthRoutes(v1, cfg)
}

func setupAuthRoutes(g *echo.Group, cfg *config.Config) {
    // Auth routes will be implemented later
}