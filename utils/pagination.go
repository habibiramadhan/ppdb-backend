// utils/pagination.go
package utils

import (
	"strconv"

	"github.com/labstack/echo/v4"
)

const (
	DefaultPage  = 1
	DefaultLimit = 10
)

func GetPageFromQuery(c echo.Context) int {
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil || page < 1 {
		return DefaultPage
	}
	return page
}

func GetLimitFromQuery(c echo.Context) int {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < 1 {
		return DefaultLimit
	}
	return limit
}