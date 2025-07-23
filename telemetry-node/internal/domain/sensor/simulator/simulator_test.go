package simulator_test

import (
	"context"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/rate"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/sensor/simulator"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestIntervalSensor_Run(t *testing.T) {
	intervalSeconds := 1
	totalSeconds := 5
	batchSize := 1
	channelCap := 0
	wg := sync.WaitGroup{}

	t.Run("SensorSimulator.SendSensorValues() happy flow", func(t *testing.T) {
		intervalSensor, err := simulator.New(
			randomValue,
			time.Duration(intervalSeconds)*time.Second,
			batchSize,
			"",
			channelCap,
			zap.NewNop(),
		)
		assert.NoError(t, err)

		// context timeout duration has small approximation error or testing code call takes some time
		timeout := time.Duration(totalSeconds*1000+10) * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		valuesChan, err := intervalSensor.Run(ctx, &wg)
		assert.NoError(t, err)

		var values []measurement.SensorValue
		for value := range valuesChan {
			values = append(values, value...)
		}

		assert.Equal(t, totalSeconds/intervalSeconds, len(values))
	})
	t.Run("should gracefully send partial batch before context cancellation", func(t *testing.T) {
		batchSize = 5

		intervalSensor, err := simulator.New(
			randomValue,
			time.Duration(intervalSeconds)*time.Second,
			batchSize,
			"",
			channelCap,
			zap.NewNop(),
		)
		assert.NoError(t, err)

		milliseconds := 4500
		// context timeout duration has small approximation error or testing code call takes some time
		timeoutMillisecond := time.Duration(milliseconds) * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeoutMillisecond)
		defer cancel()

		valuesChan, err := intervalSensor.Run(ctx, &wg)
		assert.NoError(t, err)

		var values []measurement.SensorValue
		for value := range valuesChan {
			values = append(values, value...)
		}

		assert.Equal(t, milliseconds/intervalSeconds/1000, len(values))
	})
	t.Run("sender should close channel after generateFunc panic", func(t *testing.T) {
		intervalSensor, err := simulator.New(
			randomPanic,
			time.Duration(intervalSeconds)*time.Second,
			batchSize,
			"",
			channelCap,
			zap.NewNop(),
		)
		assert.NoError(t, err)

		valuesChan, err := intervalSensor.Run(context.Background(), &wg)
		assert.NoError(t, err)

		for range valuesChan {
		}
		// no deadlock on channel close
	})
}

func TestIntervalSensor_NewIntervalSensor(t *testing.T) {
	batchSize := 1
	channelCap := 0

	t.Run("should return nil error when interval is zero", func(t *testing.T) {
		intervalSensor, err := simulator.New(
			randomValue,
			0,
			batchSize,
			"",
			channelCap,
			zap.NewNop(),
		)

		assert.NoError(t, err)
		assert.NotEmpty(t, intervalSensor)
	})
	t.Run("should return error when interval is below zero", func(t *testing.T) {
		intervalSensor, err := simulator.New(
			randomValue,
			-5,
			batchSize,
			"",
			channelCap,
			zap.NewNop(),
		)

		assert.Error(t, err)
		assert.Nil(t, intervalSensor)
	})
	t.Run("should return an error when generateFunc is nil", func(t *testing.T) {
		intervalSensor, err := simulator.New(
			nil,
			0,
			batchSize,
			"",
			channelCap,
			zap.NewNop(),
		)

		assert.Error(t, err)
		assert.Nil(t, intervalSensor)
	})
}

func TestIntervalSensor_NewRateSensor(t *testing.T) {
	ratePerSecond, _ := rate.New(5.0, time.Second)
	channelCap := 0

	t.Run("should return nil error on positive ratePerSecond", func(t *testing.T) {
		_, err := simulator.NewWithRate(
			randomValue,
			ratePerSecond,
			"",
			channelCap,
			zap.NewNop(),
		)

		assert.NoError(t, err)
	})
	t.Run("should return an error when generateFunc is nil", func(t *testing.T) {
		rateSensor, err := simulator.NewWithRate(
			nil,
			ratePerSecond,
			"",
			channelCap,
			zap.NewNop(),
		)

		assert.Error(t, err)
		assert.Nil(t, rateSensor)
	})
}

func randomValue() int64 {
	return rand.Int64N(int64(2 << 16))
}

func randomPanic() int64 {
	panic("any reason")
}
