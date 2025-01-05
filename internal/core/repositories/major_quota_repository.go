// internal/core/repositories/major_quota_repository.go
package repositories

import (
    "ppdb-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type MajorQuotaRepository interface {
    Create(quota *models.MajorQuota) error
    FindAll(limit, offset int) ([]models.MajorQuota, int64, error)
    FindByID(id uuid.UUID) (*models.MajorQuota, error)
    FindByMajorAndYear(majorID, yearID uuid.UUID) (*models.MajorQuota, error)
    Update(quota *models.MajorQuota) error
    Delete(id uuid.UUID) error
    IncreaseFilled(id uuid.UUID) error
    DecreaseFilled(id uuid.UUID) error
    GetQuotaLogs(quotaID uuid.UUID, limit, offset int) ([]models.MajorQuotaLog, int64, error)
}

type majorQuotaRepository struct {
    db *gorm.DB
}

func NewMajorQuotaRepository(db *gorm.DB) MajorQuotaRepository {
    return &majorQuotaRepository{db}
}

func (r *majorQuotaRepository) Create(quota *models.MajorQuota) error {
    return r.db.Create(quota).Error
}

func (r *majorQuotaRepository) FindAll(limit, offset int) ([]models.MajorQuota, int64, error) {
    var quotas []models.MajorQuota
    var total int64

    // Count total data
    if err := r.db.Model(&models.MajorQuota{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get data dengan eager loading
    err := r.db.Preload("AcademicYear").
        Preload("Major").
        Preload("Creator").
        Preload("Updater").
        Limit(limit).
        Offset(offset).
        Order("created_at desc").
        Find(&quotas).Error

    if err != nil {
        return nil, 0, err
    }

    return quotas, total, nil
}

func (r *majorQuotaRepository) FindByID(id uuid.UUID) (*models.MajorQuota, error) {
    var quota models.MajorQuota

    err := r.db.Preload("AcademicYear").
        Preload("Major").
        Preload("Creator").
        Preload("Updater").
        First(&quota, "id = ?", id).Error

    if err != nil {
        return nil, err
    }

    return &quota, nil
}

func (r *majorQuotaRepository) FindByMajorAndYear(majorID, yearID uuid.UUID) (*models.MajorQuota, error) {
    var quota models.MajorQuota

    err := r.db.Preload("AcademicYear").
        Preload("Major").
        Where("major_id = ? AND academic_year_id = ?", majorID, yearID).
        First(&quota).Error

    if err != nil {
        return nil, err
    }

    return &quota, nil
}

func (r *majorQuotaRepository) Update(quota *models.MajorQuota) error {
    return r.db.Save(quota).Error
}

func (r *majorQuotaRepository) Delete(id uuid.UUID) error {
    // Pake transaction untuk hapus quota dan logs
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Hapus logs dulu
        if err := tx.Delete(&models.MajorQuotaLog{}, "major_quota_id = ?", id).Error; err != nil {
            return err
        }
        // Hapus quota
        if err := tx.Delete(&models.MajorQuota{}, "id = ?", id).Error; err != nil {
            return err
        }
        return nil
    })
}

func (r *majorQuotaRepository) IncreaseFilled(id uuid.UUID) error {
    // Increment filled_quota by 1
    return r.db.Model(&models.MajorQuota{}).
        Where("id = ?", id).
        Where("filled_quota < total_quota"). // Validasi masih ada kuota
        UpdateColumn("filled_quota", gorm.Expr("filled_quota + ?", 1)).Error
}

func (r *majorQuotaRepository) DecreaseFilled(id uuid.UUID) error {
    // Decrement filled_quota by 1
    return r.db.Model(&models.MajorQuota{}).
        Where("id = ? AND filled_quota > 0", id).
        UpdateColumn("filled_quota", gorm.Expr("filled_quota - ?", 1)).Error
}

func (r *majorQuotaRepository) GetQuotaLogs(quotaID uuid.UUID, limit, offset int) ([]models.MajorQuotaLog, int64, error) {
    var logs []models.MajorQuotaLog
    var total int64

    // Count total logs
    if err := r.db.Model(&models.MajorQuotaLog{}).
        Where("major_quota_id = ?", quotaID).
        Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get logs dengan eager loading
    err := r.db.Preload("Creator").
        Where("major_quota_id = ?", quotaID).
        Limit(limit).
        Offset(offset).
        Order("created_at desc").
        Find(&logs).Error

    if err != nil {
        return nil, 0, err
    }

    return logs, total, nil
}