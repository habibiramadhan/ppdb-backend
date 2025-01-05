// internal/models/major_file.go
package models

import (
    "time"
    "github.com/google/uuid"
)

// Enum untuk tipe file
type FileType string

const (
    Brochure   FileType = "brochure"
    Syllabus   FileType = "syllabus"
    Curriculum FileType = "curriculum"
    Other      FileType = "other"
)

type MajorFile struct {
    ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
    MajorID   uuid.UUID `gorm:"type:uuid;not null" json:"major_id"`
    Title     string    `gorm:"size:255;not null" json:"title"`
    FileType  FileType  `gorm:"type:file_type;not null" json:"file_type"`
    FilePath  string    `gorm:"size:255;not null" json:"file_path"`
    FileSize  int64     `gorm:"not null" json:"file_size"`
    MimeType  string    `gorm:"size:100;not null" json:"mime_type"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`

    Major Major `gorm:"foreignKey:MajorID" json:"-"`
}

func (mf *MajorFile) ToResponse() map[string]interface{} {
    return map[string]interface{}{
        "id": mf.ID,
        "major_id": mf.MajorID,
        "title": mf.Title,
        "file_type": mf.FileType,
        "file_path": mf.FilePath,
        "file_size": mf.FileSize,
        "mime_type": mf.MimeType,
        "created_at": mf.CreatedAt,
        "updated_at": mf.UpdatedAt,
    }
}

func IsValidFileType(fileType string) bool {
    switch FileType(fileType) {
    case Brochure, Syllabus, Curriculum, Other:
        return true
    }
    return false
}

func IsAllowedMimeType(mimeType string) bool {
    allowedTypes := map[string]bool{
        "application/pdf": true,
        "application/msword": true,
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
        "image/jpeg": true,
        "image/png": true,
    }
    return allowedTypes[mimeType]
}

// Constant untuk file
const (
    MaxFileSize = 5 * 1024 * 1024 // 5MB
    UploadPath = "public/uploads/majors"
)