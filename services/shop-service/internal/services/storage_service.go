package services

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"io"
	"path/filepath"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageService struct {
	client     *minio.Client
	bucket     string
	publicURL  string
}

func NewStorageService(endpoint, accessKey, secretKey, bucket string, useSSL bool, publicURL string) (*StorageService, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create minio client: %w", err)
	}
	
	// Ensure bucket exists
	ctx := context.Background()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}
	
	return &StorageService{
		client:    client,
		bucket:    bucket,
		publicURL: publicURL,
	}, nil
}

func (s *StorageService) UploadFile(ctx context.Context, data []byte, filename, contentType string) (string, error) {
	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	
	reader := bytes.NewReader(data)
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	
	url := fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, objectName)
	return url, nil
}

func (s *StorageService) UploadFromReader(ctx context.Context, reader io.Reader, size int64, filename, contentType string) (string, error) {
	objectName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), filename)
	
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}
	
	url := fmt.Sprintf("%s/%s/%s", s.publicURL, s.bucket, objectName)
	return url, nil
}

func (s *StorageService) UploadQRCode(ctx context.Context, data []byte, entityType, entityID string) (string, error) {
	filename := fmt.Sprintf("qr_%s_%s.png", entityType, entityID)
	return s.UploadFile(ctx, data, filename, "image/png")
}

func (s *StorageService) DeleteFile(ctx context.Context, url string) error {
	objectName := filepath.Base(url)
	return s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
}

func (s *StorageService) PresignGet(ctx context.Context, fileURL string, expiry time.Duration) (string, error) {
	if s == nil || s.client == nil {
		return "", fmt.Errorf("storage service not available")
	}

	parsed, err := url.Parse(fileURL)
	if err != nil {
		return "", fmt.Errorf("invalid file url: %w", err)
	}
	objectName := filepath.Base(parsed.Path)
	if objectName == "" {
		return "", fmt.Errorf("invalid object name")
	}

	presigned, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to presign url: %w", err)
	}

	return presigned.String(), nil
}
