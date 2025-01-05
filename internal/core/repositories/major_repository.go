// internal/core/repositories/major_repository.go
package repositories

import (
    "ppdb-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type MajorRepository interface {
    Create(major *models.Major) error
    FindAll(limit, offset int) ([]models.Major, int64, error)
    FindByID(id uuid.UUID) (*models.Major, error)
    FindByCode(code string) (*models.Major, error)
    Update(major *models.Major) error
    Delete(id uuid.UUID) error
    SetStatus(id uuid.UUID, isActive bool) error
    UpdateIcon(id uuid.UUID, iconURL string) error
    // Method tambahan buat search & filter
    SearchMajors(keyword string, status *bool, limit, offset int) ([]models.Major, int64, error)
}

type majorRepository struct {
    db *gorm.DB
}

func NewMajorRepository(db *gorm.DB) MajorRepository {
    return &majorRepository{db}
}

func (r *majorRepository) Create(major *models.Major) error {
    return r.db.Create(major).Error
}

func (r *majorRepository) FindAll(limit, offset int) ([]models.Major, int64, error) {
    var majors []models.Major
    var total int64

    // Hitung total data
    if err := r.db.Model(&models.Major{}).Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Ambil data dengan pagination
    if err := r.db.Preload("Files").
        Limit(limit).
        Offset(offset).
        Order("created_at desc").
        Find(&majors).Error; err != nil {
        return nil, 0, err
    }

    return majors, total, nil
}

func (r *majorRepository) FindByID(id uuid.UUID) (*models.Major, error) {
    var major models.Major
    
    if err := r.db.Preload("Files").First(&major, "id = ?", id).Error; err != nil {
        return nil, err
    }
    
    return &major, nil
}

func (r *majorRepository) FindByCode(code string) (*models.Major, error) {
    var major models.Major
    
    if err := r.db.Where("code = ?", code).First(&major).Error; err != nil {
        return nil, err
    }
    
    return &major, nil
}

func (r *majorRepository) Update(major *models.Major) error {
    return r.db.Save(major).Error
}

func (r *majorRepository) Delete(id uuid.UUID) error {
    // Gunakan transaction untuk hapus major dan files
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Hapus files dulu
        if err := tx.Where("major_id = ?", id).Delete(&models.MajorFile{}).Error; err != nil {
            return err
        }
        
        // Hapus major
        if err := tx.Delete(&models.Major{}, id).Error; err != nil {
            return err
        }
        
        return nil
    })
}

func (r *majorRepository) SetStatus(id uuid.UUID, isActive bool) error {
    return r.db.Model(&models.Major{}).
        Where("id = ?", id).
        Update("is_active", isActive).Error
}

func (r *majorRepository) UpdateIcon(id uuid.UUID, iconURL string) error {
    return r.db.Model(&models.Major{}).
        Where("id = ?", id).
        Update("icon_url", iconURL).Error
}

func (r *majorRepository) SearchMajors(keyword string, status *bool, limit, offset int) ([]models.Major, int64, error) {
    var majors []models.Major
    var total int64
    
    query := r.db.Model(&models.Major{})

    // Filter by keyword
    if keyword != "" {
        query = query.Where("name ILIKE ? OR code ILIKE ? OR description ILIKE ?", 
            "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
    }

    // Filter by status
    if status != nil {
        query = query.Where("is_active = ?", *status)
    }

    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }

    // Get data with pagination
    if err := query.Preload("Files").
        Limit(limit).
        Offset(offset).
        Order("created_at desc").
        Find(&majors).Error; err != nil {
        return nil, 0, err
    }

    return majors, total, nil
}