package services

import (
	"context"
	"errors"
	"io"
	"time"

	"example/pdfgenerator/initializers"

	"github.com/minio/minio-go/v7"
)

// UploadFile uploads a file to MinIO.
func UploadFile(bucketName, objectName string, file io.Reader) error {
	// Check if MinioClient is initialized
	if initializers.MinioClient == nil {
		return errors.New("MinioClient is not initialized")
	}

	// Validate bucketName and objectName
	if bucketName == "" || objectName == "" {
		return errors.New("bucketName and objectName cannot be empty")
	}

	// Upload the file to MinIO
	_, err := initializers.MinioClient.PutObject(context.Background(), bucketName, objectName, file, -1, minio.PutObjectOptions{
		ContentType: "application/pdf",
	})
	if err != nil {
		return err
	}

	return nil
}

// UploadFile uploads a file to MinIO.
func UploadTemplate(bucketName2, objectName2 string, file io.Reader) error {
	// Check if MinioClient is initialized
	if initializers.MinioClient == nil {
		return errors.New("MinioClient is not initialized")
	}

	// Validate bucketName and objectName
	if bucketName2 == "" || objectName2 == "" {
		return errors.New("bucketName and objectName cannot be empty")
	}

	// Upload the file to MinIO
	_, err := initializers.MinioClient.PutObject(context.Background(), bucketName2, objectName2, file, -1, minio.PutObjectOptions{
		ContentType: "text/html",
	})
	if err != nil {
		return err
	}

	return nil
}

// GenerateFileURL generates a presigned URL for accessing a file.
func GenerateFileURL(bucketName, objectName string) string {
	// Check if MinioClient is initialized
	if initializers.MinioClient == nil {
		return ""
	}

	// Validate bucketName and objectName
	if bucketName == "" || objectName == "" {
		return ""
	}

	// Generate a presigned URL for the object
	presignedURL, err := initializers.MinioClient.PresignedGetObject(context.Background(), bucketName, objectName, 24*time.Hour, nil)
	if err != nil {
		return ""
	}
	return presignedURL.String()

}
