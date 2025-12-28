package services

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageService handles generating presigned URLs for secure document access
type StorageService struct {
	client *minio.Client
	bucket string
}

// NewStorageService creates a new storage service connected to Minio
func NewStorageService(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*StorageService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	log.Printf("Connected to Minio at %s for presigned URL generation", endpoint)

	return &StorageService{
		client: client,
		bucket: bucket,
	}, nil
}

// GetPresignedURL generates a presigned URL for temporary access to a document
// The objectPath should be in format: bucket/path/to/file.pdf
// Returns a URL valid for the specified duration (default 15 minutes)
func (s *StorageService) GetPresignedURL(ctx context.Context, objectPath string, expiry time.Duration) (string, error) {
	// Parse object path - format is "bucket/object/path"
	// Remove bucket prefix if present
	objectName := objectPath
	if strings.HasPrefix(objectPath, s.bucket+"/") {
		objectName = strings.TrimPrefix(objectPath, s.bucket+"/")
	}

	// Generate presigned URL
	url, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}

// GetPresignedURLWithDownload generates a presigned URL that forces download with specified filename
func (s *StorageService) GetPresignedURLWithDownload(ctx context.Context, objectPath string, expiry time.Duration, downloadName string) (string, error) {
	objectName := objectPath
	if strings.HasPrefix(objectPath, s.bucket+"/") {
		objectName = strings.TrimPrefix(objectPath, s.bucket+"/")
	}

	// Request parameters for content-disposition header
	reqParams := make(map[string]string)
	if downloadName != "" {
		reqParams["response-content-disposition"] = fmt.Sprintf("attachment; filename=\"%s\"", downloadName)
	}

	url, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return url.String(), nil
}
