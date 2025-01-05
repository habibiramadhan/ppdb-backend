// internal/core/services/major_quota_service.go
package services

import (
    "errors"
    "fmt"
    "ppdb-backend/internal/core/repositories"
    "ppdb-backend/internal/models"
    "ppdb-backend/utils"
    "github.com/google/uuid"
)

type MajorQuotaService interface {
    Create(input CreateMajorQuotaInput, userID uuid.UUID) error
    GetAll(page, limit int) ([]models.MajorQuota, *utils.PaginationMeta, error)
    GetByID(id uuid.UUID) (*models.MajorQuota, error)
    Update(id uuid.UUID, input UpdateMajorQuotaInput, userID uuid.UUID) error
    Delete(id uuid.UUID) error
    IncreaseFilled(id uuid.UUID) error
    DecreaseFilled(id uuid.UUID) error
    GetQuotaLogs(quotaID uuid.UUID, page, limit int) ([]models.MajorQuotaLog, *utils.PaginationMeta, error)
    ValidateAvailableQuota(majorID, yearID uuid.UUID) error
}

type majorQuotaService struct {
    quotaRepo       repositories.MajorQuotaRepository
    academicYearRepo repositories.AcademicYearRepository
    majorRepo       repositories.MajorRepository
}

type CreateMajorQuotaInput struct {
    AcademicYearID uuid.UUID `json:"academic_year_id" validate:"required"`
    MajorID        uuid.UUID `json:"major_id" validate:"required"`
    TotalQuota     int      `json:"total_quota" validate:"required,min=0"`
    Notes          string   `json:"notes"`
}

type UpdateMajorQuotaInput struct {
    TotalQuota   int    `json:"total_quota" validate:"required,min=0"`
    Notes        string `json:"notes"`
}

func NewMajorQuotaService(
    quotaRepo repositories.MajorQuotaRepository,
    academicYearRepo repositories.AcademicYearRepository,
    majorRepo repositories.MajorRepository,
) MajorQuotaService {
    return &majorQuotaService{
        quotaRepo: quotaRepo,
        academicYearRepo: academicYearRepo,
        majorRepo: majorRepo,
    }
}

// Bikin kuota baru
func (s *majorQuotaService) Create(input CreateMajorQuotaInput, userID uuid.UUID) error {
    // Validasi tahun ajaran exists dan aktif
    academicYear, err := s.academicYearRepo.FindByID(input.AcademicYearID)
    if err != nil {
        return errors.New("tahun ajaran ga ketemu")
    }
    if !academicYear.IsActive {
        return errors.New("tahun ajaran ga aktif")
    }

    // Validasi jurusan exists dan aktif
    major, err := s.majorRepo.FindByID(input.MajorID)
    if err != nil {
        return errors.New("jurusan ga ketemu")
    }
    if !major.IsActive {
        return errors.New("jurusan ga aktif")
    }

    // Cek apakah kuota sudah ada
    existing, _ := s.quotaRepo.FindByMajorAndYear(input.MajorID, input.AcademicYearID)
    if existing != nil {
        return errors.New("kuota untuk jurusan dan tahun ajaran ini udah ada")
    }

    quota := &models.MajorQuota{
        AcademicYearID: input.AcademicYearID,
        MajorID:        input.MajorID,
        TotalQuota:     input.TotalQuota,
        FilledQuota:    0,
        Notes:          input.Notes,
        CreatedBy:      userID,
    }

    return s.quotaRepo.Create(quota)
}

// Ambil semua kuota
func (s *majorQuotaService) GetAll(page, limit int) ([]models.MajorQuota, *utils.PaginationMeta, error) {
    offset := (page - 1) * limit
    
    quotas, total, err := s.quotaRepo.FindAll(limit, offset)
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

    return quotas, pagination, nil
}

// Ambil detail kuota
func (s *majorQuotaService) GetByID(id uuid.UUID) (*models.MajorQuota, error) {
    return s.quotaRepo.FindByID(id)
}

// Update kuota
func (s *majorQuotaService) Update(id uuid.UUID, input UpdateMajorQuotaInput, userID uuid.UUID) error {
    quota, err := s.quotaRepo.FindByID(id)
    if err != nil {
        return errors.New("kuota ga ketemu")
    }

    // Validasi total kuota baru harus >= filled_quota
    if input.TotalQuota < quota.FilledQuota {
        return errors.New("total kuota ga boleh lebih kecil dari kuota yang udah keisi")
    }

    quota.TotalQuota = input.TotalQuota
    quota.Notes = input.Notes
    quota.UpdatedBy = &userID

    return s.quotaRepo.Update(quota)
}

// Hapus kuota
func (s *majorQuotaService) Delete(id uuid.UUID) error {
    quota, err := s.quotaRepo.FindByID(id)
    if err != nil {
        return errors.New("kuota ga ketemu")
    }

    if quota.FilledQuota > 0 {
        return errors.New("ga bisa hapus kuota yang udah keisi")
    }

    return s.quotaRepo.Delete(id)
}

// Naikin filled quota (+1)
func (s *majorQuotaService) IncreaseFilled(id uuid.UUID) error {
    if err := s.quotaRepo.IncreaseFilled(id); err != nil {
        return errors.New("kuota udah penuh")
    }
    return nil
}

// Kurangin filled quota (-1)
func (s *majorQuotaService) DecreaseFilled(id uuid.UUID) error {
    if err := s.quotaRepo.DecreaseFilled(id); err != nil {
        return errors.New("filled quota udah 0")
    }
    return nil
}

// Ambil history logs
func (s *majorQuotaService) GetQuotaLogs(quotaID uuid.UUID, page, limit int) ([]models.MajorQuotaLog, *utils.PaginationMeta, error) {
    offset := (page - 1) * limit

    logs, total, err := s.quotaRepo.GetQuotaLogs(quotaID, limit, offset)
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

    return logs, pagination, nil
}

// Validasi masih ada kuota available
func (s *majorQuotaService) ValidateAvailableQuota(majorID, yearID uuid.UUID) error {
    quota, err := s.quotaRepo.FindByMajorAndYear(majorID, yearID)
    if err != nil {
        return errors.New("kuota belum diset")
    }

    if quota.FilledQuota >= quota.TotalQuota {
        return fmt.Errorf("kuota jurusan %s sudah penuh", quota.Major.Name)
    }

    return nil
}