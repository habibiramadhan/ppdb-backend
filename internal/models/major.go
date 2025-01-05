// internal/models/major.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type Major struct {
    ID          uuid.UUID    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    Name        string       `gorm:"size:100;not null" json:"name"`
    Code        string       `gorm:"size:20;not null;unique" json:"code"`
    Description string       `gorm:"type:text" json:"description"`
    IsActive    bool        `gorm:"default:true" json:"is_active"`
    IconURL     string       `gorm:"size:255" json:"icon_url"`
    CreatedAt   time.Time    `json:"created_at"`
    UpdatedAt   time.Time    `json:"updated_at"`
    
    Files []MajorFile `gorm:"foreignKey:MajorID" json:"files,omitempty"`
}

func (m *Major) ToResponse() map[string]interface{} {
    return map[string]interface{}{
        "id": m.ID,
        "name": m.Name,
        "code": m.Code,
        "description": m.Description,
        "is_active": m.IsActive,
        "icon_url": m.IconURL,
        "created_at": m.CreatedAt,
        "updated_at": m.UpdatedAt,
    }
}