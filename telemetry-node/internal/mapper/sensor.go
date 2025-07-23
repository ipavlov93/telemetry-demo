package mapper

import (
	sensorpb "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/measurement"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// ToProtoRequest sets []domain.SensorValue to *sensorpb.SensorValuesRequest.
// sensorpb.SensorValuesRequest wraps []domain.SensorValue.
func ToProtoRequest(sensorValues []measurement.SensorValue) *sensorpb.SensorValuesRequest {
	valueBatches := make([]*sensorpb.SensorValue, len(sensorValues))
	for i, v := range sensorValues {
		valueBatches[i] = ToProtoMessage(v)
	}
	return &sensorpb.SensorValuesRequest{Items: valueBatches}
}

// ToProtoMessage converts domain.SensorValue entity to *sensorpb.SensorValue DTO.
func ToProtoMessage(sensorValue measurement.SensorValue) *sensorpb.SensorValue {
	return &sensorpb.SensorValue{
		SensorName: sensorValue.SensorName,
		Measurement: &sensorpb.SensorValue_Measurement{
			SensorValue: &wrapperspb.Int64Value{Value: sensorValue.Value},
			CreatedAt:   timestamppb.New(sensorValue.Timestamp),
		},
	}
}
