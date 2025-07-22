package logger

import "go.uber.org/zap/zapcore"

// ParseLevel parses given level or set default level.
func ParseLevel(level string, fallback zapcore.Level) zapcore.Level {
	if level == "" {
		return fallback
	}
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		// set default log level
		logLevel = fallback
	}
	return logLevel
}
