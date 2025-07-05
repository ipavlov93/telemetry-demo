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

// ZapLogger is simple wrapper for zap.Logger
// ZapLogger implements Logger to make it possible swap logger in the future.
type ZapLogger struct {
	logger *zap.Logger
}

// New constructor creates ZapLogger wrapper with configured underlined logger.
func New(w io.Writer, logLevel zapcore.Level) ZapLogger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(w), // Destination where logs are written
		logLevel,           // Minimum level to log
	)

	logger := zap.New(core, zap.AddCaller()).WithOptions(zap.AddCallerSkip(1))

	return ZapLogger{logger: logger}
}

// Sync calls the underlying Sync method, flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func (z *ZapLogger) Sync() error { return z.logger.Sync() }

func (z *ZapLogger) Info(msg string, fields ...zap.Field)  { z.logger.Info(msg, fields...) }
func (z *ZapLogger) Warn(msg string, fields ...zap.Field)  { z.logger.Warn(msg, fields...) }
func (z *ZapLogger) Debug(msg string, fields ...zap.Field) { z.logger.Debug(msg, fields...) }
func (z *ZapLogger) Error(msg string, fields ...zap.Field) { z.logger.Error(msg, fields...) }
