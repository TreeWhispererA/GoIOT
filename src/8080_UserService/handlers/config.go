package middlewares

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// DotEnvVariable -> get .env
func DotEnvVariable(key string) string {

	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Println("Getting value: " + key + ": " + os.Getenv(key))
	return os.Getenv(key)
}
