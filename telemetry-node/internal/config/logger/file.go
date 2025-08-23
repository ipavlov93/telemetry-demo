package logger

type FileLoggerConfig struct {
	FilePath           string `json:"filePath"`
	FileName           string `json:"fileName"`
	FileRotationPeriod string `json:"fileRotationPeriod"`
}
