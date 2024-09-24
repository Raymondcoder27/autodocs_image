package initializers

import (
	"example/pdfgenerator/models"
	"fmt"
	"log"
	"os"

	// "gorm.io/driver/mysql"
	// _ "github.com/lib/pq"
	"gorm.io/driver/postgres"

	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	// dsn := os.Getenv("DB")

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

}

func MigrateDB() {
	err := DB.AutoMigrate(&models.Document{})
	if err != nil {
		log.Printf("Error migrating database: %v", err)
	}

	err2 := DB.AutoMigrate(&models.Template{})
	if err2 != nil {
		log.Printf("Error migrating database: %v", err)
	}

	err3 := DB.AutoMigrate(&models.Logs{})
	if err3 != nil {
		log.Printf("Error migrating database: %v", err)
	}
}
