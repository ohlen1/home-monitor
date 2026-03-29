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
		Enabled      bool     `yaml:"enabled"`
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

	apiToken, _ := os.LookupEnv("TIBBER_API_TOKEN")
	if apiToken == "" {
		slog.Warn("TIBBER_API_TOKEN env variable not set")
	} else {
		slog.Info("Tibber API token provided via env")
		cfg.Tibber.ApiToken = apiToken
	}

	homeId, _ := os.LookupEnv("TIBBER_HOME_ID")
	if homeId == "" {
		slog.Warn("TIBBER_HOME_ID env variable not set")
	} else {
		slog.Info("Tibber Home ID provided via env")
		cfg.Tibber.HomeId = homeId
	}

	if len(cfg.Tibber.ApiToken) == 0 || len(cfg.Tibber.HomeId) == 0 {
		panic("Tibber configuration incomplete")
	}

	slog.Info("Config initialized", "config", cfg)
	return &cfg
}
