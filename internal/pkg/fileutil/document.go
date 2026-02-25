package fileutil

import (
	"mime/multipart"

	"cms-api/internal/pkg/apperror"
)

const (
	MaxDocumentSize = 50 * 1024 * 1024
)

var allowedDocumentTypes = map[string]bool{
	"application/pdf":    true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
	"application/vnd.ms-powerpoint":                                             true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
}

var (
	ErrUnsupportedDocumentType = apperror.NewAppError(nil, "Unsupported document type", 422)
	ErrDocumentTooLarge        = apperror.NewAppError(nil, "Document exceeds maximum size of 50MB", 422)
	ErrDocumentEmpty           = apperror.NewAppError(nil, "Document file is empty", 422)
	ErrDocumentRequired        = apperror.NewAppError(nil, "Document file is required", 422)
)

func ValidateDocument(file *multipart.FileHeader) error {
	if file == nil {
		return ErrDocumentRequired
	}

	if file.Size == 0 {
		return ErrDocumentEmpty
	}

	if file.Size > MaxDocumentSize {
		return ErrDocumentTooLarge
	}

	mimeType := file.Header.Get("Content-Type")
	if !IsAllowedDocumentType(mimeType) {
		return ErrUnsupportedDocumentType
	}

	return nil
}

func IsAllowedDocumentType(mimeType string) bool {
	return allowedDocumentTypes[mimeType]
}
