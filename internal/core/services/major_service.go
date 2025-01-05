// internal/core/services/major_service.go
package services

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"ppdb-backend/internal/core/repositories"
	"ppdb-backend/internal/models"
	"ppdb-backend/utils"

	"github.com/google/uuid"
)

type MajorService interface {
	Create(input CreateMajorInput) error
	GetAll(page, limit int) ([]models.Major, *utils.PaginationMeta, error)
	GetByID(id uuid.UUID) (*models.Major, error)
	Update(id uuid.UUID, input UpdateMajorInput) error
	Delete(id uuid.UUID) error
	SetStatus(id uuid.UUID, isActive bool) error
	UpdateIcon(id uuid.UUID, file *multipart.FileHeader) error
	Search(keyword string, status *bool, page, limit int) ([]models.Major, *utils.PaginationMeta, error)
	// Method untuk files
	UploadFiles(majorID uuid.UUID, files []*multipart.FileHeader, fileType models.FileType) error
	DeleteFile(fileID uuid.UUID) error
	GetMajorFiles(majorID uuid.UUID) ([]models.MajorFile, error)
}

type majorService struct {
	majorRepo     repositories.MajorRepository
	majorFileRepo repositories.MajorFileRepository
}

type CreateMajorInput struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description"`
	IsActive    bool   `json:"is_active"`
}

type UpdateMajorInput struct {
	Name        string `json:"name" validate:"required"`
	Code        string `json:"code" validate:"required"`
	Description string `json:"description"`
}

func NewMajorService(majorRepo repositories.MajorRepository, majorFileRepo repositories.MajorFileRepository) MajorService {
	return &majorService{
		majorRepo:     majorRepo,
		majorFileRepo: majorFileRepo,
	}
}

// Bikin jurusan baru
func (s *majorService) Create(input CreateMajorInput) error {
	// Cek kode jurusan udah ada belum
	existing, _ := s.majorRepo.FindByCode(input.Code)
	if existing != nil {
		return errors.New("waduh kode jurusan udah dipake nih")
	}

	major := &models.Major{
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
		IsActive:    input.IsActive,
	}

	return s.majorRepo.Create(major)
}

// Ambil semua jurusan dengan pagination
func (s *majorService) GetAll(page, limit int) ([]models.Major, *utils.PaginationMeta, error) {
	offset := (page - 1) * limit

	majors, total, err := s.majorRepo.FindAll(limit, offset)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (int(total) + limit - 1) / limit
	pagination := &utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPages,
	}

	return majors, pagination, nil
}

// Ambil detail jurusan by ID
func (s *majorService) GetByID(id uuid.UUID) (*models.Major, error) {
	major, err := s.majorRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("jurusan ga ketemu nih")
	}
	return major, nil
}

// Update data jurusan
func (s *majorService) Update(id uuid.UUID, input UpdateMajorInput) error {
	major, err := s.majorRepo.FindByID(id)
	if err != nil {
		return errors.New("jurusan ga ketemu nih")
	}

	// Kalo ganti kode, cek duplikat
	if major.Code != input.Code {
		existing, _ := s.majorRepo.FindByCode(input.Code)
		if existing != nil {
			return errors.New("waduh kode jurusan udah dipake nih")
		}
	}

	major.Name = input.Name
	major.Code = input.Code
	major.Description = input.Description

	return s.majorRepo.Update(major)
}

// Hapus jurusan
func (s *majorService) Delete(id uuid.UUID) error {
	major, err := s.majorRepo.FindByID(id)
	if err != nil {
		return errors.New("jurusan ga ketemu nih")
	}

	if major.IsActive {
		return errors.New("ga bisa hapus jurusan yang masih aktif, non-aktifin dulu ya")
	}

	// Hapus semua file terkait
	files, err := s.majorFileRepo.FindByMajorID(id)
	if err == nil {
		for _, file := range files {
			_ = utils.DeleteFile(file.FilePath) // Hapus file fisik
		}
	}

	return s.majorRepo.Delete(id)
}

// Update status aktif/non-aktif
func (s *majorService) SetStatus(id uuid.UUID, isActive bool) error {
	_, err := s.majorRepo.FindByID(id)
	if err != nil {
		return errors.New("jurusan ga ketemu nih")
	}

	return s.majorRepo.SetStatus(id, isActive)
}

