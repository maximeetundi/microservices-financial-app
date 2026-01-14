package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// StorageService handles file uploads to Minio/S3
type StorageService struct {
	client    *minio.Client
	bucket    string
	publicURL string
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

	// Check if bucket exists, create if not
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		log.Printf("Warning: could not check if bucket exists: %v", err)
	}
	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			log.Printf("Warning: could not create bucket %s: %v", bucket, err)
		} else {
			log.Printf("Created bucket: %s", bucket)
			// Set bucket policy to public read for enterprise modules
			policy := fmt.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Principal": {"AWS": ["*"]},
						"Action": ["s3:GetObject"],
						"Resource": ["arn:aws:s3:::%s/*"]
					}
				]
			}`, bucket)
			err = client.SetBucketPolicy(ctx, bucket, policy)
			if err != nil {
				log.Printf("Warning: could not set bucket policy: %v", err)
			}
		}
	}

	return &StorageService{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}, nil
}

// UploadFile uploads a file to Minio and returns the public URL
func (s *StorageService) UploadFile(ctx context.Context, reader io.Reader, fileName string, fileSize int64, contentType string, folder string) (string, error) {
	// Generate unique object name
	ext := filepath.Ext(fileName)
	objectName := fmt.Sprintf("%s/%s_%s%s", folder, time.Now().Format("20060102150405"), uuid.New().String()[:8], ext)

	// Upload the file
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, fileSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return public URL
	return s.GetPublicURL(objectName), nil
}

// GetPublicURL returns the public URL for an object
func (s *StorageService) GetPublicURL(objectName string) string {
	publicURL := strings.TrimSuffix(s.publicURL, "/")
	return fmt.Sprintf("%s/%s/%s", publicURL, s.bucket, objectName)
}

// DeleteFile deletes a file from Minio
func (s *StorageService) DeleteFile(ctx context.Context, objectName string) error {
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}
