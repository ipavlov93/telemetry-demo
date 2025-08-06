package factory

import (
	"time"

	sensorapi "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor_service"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/service"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/strategy/channel"
	rps "github.com/ipavlov93/telemetry-demo/telemetry-node/pkg/utils/rate"
	"go.uber.org/zap"
	ratelimiter "golang.org/x/time/rate"
	"google.golang.org/grpc"
)

const (
	defaultRatePerSecond = 1
	shutdownDuration     = 5 * time.Second
)

// NewSensorServiceRPS factory constructs sensorService that has max sending rate once per second.
func NewSensorServiceRPS(
	requestRatePerSecond float32,
	clientConn grpc.ClientConnInterface,
	lg logger.Logger,
) (*service.SensorService, error) {
	// rps replaced with actualRPS to make rate = burst.
	actualRPS := rps.RoundOrDefaultRPS(
		float64(requestRatePerSecond),
		defaultRatePerSecond,
	)
	limiter := ratelimiter.NewLimiter(ratelimiter.Limit(actualRPS), actualRPS)

	sensorClient := sensorapi.NewSensorServiceClient(clientConn)
	sensorService, err := service.NewSensorService(
		sensorClient,
		limiter,
		channel.NewDrainLastStrategy(),
		shutdownDuration,
		lg,
	)
	if err != nil {
		return nil, err
	}

	lg.Info("SensorServiceClient configured to send requests with RPS",
		zap.Float64("limit", float64(limiter.Limit())),
		zap.Int("burst", limiter.Burst()),
	)

	return sensorService, nil
}
