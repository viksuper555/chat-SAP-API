package config

import (
	"github.com/joho/godotenv"
	"github.com/vrischmann/envconfig"
	"log"
)

type Config struct {
	USER2    string `envconfig:"default=postgres"`
	Database struct {
		Dialect  string `envconfig:"default=postgres"`
		User     string `envconfig:"default=postgres"`
		Password string `envconfig:"default=postgres"`
		Name     string `envconfig:"default=notifications"`
		IP       string `envconfig:"optional"`
	} `env:"DB_"`
}

// NewFromEnv creates a Config from environment values
func NewFromEnv() Config {
	appConfig := Config{}
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading env file: %v\n", err)
	}
	err = envconfig.Init(&appConfig)
	if err != nil {
		log.Printf("Error loading env file: %v\n", err)
	}
	return appConfig
}
