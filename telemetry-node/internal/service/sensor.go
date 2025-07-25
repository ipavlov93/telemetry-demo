package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	sensorapi "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor_service"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/mapper"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

const gracefulShutdownMinDuration = 100 * time.Millisecond

// SensorService represents component that is responsible for delivery to destination server.
type SensorService struct {
	sensorClient     sensorapi.SensorServiceClient
	limiter          RateLimiter
	drainStrategy    DrainStrategy
	gracefulShutdown time.Duration
	logger           logger.Logger
}

// NewSensorService constructor returns error if gracefulShutdown is invalid.
func NewSensorService(
	sensorClient sensorapi.SensorServiceClient,
	limiter RateLimiter,
	drainStrategy DrainStrategy,
	gracefulShutdown time.Duration,
	lg logger.Logger,
) (*SensorService, error) {
	if gracefulShutdown < gracefulShutdownMinDuration {
		return nil, fmt.Errorf("graceful shutdown must be at least %v", gracefulShutdownMinDuration)
	}
	return &SensorService{
		sensorClient:     sensorClient,
		limiter:          limiter,
		drainStrategy:    drainStrategy,
		gracefulShutdown: gracefulShutdown,
		logger:           lg,
	}, nil
}

// Run starts send requests using corresponding client in a separate goroutine.
// It will start graceful shutdown:
// - when parent context is done (via <-ctx.Done())
// - if gracefulShutdown is bigger than config.totalTimeoutRPC.
func (s *SensorService) Run(parentCtx context.Context, config *RunConfig) error {
	if config == nil {
		return fmt.Errorf("config must be not nil")
	}
	wg := config.wg
	if wg != nil {
		wg.Add(1)
	}

	go func() {
		defer func() {
			if wg == nil {
				return
			}
			wg.Done()
		}()

		for {
			select {
			case <-parentCtx.Done():
				if config.totalTimeoutRPC <= s.gracefulShutdown {
					s.shutdown(config)
				}
				return
			case sensorValues, ok := <-config.valuesChan:
				if !ok {
					return
				}
				s.sendRequest(parentCtx, sensorValues, config.totalTimeoutRPC)
			}
		}
	}()
	return nil
}

// sendRequest sends request using corresponding client and ignores errors.
// timeoutDuration is total timeout per RPC call (including retry attempts).
func (s *SensorService) sendRequest(
	ctx context.Context,
	sensorValues []measurement.SensorValue,
	timeout time.Duration,
) {
	// Block until the rate limiter allows sending the next request
	err := s.limiter.Wait(context.Background())
	if err != nil {
		s.logger.Error("Rate limiter wait interrupted", zap.Error(err))
		return
	}

	reqCtx, reqCancel := context.WithTimeout(ctx, timeout)
	func() {
		defer reqCancel()
		_, err = s.sensorClient.SendSensorValues(reqCtx,
			mapper.ToProtoRequest(sensorValues),
		)
		if err != nil {
			statusCode := status.Convert(err)
			s.logger.Debug("SensorServiceClient.SendSensorValues failed",
				zap.Uint32("status_code", uint32(statusCode.Code())),
				zap.String("message", statusCode.Message()),
			)
		}
	}()
}

func (s *SensorService) shutdown(config *RunConfig) {
	if config == nil {
		return
	}

	var once sync.Once
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), s.gracefulShutdown)
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				return
			default:
				valuesBatch := s.drainStrategy.Receive(ctx, config.valuesChan)
				if len(valuesBatch) == 0 {
					return
				}

				s.sendRequest(ctx, valuesBatch, config.totalTimeoutRPC)
			}
		}
	})
}