// Upload/update icon jurusan
func (s *majorService) UpdateIcon(id uuid.UUID, file *multipart.FileHeader) error {
	// Cek jurusan exists
	major, err := s.majorRepo.FindByID(id)
	if err != nil {
		return errors.New("jurusan ga ketemu nih")
	}

	// Validasi file
	if err := utils.ValidateFileSize(file.Size); err != nil {
		return err
	}

	if !utils.IsImageFile(file.Filename) {
		return errors.New("format file harus gambar (jpg/png)")
	}

	if err := utils.ValidateMimeType(file.Header.Get("Content-Type")); err != nil {
		return err
	}

	// Generate nama file yang aman
	filename := fmt.Sprintf("%s_%s%s",
		major.Code,
		uuid.New().String(),
		filepath.Ext(file.Filename),
	)
	filename = utils.CleanFileName(filename)

	// Set path file
	path := filepath.Join(utils.UploadDir, "majors/icons", filename)

	// Hapus icon lama kalo ada
	if major.IconURL != "" {
		_ = utils.DeleteFile(major.IconURL)
	}

	// Upload file baru
	if err := utils.SaveUploadedFile(file, path); err != nil {
		return fmt.Errorf("gagal upload file: %v", err)
	}

	// Update database
	return s.majorRepo.UpdateIcon(id, path)
}

// Search jurusan dengan filter
func (s *majorService) Search(keyword string, status *bool, page, limit int) ([]models.Major, *utils.PaginationMeta, error) {
	offset := (page - 1) * limit

	majors, total, err := s.majorRepo.SearchMajors(keyword, status, limit, offset)
	if err != nil {
		return nil, nil, err
	}

	totalPages := (int(total) + limit - 1) / limit
	pagination := &utils.PaginationMeta{
		Page:      page,
		Limit:     limit,
		TotalData: total,
		TotalPage: totalPages,
	}

	return majors, pagination, nil
}

// Upload file pendukung jurusan
func (s *majorService) UploadFiles(majorID uuid.UUID, files []*multipart.FileHeader, fileType models.FileType) error {
	// Validasi jurusan exists
	if _, err := s.majorRepo.FindByID(majorID); err != nil {
		return errors.New("jurusan ga ketemu nih")
	}

	// Validasi tipe file
	if !models.IsValidFileType(string(fileType)) {
		return errors.New("tipe file ga valid")
	}

	for _, file := range files {
		// Validasi tiap file
		if err := utils.ValidateFileSize(file.Size); err != nil {
			return fmt.Errorf("file %s: %v", file.Filename, err)
		}

		if !utils.IsDocumentFile(file.Filename) {
			return fmt.Errorf("file %s: format file harus pdf/doc/docx", file.Filename)
		}

		if err := utils.ValidateMimeType(file.Header.Get("Content-Type")); err != nil {
			return fmt.Errorf("file %s: %v", file.Filename, err)
		}

		// Generate nama file yang aman
		filename := fmt.Sprintf("%s_%s_%s%s",
			fileType,
			uuid.New().String(),
			utils.CleanFileName(filepath.Base(file.Filename)),
			filepath.Ext(file.Filename),
		)

		// Set path file
		path := filepath.Join(utils.UploadDir, "majors/files", filename)

		// Upload file
		if err := utils.SaveUploadedFile(file, path); err != nil {
			return fmt.Errorf("gagal upload file %s: %v", file.Filename, err)
		}

		// Simpan ke database
		majorFile := &models.MajorFile{
			MajorID:  majorID,
			Title:    file.Filename,
			FileType: fileType,
			FilePath: path,
			FileSize: file.Size,
			MimeType: file.Header.Get("Content-Type"),
		}

		if err := s.majorFileRepo.Create(majorFile); err != nil {
			// Hapus file yang udah terupload kalo gagal
			_ = utils.DeleteFile(path)
			return fmt.Errorf("gagal simpan data file %s: %v", file.Filename, err)
		}
	}

	return nil
}

// Hapus file pendukung
func (s *majorService) DeleteFile(fileID uuid.UUID) error {
	file, err := s.majorFileRepo.FindByID(fileID)
	if err != nil {
		return errors.New("file ga ketemu nih")
	}

	// Hapus file fisik
	if err := utils.DeleteFile(file.FilePath); err != nil {
		return fmt.Errorf("gagal hapus file: %v", err)
	}

	// Hapus dari database
	return s.majorFileRepo.Delete(fileID)
}

// Ambil semua file jurusan
func (s *majorService) GetMajorFiles(majorID uuid.UUID) ([]models.MajorFile, error) {
	// Validasi jurusan exists
	if _, err := s.majorRepo.FindByID(majorID); err != nil {
		return nil, errors.New("jurusan ga ketemu nih")
	}

	return s.majorFileRepo.FindByMajorID(majorID)
}
