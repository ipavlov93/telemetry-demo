package factory

import (
	"github.com/ipavlov93/telemetry-demo/pkg/env"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
)

func NewDefaultLoggerConfig(key string) logger.ConfigMap {
	return logger.ConfigMap{
		key: NewStdoutLoggerConfig(),
	}
}

func NewStdoutLoggerConfig() logger.Configuration {
	return logger.Configuration{
		Enabled:     env.ParseBoolEnv("LOGGER_ENABLED", true),
		MinLevel:    env.EnvironmentVariable("LOGGER_MIN_LOG_LEVEL", "info"),
		Destination: logger.LogOutputStdout,
	}
}
