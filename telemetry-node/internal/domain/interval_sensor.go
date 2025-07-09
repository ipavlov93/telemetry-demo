package domain

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"go.uber.org/zap"
)

const (
	defaultBufferSize = 1 // turn off batch support by default
	defaultInterval   = time.Second
)

// IntervalSensor simulates real sensor that emits values with given interval.
type IntervalSensor struct {
	sensorName string
	// delay duration between send operation
	interval time.Duration
	// number of values to send per tick
	bufferSize int
	// function that calculates integers that are set to SensorValue.Value
	generateFunc func() int64
	logger       logger.Logger
}

// NewRateSensor returns pointer to created instance of IntervalSensor.
// This constructor is used to create instance only with defined ratePerSecond (with delay).
// Constructor will return error if:
// - ratePerSecond is below or equals zero;
// - generateFunc is nil.
// If any optional parameters are zero-valued, the constructor will assign default values to the corresponding fields.
// It skips logger validation.
func NewRateSensor(
	generateFunc func() int64,
	ratePerSecond float32,
	name string,
	lg logger.Logger,
) (*IntervalSensor, error) {
	if ratePerSecond <= 0 {
		return nil, fmt.Errorf("can't init IntervalSensor, ratePerSecond is invalid")
	}
	interval := ratePerSecondToInterval(ratePerSecond, defaultInterval)
	return NewIntervalSensor(generateFunc, interval, 1, name, lg)
}

// NewIntervalSensor returns pointer to created instance of IntervalSensor.
// This constructor allows to create instance with zero interval (without delay).
// Constructor will return error if:
// - interval is below zero;
// - generateFunc is nil.
// If any optional parameters are zero-valued, the constructor will assign default values to the corresponding fields.
// It skips logger validation.
func NewIntervalSensor(
	generateFunc func() int64,
	interval time.Duration,
	batchSize int,
	name string,
	lg logger.Logger,
) (*IntervalSensor, error) {
	// required parameters
	if interval < 0 {
		return nil, fmt.Errorf("can't init IntervalSensor, interval is invalid")
	}
	if generateFunc == nil {
		return nil, fmt.Errorf("can't init IntervalSensor, generateFunc is nil")
	}

	// optional parameters
	actualBatchSize := batchSize
	if batchSize == 0 {
		batchSize = defaultBufferSize
	}
	sensorName := name
	if sensorName == "" {
		sensorName = fmt.Sprintf("sensor-%d", time.Now().Unix())
	}

	return &IntervalSensor{
		sensorName:   sensorName,
		interval:     interval,
		bufferSize:   actualBatchSize,
		generateFunc: generateFunc,
		logger:       lg,
	}, nil
}

// Run starts producing data at a constant rate in a separate goroutine.
// Method returns error if IntervalSensor hasn't been set properly.
// The measurement process stops when the context is done (via <-ctx.Done()).
// Notice: valuesChan channel is passed as sender-only/receive-only to avoid possible deadlocks.
func (s *IntervalSensor) Run(ctx context.Context, wg *sync.WaitGroup) (<-chan []SensorValue, error) {
	if s.generateFunc == nil {
		return nil, fmt.Errorf("can't generate value, generateFunc is nil")
	}

	// buffered channels is used to prevent immediate block on channel send
	valuesChan := make(chan []SensorValue, 100)

	if wg != nil {
		wg.Add(1)
	}

	go func(ctx context.Context, valuesChan chan<- []SensorValue) {
		defer func() {
			if r := recover(); r != nil {
				// recovered generateFunc() panic
				s.logger.Error("generateFunc() panic recovered",
					zap.String("sensor", s.sensorName),
					zap.Any("panic", r),
				)
			}
			// ensure that it's single closer
			// receivers will not wait forever on channel close
			close(valuesChan)

			if wg == nil {
				return
			}
			wg.Done()
		}()

		// create buffer with bufferSize capacity
		buffer := make([]SensorValue, 0, s.bufferSize)

		for {
			select {
			case <-ctx.Done():
				// try to drain buffer gracefully to send all generated values
				if len(buffer) > 0 {
					valuesChan <- buffer
				}
				s.logger.Debug("IntervalSensor received context done, returning")
				return
			case <-time.After(s.interval):
				//for range s.bufferSize {
				//	buffer = append(buffer, SensorValue{
				//		SensorName: s.sensorName,
				//		Value:      s.generateFunc(),
				//		Timestamp:  time.Now(),
				//	})
				//}
				//valuesChan <- buffer

				buffer = append(buffer, SensorValue{
					SensorName: s.sensorName,
					Value:      s.generateFunc(),
					Timestamp:  time.Now(),
				})

				// send data when batchSize is reached
				if len(buffer) >= s.bufferSize {
					valuesChan <- buffer
					// clean the buffer
					buffer = make([]SensorValue, 0, s.bufferSize)
				}
			}
		}
	}(ctx, valuesChan)

	return valuesChan, nil
}

// ratePerSecondToInterval is utility function that calculates time interval (truncated to ms precision) based on given rate per second.
// It will return defaultInterval if rps is below zero.
func ratePerSecondToInterval(
	rps float32,
	defaultInterval time.Duration,
) time.Duration {
	if rps <= 0 {
		return defaultInterval
	}
	raw := float32(time.Second) / rps
	interval := time.Duration(raw).Truncate(time.Millisecond)
	return interval
}
