//utils/response.go
package utils

import (
    "net/http"
    "github.com/labstack/echo/v4"
)

type Response struct {
    Status  bool        `json:"status"`
    Message string      `json:"message"`
    Errors  interface{} `json:"errors,omitempty"`
    Data    interface{} `json:"data,omitempty"`
}

type PaginationResponse struct {
    Status     bool        `json:"status"`
    Message    string      `json:"message"`
    Data       interface{} `json:"data,omitempty"`
    Pagination interface{} `json:"pagination,omitempty"`
    Errors     interface{} `json:"errors,omitempty"`
}

type PaginationMeta struct {
    Page      int   `json:"page"`
    Limit     int   `json:"limit"`
    TotalData int64 `json:"total_data"`
    TotalPage int   `json:"total_page"`
}

func SuccessResponse(c echo.Context, message string, data interface{}) error {
    return c.JSON(http.StatusOK, Response{
        Status:  true,
        Message: message,
        Data:    data,
    })
}

func CreatedResponse(c echo.Context, message string, data interface{}) error {
    return c.JSON(http.StatusCreated, Response{
        Status:  true,
        Message: message,
        Data:    data,
    })
}

func ErrorResponse(c echo.Context, statusCode int, message string, err interface{}) error {
    return c.JSON(statusCode, Response{
        Status:  false,
        Message: message,
        Errors:  err,
    })
}

func ValidationErrorResponse(c echo.Context, message string, err interface{}) error {
    return c.JSON(http.StatusUnprocessableEntity, Response{
        Status:  false,
        Message: message,
        Errors:  err,
    })
}

func PaginationSuccessResponse(c echo.Context, message string, data interface{}, pagination interface{}) error {
    return c.JSON(http.StatusOK, PaginationResponse{
        Status:     true,
        Message:    message,
        Data:       data,
        Pagination: pagination,
    })
}