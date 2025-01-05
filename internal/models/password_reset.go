// internal/models/password_reset.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type PasswordReset struct {
    ID        uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
    Token     string     `gorm:"unique;not null" json:"token"`
    ExpiresAt time.Time  `json:"expires_at"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`

    User User `gorm:"foreignkey:UserID" json:"user"`
}