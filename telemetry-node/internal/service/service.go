package service

import (
	"context"
	"time"

	sensorapi "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor_service"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/mapper"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
)

// SensorService represents component that is responsible for delivery to destination server.
type SensorService struct {
	sensorClient sensorapi.SensorServiceClient
	logger       logger.Logger
}

func NewSensorService(
	sensorClient sensorapi.SensorServiceClient,
	lg logger.Logger,
) *SensorService {
	return &SensorService{
		sensorClient: sensorClient,
		logger:       lg,
	}
}

// SendSensorValues consumes sensor values (single slice) from the input channel and sends them using corresponding client.
// timeoutDuration is total timeout per RPC call (including retry attempts).
// Notice: it won't drain remaining messages from channel if context is cancelled.
func (s *SensorService) SendSensorValues(
	parentCtx context.Context,
	timeoutDuration time.Duration,
	input <-chan []measurement.SensorValue,
) {
	for {
		select {
		case <-parentCtx.Done():
			s.logger.Debug("SensorService received context done, returning")
			return
		case sensorValues, ok := <-input:
			if !ok {
				s.logger.Debug("SensorService input channel closed")
				return
			}
			ctx, cancel := context.WithTimeout(parentCtx, timeoutDuration)
			_, err := s.sensorClient.SendSensorValues(ctx,
				mapper.ToProtoRequest(sensorValues),
			)
			if err != nil {
				statusCode := status.Convert(err)
				s.logger.Debug("SensorServiceClient.SendSensorValues failed",
					zap.Uint32("status_code", uint32(statusCode.Code())),
					zap.String("message", statusCode.Message()),
				)
			}
			cancel()
		}
	}
}
