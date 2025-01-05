// internal/api/handlers/major_quota_handler.go
package handlers

import (
    "net/http"
    "ppdb-backend/internal/core/services"
    "ppdb-backend/utils"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
    "github.com/golang-jwt/jwt/v4"
)

type MajorQuotaHandler struct {
    majorQuotaService services.MajorQuotaService
}

func NewMajorQuotaHandler(majorQuotaService services.MajorQuotaService) *MajorQuotaHandler {
    return &MajorQuotaHandler{majorQuotaService}
}

func (h *MajorQuotaHandler) Create(c echo.Context) error {
    var input services.CreateMajorQuotaInput
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

    if err := h.majorQuotaService.Create(input, userID); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal bikin kuota", err.Error())
    }

    return utils.CreatedResponse(c, "Sukses bikin kuota jurusan", nil)
}

func (h *MajorQuotaHandler) GetAll(c echo.Context) error {
    page := utils.GetPageFromQuery(c)
    limit := utils.GetLimitFromQuery(c)

    quotas, pagination, err := h.majorQuotaService.GetAll(page, limit)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil data kuota", err.Error())
    }

    return utils.PaginationSuccessResponse(c, "Sukses ambil data kuota", quotas, pagination)
}

func (h *MajorQuotaHandler) GetByID(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID kuota ga valid", err.Error())
    }

    quota, err := h.majorQuotaService.GetByID(id)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusNotFound, "Kuota ga ketemu", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses ambil detail kuota", quota)
}

func (h *MajorQuotaHandler) Update(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID kuota ga valid", err.Error())
    }

    var input services.UpdateMajorQuotaInput
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

    if err := h.majorQuotaService.Update(id, input, userID); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal update kuota", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses update kuota", nil)
}

func (h *MajorQuotaHandler) Delete(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID kuota ga valid", err.Error())
    }

    if err := h.majorQuotaService.Delete(id); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal hapus kuota", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses hapus kuota", nil)
}

func (h *MajorQuotaHandler) GetLogs(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID kuota ga valid", err.Error())
    }

    page := utils.GetPageFromQuery(c)
    limit := utils.GetLimitFromQuery(c)

    logs, pagination, err := h.majorQuotaService.GetQuotaLogs(id, page, limit)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil history kuota", err.Error())
    }

    return utils.PaginationSuccessResponse(c, "Sukses ambil history kuota", logs, pagination)
}