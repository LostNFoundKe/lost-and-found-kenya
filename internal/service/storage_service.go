package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"lostnfound-api/internal/models"
	"lostnfound-api/internal/repository"
	"lostnfound-api/internal/util/storage"
	"mime/multipart"
	"path/filepath"
	"strings"
)

// StorageService handles file storage operations
type StorageService struct {
	storage   *storage.GoogleCloudStorage
	imageRepo *repository.ImageRepository
}

// NewStorageService creates a new StorageService
func NewStorageService(storage *storage.GoogleCloudStorage, imageRepo *repository.ImageRepository) *StorageService {
	return &StorageService{
		storage:   storage,
		imageRepo: imageRepo,
	}
}

// UploadItemImage uploads an image for an item and creates a database record
func (s *StorageService) UploadItemImage(ctx context.Context, itemID uuid.UUID, file multipart.File, fileHeader *multipart.FileHeader) (*models.Image, error) {
	// Generate a unique filename
	filename := generateUniqueFilename(fileHeader.Filename)

	// Determine content type from file extension
	contentType := getContentTypeFromFileName(filename)

	// Define the object path in GCS
	objectName := fmt.Sprintf("items/%d/%s", itemID, filename)

	// Upload file to Google Cloud Storage
	url, err := s.storage.UploadFile(ctx, objectName, file, contentType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file: %w", err)
	}

	// Create image record in database
	image := &models.Image{
		URL:    url,
		ItemID: itemID,
	}

	if err := s.imageRepo.Create(image); err != nil {
		// Try to delete the uploaded file if database operation fails
		_ = s.storage.DeleteFile(ctx, objectName)
		return nil, fmt.Errorf("failed to save image record: %w", err)
	}

	return image, nil
}

// DeleteItemImage deletes an image from storage and database
func (s *StorageService) DeleteItemImage(ctx context.Context, imageID uint) error {
	// Fetch image record
	image, err := s.imageRepo.GetByID(imageID)
	if err != nil {
		return fmt.Errorf("image not found: %w", err)
	}

	// Extract object name from URL
	// URL format: https://storage.googleapis.com/bucket-name/object-name
	urlParts := strings.Split(image.URL, "/")
	if len(urlParts) < 5 {
		return fmt.Errorf("invalid image URL format")
	}

	objectName := strings.Join(urlParts[4:], "/")

	// Delete file from Google Cloud Storage
	if err := s.storage.DeleteFile(ctx, objectName); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// Delete record from database
	if err := s.imageRepo.Delete(imageID); err != nil {
		return fmt.Errorf("failed to delete image record: %w", err)
	}

	return nil
}

// GenerateSignedUploadURL generates a signed URL for direct file upload
func (s *StorageService) GenerateSignedUploadURL(ctx context.Context, itemID uint, filename string) (string, string, error) {
	// Generate a unique filename
	uniqueFilename := generateUniqueFilename(filename)

	// Determine content type from file extension
	contentType := getContentTypeFromFileName(uniqueFilename)

	// Define the object path in GCS
	objectName := fmt.Sprintf("items/%d/%s", itemID, uniqueFilename)

	// Generate signed URL
	signedURL, err := s.storage.GenerateSignedURL(ctx, objectName, contentType)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate signed URL: %w", err)
	}

	// Get the eventual public URL
	publicURL := s.storage.GetPublicURL(objectName)

	return signedURL, publicURL, nil
}

// Helper functions

// generateUniqueFilename generates a unique filename by prefixing a UUID
func generateUniqueFilename(originalFilename string) string {
	extension := filepath.Ext(originalFilename)
	return fmt.Sprintf("%s%s", uuid.New().String(), extension)
}

// getContentTypeFromFileName determines the content type based on file extension
func getContentTypeFromFileName(filename string) string {
	extension := strings.ToLower(filepath.Ext(filename))

	switch extension {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".heif", ".heic":
		return "image/heif"
	default:
		return "application/octet-stream"
	}
}
