package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
)

const (
	hostKey       = "WEB_HOST"
	portKey       = "WEB_PORT"
	dbDSNKey      = "DB_DSN"
	jwtSigningKey = "JWT_SIGNING_KEY"
)

type Config struct {
	Host          string
	Port          int
	DbDSN         string
	JWTSigningKey string
}

func (c *Config) setDefaults() {
	if c.Host == "" {
		c.Host = "0.0.0.0"
	}
	if c.Port == 0 {
		c.Port = 8080
	}
}

func (c *Config) validate() error {
	if c.DbDSN == "" {
		return errors.New("DB_DSN must be present")
	}

	if c.JWTSigningKey == "" {
		return errors.New("JWT_SIGNING_KEY must be present")
	}

	return nil
}

func (c *Config) log() {
	slog.Info(
		"app config",
		slog.String("host", c.Host),
		slog.String("port", strconv.Itoa(c.Port)),
	)
}

func New() (*Config, error) {
	portString := os.Getenv(portKey)
	if portString == "" {
		portString = "8080"
	}

	portInt, err := strconv.Atoi(portString)
	if err != nil {
		return nil, err
	}

	c := Config{
		Host:          os.Getenv(hostKey),
		Port:          portInt,
		DbDSN:         os.Getenv(dbDSNKey),
		JWTSigningKey: os.Getenv(jwtSigningKey),
	}

	c.setDefaults()
	if err = c.validate(); err != nil {
		return nil, err
	}

	c.log()

	return &c, nil
}
