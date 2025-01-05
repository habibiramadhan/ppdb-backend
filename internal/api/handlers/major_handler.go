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

// @Summary Bikin jurusan baru
// @Tags Majors
// @Accept json
// @Produce json
// @Param input body services.CreateMajorInput true "Data Jurusan"
// @Success 201 {object} utils.Response
// @Router /admin/majors [post]
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

// @Summary Ambil semua jurusan
// @Tags Majors
// @Produce json
// @Param page query int false "Halaman"
// @Param limit query int false "Jumlah data per halaman"
// @Success 200 {object} utils.PaginationResponse
// @Router /majors [get]
func (h *MajorHandler) GetAll(c echo.Context) error {
    page := utils.GetPageFromQuery(c)
    limit := utils.GetLimitFromQuery(c)

    majors, pagination, err := h.majorService.GetAll(page, limit)
    if err != nil {
        return utils.ErrorResponse(c, http.StatusInternalServerError, "Gagal ambil data jurusan", err.Error())
    }

    return utils.PaginationSuccessResponse(c, "Sukses ambil data jurusan", majors, pagination)
}

// @Summary Ambil detail jurusan
// @Tags Majors
// @Produce json
// @Param id path string true "ID Jurusan"
// @Success 200 {object} utils.Response
// @Router /majors/{id} [get]
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

// @Summary Update data jurusan
// @Tags Majors
// @Accept json
// @Produce json
// @Param id path string true "ID Jurusan"
// @Param input body services.UpdateMajorInput true "Data Update"
// @Success 200 {object} utils.Response
// @Router /admin/majors/{id} [put]
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

// @Summary Delete jurusan
// @Tags Majors
// @Produce json
// @Param id path string true "ID Jurusan"
// @Success 200 {object} utils.Response
// @Router /admin/majors/{id} [delete]
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

// @Summary Set status jurusan
// @Tags Majors
// @Accept json
// @Produce json
// @Param id path string true "ID Jurusan"
// @Success 200 {object} utils.Response
// @Router /admin/majors/{id}/status [patch]
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

// @Summary Upload icon jurusan
// @Tags Majors
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID Jurusan"
// @Param icon formData file true "File Icon"
// @Success 200 {object} utils.Response
// @Router /admin/majors/{id}/icon [post]
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

// @Summary Upload file pendukung jurusan
// @Tags Majors
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "ID Jurusan"
// @Param type formData string true "Tipe File (brochure/syllabus/curriculum/other)"
// @Param files formData file true "Files"
// @Success 200 {object} utils.Response
// @Router /admin/majors/{id}/files [post]
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

// @Summary Hapus file pendukung
// @Tags Majors
// @Produce json
// @Param id path string true "ID File"
// @Success 200 {object} utils.Response
// @Router /admin/majors/files/{id} [delete]
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

// @Summary Ambil file jurusan
// @Tags Majors
// @Produce json
// @Param id path string true "ID Jurusan"
// @Success 200 {object} utils.Response
// @Router /majors/{id}/files [get]
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