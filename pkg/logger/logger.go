package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(logFile string) *zap.Logger {
	cfg := zap.NewProductionEncoderConfig()
	cfg.EncodeTime = zapcore.RFC3339TimeEncoder

	//core := zapcore.NewCore(
	//	zapcore.NewJSONEncoder(cfg),
	//zapcore.AddSync(w), // Where logs are written
	//	zapcore.InfoLevel, // Minimum level to log (Info and above)
	//)

	// Write to stdout
	consoleSync := zapcore.AddSync(os.Stdout)

	// Write to file
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return nil
	}
	fileSync := zapcore.AddSync(file)

	// Combine outputs using zapcore.NewTee
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(cfg), consoleSync, zapcore.DebugLevel),
		zapcore.NewCore(zapcore.NewJSONEncoder(cfg), fileSync, zapcore.InfoLevel),
	)

	return zap.New(core, zap.AddCaller())
}
