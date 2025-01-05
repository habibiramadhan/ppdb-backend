// internal/models/email_verification.go
package models

import (
	"time"
	"github.com/google/uuid"
)

type EmailVerification struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID      uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Token       string     `gorm:"unique;not null" json:"token"`
	SentAt      time.Time  `json:"sent_at"`
	VerifiedAt  *time.Time `json:"verified_at"`
	ExpiresAt   time.Time  `json:"expires_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`

	// Relations
	User User `gorm:"foreignkey:UserID" json:"user"`
}