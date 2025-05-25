// Package config provides main application configuration capabilities.
package config

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
)

// Config represents application configuration.
type Config struct {
	GRPC struct {
		Port int `koanf:"port"`
	} `koanf:"grpc"`
	Redis struct {
		Host     string        `koanf:"host"`
		Port     int           `koanf:"port"`
		Username string        `koanf:"username"`
		Password string        `koanf:"password"` // FIXME: replace with Vault-readed value.
		TTL      time.Duration `koanf:"ttl"`
	} `koanf:"redis"`
}

// New creates new instance of config with default values.
//
// Depends on application mode parses different config files.
func New(mode Mode, log *slog.Logger) (*Config, error) {
	c := &Config{}

	c.loadDefaults()

	var configPath = "configs/config-dev.yml"
	if mode == ModeProduction {
		configPath = "configs/config-prod.yml"
	}

	const delim = ":"
	k := koanf.New(delim)
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("load yaml config: %w", err)
	}

	if err := k.Load(env.Provider("HASHER_", delim, func(s string) string {
		return strings.Replace(strings.ToLower(strings.TrimPrefix(s, "HASHER_")), "_", delim, -1)
	}), nil); err != nil {
		return nil, fmt.Errorf("load env config: %w", err)
	}

	if err := k.Unmarshal("", c); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	log.Info("config initialized", slog.String("mode", string(mode)))

	return c, nil
}

// Mode represents application mode.
type Mode string

// Application modes.
const (
	ModeDevelopment Mode = "dev"
	ModeProduction  Mode = "prod"
)

func (c *Config) loadDefaults() {
	c.GRPC.Port = 6969
	c.Redis.Host = "127.0.0.1"
	c.Redis.Port = 6379
	c.Redis.TTL = 5 * time.Minute
	c.Redis.Username = "default"
	c.Redis.Password = "1234qwerASDF"
}
