package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnv(key string) string {

	err := godotenv.Load()

	if err != nil {
		log.Println("Error loading .env file, will use system environment variables.")
	}

	return os.Getenv(key)
}
