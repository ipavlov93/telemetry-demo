package logger

import (
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
	options = append([]zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}, options...)

	return &ZapLogger{
		logger: zap.New(core, options...),
	}
}

// NewNopLogger never writes out logs or internal errors.
func NewNopLogger() *ZapLogger {
	return &ZapLogger{logger: zap.NewNop()}
}
