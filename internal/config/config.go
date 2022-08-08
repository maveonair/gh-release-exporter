package config

import (
	"time"

	"github.com/BurntSushi/toml"
	"github.com/maveonair/gh-release-exporter/internal/release"
)

const (
	defaultInterval = 24 * time.Hour
)

type Config struct {
	Interval      time.Duration
	ListeningAddr string                     `toml:"listening_addr"`
	Releases      map[string]release.Release `toml:"releases"`
}

func NewConfig(configFilePath string) (*Config, error) {
	var config Config
	if _, err := toml.DecodeFile(configFilePath, &config); err != nil {
		return nil, err
	}

	config.Interval = defaultInterval

	return &config, nil
}
