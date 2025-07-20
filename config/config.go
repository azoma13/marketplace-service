package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		APP
		HTTP
		PG
		JWT
		Hasher
	}

	APP struct {
		Name                  string   `env:"APP_NAME"`
		Version               string   `env:"APP_VERSION"`
		AllowedFileExtensions []string `env:"APP_ALLOWED_FILE_EXTENSIONS"`
		MaxImageSize          int64    `env:"APP_MAX_IMAGE_SIZE"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT"`
	}

	PG struct {
		MaxPoolSize int    `env:"PG_MAX_POOL_SIZE"`
		URL         string `env:"PG_URL"`
	}

	JWT struct {
		SignKey  string        `env:"JWT_SIGN_KEY"`
		TokenTTL time.Duration `env:"JWT_TOKEN_TTL"`
	}

	Hasher struct {
		Salt string `env:"HASHER_SALT"`
	}
)

var Cfg *Config

func NewConfig() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("unable to load .env file: %w", err)
	}

	Cfg = &Config{}
	if err := env.Parse(Cfg); err != nil {
		return fmt.Errorf("error parse env: %w", err)
	}

	return nil
}
