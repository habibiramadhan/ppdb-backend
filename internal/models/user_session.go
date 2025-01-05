// internal/models/user_session.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type UserSession struct {
    ID           uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    UserID       uuid.UUID `gorm:"type:uuid;not null" json:"user_id"`
    Token        string    `gorm:"unique;not null" json:"token"`
    DeviceInfo   string    `json:"device_info"`
    IPAddress    string    `json:"ip_address"`
    LastActivity time.Time `json:"last_activity"`
    IsRevoked    bool      `gorm:"default:false" json:"is_revoked"`
    ExpiresAt    time.Time `json:"expires_at"`
    CreatedAt    time.Time `json:"created_at"`
    UpdatedAt    time.Time `json:"updated_at"`

    User User `gorm:"foreignkey:UserID" json:"user"`
}