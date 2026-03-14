package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DatabaseHost     string `envconfig:"DATABASE_HOST" default:"localhost"`
	DatabasePort     int    `envconfig:"DATABASE_PORT" default:"5432"`
	DatabaseUser     string `envconfig:"DATABASE_USER" default:"postgres"`
	DatabasePassword string `envconfig:"DATABASE_PASSWORD" default:"postgres"`
	DatabaseName     string `envconfig:"DATABASE_NAME" default:"juwanna_fantasy"`
	ServerPort       int    `envconfig:"SERVER_PORT" default:"8080"`
	SleeperLeagueID  string `envconfig:"SLEEPER_LEAGUE_ID"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (c *Config) DatabaseDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.DatabaseUser, c.DatabasePassword,
		c.DatabaseHost, c.DatabasePort, c.DatabaseName,
	)
}
