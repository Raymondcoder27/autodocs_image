package initializers

import (
	"example/pdfgenerator/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	// dsn := os.Getenv("DB")

	dbHost := os.Getenv("AUTODOCS_DB_HOST")
	dbUser := os.Getenv("AUTODOCS_DB_USER")
	dbPassword := os.Getenv("AUTODOCS_DB_PASSWORD")
	dbName := os.Getenv("AUTODOCS_DB_NAME")
	dbSchema := os.Getenv("AUTODOCS_DB_SCHEMA")

	if dbSchema == "" {
		dbSchema = "public"
	}

	// dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbName)

	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Printf("Failed to connect to database: %v", err)
	// }

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s.", dbSchema),
			SingularTable: false,
		},
	})

	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Fatal("Failed to connect to database")
	}

	log.Printf("Connected to database: %v", DB)

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

	err4 := DB.AutoMigrate(&models.FailedGenerations{})
	if err4 != nil {
		log.Printf("Error migrating database: %v", err)
	}
}
