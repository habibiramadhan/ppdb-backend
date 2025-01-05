// internal/models/major_quota.go
package models

import (
	"time"

	"github.com/google/uuid"
)

type MajorQuota struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	AcademicYearID uuid.UUID  `gorm:"type:uuid;not null" json:"academic_year_id"`
	MajorID        uuid.UUID  `gorm:"type:uuid;not null" json:"major_id"`
	TotalQuota     int        `gorm:"not null" json:"total_quota"`
	FilledQuota    int        `gorm:"not null;default:0" json:"filled_quota"`
	RemainingQuota int        `gorm:"<-:false" json:"remaining_quota"`
	Notes          string     `gorm:"type:text" json:"notes"`
	CreatedBy      uuid.UUID  `gorm:"type:uuid;not null" json:"created_by"`
	UpdatedBy      *uuid.UUID `gorm:"type:uuid" json:"updated_by"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	AcademicYear AcademicYear    `gorm:"foreignKey:AcademicYearID" json:"academic_year"`
	Major        Major           `gorm:"foreignKey:MajorID" json:"major"`
	Creator      User            `gorm:"foreignKey:CreatedBy" json:"creator"`
	Updater      *User           `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
	Logs         []MajorQuotaLog `gorm:"foreignKey:MajorQuotaID" json:"logs,omitempty"`
}

func (mq *MajorQuota) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":            mq.ID,
		"academic_year": mq.AcademicYear.GetFormattedName(),
		"major": map[string]interface{}{
			"id":   mq.Major.ID,
			"name": mq.Major.Name,
			"code": mq.Major.Code,
		},
		"total_quota":     mq.TotalQuota,
		"filled_quota":    mq.FilledQuota,
		"remaining_quota": mq.RemainingQuota,
		"notes":           mq.Notes,
		"created_at":      mq.CreatedAt,
		"updated_at":      mq.UpdatedAt,
	}
}

func (MajorQuota) TableName() string {
	return "major_quotas"
}
