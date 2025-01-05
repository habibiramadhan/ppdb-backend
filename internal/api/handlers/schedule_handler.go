// internal/api/handlers/schedule_handler.go
package handlers

import (
    "net/http"
    "ppdb-backend/internal/core/services"
    "ppdb-backend/utils"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v4"
)

type ScheduleHandler struct {
    scheduleService services.ScheduleService
}

func NewScheduleHandler(scheduleService services.ScheduleService) *ScheduleHandler {
    return &ScheduleHandler{scheduleService}
}

func (h *ScheduleHandler) Create(c echo.Context) error {
    var input services.CreateScheduleInput
    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid", err.Error())
    }

    if err := c.Validate(input); err != nil {
        return utils.ValidationErrorResponse(c, "Data wajib ada yang belum diisi", err.Error())
    }

    user := c.Get("user").(jwt.MapClaims)
    userID, err := uuid.Parse(user["user_id"].(string))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "User ID ga valid", err.Error())
    }

    if err := h.scheduleService.Create(input, userID); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal bikin jadwal", err.Error())
    }

    return utils.CreatedResponse(c, "Sukses bikin jadwal baru", nil)
}

func (h *ScheduleHandler) GetAll(c echo.Context) error {
    page := utils.GetPageFromQuery(c)
    limit := utils.GetLimitFromQuery(c)

    schedules, pagination, err := h.scheduleService.GetAll(page, limit)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil data jadwal", err.Error())
    }

    return utils.PaginationSuccessResponse(c, "Sukses ambil data jadwal", schedules, pagination)
}

func (h *ScheduleHandler) GetByAcademicYear(c echo.Context) error {
    yearID, err := uuid.Parse(c.Param("yearId"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID tahun ajaran ga valid", err.Error())
    }

    page := utils.GetPageFromQuery(c)
    limit := utils.GetLimitFromQuery(c)

    schedules, pagination, err := h.scheduleService.GetByAcademicYear(yearID, page, limit)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil data jadwal", err.Error())
    }

    return utils.PaginationSuccessResponse(c, "Sukses ambil data jadwal", schedules, pagination)
}

func (h *ScheduleHandler) GetByID(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jadwal ga valid", err.Error())
    }

    schedule, err := h.scheduleService.GetByID(id)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusNotFound, "Jadwal ga ketemu", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses ambil detail jadwal", schedule)
}

func (h *ScheduleHandler) Update(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jadwal ga valid", err.Error())
    }

    var input services.UpdateScheduleInput
    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid", err.Error())
    }

    if err := c.Validate(input); err != nil {
        return utils.ValidationErrorResponse(c, "Data wajib ada yang belum diisi", err.Error())
    }

    user := c.Get("user").(jwt.MapClaims)
    userID, err := uuid.Parse(user["user_id"].(string))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "User ID ga valid", err.Error())
    }

    if err := h.scheduleService.Update(id, input, userID); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal update jadwal", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses update jadwal", nil)
}

func (h *ScheduleHandler) Delete(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jadwal ga valid", err.Error())
    }

    if err := h.scheduleService.Delete(id); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal hapus jadwal", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses hapus jadwal", nil)
}

func (h *ScheduleHandler) SetStatus(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jadwal ga valid", err.Error())
    }

    var input struct {
        IsActive bool `json:"is_active"`
    }

    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid", err.Error())
    }

    if err := h.scheduleService.SetStatus(id, input.IsActive); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal update status jadwal", err.Error())
    }

    status := "non-aktif"
    if input.IsActive {
        status = "aktif"
    }

    return utils.SuccessResponse(c, "Sukses update status jadwal jadi "+status, nil)
}

func (h *ScheduleHandler) GetUpcoming(c echo.Context) error {
    schedules, err := h.scheduleService.GetUpcomingSchedules(10)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil jadwal", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses ambil jadwal yang akan datang", schedules)
}