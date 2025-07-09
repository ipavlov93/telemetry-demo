package service

import (
	"context"
	"strconv"
	"time"

	sensorpb "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor"
	sensorapi "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor_service"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// SensorService represents component that is responsible for delivery to destination server.
type SensorService struct {
	sensorClient sensorapi.SensorServiceClient
	logger       logger.Logger
}

// NewSensorService returns pointer to created instance of SensorService.
// It skips logger validation.
func NewSensorService(
	sensorClient sensorapi.SensorServiceClient,
	lg logger.Logger,
) *SensorService {
	return &SensorService{
		sensorClient: sensorClient,
		logger:       lg,
	}
}

// SendSensorValues consumes sensor values (single slice) from the channel and sends them using corresponding client.
// Notice: it won't drain remaining messages from channel if context is cancelled.
func (s *SensorService) SendSensorValues(parentCtx context.Context, timeoutDuration time.Duration, input <-chan []domain.SensorValue) {
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

			// prepare request
			req := toProtoRequest(sensorValues)

			// Notice:
			// individual context is passed per RPC call
			ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
			defer cancel()

			// send attempt with retry
			// there won't be any retry attempts if server is totally unavailable due network issues
			_, err := s.sensorClient.SendSensorValues(ctx, req)
			if err != nil {
				st := status.Convert(err)
				s.logger.Debug("SensorServiceClient.SendSensorValues failed",
					zap.String(strconv.Itoa(int(status.Code(err))), st.Message()))

				continue
			}
			s.logger.Debug("SensorServiceClient.SendSensorValues successful",
				zap.Int("messageCount", len(sensorValues)),
			)
		}
	}
}

// toProtoRequest utility function set []domain.SensorValue to *sensorpb.SensorValuesRequest.
func toProtoRequest(sensorValues []domain.SensorValue) *sensorpb.SensorValuesRequest {
	valueBatches := make([]*sensorpb.SensorValue, len(sensorValues))
	for i, v := range sensorValues {
		valueBatches[i] = toProtoMessage(v)
	}
	return &sensorpb.SensorValuesRequest{Items: valueBatches}
}

// toProtoMessage utility function convert domain.SensorValue entity to *sensorpb.SensorValue DTO.
func toProtoMessage(sensorValue domain.SensorValue) *sensorpb.SensorValue {
	return &sensorpb.SensorValue{
		SensorName: sensorValue.SensorName,
		Measurement: &sensorpb.SensorValue_Measurement{
			SensorValue: &wrapperspb.Int64Value{Value: sensorValue.Value},
			CreatedAt:   timestamppb.New(sensorValue.Timestamp),
		},
	}
}
