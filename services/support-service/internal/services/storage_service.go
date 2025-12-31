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

// StorageService handles file uploads and presigned URL generation
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
	
	// Create bucket if not exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		log.Printf("Warning: Failed to check if bucket exists: %v", err)
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

// GetPresignedURL generates a presigned URL for temporary access to a document
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

// UploadFile uploads a file to MinIO and returns the object path (and optional presigned URL)
func (s *StorageService) UploadFile(ctx context.Context, objectName string, filePath string, contentType string) (string, error) {
	// Upload the file
	info, err := s.client.FPutObject(ctx, s.bucket, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	log.Printf("Successfully uploaded %s of size %d", objectName, info.Size)
	return fmt.Sprintf("%s/%s", s.bucket, objectName), nil
}

// UploadFileStream uploads a file from a stream (reader)
func (s *StorageService) UploadFileStream(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	// Upload the file from stream
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file stream: %w", err)
	}

	return fmt.Sprintf("%s/%s", s.bucket, objectName), nil
}
