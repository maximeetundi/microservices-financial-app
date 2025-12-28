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
	client       *minio.Client // Client for internal operations (upload/download if needed)
	signerClient *minio.Client // Client for generating public presigned URLs
	bucket       string
}

// NewStorageService creates a new storage service connected to Minio
func NewStorageService(endpoint, accessKey, secretKey, bucket string, useSSL bool, publicURL string) (*StorageService, error) {
	// 1. Initialize internal client (backend-to-backend)
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// 2. Initialize signer client (for public URLs)
	// If publicURL is provided, we use it to configure the signer client
	// so that generated signatures match the public hostname.
	var signerClient *minio.Client
	if publicURL != "" {
		u, err := url.Parse(publicURL)
		if err != nil {
			log.Printf("Warning: Invalid public URL %s, falling back to internal endpoint for signing: %v", publicURL, err)
			signerClient = client
		} else {
			publicHost := u.Host
			publicSecure := u.Scheme == "https"
			
			signerClient, err = minio.New(publicHost, &minio.Options{
				Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
				Secure: publicSecure,
			})
			if err != nil {
				log.Printf("Warning: Failed to create public minio client, falling back to internal: %v", err)
				signerClient = client
			} else {
				log.Printf("Initialized Minio signer for public URL: %s (Secure: %v)", publicHost, publicSecure)
			}
		}
	} else {
		signerClient = client // Fallback to internal if no public URL
	}

	log.Printf("Connected to Minio at %s (Internal)", endpoint)

	return &StorageService{
		client:       client,
		signerClient: signerClient,
		bucket:       bucket,
	}, nil
}

// GetPresignedURL generates a presigned URL for temporary access to a document
func (s *StorageService) GetPresignedURL(ctx context.Context, objectPath string, expiry time.Duration) (string, error) {
	objectName := objectPath
	if strings.HasPrefix(objectPath, s.bucket+"/") {
		objectName = strings.TrimPrefix(objectPath, s.bucket+"/")
	}

	// Use signerClient to generate URL with correct host and signature
	presignedURL, err := s.signerClient.PresignedGetObject(ctx, s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
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

	// Use signerClient to generate URL with correct host and signature
	presignedURL, err := s.signerClient.PresignedGetObject(ctx, s.bucket, objectName, expiry, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

