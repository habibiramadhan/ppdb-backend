// internal/core/services/session_service.go
package services

import (
    "crypto/rand"
    "encoding/hex"
    "errors"
    "time"
    "ppdb-backend/internal/models"
    "ppdb-backend/internal/core/repositories"
    "github.com/google/uuid"
)

type SessionService interface {
    CreateSession(userID uuid.UUID, deviceInfo, ipAddress string) (*models.UserSession, error)
    ValidateSession(token string) (*models.UserSession, error)
    UpdateActivity(sessionID uuid.UUID) error
    RevokeSession(sessionID uuid.UUID) error
    RevokeAllSessions(userID uuid.UUID) error
    GetActiveSessions(userID uuid.UUID) ([]models.UserSession, error)
    CleanupExpiredSessions() error
}

type sessionService struct {
    sessionRepo repositories.SessionRepository
    userRepo    repositories.UserRepository
}

func NewSessionService(
    sessionRepo repositories.SessionRepository,
    userRepo repositories.UserRepository,
) SessionService {
    return &sessionService{
        sessionRepo: sessionRepo,
        userRepo:    userRepo,
    }
}

func generateSessionToken() (string, error) {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil
}

func (s *sessionService) CreateSession(userID uuid.UUID, deviceInfo, ipAddress string) (*models.UserSession, error) {
    // Generate session token
    token, err := generateSessionToken()
    if err != nil {
        return nil, err
    }

    // Create new session
    session := &models.UserSession{
        UserID:       userID,
        Token:        token,
        DeviceInfo:   deviceInfo,
        IPAddress:    ipAddress,
        LastActivity: time.Now(),
        ExpiresAt:    time.Now().Add(24 * time.Hour), // Session expires in 24 hours
    }

    if err := s.sessionRepo.Create(session); err != nil {
        return nil, err
    }

    return session, nil
}

func (s *sessionService) ValidateSession(token string) (*models.UserSession, error) {
    session, err := s.sessionRepo.FindByToken(token)
    if err != nil {
        return nil, errors.New("invalid session")
    }

    if time.Now().After(session.ExpiresAt) {
        return nil, errors.New("session expired")
    }

    if session.IsRevoked {
        return nil, errors.New("session revoked")
    }

    return session, nil
}

func (s *sessionService) UpdateActivity(sessionID uuid.UUID) error {
    return s.sessionRepo.UpdateLastActivity(sessionID, time.Now())
}

func (s *sessionService) RevokeSession(sessionID uuid.UUID) error {
    return s.sessionRepo.RevokeSession(sessionID)
}

func (s *sessionService) RevokeAllSessions(userID uuid.UUID) error {
    return s.sessionRepo.RevokeAllUserSessions(userID)
}

func (s *sessionService) GetActiveSessions(userID uuid.UUID) ([]models.UserSession, error) {
    return s.sessionRepo.FindByUserID(userID)
}

func (s *sessionService) CleanupExpiredSessions() error {
    return s.sessionRepo.DeleteExpiredSessions()
}