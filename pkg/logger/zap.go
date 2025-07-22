package logger

import (
	"io"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.Logger
}

// Sync flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func (z *ZapLogger) Sync() error { return z.logger.Sync() }

func (z *ZapLogger) Info(msg string, fields ...zap.Field)  { z.logger.Info(msg, fields...) }
func (z *ZapLogger) Warn(msg string, fields ...zap.Field)  { z.logger.Warn(msg, fields...) }
func (z *ZapLogger) Debug(msg string, fields ...zap.Field) { z.logger.Debug(msg, fields...) }
func (z *ZapLogger) Error(msg string, fields ...zap.Field) { z.logger.Error(msg, fields...) }
func (z *ZapLogger) Fatal(msg string, fields ...zap.Field) { z.logger.Fatal(msg, fields...) }

func NewWithCore(core zapcore.Core, options ...zap.Option) *ZapLogger {
	return &ZapLogger{
		logger: zap.New(core, options...).
			WithOptions(zap.AddCallerSkip(1)),
	}
}

// New logger uses JSON encoding with RFC3339 timestamps.
func New(w io.Writer, minLogLevel zapcore.Level) *ZapLogger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(cfg),
		zapcore.AddSync(w),
		minLogLevel,
	)

	logger := zap.New(core, zap.AddCaller()).
		WithOptions(zap.AddCallerSkip(1))

	return &ZapLogger{logger: logger}
}

// NewNopLogger never writes out logs or internal errors.
func NewNopLogger() *ZapLogger {
	return &ZapLogger{logger: zap.NewNop()}
}
