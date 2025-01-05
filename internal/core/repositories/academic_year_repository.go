// internal/core/repositories/academic_year_repository.go
package repositories

import (
    "ppdb-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type AcademicYearRepository interface {
    Create(year *models.AcademicYear) error
    FindAll(limit, offset int) ([]models.AcademicYear, int64, error)
    FindByID(id uuid.UUID) (*models.AcademicYear, error)
    Update(year *models.AcademicYear) error
    Delete(id uuid.UUID) error
    FindActive() (*models.AcademicYear, error)
    SetActive(id uuid.UUID) error
    SetInactive(id uuid.UUID) error
}

type academicYearRepository struct {
    db *gorm.DB
}

func NewAcademicYearRepository(db *gorm.DB) AcademicYearRepository {
    return &academicYearRepository{db}
}

func (r *academicYearRepository) Create(year *models.AcademicYear) error {
    return r.db.Create(year).Error
}

func (r *academicYearRepository) FindAll(limit, offset int) ([]models.AcademicYear, int64, error) {
    var years []models.AcademicYear
    var total int64

    if err := r.db.Model(&models.AcademicYear{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    if err := r.db.Limit(limit).Offset(offset).Order("year_start desc").Find(&years).Error; err != nil {
        return nil, 0, err
    }

    return years, total, nil
}

func (r *academicYearRepository) FindByID(id uuid.UUID) (*models.AcademicYear, error) {
    var year models.AcademicYear
    err := r.db.First(&year, "id = ?", id).Error
    if err != nil {
        return nil, err
    }
    return &year, nil
}

func (r *academicYearRepository) Update(year *models.AcademicYear) error {
    return r.db.Save(year).Error
}

func (r *academicYearRepository) Delete(id uuid.UUID) error {
    return r.db.Delete(&models.AcademicYear{}, "id = ?", id).Error
}

func (r *academicYearRepository) FindActive() (*models.AcademicYear, error) {
    var year models.AcademicYear
    err := r.db.Where("is_active = ?", true).First(&year).Error
    if err != nil {
        return nil, err
    }
    return &year, nil
}

func (r *academicYearRepository) SetActive(id uuid.UUID) error {
    tx := r.db.Begin()

    if err := tx.Model(&models.AcademicYear{}).Update("is_active", false).Error; err != nil {
        tx.Rollback()
        return err
    }

    if err := tx.Model(&models.AcademicYear{}).Where("id = ?", id).Update("is_active", true).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}

func (r *academicYearRepository) SetInactive(id uuid.UUID) error {
    return r.db.Model(&models.AcademicYear{}).Where("id = ?", id).Update("is_active", false).Error
}