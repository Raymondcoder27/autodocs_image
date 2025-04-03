package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	// err := godotenv.Load()
	godotenv.Load()

	// if err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	// if err != nil {
	// 	fmt.Print("Error loading .env file")
	// }

	TemplateBucket := os.Getenv("TEMPLATE_BUCKET")
	PdfBucket := os.Getenv("PDF_BUCKET")

	// Optional: log the values for confirmation
	log.Printf("Loaded Template Bucket: %s", TemplateBucket)
	log.Printf("Loaded PDF Bucket: %s", PdfBucket)
}
