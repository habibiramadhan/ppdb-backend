// internal/models/schedule_notification.go
package models

import (
    "time"
    "github.com/google/uuid"
)

const (
    NotificationPending = "pending"
    NotificationSent   = "sent"
    NotificationFailed = "failed"
)

type ScheduleNotification struct {
    ID            uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    ScheduleID    uuid.UUID  `gorm:"type:uuid;not null" json:"schedule_id"`
    UserID        uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
    Type          string     `gorm:"size:50;not null" json:"type"`
    Status        string     `gorm:"size:50;not null;default:'pending'" json:"status"`
    SentAt        *time.Time `json:"sent_at,omitempty"`
    ErrorMessage  string     `gorm:"type:text" json:"error_message,omitempty"`
    CreatedAt     time.Time  `json:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at"`

    Schedule    Schedule `gorm:"foreignKey:ScheduleID" json:"schedule"`
    User        User    `gorm:"foreignKey:UserID" json:"user"`
}

func (sn *ScheduleNotification) IsSent() bool {
    return sn.Status == NotificationSent
}

func (sn *ScheduleNotification) SetFailed(errorMsg string) {
    sn.Status = NotificationFailed
    sn.ErrorMessage = errorMsg
}

func (sn *ScheduleNotification) SetSent() {
    now := time.Now()
    sn.Status = NotificationSent
    sn.SentAt = &now
    sn.ErrorMessage = ""
}