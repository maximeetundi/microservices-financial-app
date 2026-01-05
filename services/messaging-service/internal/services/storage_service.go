package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageService handles file uploads to MinIO for chat attachments
type StorageService struct {
	client       *minio.Client
	signerClient *minio.Client
	bucket       string
}

// NewStorageService creates a new storage service connected to MinIO
func NewStorageService(endpoint, accessKey, secretKey, bucket string, useSSL bool, publicURL string) (*StorageService, error) {
	// Initialize internal client
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}

	// Initialize signer client for public URLs
	var signerClient *minio.Client
	if publicURL != "" {
		u, err := url.Parse(publicURL)
		if err != nil {
			log.Printf("Warning: Invalid public URL %s: %v", publicURL, err)
			signerClient = client
		} else {
			publicHost := u.Host
			publicSecure := u.Scheme == "https"
			
			signerClient, err = minio.New(publicHost, &minio.Options{
				Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
				Secure: publicSecure,
			})
			if err != nil {
				log.Printf("Warning: Failed to create public minio client: %v", err)
				signerClient = client
			} else {
				log.Printf("Initialized MinIO signer for public URL: %s", publicHost)
			}
		}
	} else {
		signerClient = client
	}
	
	// Create bucket if not exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		log.Printf("Warning: Failed to check bucket: %v", err)
	} else if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Warning: Failed to create bucket: %v", err)
		} else {
			log.Printf("Created bucket %s", bucket)
		}
	}

	return &StorageService{
		client:       client,
		signerClient: signerClient,
		bucket:       bucket,
	}, nil
}

// UploadFileStream uploads a file from a stream (for multipart uploads)
func (s *StorageService) UploadFileStream(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("Successfully uploaded %s (size: %d, type: %s)", objectName, objectSize, contentType)
	return fmt.Sprintf("%s/%s", s.bucket, objectName), nil
}

// GetPresignedURL generates a presigned URL for file access
func (s *StorageService) GetPresignedURL(ctx context.Context, objectPath string, expiry time.Duration) (string, error) {
	objectName := objectPath
	if strings.HasPrefix(objectPath, s.bucket+"/") {
		objectName = strings.TrimPrefix(objectPath, s.bucket+"/")
	}

	presignedURL, err := s.signerClient.PresignedGetObject(ctx, s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return presignedURL.String(), nil
}

// DeleteFile removes a file from MinIO
func (s *StorageService) DeleteFile(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}
