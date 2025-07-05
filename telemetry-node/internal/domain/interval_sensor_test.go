package domain_test

import (
	"context"
	"math/rand/v2"
	"testing"
	"time"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func randomValue() int64 {
	return rand.Int64N(int64(2 << 16))
}

func randomPanic() int64 {
	panic("any reason")
}

func TestIntervalSensor_Value(t *testing.T) {
	intervalSeconds := 1
	totalSeconds := 5

	t.Run("IntervalSensor.Run() happy flow", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			"",
			time.Duration(intervalSeconds)*time.Second,
			randomValue,
			zap.NewNop(),
		)
		assert.NoError(t, err)

		// context timeout duration has small approximation error or testing code call takes some time
		totalMilliseconds := time.Duration(totalSeconds*1000+10) * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), totalMilliseconds)
		defer cancel()

		valuesChan, err := intervalSensor.Run(ctx)
		assert.NoError(t, err)

		var values []*domain.SensorValue
		for value := range valuesChan {
			values = append(values, value)
		}

		assert.Equal(t, totalSeconds/intervalSeconds, len(values))
	})
	t.Run("sender should close channel after generateFunc panic", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			"",
			time.Duration(intervalSeconds)*time.Second,
			randomPanic,
			zap.NewNop(),
		)
		assert.NoError(t, err)

		valuesChan, err := intervalSensor.Run(context.Background())
		assert.NoError(t, err)

		for range valuesChan {
		}
		// no deadlock on channel close
	})
}

func TestIntervalSensor_NewIntervalSensor(t *testing.T) {
	t.Run("should return an error when generateFunc is nil", func(t *testing.T) {
		intervalSensor, err := domain.NewIntervalSensor(
			"",
			time.Second,
			nil,
			zap.NewNop(),
		)
		assert.Error(t, err)
		assert.Nil(t, intervalSensor)
	})
}
