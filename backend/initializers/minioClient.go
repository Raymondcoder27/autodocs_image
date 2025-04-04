package initializers

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinioClient() {
	var err error
	minioURL := os.Getenv("AUTODOCS_MINIO_URL")
	minioAccessKey := os.Getenv("AUTODOCS_MINIO_ACCESS_KEY")
	minioSecretKey := os.Getenv("AUTODOCS_MINIO_SECRET_KEY")
	if minioURL == "" || minioAccessKey == "" || minioSecretKey == "" {
		log.Fatalf("AUTODOCS_MINIO_URL, AUTODOCS_MINIO_ACCESS_KEY, or AUTODOCS_MINIO_SECRET_KEY environment variable not set")
	}
	MinioClient, err = minio.New(minioURL, &minio.Options{
		Creds:  credentials.NewStaticV4(minioAccessKey, minioSecretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to create MinIO client: %v", err)
	}
}
