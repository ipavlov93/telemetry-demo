package writer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	pb "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const defaultFileNameFormat = "telemetry-sink-%d.log.json"

// JsonWriter represents component that writes telemetry messages to a JSON file.
type JsonWriter struct {
	fileName string
	filePath string
	logger   logger.Logger
}

// NewJsonWriter returns pointer to created instance of JsonWriter.
func NewJsonWriter(
	fName string,
	fPath string,
	lg logger.Logger,
) *JsonWriter {
	fileName := fName
	if fileName == "" {
		fileName = fmt.Sprintf(defaultFileNameFormat, time.Now().Unix())
	}

	filePath := fPath
	if filePath == "" {
		filePath = "." // default path is current directory
	}

	return &JsonWriter{
		fileName: fileName,
		filePath: filePath,
		logger:   lg,
	}
}

// Run starts process messages in a separate goroutine.
// It writes messages to json file on every channel receive.
// Run respects context cancellation (e.g., via <-ctx.Done()) and wait group by design.
// It writes each message batch as a separate log line with timestamp and values.
// Notice: actual logs format is different from JSON.
func (w *JsonWriter) Run(ctx context.Context, inputChan <-chan []*pb.SensorValue, wg *sync.WaitGroup) {
	if wg != nil {
		wg.Add(1)
	}

	go func() {
		defer func() {
			if wg != nil {
				wg.Done()
			}
		}()

		// Create or open the file for appending
		fullPath := filepath.Join(w.filePath, w.fileName)
		file, err := os.OpenFile(fullPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			w.logger.Error("failed to open file",
				zap.String("full_file_path", fullPath),
				zap.Error(err))
			return
		}
		defer file.Close()

		// fast solution
		// create separate logger to write to file
		cfg := zap.NewProductionEncoderConfig()
		cfg.EncodeTime = zapcore.RFC3339TimeEncoder

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(cfg),
			zapcore.AddSync(file), // Destination where logs are written
			zap.InfoLevel,
		)

		fileLogger := logger.NewWithCore(core)
		defer fileLogger.Sync()

		for {
			select {
			case <-ctx.Done():
				return

			case messages, ok := <-inputChan:
				if !ok {
					// channel closed
					return
				}

				if len(messages) == 0 {
					continue // skip empty batches
				}

				// write to file using logger
				fileLogger.Info("sensor values batch received",
					zap.Int64("timestamp", time.Now().Unix()),
					zap.Int("messageCount", len(messages)),
					zap.Any("values", messages),
				)
			}
		}
	}()

	return
}
