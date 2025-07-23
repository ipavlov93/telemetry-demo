package service

import (
	"context"
	"sync"
	"time"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
)

const gracefulShutdownMinDuration = 50 * time.Millisecond

type RateLimiter interface {
	Wait(ctx context.Context) error
}

type RunConfig struct {
	valuesChan       <-chan []measurement.SensorValue
	totalTimeoutRPC  time.Duration
	limiter          RateLimiter
	gracefulShutdown time.Duration
	wg               *sync.WaitGroup
}

func NewRunConfig(
	ch <-chan []measurement.SensorValue,
	totalTimeoutRPC time.Duration,
	limiter RateLimiter,
	gracefulShutdown time.Duration,
	wg *sync.WaitGroup,
) *RunConfig {
	return &RunConfig{
		valuesChan:       ch,
		totalTimeoutRPC:  totalTimeoutRPC,
		limiter:          limiter,
		gracefulShutdown: gracefulShutdown,
		wg:               wg,
	}
}

func (r *RunConfig) Valid() bool {
	if r == nil {
		return false
	}
	if r.valuesChan == nil {
		return false
	}
	if r.totalTimeoutRPC <= 0 {
		return false
	}
	if r.gracefulShutdown < gracefulShutdownMinDuration {
		return false
	}
	return r.wg != nil
}
