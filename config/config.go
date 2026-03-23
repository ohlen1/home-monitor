package config

import (
	"log/slog"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Config struct {
	ServiceName    string `yaml:"serviceName"`
	LogLevel       string `yaml:"logLevel"`
	MetricsEnabled bool   `yaml:"metricsEnabled"`
	Parallelism    int    `yaml:"parallelism"`
	Tibber         struct {
		Url          string   `yaml:"url"`
		HomeId       string   `yaml:"homeId"`
		ApiToken     string   `yaml:"apiToken"`
		Measurements []string `yaml:"measurements"`
	} `yaml:"tibber"`
}

func LoadConfig(path string) *Config {
	data, err := os.ReadFile(path)
	if err != nil {
		slog.Error("Failed to read config", "path", path, "error", err)
		os.Exit(1)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		slog.Error("Failed to unmarshal config", "path", path, "error", err)
		os.Exit(1)
	}

	slog.Info("Config initialized", "config", cfg)
	return &cfg
}
