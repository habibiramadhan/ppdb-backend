// internal/api/handlers/academic_year_handler.go 
package handlers
import (
    "net/http"
    "ppdb-backend/internal/core/services"
    "ppdb-backend/utils"
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
)

type AcademicYearHandler struct {
    academicYearService services.AcademicYearService
}

func NewAcademicYearHandler(academicYearService services.AcademicYearService) *AcademicYearHandler {
    return &AcademicYearHandler{academicYearService}
}

func (h *AcademicYearHandler) Create(c echo.Context) error {
    var input services.CreateAcademicYearInput
    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid nih", err.Error())
    }

    if err := c.Validate(input); err != nil {
        return utils.ValidationErrorResponse(c, "Data wajib ada yang belum diisi nih", err.Error())
    }

    if err := h.academicYearService.Create(input); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal bikin tahun ajaran baru", err.Error())
    }

    return utils.CreatedResponse(c, "Sukses bikin tahun ajaran baru", nil)
}

func (h *AcademicYearHandler) GetAll(c echo.Context) error {
    page := utils.GetPageFromQuery(c)
    limit := utils.GetLimitFromQuery(c)

    academicYears, pagination, err := h.academicYearService.GetAll(page, limit)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil data tahun ajaran", err.Error())
    }

    return utils.PaginationSuccessResponse(c, "Sukses ambil data tahun ajaran", academicYears, pagination)
}

func (h *AcademicYearHandler) GetByID(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID tahun ajaran ga valid", err.Error())
    }

    academicYear, err := h.academicYearService.GetByID(id)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusNotFound, "Tahun ajaran ga ketemu", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses ambil detail tahun ajaran", academicYear)
}

func (h *AcademicYearHandler) Update(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID tahun ajaran ga valid", err.Error())
    }

    var input services.UpdateAcademicYearInput
    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid", err.Error())
    }

    if err := c.Validate(input); err != nil {
        return utils.ValidationErrorResponse(c, "Data wajib ada yang belum diisi nih", err.Error())
    }

    if err := h.academicYearService.Update(id, input); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal update tahun ajaran", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses update tahun ajaran", nil)
}

func (h *AcademicYearHandler) Delete(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID tahun ajaran ga valid", err.Error())
    }

    if err := h.academicYearService.Delete(id); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal hapus tahun ajaran", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses hapus tahun ajaran", nil)
}

func (h *AcademicYearHandler) SetStatus(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID tahun ajaran ga valid", err.Error())
    }

    var input struct {
        IsActive bool `json:"is_active"`
    }

    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid", err.Error())
    }

    if err := h.academicYearService.SetStatus(id, input.IsActive); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal update status", err.Error())
    }

    status := "non-aktif"
    if input.IsActive {
        status = "aktif"
    }

    return utils.SuccessResponse(c, "Sukses update status jadi "+status, nil)
}

func (h *AcademicYearHandler) GetActive(c echo.Context) error {
    academicYear, err := h.academicYearService.GetActive()
    if err != nil {
        return utils.ErrorResponse(c, http.StatusNotFound, "Belum ada tahun ajaran aktif", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses ambil tahun ajaran aktif", academicYear)
}