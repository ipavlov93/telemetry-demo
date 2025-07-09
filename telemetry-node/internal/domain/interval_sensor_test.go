package domain_test

import (
	"context"
	"math/rand/v2"
	"sync"
	"testing"
	"time"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestIntervalSensor_Run(t *testing.T) {
	intervalSeconds := 1
	totalSeconds := 5
	batchSize := 1
	wg := sync.WaitGroup{}

	t.Run("IntervalSensor.SendSensorValues() happy flow", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			randomValue,
			time.Duration(intervalSeconds)*time.Second,
			batchSize,
			"",
			zap.NewNop(),
		)
		assert.NoError(t, err)

		// context timeout duration has small approximation error or testing code call takes some time
		timeout := time.Duration(totalSeconds*1000+10) * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		valuesChan, err := intervalSensor.Run(ctx, &wg)
		assert.NoError(t, err)

		var values []domain.SensorValue
		for value := range valuesChan {
			values = append(values, value...)
		}

		assert.Equal(t, totalSeconds/intervalSeconds, len(values))
	})
	t.Run("should gracefully send partial batch before context cancellation", func(t *testing.T) {
		batchSize = 5

		intervalSensor, err := domain.NewIntervalSensor(
			randomValue,
			time.Duration(intervalSeconds)*time.Second,
			batchSize,
			"",
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

		var values []domain.SensorValue
		for value := range valuesChan {
			values = append(values, value...)
		}

		assert.Equal(t, milliseconds/intervalSeconds/1000, len(values))
	})
	t.Run("sender should close channel after generateFunc panic", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			randomPanic,
			time.Duration(intervalSeconds)*time.Second,
			batchSize,
			"",
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

	t.Run("should return nil error when interval is zero", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			randomValue,
			0,
			batchSize,
			"",
			zap.NewNop(),
		)

		assert.NoError(t, err)
		assert.NotEmpty(t, intervalSensor)
	})
	t.Run("should return error when interval is below zero", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			randomValue,
			-5,
			batchSize,
			"",
			zap.NewNop(),
		)

		assert.Error(t, err)
		assert.Nil(t, intervalSensor)
	})
	t.Run("should return an error when generateFunc is nil", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			nil,
			0,
			batchSize,
			"",
			zap.NewNop(),
		)

		assert.Error(t, err)
		assert.Nil(t, intervalSensor)
	})
}

func TestIntervalSensor_NewRateSensor(t *testing.T) {
	ratePerSecond := float32(5.0)

	t.Run("should return nil error on positive ratePerSecond", func(t *testing.T) {
		_, err := domain.NewRateSensor(
			randomValue,
			ratePerSecond,
			"",
			zap.NewNop(),
		)

		assert.NoError(t, err)
	})
	t.Run("should return error on zero ratePerSecond", func(t *testing.T) {
		_, err := domain.NewRateSensor(
			randomValue,
			0,
			"",
			zap.NewNop(),
		)

		assert.Error(t, err)
	})
	t.Run("should return error when ratePerSecond is below zero", func(t *testing.T) {
		rateSensor, err := domain.NewRateSensor(
			randomValue,
			-5,
			"",
			zap.NewNop(),
		)

		assert.Error(t, err)
		assert.Empty(t, rateSensor)
	})
	t.Run("should return an error when generateFunc is nil", func(t *testing.T) {
		rateSensor, err := domain.NewRateSensor(
			nil,
			0,
			"",
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
