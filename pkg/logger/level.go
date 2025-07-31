package logger

import (
	"fmt"

	"go.uber.org/zap/zapcore"
)

// ParseLevelOrDefault parses given level or set default level.
func ParseLevelOrDefault(level string, fallback zapcore.Level) zapcore.Level {
	logLevel, err := ParseLevel(level)
	if err != nil {
		return fallback
	}
	return logLevel
}

func ParseLevel(level string) (zapcore.Level, error) {
	if level == "" {
		return 0, fmt.Errorf("empty log level")
	}
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return 0, err
	}
	return logLevel, nil
}
