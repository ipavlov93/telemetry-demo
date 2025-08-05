package writer

import (
	"fmt"
	"io"
	"os"

	config "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
)

type Factory func(config.Configuration) (io.Writer, error)

func NewLogWriter(cfg config.Configuration) (io.Writer, error) {
	switch cfg.Destination {
	case config.LogOutputStdout:
		return os.Stdout, nil

	case config.LogOutputFile:
		if cfg.File == nil {
			return nil, fmt.Errorf("file config missing for file output")
		}
		return createFileWriter(*cfg.File)

	default:
		return nil, fmt.Errorf("unsupported log output type: %s", cfg.Destination)
	}
}
