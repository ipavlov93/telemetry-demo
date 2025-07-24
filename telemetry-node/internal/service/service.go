package service

import (
	"context"
	"sync"
	"time"

	sensorapi "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor_service"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/mapper"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

const gracefulShutdownMinDuration = 50 * time.Millisecond

// SensorService represents component that is responsible for delivery to destination server.
type SensorService struct {
	sensorClient     sensorapi.SensorServiceClient
	gracefulShutdown time.Duration
	logger           logger.Logger
}

func NewSensorService(
	sensorClient sensorapi.SensorServiceClient,
	gracefulShutdown time.Duration,
	lg logger.Logger,
) *SensorService {
	return &SensorService{
		sensorClient:     sensorClient,
		gracefulShutdown: gracefulShutdown,
		logger:           lg,
	}
}

// sendRequest sends request using corresponding client and ignores errors.
// timeoutDuration is total timeout per RPC call (including retry attempts).
func (s *SensorService) sendRequest(
	ctx context.Context,
	sensorValues []measurement.SensorValue,
	limiter RateLimiter,
	timeout time.Duration,
) {
	// Block until the rate limiter allows sending the next request
	err := limiter.Wait(ctx)
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

// Run starts send requests using corresponding client in a separate goroutine.
// It will start graceful shutdown when the parent context is done (via <-ctx.Done()).
func (s *SensorService) Run(parentCtx context.Context, config *RunConfig) {
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
				if s.gracefulShutdown < gracefulShutdownMinDuration {
					return
				}
				s.shutdown(config)
				return
			case sensorValues, ok := <-config.valuesChan:
				if !ok {
					return
				}
				s.sendRequest(parentCtx, sensorValues, config.limiter, config.totalTimeoutRPC)
			}
		}
	}()
}

func (s *SensorService) shutdown(config *RunConfig) {
	if config == nil {
		return
	}

	var once sync.Once
	once.Do(func() {
		s.logger.Debug("SensorService has started shutdown")

		ctx, cancel := context.WithTimeout(context.Background(), s.gracefulShutdown)
		defer cancel()

		valuesBatch := drainChannel(ctx, config.valuesChan)
		if len(valuesBatch) == 0 {
			return
		}

		reqCtx, reqCancel := context.WithTimeout(context.Background(), config.totalTimeoutRPC)
		s.sendRequest(reqCtx, valuesBatch, config.limiter, config.totalTimeoutRPC)
		reqCancel()
	})
}

func drainChannel(ctx context.Context, valuesChan <-chan []measurement.SensorValue) []measurement.SensorValue {
	var valuesBatch []measurement.SensorValue

	for {
		select {
		case <-ctx.Done():
			return valuesBatch
		case sensorValues, ok := <-valuesChan:
			if !ok {
				return valuesBatch
			}
			valuesBatch = append(valuesBatch, sensorValues...)
		default:
			return valuesBatch
		}
	}
}
