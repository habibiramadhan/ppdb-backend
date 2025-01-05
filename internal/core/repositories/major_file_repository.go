// internal/core/repositories/major_file_repository.go
package repositories

import (
    "ppdb-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type MajorFileRepository interface {
    Create(file *models.MajorFile) error
    FindByID(id uuid.UUID) (*models.MajorFile, error)
    FindByMajorID(majorID uuid.UUID) ([]models.MajorFile, error)
    Delete(id uuid.UUID) error
    DeleteByMajorID(majorID uuid.UUID) error
}

type majorFileRepository struct {
    db *gorm.DB
}

func NewMajorFileRepository(db *gorm.DB) MajorFileRepository {
    return &majorFileRepository{db}
}

func (r *majorFileRepository) Create(file *models.MajorFile) error {
    return r.db.Create(file).Error
}

func (r *majorFileRepository) FindByID(id uuid.UUID) (*models.MajorFile, error) {
    var file models.MajorFile
    
    if err := r.db.First(&file, "id = ?", id).Error; err != nil {
        return nil, err
    }
    
    return &file, nil
}

func (r *majorFileRepository) FindByMajorID(majorID uuid.UUID) ([]models.MajorFile, error) {
    var files []models.MajorFile
    
    if err := r.db.Where("major_id = ?", majorID).Find(&files).Error; err != nil {
        return nil, err
    }
    
    return files, nil
}

func (r *majorFileRepository) Delete(id uuid.UUID) error {
    return r.db.Delete(&models.MajorFile{}, "id = ?", id).Error
}

func (r *majorFileRepository) DeleteByMajorID(majorID uuid.UUID) error {
    return r.db.Delete(&models.MajorFile{}, "major_id = ?", majorID).Error
}