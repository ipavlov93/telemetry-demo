package logger

import "time"

type FileLoggerConfig struct {
	FilePath           string
	FileName           string
	FileRotationPeriod time.Duration
}
