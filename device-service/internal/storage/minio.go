package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

// MinioConfig holds MinIO connection settings
type MinioConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// MinioClient wraps minio-go for firmware object storage
type MinioClient struct {
	client *minio.Client
	cfg    MinioConfig
	logger *zap.Logger
}

// NewMinioClient creates a new MinIO client and ensures the bucket exists
func NewMinioClient(cfg MinioConfig, logger *zap.Logger) (*MinioClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	mc := &MinioClient{client: client, cfg: cfg, logger: logger}

	// Ensure bucket exists
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket %s: %w", cfg.Bucket, err)
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket %s: %w", cfg.Bucket, err)
		}
		logger.Info("MinIO bucket created", zap.String("bucket", cfg.Bucket))
	}

	logger.Info("Connected to MinIO",
		zap.String("endpoint", cfg.Endpoint),
		zap.String("bucket", cfg.Bucket),
	)

	return mc, nil
}

// UploadFirmware uploads a firmware binary to MinIO
// Returns the object path within the bucket and the SHA256 checksum
func (m *MinioClient) UploadFirmware(ctx context.Context, reader io.Reader, objectName string, fileSize int64) (string, error) {
	info, err := m.client.PutObject(ctx, m.cfg.Bucket, objectName, reader, fileSize, minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload firmware: %w", err)
	}

	m.logger.Info("Firmware uploaded to MinIO",
		zap.String("path", objectName),
		zap.Int64("size", info.Size),
		zap.String("etag", info.ETag),
	)

	return objectName, nil
}

// DownloadFirmware retrieves a firmware binary from MinIO
func (m *MinioClient) DownloadFirmware(ctx context.Context, objectName string) (io.ReadCloser, error) {
	obj, err := m.client.GetObject(ctx, m.cfg.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to download firmware %s: %w", objectName, err)
	}
	return obj, nil
}

// DeleteFirmware removes a firmware binary from MinIO
func (m *MinioClient) DeleteFirmware(ctx context.Context, objectName string) error {
	if objectName == "" {
		return nil
	}

	err := m.client.RemoveObject(ctx, m.cfg.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete firmware %s: %w", objectName, err)
	}

	m.logger.Info("Firmware deleted from MinIO", zap.String("path", objectName))
	return nil
}

// GeneratePresignedURL creates a temporary download URL for a firmware file
func (m *MinioClient) GeneratePresignedURL(ctx context.Context, objectName string, expiry time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, m.cfg.Bucket, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}
