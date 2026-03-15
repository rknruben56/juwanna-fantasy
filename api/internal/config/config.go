package config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseHost     string `envconfig:"DATABASE_HOST" default:"localhost"`
	DatabasePort     int    `envconfig:"DATABASE_PORT" default:"5432"`
	DatabaseUser     string `envconfig:"DATABASE_USER" default:"postgres"`
	DatabasePassword string `envconfig:"DATABASE_PASSWORD" default:"postgres"`
	DatabaseName     string `envconfig:"DATABASE_NAME" default:"juwanna_fantasy"`
	DatabaseSSLMode  string `envconfig:"DATABASE_SSL_MODE" default:"disable"`
	DatabaseMaxConns int    `envconfig:"DATABASE_MAX_CONNS" default:"10"`
	DatabaseMinConns int    `envconfig:"DATABASE_MIN_CONNS" default:"2"`
	ServerPort       int    `envconfig:"SERVER_PORT" default:"8080"`
	SleeperLeagueID  string `envconfig:"SLEEPER_LEAGUE_ID"`
	LogLevel         string `envconfig:"LOG_LEVEL" default:"info"`
	CORSAllowedOrigins string `envconfig:"CORS_ALLOWED_ORIGINS" default:"*"`
	ReadTimeout      int    `envconfig:"READ_TIMEOUT" default:"10"`
	WriteTimeout     int    `envconfig:"WRITE_TIMEOUT" default:"10"`
	ShutdownTimeout  int    `envconfig:"SHUTDOWN_TIMEOUT" default:"10"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		c.DatabaseUser, c.DatabasePassword,
		c.DatabaseHost, c.DatabasePort, c.DatabaseName,
		c.DatabaseSSLMode,
	)
}

// CORSOrigins returns the allowed origins as a slice.
func (c *Config) CORSOrigins() []string {
	origins := strings.Split(c.CORSAllowedOrigins, ",")
	for i := range origins {
		origins[i] = strings.TrimSpace(origins[i])
	}
	return origins
}
