package others

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() error {
	err := godotenv.Load()
	return err
}

func ENV(key string) string {
	return os.Getenv(key)
}
