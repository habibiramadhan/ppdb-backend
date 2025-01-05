// internal/models/academic_year.go
package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type AcademicYear struct {
	ID                uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	YearStart         int       `gorm:"not null" json:"year_start"`
	YearEnd           int       `gorm:"not null" json:"year_end"`
	IsActive          bool      `gorm:"default:false" json:"is_active"`
	RegistrationStart time.Time `json:"registration_start"`
	RegistrationEnd   time.Time `json:"registration_end"`
	Description       string    `json:"description"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (ay *AcademicYear) IsRegistrationOpen() bool {
	now := time.Now()
	return ay.IsActive &&
		now.After(ay.RegistrationStart) &&
		now.Before(ay.RegistrationEnd)
}

func (ay *AcademicYear) GetFormattedName() string {
	return fmt.Sprintf("%d/%d", ay.YearStart, ay.YearEnd)
}
