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
    // Base directory untuk upload
    UploadDir = "public/uploads"
    
    // Max file size (5MB)
    MaxFileSize = 5 * 1024 * 1024
    
    // Permission untuk folder & file
    DirPermission  = 0755 // rwxr-xr-x
    FilePermission = 0644 // rw-r--r--
)

var (
    // Allowed image extensions
    AllowedImageExt = []string{".jpg", ".jpeg", ".png"}
    
    // Allowed document extensions
    AllowedDocExt = []string{".pdf", ".doc", ".docx"}
    
    // Allowed mime types
    AllowedMimeTypes = map[string]bool{
        "image/jpeg":                                           true,
        "image/png":                                            true,
        "application/pdf":                                      true,
        "application/msword":                                   true,
        "application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
    }
)

// Setup folder upload
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

// Cek file size
func ValidateFileSize(size int64) error {
    if size > MaxFileSize {
        return errors.New("ukuran file terlalu besar (max 5MB)")
    }
    return nil
}

// Cek extension file gambar
func IsImageFile(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    for _, allowedExt := range AllowedImageExt {
        if ext == allowedExt {
            return true
        }
    }
    return false
}

// Cek extension file dokumen
func IsDocumentFile(filename string) bool {
    ext := strings.ToLower(filepath.Ext(filename))
    for _, allowedExt := range AllowedDocExt {
        if ext == allowedExt {
            return true
        }
    }
    return false
}

// Validasi mime type
func ValidateMimeType(mimeType string) error {
    if !AllowedMimeTypes[mimeType] {
        return errors.New("tipe file tidak diizinkan")
    }
    return nil
}

// Save uploaded file
func SaveUploadedFile(file *multipart.FileHeader, dst string) error {
    // Buat folder jika belum ada
    if err := os.MkdirAll(filepath.Dir(dst), DirPermission); err != nil {
        return err
    }

    // Open source file
    src, err := file.Open()
    if err != nil {
        return err
    }
    defer src.Close()

    // Create destination file
    out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, FilePermission)
    if err != nil {
        return err
    }
    defer out.Close()

    // Copy file
    _, err = io.Copy(out, src)
    return err
}

// Delete file
func DeleteFile(path string) error {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return nil // File udah ga ada, anggap sukses
    }
    return os.Remove(path)
}

// Clean filename 
func CleanFileName(filename string) string {
    // Hapus karakter berbahaya
    filename = strings.Map(func(r rune) rune {
        if strings.ContainsRune(`<>:"/\|?*`, r) {
            return '_'
        }
        return r
    }, filename)

    // Max length 100 karakter
    if len(filename) > 100 {
        ext := filepath.Ext(filename)
        filename = filename[:100-len(ext)] + ext
    }

    return filename
}