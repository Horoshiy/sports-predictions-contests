package handlers

import (
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	_ "image/gif"  // Import for gif support
	_ "image/jpeg" // Import for jpeg support
	_ "image/png"  // Import for png support
)

// UploadHandler handles file uploads
type UploadHandler struct {
	maxFileSize  int64
	allowedTypes []string
	uploadDir    string
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(maxFileSize int64, uploadDir string) *UploadHandler {
	return &UploadHandler{
		maxFileSize:  maxFileSize,
		allowedTypes: []string{"image/jpeg", "image/png", "image/gif"},
		uploadDir:    uploadDir,
	}
}

// ValidateImageFile validates an uploaded image file
func (h *UploadHandler) ValidateImageFile(file multipart.File, header *multipart.FileHeader) error {
	// Check file size
	if header.Size > h.maxFileSize {
		return fmt.Errorf("file size %d exceeds maximum allowed size %d", header.Size, h.maxFileSize)
	}

	// Read first 512 bytes for content type detection
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Reset file pointer
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %v", err)
	}

	// Detect actual content type from file content
	detectedType := http.DetectContentType(buffer[:n])
	if !h.isAllowedType(detectedType) {
		return fmt.Errorf("file type %s is not allowed (detected from content)", detectedType)
	}

	// Validate that it's actually an image by trying to decode it
	_, _, err = image.Decode(file)
	if err != nil {
		return fmt.Errorf("invalid image file: %v", err)
	}

	// Reset file pointer again
	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("failed to reset file pointer: %v", err)
	}

	return nil
}

// GenerateSecureFilename generates a secure filename for uploaded files
func (h *UploadHandler) GenerateSecureFilename(originalFilename string, userID uint) string {
	ext := strings.ToLower(filepath.Ext(originalFilename))
	timestamp := time.Now().Unix()

	// Generate secure filename: user_id_timestamp.ext
	return fmt.Sprintf("avatar_%d_%d%s", userID, timestamp, ext)
}

// ProcessAvatarUpload processes an avatar upload
// NOTE: This is a placeholder implementation. In production, you should:
// 1. Save the file to a secure location (local storage, S3, CDN, etc.)
// 2. Resize/optimize the image
// 3. Generate thumbnails
// 4. Return the actual public URL
func (h *UploadHandler) ProcessAvatarUpload(file multipart.File, header *multipart.FileHeader, userID uint) (string, error) {
	// Validate the file
	if err := h.ValidateImageFile(file, header); err != nil {
		return "", err
	}

	// Generate secure filename
	filename := h.GenerateSecureFilename(header.Filename, userID)

	// TODO: Implement actual file storage
	// Example implementations:
	// - Local: ioutil.WriteFile(filepath.Join(h.uploadDir, filename), fileBytes, 0644)
	// - S3: s3Client.PutObject(...)
	// - CDN: Upload to CDN service

	// For now, return a placeholder URL
	// WARNING: Files are not actually stored!
	avatarURL := fmt.Sprintf("/uploads/avatars/%s", filename)

	return avatarURL, nil
}

// isAllowedType checks if the content type is allowed
func (h *UploadHandler) isAllowedType(contentType string) bool {
	for _, allowedType := range h.allowedTypes {
		if contentType == allowedType {
			return true
		}
	}
	return false
}

// ValidateFileExtension validates file extension
func (h *UploadHandler) ValidateFileExtension(filename string) error {
	ext := strings.ToLower(filepath.Ext(filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".gif"}

	for _, allowedExt := range allowedExts {
		if ext == allowedExt {
			return nil
		}
	}

	return errors.New("file extension not allowed")
}
