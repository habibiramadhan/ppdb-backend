// internal/core/repositories/session_repository.go
package repositories

import (
    "ppdb-backend/internal/models"
    "github.com/google/uuid"
    "gorm.io/gorm"
    "time"
)

type SessionRepository interface {
    Create(session *models.UserSession) error
    FindByToken(token string) (*models.UserSession, error)
    FindByUserID(userID uuid.UUID) ([]models.UserSession, error)
    UpdateLastActivity(id uuid.UUID, lastActivity time.Time) error
    RevokeSession(id uuid.UUID) error
    RevokeAllUserSessions(userID uuid.UUID) error
    DeleteExpiredSessions() error
    Update(session *models.UserSession) error
}

type sessionRepository struct {
    db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) SessionRepository {
    return &sessionRepository{db}
}

func (r *sessionRepository) Create(session *models.UserSession) error {
    return r.db.Create(session).Error
}

func (r *sessionRepository) FindByToken(token string) (*models.UserSession, error) {
    var session models.UserSession
    err := r.db.Where("token = ? AND is_revoked = ?", token, false).First(&session).Error
    if err != nil {
        return nil, err
    }
    return &session, nil
}

func (r *sessionRepository) FindByUserID(userID uuid.UUID) ([]models.UserSession, error) {
    var sessions []models.UserSession
    err := r.db.Where("user_id = ? AND is_revoked = ?", userID, false).Find(&sessions).Error
    return sessions, err
}

func (r *sessionRepository) UpdateLastActivity(id uuid.UUID, lastActivity time.Time) error {
    return r.db.Model(&models.UserSession{}).
        Where("id = ?", id).
        Update("last_activity", lastActivity).Error
}

func (r *sessionRepository) RevokeSession(id uuid.UUID) error {
    return r.db.Model(&models.UserSession{}).
        Where("id = ?", id).
        Update("is_revoked", true).Error
}

func (r *sessionRepository) RevokeAllUserSessions(userID uuid.UUID) error {
    return r.db.Model(&models.UserSession{}).
        Where("user_id = ? AND is_revoked = ?", userID, false).
        Update("is_revoked", true).Error
}

func (r *sessionRepository) DeleteExpiredSessions() error {
    return r.db.Where("expires_at < ?", time.Now()).Delete(&models.UserSession{}).Error
}

func (r *sessionRepository) Update(session *models.UserSession) error {
    return r.db.Save(session).Error
}