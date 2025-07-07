package domain

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"go.uber.org/zap"
)

const defaultBatchSize = 1 // turn off batch support by default

// IntervalSensor simulates real sensor that emits value with given interval.
type IntervalSensor struct {
	sensorName   string
	interval     time.Duration
	batchSize    int64
	generateFunc func() int64
	logger       logger.Logger
}

// NewIntervalSensor returns pointer to created instance of IntervalSensor.
// Constructor will return error if:
// - interval is non-positive;
// - generateFunc is nil.
// If any optional parameters are zero-valued, the constructor will assign default values to the corresponding fields.
func NewIntervalSensor(
	generateFunc func() int64,
	interval time.Duration,
	batchSize int64,
	name string,
	inputLogger logger.Logger,
) (*IntervalSensor, error) {
	// required parameters
	if interval == 0 {
		return nil, fmt.Errorf("can't init IntervalSensor, interval is zero")
	}
	if generateFunc == nil {
		return nil, fmt.Errorf("can't init IntervalSensor, generateFunc is nil")
	}

	// optional parameters
	actualBatchSize := batchSize
	if batchSize == 0 {
		batchSize = defaultBatchSize
	}
	sensorName := name
	if sensorName == "" {
		sensorName = fmt.Sprintf("sensor-%d", time.Now().Unix())
	}
	actualLogger := inputLogger
	if inputLogger == nil || reflect.ValueOf(inputLogger).IsNil() {
		actualLogger = logger.NewNopLogger()
	}

	return &IntervalSensor{
		sensorName:   sensorName,
		interval:     interval,
		batchSize:    actualBatchSize,
		generateFunc: generateFunc,
		logger:       actualLogger,
	}, nil
}

// Run starts producing data at a constant rate in a separate goroutine.
// Method returns error if IntervalSensor hasn't been set properly.
// The measurement process stops when the context is done (via <-ctx.Done()).
// Notice: valuesChan channel is passed as sender-only/receive-only to avoid possible deadlocks.
func (s *IntervalSensor) Run(ctx context.Context) (<-chan []SensorValue, error) {
	if s.generateFunc == nil {
		return nil, fmt.Errorf("can't generate value, generateFunc is nil")
	}

	// buffered channels is used to prevent immediate block on channel send
	valuesChan := make(chan []SensorValue, 100)

	go func(ctx context.Context, valuesChan chan<- []SensorValue) {
		defer func() {
			if r := recover(); r != nil {
				// recovered generateFunc() panic
				s.logger.Error("generateFunc() panic recovered",
					zap.String("sensor", s.sensorName),
					zap.Any("panic", r),
				)
			}
			// ensure receivers will not wait forever
			close(valuesChan)
		}()

		// create buffer with batchSize capacity
		buffer := make([]SensorValue, 0, s.batchSize)

		for {
			select {
			case <-ctx.Done():
				// try to drain buffer gracefully to send all generated values
				if len(buffer) > 0 {
					valuesChan <- buffer
				}
				return
			case <-time.After(s.interval):
				buffer = append(buffer, SensorValue{
					SensorName: s.sensorName,
					Value:      s.generateFunc(),
					Timestamp:  time.Now(),
				})

				// send data when batchSize is reached
				if int64(len(buffer)) >= s.batchSize {
					valuesChan <- buffer
					// clean the buffer
					buffer = make([]SensorValue, 0, s.batchSize)
				}
			}
		}
	}(ctx, valuesChan)

	return valuesChan, nil
}
