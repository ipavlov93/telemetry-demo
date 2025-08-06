package factory

import (
	"github.com/ipavlov93/telemetry-demo/pkg/env"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
)

func NewDefaultLoggerConfig() logger.ConfigMap {
	return logger.ConfigMap{
		"APP_LOGGER": logger.Configuration{
			Enabled:     env.ParseBoolEnv("LOGGER_ENABLED", true),
			MinLevel:    env.EnvironmentVariable("LOGGER_MIN_LOG_LEVEL", "info"),
			Destination: logger.LogOutputStdout,
		},
	}
}
