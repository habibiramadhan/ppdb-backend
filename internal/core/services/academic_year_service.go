// internal/core/services/academic_year_service.go
package services

import (
    "errors"
    "time"
    "ppdb-backend/internal/core/repositories"
    "ppdb-backend/internal/models"
    "ppdb-backend/utils"
    "github.com/google/uuid"
)

type AcademicYearService interface {
    Create(input CreateAcademicYearInput) error
    GetAll(page, limit int) ([]models.AcademicYear, *utils.PaginationMeta, error)
    GetByID(id uuid.UUID) (*models.AcademicYear, error)
    Update(id uuid.UUID, input UpdateAcademicYearInput) error
    Delete(id uuid.UUID) error
    SetStatus(id uuid.UUID, isActive bool) error
    GetActive() (*models.AcademicYear, error)
}

type academicYearService struct {
    academicYearRepo repositories.AcademicYearRepository
}

type CreateAcademicYearInput struct {
    YearStart         int       `json:"year_start" validate:"required"`
    YearEnd           int       `json:"year_end" validate:"required,gtfield=YearStart"`
    RegistrationStart time.Time `json:"registration_start" validate:"required"`
    RegistrationEnd   time.Time `json:"registration_end" validate:"required,gtfield=RegistrationStart"`
    Description       string    `json:"description"`
    IsActive         bool      `json:"is_active"`
}

type UpdateAcademicYearInput struct {
    YearStart         int       `json:"year_start" validate:"required"`
    YearEnd           int       `json:"year_end" validate:"required,gtfield=YearStart"`
    RegistrationStart time.Time `json:"registration_start" validate:"required"`
    RegistrationEnd   time.Time `json:"registration_end" validate:"required,gtfield=RegistrationStart"`
    Description       string    `json:"description"`
}

func NewAcademicYearService(academicYearRepo repositories.AcademicYearRepository) AcademicYearService {
    return &academicYearService{
        academicYearRepo: academicYearRepo,
    }
}

func (s *academicYearService) Create(input CreateAcademicYearInput) error {
    if input.IsActive {
        active, _ := s.academicYearRepo.FindActive()
        if active != nil {
            return errors.New("udah ada tahun ajaran yang aktif nih, non-aktifin dulu ya yang lama")
        }
    }

    academicYear := &models.AcademicYear{
        YearStart:         input.YearStart,
        YearEnd:           input.YearEnd,
        RegistrationStart: input.RegistrationStart,
        RegistrationEnd:   input.RegistrationEnd,
        Description:       input.Description,
        IsActive:         input.IsActive,
    }

    return s.academicYearRepo.Create(academicYear)
}

func (s *academicYearService) GetAll(page, limit int) ([]models.AcademicYear, *utils.PaginationMeta, error) {
    offset := (page - 1) * limit
    
    academicYears, total, err := s.academicYearRepo.FindAll(limit, offset)
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

    return academicYears, pagination, nil
}

func (s *academicYearService) GetByID(id uuid.UUID) (*models.AcademicYear, error) {
    return s.academicYearRepo.FindByID(id)
}

func (s *academicYearService) Update(id uuid.UUID, input UpdateAcademicYearInput) error {
    academicYear, err := s.academicYearRepo.FindByID(id)
    if err != nil {
        return errors.New("tahun ajaran ga ketemu nih")
    }

    academicYear.YearStart = input.YearStart
    academicYear.YearEnd = input.YearEnd
    academicYear.RegistrationStart = input.RegistrationStart
    academicYear.RegistrationEnd = input.RegistrationEnd
    academicYear.Description = input.Description

    return s.academicYearRepo.Update(academicYear)
}

func (s *academicYearService) Delete(id uuid.UUID) error {
    academicYear, err := s.academicYearRepo.FindByID(id)
    if err != nil {
        return errors.New("tahun ajaran ga ketemu nih")
    }

    if academicYear.IsActive {
        return errors.New("gabisa hapus tahun ajaran yang masih aktif")
    }

    return s.academicYearRepo.Delete(id)
}

func (s *academicYearService) SetStatus(id uuid.UUID, isActive bool) error {
    if isActive {
        return s.academicYearRepo.SetActive(id)
    }
    return s.academicYearRepo.SetInactive(id)
}

func (s *academicYearService) GetActive() (*models.AcademicYear, error) {
    active, err := s.academicYearRepo.FindActive()
    if err != nil {
        return nil, errors.New("belum ada tahun ajaran yang aktif nih")
    }
    return active, nil
}