package services

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ConfigService interface {
	Get(key string) (string, error)
}

type configService struct{}

func NewConfigService() ConfigService {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file found, using system environment variables")
	}
	return &configService{}
}

func (c *configService) Get(key string) (string, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return "", errors.New("value not found for key: " + key)
	}
	return value, nil
}
