package config

import (
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
	"os"
)

// QbotConfig represents the configuration structure
type QbotConfig struct {
	AppID     uint64 `yaml:"app_id"`
	Token     string `yaml:"token"`
	AppSecret string
}

type ApexConfig struct {
	AuthKey string `yaml:"auth_key"`
}

type MarvinConfig struct {
	Marvin   QbotConfig    `yaml:"marvin"`
	Apex     ApexConfig    `yaml:"apex"`
	ShortKey string        `yaml:"short_key"`
	Redis    redis.Options `yaml:"redis"`
}

var config *MarvinConfig = nil

func GetConfig() *MarvinConfig {
	if config == nil {
		var err error = nil
		config, err = LoadConfig()
		if err != nil {
			panic(err)
		}
	}
	return config
}

// LoadConfig reads the configuration from a YAML file
func LoadConfig() (*MarvinConfig, error) {
	cfgPath := "./marvin.yml"

	if len(os.Args) > 1 {
		cfgPath = os.Args[1]
	}

	data, err := os.ReadFile(cfgPath)
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
