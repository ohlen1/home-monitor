package main

import (
	"log/slog"
	"os"

	"codingminds.com/homemmonitor/config"
	"codingminds.com/homemmonitor/metrics"
	"codingminds.com/homemmonitor/tibber"
)

var cfg *config.Config

func main() {
	cfg = config.LoadConfig("config.yaml")

	initLogger(cfg.LogLevel)
	metrics.Init(*cfg)
	tibber.Init(*cfg)
}

func initLogger(level string) {
	var slogLevel slog.Level
	switch level {
	case "DEBUG":
		slogLevel = slog.LevelDebug
	case "INFO":
		slogLevel = slog.LevelInfo
	case "WARN":
		slogLevel = slog.LevelWarn
	case "ERROR":
		slogLevel = slog.LevelError
	default:
		slog.Warn("Log level not recognized. Defaulting to INFO")
		slogLevel = slog.LevelInfo
	}
	logger := *slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel}))
	slog.SetDefault(&logger)
	slog.Info("Log level set", "level", slogLevel)
}
