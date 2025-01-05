// internal/core/repositories/password_reset_repository.go
package repositories

import (
    "ppdb-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
)

type PasswordResetRepository interface {
    Create(reset *models.PasswordReset) error
    FindByToken(token string) (*models.PasswordReset, error)
    FindByUserID(userID uuid.UUID) (*models.PasswordReset, error)
    Delete(token string) error
}

type passwordResetRepository struct {
    db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) PasswordResetRepository {
    return &passwordResetRepository{db}
}

func (r *passwordResetRepository) Create(reset *models.PasswordReset) error {
    return r.db.Create(reset).Error
}

func (r *passwordResetRepository) FindByToken(token string) (*models.PasswordReset, error) {
    var reset models.PasswordReset
    err := r.db.Where("token = ?", token).First(&reset).Error
    if err != nil {
        return nil, err
    }
    return &reset, nil
}

func (r *passwordResetRepository) FindByUserID(userID uuid.UUID) (*models.PasswordReset, error) {
    var reset models.PasswordReset
    err := r.db.Where("user_id = ?", userID).First(&reset).Error
    if err != nil {
        return nil, err
    }
    return &reset, nil
}

func (r *passwordResetRepository) Delete(token string) error {
    return r.db.Where("token = ?", token).Delete(&models.PasswordReset{}).Error
}