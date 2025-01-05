// internal/models/schedule.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type ScheduleType string
type PriorityLevel string

const (
    ScheduleRegistration ScheduleType = "registration"
    ScheduleTest        ScheduleType = "test"
    ScheduleInterview   ScheduleType = "interview"
    ScheduleAnnouncement ScheduleType = "announcement"
    ScheduleEnrollment  ScheduleType = "enrollment"
    ScheduleOther      ScheduleType = "other"

    PriorityHigh   PriorityLevel = "high"
    PriorityMedium PriorityLevel = "medium"
    PriorityLow    PriorityLevel = "low"
)

type Schedule struct {
    ID             uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    AcademicYearID uuid.UUID     `gorm:"type:uuid;not null" json:"academic_year_id"`
    Title          string        `gorm:"size:255;not null" json:"title"`
    Description    string        `gorm:"type:text" json:"description"`
    StartDate      time.Time     `gorm:"not null" json:"start_date"`
    EndDate        time.Time     `gorm:"not null" json:"end_date"`
    ScheduleType   ScheduleType  `gorm:"type:schedule_type;not null" json:"schedule_type"`
    Priority       PriorityLevel `gorm:"type:priority_level;not null;default:'medium'" json:"priority"`
    IsActive       bool          `gorm:"default:true" json:"is_active"`
    RemindBefore   *int          `json:"remind_before"`
    Location       string        `gorm:"size:255" json:"location"`
    CreatedBy      uuid.UUID     `gorm:"type:uuid;not null" json:"created_by"`
    UpdatedBy      *uuid.UUID    `gorm:"type:uuid" json:"updated_by,omitempty"`
    CreatedAt      time.Time     `json:"created_at"`
    UpdatedAt      time.Time     `json:"updated_at"`

    AcademicYear   AcademicYear  `gorm:"foreignKey:AcademicYearID" json:"academic_year"`
    Creator        User          `gorm:"foreignKey:CreatedBy" json:"creator"`
    Updater        *User         `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
    Notifications  []ScheduleNotification `gorm:"foreignKey:ScheduleID" json:"notifications,omitempty"`
}

func (s *Schedule) IsOngoing() bool {
    now := time.Now()
    return now.After(s.StartDate) && now.Before(s.EndDate)
}

func (s *Schedule) HasStarted() bool {
    return time.Now().After(s.StartDate)
}

func (s *Schedule) HasEnded() bool {
    return time.Now().After(s.EndDate)
}

func (s *Schedule) DaysUntilStart() int {
    return int(time.Until(s.StartDate).Hours() / 24)
}

func IsValidScheduleType(scheduleType string) bool {
    switch ScheduleType(scheduleType) {
    case ScheduleRegistration, ScheduleTest, ScheduleInterview, 
         ScheduleAnnouncement, ScheduleEnrollment, ScheduleOther:
        return true
    }
    return false
}

func IsValidPriorityLevel(priority string) bool {
    switch PriorityLevel(priority) {
    case PriorityHigh, PriorityMedium, PriorityLow:
        return true
    }
    return false
}