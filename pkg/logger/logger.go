package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger contains set of method to log a messages.
type Logger interface {
	// Sync flushing any buffered log entries.
	// Applications should take care to call Sync before exiting.
	Sync() error

	// Info , Warn, Debug, Error method names indicate to log level message will be logged.
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
}

// ZapLogger is simple wrapper around underlying zap.Logger.
// ZapLogger implements Logger to make it possible swap logger in the future.
type ZapLogger struct {
	logger *zap.Logger
}

// NewNopLogger returns ZapLogger with a no-op underlined logger that never writes out logs or internal errors.
func NewNopLogger() *ZapLogger {
	return &ZapLogger{logger: zap.NewNop()}
}

// New creates ZapLogger with a configured underlying zap.Logger.
// The logLevel parameter specifies the minimum level to log; messages below this level will be ignored.
// Logs are written to the given io.Writer using JSON encoding with RFC3339 timestamps.
func New(w io.Writer, logLevel zapcore.Level) *ZapLogger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(w), // Destination where logs are written
		logLevel,           // Minimum level to log
	)

	logger := zap.New(core, zap.AddCaller()).WithOptions(zap.AddCallerSkip(1))

	return &ZapLogger{logger: logger}
}

// NewWithCore creates a ZapLogger with given zapcore.Core and zap.Option slice.
func NewWithCore(core zapcore.Core, options ...zap.Option) *ZapLogger {
	logger := zap.New(core, options...).WithOptions(zap.AddCallerSkip(1))
	return &ZapLogger{logger: logger}
}

// Sync calls the underlying Sync method, flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func (z *ZapLogger) Sync() error { return z.logger.Sync() }

func (z *ZapLogger) Info(msg string, fields ...zap.Field)  { z.logger.Info(msg, fields...) }
func (z *ZapLogger) Warn(msg string, fields ...zap.Field)  { z.logger.Warn(msg, fields...) }
func (z *ZapLogger) Debug(msg string, fields ...zap.Field) { z.logger.Debug(msg, fields...) }
func (z *ZapLogger) Error(msg string, fields ...zap.Field) { z.logger.Error(msg, fields...) }

// ParseLevel utility function tries to parse given level or set default level.
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
