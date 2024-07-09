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

type ApexConfig struct {
	AuthKey string `yaml:"auth_key"`
}

type MarvinConfig struct {
	Marvin Config     `yaml:"marvin"`
	Apex   ApexConfig `yaml:"apex"`
}

// LoadConfig reads the configuration from a YAML file
func LoadConfig(filename string) (*MarvinConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg MarvinConfig
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
