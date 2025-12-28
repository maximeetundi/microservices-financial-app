package services

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageService handles generating presigned URLs for secure document access
type StorageService struct {
	client       *minio.Client
	bucket       string
	internalHost string // Internal Docker endpoint (minio:9000)
	publicURL    string // Public accessible URL (https://minio.example.com)
}

// NewStorageService creates a new storage service connected to Minio
func NewStorageService(endpoint, accessKey, secretKey, bucket string, useSSL bool, publicURL string) (*StorageService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	log.Printf("Connected to Minio at %s for presigned URL generation (public: %s)", endpoint, publicURL)

	return &StorageService{
		client:       client,
		bucket:       bucket,
		internalHost: endpoint,
		publicURL:    publicURL,
	}, nil
}

// replaceInternalWithPublic replaces the internal Docker address with the public URL
func (s *StorageService) replaceInternalWithPublic(presignedURL string) string {
	if s.publicURL == "" {
		return presignedURL
	}

	// Parse the presigned URL
	u, err := url.Parse(presignedURL)
	if err != nil {
		return presignedURL
	}

	// Parse the public URL
	pub, err := url.Parse(s.publicURL)
	if err != nil {
		return presignedURL
	}

	// Replace host with public host
	u.Scheme = pub.Scheme
	u.Host = pub.Host

	return u.String()
}

// GetPresignedURL generates a presigned URL for temporary access to a document
func (s *StorageService) GetPresignedURL(ctx context.Context, objectPath string, expiry time.Duration) (string, error) {
	objectName := objectPath
	if strings.HasPrefix(objectPath, s.bucket+"/") {
		objectName = strings.TrimPrefix(objectPath, s.bucket+"/")
	}

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Replace internal address with public URL
	return s.replaceInternalWithPublic(presignedURL.String()), nil
}

// GetPresignedURLWithDownload generates a presigned URL that forces download with specified filename
func (s *StorageService) GetPresignedURLWithDownload(ctx context.Context, objectPath string, expiry time.Duration, downloadName string) (string, error) {
	objectName := objectPath
	if strings.HasPrefix(objectPath, s.bucket+"/") {
		objectName = strings.TrimPrefix(objectPath, s.bucket+"/")
	}

	reqParams := make(url.Values)
	if downloadName != "" {
		reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", downloadName))
	}

	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	// Replace internal address with public URL
	return s.replaceInternalWithPublic(presignedURL.String()), nil
}

