// utils/file.go
package utils

import (
    "errors"
    "io"
    "mime/multipart" 
    "os"
    "path/filepath"
    "strings"
)

const (
    UploadDir = "public/uploads"
    MaxFileSize = 5 * 1024 * 1024
    DirPermission  = 0755
    FilePermission = 0644
)

var (
    AllowedImageExt = []string{".jpg", ".jpeg", ".png"}
    
    AllowedDocExt = []string{".pdf", ".doc", ".docx"}
    
    AllowedMimeTypes = map[string]bool{
        "image/jpeg":                                           true,
        "image/png":                                            true,
        "application/pdf":                                      true,
        "application/msword":                                   true,
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
    }
)

func SetupUploadDir() error {
    paths := []string{
        UploadDir,
        filepath.Join(UploadDir, "majors/icons"),
        filepath.Join(UploadDir, "majors/files"),
    }

    for _, path := range paths {
        if err := os.MkdirAll(path, DirPermission); err != nil {
            return err
        }
    }

    return nil
}

func ValidateFileSize(size int64) error {
    if size > MaxFileSize {
        return errors.New("ukuran file terlalu besar (max 5MB)")
    }
    return nil
}

func IsImageFile(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    for _, allowedExt := range AllowedImageExt {
        if ext == allowedExt {
            return true
        }
    }
    return false
}

func IsDocumentFile(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    for _, allowedExt := range AllowedDocExt {
        if ext == allowedExt {
            return true
        }
    }
    return false
}

func ValidateMimeType(mimeType string) error {
    if !AllowedMimeTypes[mimeType] {
        return errors.New("tipe file tidak diizinkan")
    }
    return nil
}

func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
    if err := os.MkdirAll(filepath.Dir(dst), DirPermission); err != nil {
        return err
    }

    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, FilePermission)
    if err != nil {
        return err
    }
    defer out.Close()

    _, err = io.Copy(out, src)
    return err
}

func DeleteFile(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil
    }
    return os.Remove(path)
}

func CleanFileName(filename string) string {
    filename = strings.Map(func(r rune) rune {
        if strings.ContainsRune(`<>:"/\|?*`, r) {
            return '_'
        }
        return r
    }, filename)

    if len(filename) > 100 {
        ext := filepath.Ext(filename)
        filename = filename[:100-len(ext)] + ext
    }

    return filename
}