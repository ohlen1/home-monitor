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
	data, _ := os.ReadFile(path)
	y := Config{}
	err := yaml.Unmarshal([]byte(data), &y)
	if err != nil {
		slog.Error("Error umarshalling config", "error", err)
	}
	slog.Info("Config initialized", "config", y)
	return &y
}
