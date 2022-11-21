package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		Server  `yaml:"server"`
		Log     `yaml:"logger"`
		Storage `yaml:"storage"`
	}

	Server struct {
		GrpcServerPath string `yaml:"grpcserverpath" env:"GRPC_SERVER_PATH"`
	}

	Log struct {
		Level string `yaml:"log_level"   env:"LOG_LEVEL"`
	}

	Storage struct {
		UsersStoragePath   string `yaml:"users_storage_path" env:"USERS_STORAGE_PATH"`
		ObjectsStoragePath string `yaml:"objects_storage_path" env:"OBJECTS_STORAGE_PATH"`
	}
)

// LoadConfig returns app config.
func LoadConfig() (*Config, error) {
	cfg := Config{}
	err := cleanenv.ReadConfig("./internal/client/config/client_config.yml", &cfg)
	if err != nil {
		return nil, err
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
