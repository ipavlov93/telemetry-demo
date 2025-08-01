package writer

import (
	"fmt"
	"io"
	"os"

	config "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
)

type Factory func(config.LogOutput) (io.Writer, error)

func NewLogWriter(dest config.LogOutput) (io.Writer, error) {
	switch dest.Type {
	case config.LogOutputStdout:
		return os.Stdout, nil

	case config.LogOutputFile:
		if dest.File == nil {
			return nil, fmt.Errorf("file config missing for file output")
		}
		return createFileWriter(*dest.File)

	default:
		return nil, fmt.Errorf("unsupported log output type: %s", dest.Type)
	}
}
