package simulator

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/sensor"
	"go.uber.org/zap"
)

const (
	defaultBufferSize = 1 // turn off buffering by default

	defaultChannelCapacity = 100
)

// SensorSimulator simulates real sensor that emits values with given delay.
type SensorSimulator struct {
	Name         string
	generateFunc func() int64
	// delay duration between generateFunc call
	delay time.Duration
	// number of values to store before sending
	bufferSize int
	channelCap int
	logger     logger.Logger
}

// NewWithRate constructor will return error if:
// - ratePerSecond is below or equals zero;
// - generateFunc is nil.
// If any optional parameters are zero-valued, the constructor will assign default values to the corresponding fields.
func NewWithRate(
	generateFunc func() int64,
	rate sensor.SamplingRate,
	name string,
	channelCap int,
	lg logger.Logger,
) (*SensorSimulator, error) {
	return New(
		generateFunc,
		rate.Interval(),
		defaultBufferSize,
		name,
		channelCap,
		lg,
	)
}

// New constructor allows to create instance with zero delay (without delay).
// Constructor will return error if:
// - delay is below zero;
// - generateFunc is nil.
func New(
	generateFunc func() int64,
	delay time.Duration,
	bufferSize int,
	name string,
	channelCap int,
	lg logger.Logger,
) (*SensorSimulator, error) {
	if delay < 0 {
		return nil, fmt.Errorf("can't init SensorSimulator, delay is invalid")
	}
	if generateFunc == nil {
		return nil, fmt.Errorf("can't init SensorSimulator, generateFunc is nil")
	}

	// optional parameters
	actualSize := bufferSize
	if bufferSize <= 0 {
		bufferSize = defaultBufferSize
	}
	actualChannelCap := defaultChannelCapacity
	if channelCap < 0 {
		actualChannelCap = defaultChannelCapacity
	}
	sensorName := name
	if sensorName == "" {
		sensorName = fmt.Sprintf("sensor-%d", time.Now().Unix())
	}

	return &SensorSimulator{
		Name:         sensorName,
		delay:        delay,
		bufferSize:   actualSize,
		channelCap:   actualChannelCap,
		generateFunc: generateFunc,
		logger:       lg,
	}, nil
}

// Run starts producing data at a constant rate in a separate goroutine.
// The measurement process stops when the context is done (via <-ctx.Done()).
func (s *SensorSimulator) Run(ctx context.Context, wg *sync.WaitGroup) (<-chan []measurement.SensorValue, error) {
	if s.generateFunc == nil {
		return nil, fmt.Errorf("can't generate value, generateFunc is nil")
	}

	valuesChan := make(chan []measurement.SensorValue, s.channelCap)

	if wg != nil {
		wg.Add(1)
	}

	go func(ctx context.Context, valuesChan chan<- []measurement.SensorValue) {
		defer func() {
			if r := recover(); r != nil {
				s.logger.Error("generateFunc() panic recovered",
					zap.String("sensor", s.Name),
					zap.Any("panic", r),
				)
			}

			// receivers will not wait forever on channel close
			close(valuesChan)

			if wg == nil {
				return
			}
			wg.Done()
		}()

		buffer := make([]measurement.SensorValue, 0, s.bufferSize)

		for {
			select {
			case <-ctx.Done():
				if len(buffer) > 0 {
					valuesChan <- buffer
				}
				s.logger.Debug("SensorSimulator received context done, returning")
				return
			case <-time.After(s.delay):
				buffer = append(buffer, measurement.SensorValue{
					SensorName: s.Name,
					Value:      s.generateFunc(),
					Timestamp:  time.Now(),
				})

				// send data when bufferSize is reached
				if len(buffer) >= s.bufferSize {
					valuesChan <- buffer
					buffer = make([]measurement.SensorValue, 0, s.bufferSize)
				}
			}
		}
	}(ctx, valuesChan)

	return valuesChan, nil
}
