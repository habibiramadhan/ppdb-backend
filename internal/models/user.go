package models

import (
    "time"

    "github.com/google/uuid"
)

type User struct {
    ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    Email           string     `gorm:"uniqueIndex;not null" json:"email"`
    Password        string     `gorm:"not null" json:"-"`
    Name            string     `gorm:"not null" json:"name"`
    Role            string     `gorm:"type:user_role;not null" json:"role"`
    Status          string     `gorm:"type:user_status;default:inactive" json:"status"`
    Phone           string     `json:"phone"`
    EmailVerifiedAt *time.Time `json:"email_verified_at"`
    CreatedAt       time.Time  `json:"created_at"`
    UpdatedAt       time.Time  `json:"updated_at"`
}