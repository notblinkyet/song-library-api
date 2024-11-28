package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DbPort, ServerPort                                                        int
	DbUser, DbName, DbHost, DbPassword, ApiAddrURL, MigrationPath, ServerHost string
	Timeout, IdleTimeout                                                      time.Duration
}

func LoadConfig() (*Config, error) {
	const op = "config.LoadConfig"

	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	timeOut, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	idleTimeout, err := strconv.Atoi(os.Getenv("IDLE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Config{
		DbPort:        dbPort,
		ServerPort:    serverPort,
		DbUser:        os.Getenv("DB_USER"),
		DbName:        os.Getenv("DB_NAME"),
		DbHost:        os.Getenv("DB_HOST"),
		DbPassword:    os.Getenv("DB_PASSWORD"),
		ApiAddrURL:    os.Getenv("API_ADDR_URL"),
		MigrationPath: os.Getenv("MIGRATION_PATH"),
		ServerHost:    os.Getenv("SERVER_HOST"),
		Timeout:       time.Duration(timeOut) * time.Second,
		IdleTimeout:   time.Duration(idleTimeout) * time.Second,
	}, nil
}

func MustLoadConfig() *Config {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	return config
}
