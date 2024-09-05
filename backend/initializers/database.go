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
}

// func ConnectToDB() {
// 	var err error
// 	user := os.Getenv("DB_USER")
// 	password := os.Getenv("DB_PASS")
// 	host := os.Getenv("DB_HOST")
// 	port := os.Getenv("DB_PORT")
// 	dbName := os.Getenv("DB_NAME")

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		user, password, host, port, dbName)
// 	// }

// 	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		log.Fatalf("Failed to connect to database: %v", err)
// 	}
// }

// func ConnectToDB() {
// 	var err error
// 	dsn := os.Getenv("DB")

// 	// Retry logic to handle database initialization
// 	for i := 0; i < 15; i++ { // Increased retries
// 		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
// 		if err == nil {
// 			// If connection is successful, break out of loop
// 			log.Println("Successfully connected to database.")
// 			break
// 		}

// 		log.Printf("Failed to connect to database, retrying in 10 seconds... (%d/15)", i+1)
// 		time.Sleep(10 * time.Second) // Increased wait time
// 	}

// 	if err != nil {
// 		log.Fatalf("Failed to connect to database after retries: %v", err)
// 	}
// }
