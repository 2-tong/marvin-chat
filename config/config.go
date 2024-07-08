package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

// Config represents the configuration structure
type Config struct {
	AppID     uint64 `yaml:"app_id"`
	Token     string `yaml:"token"`
	AppSecret string
}

type MarvinConfig struct {
	Marvin Config `json:"marvin"`
}

// LoadConfig reads the configuration from a YAML file
func LoadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg MarvinConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg.Marvin, nil
}
