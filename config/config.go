package config

import (
	"os"

	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v3"
)

// QbotConfig represents the configuration structure
type QbotConfig struct {
	AppID     string `yaml:"app_id"`
	Token     string `yaml:"token"`
	AppSecret string `yaml:"secret"`
}

type ApexConfig struct {
	AuthKey string `yaml:"auth_key"`
}

type MarvinConfig struct {
	Marvin       QbotConfig    `yaml:"marvin"`
	Apex         ApexConfig    `yaml:"apex"`
	ShortKey     string        `yaml:"short_key"`
	Redis        redis.Options `yaml:"redis"`
	MapHtml      string        `yaml:"map_html"`
	ImgPath      string        `yaml:"img_path"`
	ImgUrlPrefix string        `yaml:"img_url_prefix"`
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
