// internal/core/services/password_reset_service.go
package services

import (
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"
    "ppdb-backend/internal/core/repositories"
    "ppdb-backend/internal/models"
    "ppdb-backend/utils"
)

type PasswordResetService interface {
    RequestReset(email string) error
    ValidateToken(token string) error
    ResetPassword(token string, newPassword string) error
}

type passwordResetService struct {
    passwordResetRepo repositories.PasswordResetRepository
    userRepo         repositories.UserRepository
    emailService     EmailService
}

func NewPasswordResetService(
    passwordResetRepo repositories.PasswordResetRepository,
    userRepo repositories.UserRepository,
    emailService EmailService,
) PasswordResetService {
    return &passwordResetService{
        passwordResetRepo: passwordResetRepo,
        userRepo:         userRepo,
        emailService:     emailService,
    }
}

func generateResetToken() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

func (s *passwordResetService) RequestReset(email string) error {
    user, err := s.userRepo.FindByEmail(email)
    if err != nil {
        return errors.New("email not registered")
    }

    token, err := generateResetToken()
    if err != nil {
        return err
    }

    reset := &models.PasswordReset{
        UserID:    user.ID,
        Token:     token,
        ExpiresAt: time.Now().Add(1 * time.Hour),
    }

    if err := s.passwordResetRepo.Create(reset); err != nil {
        return err
    }

    return s.emailService.SendPasswordResetEmail(email, token, user.Name)
}

func (s *passwordResetService) ValidateToken(token string) error {
    reset, err := s.passwordResetRepo.FindByToken(token)
    if err != nil {
        return errors.New("invalid reset token")
    }

    if time.Now().After(reset.ExpiresAt) {
        return errors.New("reset token has expired")
    }

    return nil
}

func (s *passwordResetService) ResetPassword(token string, newPassword string) error {
    reset, err := s.passwordResetRepo.FindByToken(token)
    if err != nil {
        return errors.New("invalid reset token")
    }

    if time.Now().After(reset.ExpiresAt) {
        return errors.New("reset token has expired")
    }

    user, err := s.userRepo.FindByID(reset.UserID)
    if err != nil {
        return errors.New("user not found")
    }

    hashedPassword, err := utils.HashPassword(newPassword)
    if err != nil {
        return err
    }

    user.Password = hashedPassword
    if err := s.userRepo.Update(user); err != nil {
        return err
    }

    return s.passwordResetRepo.Delete(token)
}