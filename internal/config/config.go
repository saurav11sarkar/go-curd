package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI   string
	MongoDB    string
	ServerPort string
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, fmt.Errorf("error loading .env file")
	}

	mongoURI, err := extractEnv("MONGO_URI")
	if err != nil {
		return Config{}, err
	}
	serverPort, err := extractEnv("PORT")
	if err != nil {
		return Config{}, err
	}

	mongoDB, err := extractEnv("MONGO_DB")
	if err != nil {
		return Config{}, err
	}

	return Config{mongoURI, mongoDB, serverPort}, nil
}

func extractEnv(key string) (string, error) {
	val := os.Getenv(key)

	if val == "" {
		return "", fmt.Errorf("environment variable %s not set", key)
	}
	return val, nil
}
