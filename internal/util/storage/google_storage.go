package storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io"

	"google.golang.org/api/option"
	"lostnfound-api/internal/config"
	"time"
)

// GoogleCloudStorage implements file storage using Google Cloud Storage
type GoogleCloudStorage struct {
	client     *storage.Client
	bucketName string
	projectID  string
}

// NewGoogleCloudStorage creates a new Google Cloud Storage client
func NewGoogleCloudStorage(cfg *config.Config) (*GoogleCloudStorage, error) {
	ctx := context.Background()

	var client *storage.Client
	var err error

	if cfg.GCSCredentialsFile != "" {
		// Use credentials file if provided
		client, err = storage.NewClient(ctx, option.WithCredentialsFile(cfg.GCSCredentialsFile))
	} else {
		// Otherwise, use default credentials
		client, err = storage.NewClient(ctx)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	return &GoogleCloudStorage{
		client:     client,
		bucketName: cfg.GCSBucketName,
		projectID:  cfg.GCSProjectID,
	}, nil
}

// Close closes the storage client
func (g *GoogleCloudStorage) Close() error {
	return g.client.Close()
}

// UploadFile uploads a file to Google Cloud Storage
func (g *GoogleCloudStorage) UploadFile(ctx context.Context, objectName string, content io.Reader, contentType string) (string, error) {
	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(objectName)
	w := obj.NewWriter(ctx)
	w.ContentType = contentType

	if _, err := io.Copy(w, content); err != nil {
		return "", fmt.Errorf("io.Copy: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("Writer.Close: %w", err)
	}

	// Make the object publicly accessible
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", fmt.Errorf("ACL.Set: %w", err)
	}

	// Generate public URL
	url := fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, objectName)
	return url, nil
}

// GenerateSignedURL generates a signed URL for uploading a file directly
func (g *GoogleCloudStorage) GenerateSignedURL(ctx context.Context, objectName string, contentType string) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:         storage.SigningSchemeV4,
		Method:         "PUT",
		ContentType:    contentType,
		Expires:        time.Now().Add(15 * time.Minute),
		GoogleAccessID: "",
	}

	url, err := storage.SignedURL(g.bucketName, objectName, opts)
	if err != nil {
		return "", fmt.Errorf("storage.SignedURL: %w", err)
	}

	return url, nil
}

// DeleteFile deletes a file from Google Cloud Storage
func (g *GoogleCloudStorage) DeleteFile(ctx context.Context, objectName string) error {
	bucket := g.client.Bucket(g.bucketName)
	obj := bucket.Object(objectName)

	if err := obj.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", objectName, err)
	}

	return nil
}

// GetPublicURL returns a public URL for accessing the object
func (g *GoogleCloudStorage) GetPublicURL(objectName string) string {
	return fmt.Sprintf("https://storage.googleapis.com/%s/%s", g.bucketName, objectName)
}
