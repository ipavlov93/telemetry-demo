package writer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	config "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
)

func createFileWriter(cfg config.FileLoggerConfig) (io.Writer, error) {
	if cfg.FilePath == "" || cfg.FileName == "" {
		return nil, fmt.Errorf("invalid file logger config: missing path or name")
	}

	fullPath := filepath.Join(cfg.FilePath, cfg.FileName)
	f, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return f, nil
}
