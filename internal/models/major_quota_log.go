// internal/models/major_quota_log.go
package models

import (
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
)

type QuotaActionType string

const (
	QuotaActionCreate   QuotaActionType = "CREATE"
	QuotaActionUpdate   QuotaActionType = "UPDATE"
	QuotaActionIncrease QuotaActionType = "INCREASE"
	QuotaActionDecrease QuotaActionType = "DECREASE"
	QuotaActionReset    QuotaActionType = "RESET"
)

func (qt *QuotaActionType) Scan(value interface{}) error {
	*qt = QuotaActionType(value.(string))
	return nil
}

func (qt QuotaActionType) Value() (driver.Value, error) {
	return string(qt), nil
}

type MajorQuotaLog struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	MajorQuotaID   uuid.UUID       `gorm:"type:uuid;not null" json:"major_quota_id"`
	ActionType     QuotaActionType `gorm:"type:quota_action_type;not null" json:"action_type"`
	OldTotalQuota  *int            `json:"old_total_quota,omitempty"`
	NewTotalQuota  int             `json:"new_total_quota"`
	OldFilledQuota *int            `json:"old_filled_quota,omitempty"`
	NewFilledQuota int             `json:"new_filled_quota"`
	Notes          string          `gorm:"type:text" json:"notes"`
	CreatedBy      uuid.UUID       `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt      time.Time       `json:"created_at"`

	MajorQuota MajorQuota `gorm:"foreignKey:MajorQuotaID" json:"-"`
	Creator    User       `gorm:"foreignKey:CreatedBy" json:"creator"`
}

func (mql *MajorQuotaLog) ToResponse() map[string]interface{} {
	return map[string]interface{}{
		"id":               mql.ID,
		"action_type":      mql.ActionType,
		"old_total_quota":  mql.OldTotalQuota,
		"new_total_quota":  mql.NewTotalQuota,
		"old_filled_quota": mql.OldFilledQuota,
		"new_filled_quota": mql.NewFilledQuota,
		"notes":            mql.Notes,
		"created_by": map[string]interface{}{
			"id":   mql.Creator.ID,
			"name": mql.Creator.Name,
		},
		"created_at": mql.CreatedAt,
	}
}

func IsValidQuotaAction(action string) bool {
	switch QuotaActionType(action) {
	case QuotaActionCreate, QuotaActionUpdate, QuotaActionIncrease, QuotaActionDecrease, QuotaActionReset:
		return true
	}
	return false
}

func (MajorQuotaLog) TableName() string {
	return "major_quota_logs"
}
