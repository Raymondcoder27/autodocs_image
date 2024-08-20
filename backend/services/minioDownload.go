package services

import (
	"bytes"
	"context"
	"example/pdfgenerator/initializers"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
)

// DownloadFile downloads an object from MinIO and returns the data as a byte slice
func DownloadFile(bucketName, objectName string) ([]byte, error) {
	minioClient := initializers.MinioClient

	// Download the object from MinIO
	object, err := minioClient.GetObject(context.Background(), bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer object.Close()

	// Read the object data into a byte buffer
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, object)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// DeleteFile deletes an object from MinIO
func DeleteFile(bucketName, objectName string) error {
	minioClient := initializers.MinioClient

	// Remove the object from MinIO
	err := minioClient.RemoveObject(context.Background(), bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	fmt.Println("Deleted the file from minio")

	return nil
}
