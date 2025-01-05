// internal/core/repositories/verification_repository.go
package repositories

import (
	"ppdb-backend/internal/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VerificationRepository interface {
	Create(verification *models.EmailVerification) error
	FindByToken(token string) (*models.EmailVerification, error)
	FindByUserID(userID uuid.UUID) (*models.EmailVerification, error)
	Update(verification *models.EmailVerification) error
}

type verificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) VerificationRepository {
	return &verificationRepository{db}
}

func (r *verificationRepository) Create(verification *models.EmailVerification) error {
	return r.db.Create(verification).Error
}

func (r *verificationRepository) FindByToken(token string) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.Where("token = ?", token).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *verificationRepository) FindByUserID(userID uuid.UUID) (*models.EmailVerification, error) {
	var verification models.EmailVerification
	err := r.db.Where("user_id = ?", userID).First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

func (r *verificationRepository) Update(verification *models.EmailVerification) error {
	return r.db.Save(verification).Error
}