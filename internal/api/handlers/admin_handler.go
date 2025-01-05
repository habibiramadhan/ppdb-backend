// internal/api/handlers/admin_handler.go
package handlers

import (
	"net/http"
	"ppdb-backend/internal/core/services"
	"ppdb-backend/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminService services.AdminService
}

func NewAdminHandler(adminService services.AdminService) *AdminHandler {
	return &AdminHandler{adminService}
}

func (h *AdminHandler) GetAllUsers(c echo.Context) error {
	page := utils.GetPageFromQuery(c)
	limit := utils.GetLimitFromQuery(c)
	users, pagination, err := h.adminService.GetAllUsers(page, limit)
	
	if err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get users", err.Error())
	}

	return utils.PaginationSuccessResponse(c, "Users retrieved successfully", users, pagination)
}

func (h *AdminHandler) GetUserByID(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
	}

	user, err := h.adminService.GetUserByID(id)
	if err != nil {
		return utils.ErrorResponse(c, http.StatusNotFound, "User not found", err.Error())
	}

	return utils.SuccessResponse(c, "User retrieved successfully", user)
}

func (h *AdminHandler) UpdateUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
	}

	var input services.UpdateUserInput
	if err := c.Bind(&input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(input); err != nil {
		return utils.ValidationErrorResponse(c, "Validation failed", err.Error())
	}

	if err := h.adminService.UpdateUser(id, input); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user", err.Error())
	}

	return utils.SuccessResponse(c, "User updated successfully", nil)
}

func (h *AdminHandler) DeleteUser(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
	}

	if err := h.adminService.DeleteUser(id); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete user", err.Error())
	}

	return utils.SuccessResponse(c, "User deleted successfully", nil)
}

func (h *AdminHandler) UpdateUserStatus(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid user ID", err.Error())
	}

	var input services.UpdateStatusInput
	if err := c.Bind(&input); err != nil {
		return utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := c.Validate(input); err != nil {
		return utils.ValidationErrorResponse(c, "Validation failed", err.Error())
	}

	if err := h.adminService.UpdateUserStatus(id, input.Status); err != nil {
		return utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update user status", err.Error())
	}

	return utils.SuccessResponse(c, "User status updated successfully", nil)
}