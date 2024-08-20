package initializers

import (
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
}
