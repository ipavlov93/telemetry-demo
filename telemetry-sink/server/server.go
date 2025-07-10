package server

import (
	"context"

	pb "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor"
	sensorapi "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor_service"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"go.uber.org/zap"
)

const defaultChannelBufferSize = 1000

// Server component implements SensorServiceServer.
type Server struct {
	sensorapi.SensorServiceServer
	outChan chan []*pb.SensorValue
	logger  logger.Logger
}

// NewServer returns pointer to created instance of Server.
func NewServer(lg logger.Logger) *Server {
	return &Server{
		outChan: make(chan []*pb.SensorValue, defaultChannelBufferSize),
		logger:  lg,
	}
}

// Out return the output channel.
func (s *Server) Out() <-chan []*pb.SensorValue { return s.outChan }

// SendSensorValues gRPC handler sends request messages to outputChan.
// buffered channel is used to prevent immediate block on channel send.
func (s *Server) SendSensorValues(_ context.Context, req *pb.SensorValuesRequest) (resp *pb.SensorValuesResponse, err error) {
	if len(req.GetItems()) == 0 {
		return &pb.SensorValuesResponse{}, nil
	}

	s.outChan <- req.GetItems()

	s.logger.Debug("sent messages to channel", zap.Int("message_count", len(req.GetItems())))
	return &pb.SensorValuesResponse{}, nil
}
