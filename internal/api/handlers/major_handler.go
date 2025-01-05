// internal/api/handlers/major_handler.go
package handlers
import (
    "net/http"
    "ppdb-backend/internal/core/services"
    "ppdb-backend/internal/models"
    "ppdb-backend/utils"
    
    "github.com/google/uuid"
    "github.com/labstack/echo/v4"
)

type MajorHandler struct {
    majorService services.MajorService
}

func NewMajorHandler(majorService services.MajorService) *MajorHandler {
    return &MajorHandler{majorService}
}

func (h *MajorHandler) Create(c echo.Context) error {
    var input services.CreateMajorInput
    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid", err.Error())
    }

    if err := c.Validate(input); err != nil {
        return utils.ValidationErrorResponse(c, "Data wajib ada yang belum diisi", err.Error())
    }

    if err := h.majorService.Create(input); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal bikin jurusan baru", err.Error())
    }

    return utils.CreatedResponse(c, "Sukses bikin jurusan baru", nil)
}

func (h *MajorHandler) GetAll(c echo.Context) error {
    page := utils.GetPageFromQuery(c)
    limit := utils.GetLimitFromQuery(c)

    majors, pagination, err := h.majorService.GetAll(page, limit)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil data jurusan", err.Error())
    }

    return utils.PaginationSuccessResponse(c, "Sukses ambil data jurusan", majors, pagination)
}

func (h *MajorHandler) GetByID(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jurusan ga valid", err.Error())
    }

    major, err := h.majorService.GetByID(id)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusNotFound, "Jurusan ga ketemu", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses ambil detail jurusan", major)
}

func (h *MajorHandler) Update(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jurusan ga valid", err.Error())
    }

    var input services.UpdateMajorInput
    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid", err.Error())
    }

    if err := c.Validate(input); err != nil {
        return utils.ValidationErrorResponse(c, "Data wajib ada yang belum diisi", err.Error())
    }

    if err := h.majorService.Update(id, input); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal update jurusan", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses update jurusan", nil)
}

func (h *MajorHandler) Delete(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jurusan ga valid", err.Error())
    }

    if err := h.majorService.Delete(id); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal hapus jurusan", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses hapus jurusan", nil)
}

func (h *MajorHandler) SetStatus(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jurusan ga valid", err.Error())
    }

    var input struct {
        IsActive bool `json:"is_active"`
    }

    if err := c.Bind(&input); err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid", err.Error())
    }

    if err := h.majorService.SetStatus(id, input.IsActive); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal update status jurusan", err.Error())
    }

    status := "non-aktif"
    if input.IsActive {
        status = "aktif"
    }

    return utils.SuccessResponse(c, "Sukses update status jurusan jadi "+status, nil)
}

func (h *MajorHandler) UploadIcon(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jurusan ga valid", err.Error())
    }

    file, err := c.FormFile("icon")
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "File icon ga ada", err.Error())
    }

    if err := h.majorService.UpdateIcon(id, file); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal upload icon", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses upload icon jurusan", nil)
}

func (h *MajorHandler) UploadFiles(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jurusan ga valid", err.Error())
    }

    fileType := models.FileType(c.FormValue("type"))
    if !models.IsValidFileType(string(fileType)) {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Tipe file ga valid", nil)
    }

    form, err := c.MultipartForm()
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "Format input ga valid", err.Error())
    }

    files := form.File["files"]
    if len(files) == 0 {
        return utils.ErrorResponse(c, http.StatusBadRequest, "File ga ada", nil)
    }

    if err := h.majorService.UploadFiles(id, files, fileType); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal upload file", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses upload file pendukung", nil)
}

func (h *MajorHandler) DeleteFile(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID file ga valid", err.Error())
    }

    if err := h.majorService.DeleteFile(id); err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal hapus file", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses hapus file", nil)
}

func (h *MajorHandler) GetFiles(c echo.Context) error {
    id, err := uuid.Parse(c.Param("id"))
    if err != nil {
        return utils.ErrorResponse(c, http.StatusBadRequest, "ID jurusan ga valid", err.Error())
    }

    files, err := h.majorService.GetMajorFiles(id)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil data file", err.Error())
    }

    return utils.SuccessResponse(c, "Sukses ambil data file jurusan", files)
}