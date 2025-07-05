package domain

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"go.uber.org/zap"
)

// IntervalSensor simulates real sensor that emits value with given interval.
type IntervalSensor struct {
	sensorName   string
	interval     time.Duration
	generateFunc func() int64

	logger logger.Logger
}

// NewIntervalSensor constructor returns error when generateFunc is nil.
func NewIntervalSensor(
	name string,
	interval time.Duration,
	generateFunc func() int64,
	logger logger.Logger,
) (*IntervalSensor, error) {
	if generateFunc == nil {
		return nil, fmt.Errorf("can't init IntervalSensor, generateFunc is nil")
	}

	sensorName := name
	if sensorName == "" {
		sensorName = fmt.Sprintf("sensor-%s", rand.Text())
	}

	return &IntervalSensor{
		sensorName:   sensorName,
		interval:     interval,
		generateFunc: generateFunc,
		logger:       logger,
	}, nil
}

// Run starts producing data at a constant rate in a separate goroutine.
// Method returns error if IntervalSensor hasn't been set properly.
// The measurement process stops when the context is done (via <-ctx.Done())
func (s IntervalSensor) Run(ctx context.Context) (<-chan *SensorValue, error) {
	if s.generateFunc == nil {
		return nil, fmt.Errorf("can't generate value, generateFunc is nil")
	}

	// unbuffered channels is used to prevent immediate block on channel send
	valuesChan := make(chan *SensorValue, 100)

	go func(ctx context.Context, valuesChan chan<- *SensorValue) {
		defer func() {
			if r := recover(); r != nil {
				// recovered generateFunc() panic
				s.logger.Error("generateFunc() panic recovered", zap.Any("panic", r))
			}
			// ensure consumers/receivers will not wait forever
			close(valuesChan)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(s.interval):
				valuesChan <- &SensorValue{
					SensorName: s.sensorName,
					Value:      s.generateFunc(),
					Timestamp:  time.Now(),
				}
			}
		}
	}(ctx, valuesChan)

	return valuesChan, nil
}
