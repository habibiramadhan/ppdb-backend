// cmd/main.go
package main

import (
    "log"
    "ppdb-backend/config"
    "ppdb-backend/internal/api/routes"
    "ppdb-backend/utils"
    "os"
    
    "github.com/joho/godotenv"
    "github.com/labstack/echo/v4"
    "github.com/labstack/echo/v4/middleware"
)

func main() {
    // Load .env file
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    if err := utils.SetupUploadDir(); err != nil {
        log.Fatal("Gagal setup folder upload:", err)
    }
    
    // Initialize config
    cfg := config.NewConfig()

    // Echo instance
    e := echo.New()

    // Set custom validator
    e.Validator = utils.NewValidator()

    // Middleware
    e.Use(middleware.Logger())
    e.Use(middleware.Recover())
    e.Use(middleware.CORS())

    
    // Setup Routes
    routes.Setup(e, cfg)

    // Start server
    e.Logger.Fatal(e.Start(":" + os.Getenv("APP_PORT")))
}