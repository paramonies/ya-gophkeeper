package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server `yaml:"server"`
		Log    `yaml:"logger"`
		DB     `yaml:"db"`
	}

	Server struct {
		Address string `yaml:"address" env:"ADDRESS"`
	}

	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}

	DB struct {
		DNS            string `yaml:"dns"   env:"DNS"`
		ConnectTimeout int    `yaml:"connect_timeout" env:"CONNECTION_TIMEOUT"`
		QueryTimeout   int    `yaml:"query_timeout" env:"QUERY_TIMEOUT"`
	}
)

// LoadConfig returns app config.
func LoadConfig() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig("./internal/server/config/server_config.yml", &cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
