package config

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME"`
	ServicePort int    `envconfig:"SERVICE_PORT" default:"8080"`

	DBHost     string `envconfig:"DB_HOST" default:"localhost"`
	DBPort     int    `envconfig:"DB_PORT" default:"5432"`
	DBUsername string `envconfig:"DB_USERNAME"`
	DBPassword string `envconfig:"DB_PASSWORD"`
	DBName     string `envconfig:"DB_NAME"`
	DBSSL      string `envconfig:"DB_SSL" default:"disable"`

	BcryptCost int `envconfig:"BCRYPT_COST" default:"10"`

	JwtPrivateKey             string        `envconfig:"JWT_PRIVATE_KEY"`
	JwtPublicKey              string        `envconfig:"JWT_PUBLIC_KEY"`
	JwtAlgorithm              string        `envconfig:"JWT_ALGORITHM"`
	JwtAccessTokenExpiration  time.Duration `envconfig:"JWT_ACCESS_TOKEN_EXPIRATION"`
	JwtRefreshTokenExpiration time.Duration `envconfig:"JWT_REFRESH_TOKEN_EXPIRATION"`

	PasswordValidation string `envconfig:"PASSWORD_VALIDATION"`
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DBUsername,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSL,
	)
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("./config/.env")
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
