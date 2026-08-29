package config

import (
	"fmt"
	"os"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App AppConfig `yaml:"app"`
	HTTP HTTPConfig `yaml:"http"`
	Logging LoggingConfig `yaml:"logging"`
}

type AppConfig struct {
	Name string `yaml:"name"`
	Environment string `yaml:"environment"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port int `yaml:"port"`
}

type LoggingConfig struct {
	Level string `yaml:"level"`
	Format string `yaml:"format"`
}

func Load(path string) (Config, error) {
	fileContent, err := os.ReadFile(path)
	if err!=nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(fileContent, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse config file: %w", err)
	}

	applyDefaults(&cfg)

	return cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.App.Name == "" {
		cfg.App.Name = "StreamOps"
	}

	if cfg.App.Environment == "" {
		cfg.App.Environment = "local"
	}

	if cfg.HTTP.Host == "" {
		cfg.HTTP.Host = "127.0.0.1"
	}

	if cfg.HTTP.Port == 0 {
		cfg.HTTP.Port = 7317
	}

	if cfg.Logging.Level == "" {
		cfg.Logging.Level = "info"
	}

	if cfg.Logging.Format == "" {
		cfg.Logging.Format = "text"
	}
}
