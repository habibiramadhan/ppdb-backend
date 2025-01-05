// internal/core/services/verification_service.go
package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"ppdb-backend/internal/core/repositories"
	"ppdb-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type VerificationService interface {
	CreateVerificationToken(userID uuid.UUID) (string, error)
	VerifyEmail(token string) error
	ResendVerification(email string) error
}

type verificationService struct {
	verificationRepo repositories.VerificationRepository
	userRepo         repositories.UserRepository
	emailService     EmailService
}

func NewVerificationService(
	verificationRepo repositories.VerificationRepository,
	userRepo repositories.UserRepository,
	emailService EmailService,
) VerificationService {
	return &verificationService{
		verificationRepo: verificationRepo,
		userRepo:         userRepo,
		emailService:     emailService,
	}
}

func (s *verificationService) CreateVerificationToken(userID uuid.UUID) (string, error) {
	// Generate random token
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	// Create verification record
	verification := &models.EmailVerification{
		UserID:    userID,
		Token:     token,
		SentAt:    time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.verificationRepo.Create(verification); err != nil {
		return "", err
	}

	return token, nil
}

func (s *verificationService) VerifyEmail(token string) error {
	// Get verification record
	verification, err := s.verificationRepo.FindByToken(token)
	if err != nil {
		return errors.New("invalid verification token")
	}

	// Check if token is expired
	if time.Now().After(verification.ExpiresAt) {
		return errors.New("verification token has expired")
	}

	// Check if already verified
	if verification.VerifiedAt != nil {
		return errors.New("email already verified")
	}

	// Update verification record
	now := time.Now()
	verification.VerifiedAt = &now
	if err := s.verificationRepo.Update(verification); err != nil {
		return err
	}

	// Update user status to active
	user, err := s.userRepo.FindByID(verification.UserID)
	if err != nil {
		return err
	}

	user.Status = "active"
	user.EmailVerifiedAt = &now
	return s.userRepo.Update(user)
}

func (s *verificationService) ResendVerification(email string) error {
	// Find user
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	// Check if already verified
	if user.EmailVerifiedAt != nil {
		return errors.New("email already verified")
	}

	// Create new verification token
	token, err := s.CreateVerificationToken(user.ID)
	if err != nil {
		return err
	}

	// Send verification email
	return s.emailService.SendVerificationEmail(email, token, user.Name)
}
