package main

import (
	"context"
	"math/rand/v2"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	sensorapi "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor_service"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/rate"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/sensor/simulator"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/service"
	rps "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/utils/rate"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/utils/timeout"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ratelimiter "golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultMinLogLevel          = zapcore.InfoLevel
	defaultRequestRatePerSecond = 1

	perRetryTimeout = 100 * time.Millisecond

	gracefulShutdown = time.Second
)

func main() {
	appConfig := config.LoadConfigEnv()

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	lg := logger.New(os.Stdout,
		logger.ParseLevel(
			appConfig.LoggerMinLogLevel,
			defaultMinLogLevel,
		),
	)
	defer lg.Sync()

	samplingRate, err := rate.New(appConfig.SensorValueRatePerSecond, time.Second)
	if err != nil {
		lg.Fatal("failed to initialize sampling rate", zap.Error(err))
	}

	sensorSimulator, err := simulator.NewWithRate(
		func() int64 {
			return rand.Int64N(int64(2 << 16))
		},
		samplingRate,
		appConfig.SensorName,
		0,
		lg,
	)
	if err != nil {
		lg.Fatal("failed to initialize SensorSimulator", zap.Error(err))
	}

	// Important:
	// if totalTimeoutRPC <= perRetryTimeout then retryStrategy will never run
	totalTimeoutRPC := timeout.TotalTimeout(perRetryTimeout, appConfig.GrpcClientMaxRetryAttempts)

	grpcClientOptions := []grpc.DialOption{
		// Important: insecure is allowed to use only for local development
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			retry.UnaryClientInterceptor(
				retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
				// WithMax supplies given value minus one (first initial attempt)
				retry.WithMax(appConfig.GrpcClientMaxRetryAttempts+1),
				retry.WithPerRetryTimeout(perRetryTimeout),
				retry.WithOnRetryCallback(func(ctx context.Context, attempt uint, err error) {
					lg.Debug("", zap.Uint("retry_attempt", attempt), zap.Error(err))
				}),
			)),
	}

	clientConn, err := grpc.NewClient(appConfig.GrpcServerSocket, grpcClientOptions...)
	if err != nil {
		lg.Fatal("failed to dial gRPC server", zap.Error(err))
	}
	defer clientConn.Close()

	sensorClient := sensorapi.NewSensorServiceClient(clientConn)
	sensorService := service.NewSensorService(sensorClient, gracefulShutdown, lg)

	// rps replaced with actualRPS to make rate = burst.
	actualRPS := rps.RoundOrDefaultRPS(
		float64(appConfig.RequestRatePerSecond),
		defaultRequestRatePerSecond,
	)
	limiter := ratelimiter.NewLimiter(ratelimiter.Limit(actualRPS), actualRPS)

	lg.Info("Client configured to send requests to server socket",
		zap.String("socket", appConfig.GrpcServerSocket),
		zap.Uint("max_attempts", appConfig.GrpcClientMaxRetryAttempts),
		zap.Duration("per_retry_timeout_second", perRetryTimeout),
	)
	lg.Info("Client configured to send requests with RPS",
		zap.Float64("limit", float64(limiter.Limit())),
		zap.Int("burst", limiter.Burst()),
	)
	lg.Info("Sensor configured to produce values with RPS",
		zap.Float32("rate", appConfig.SensorValueRatePerSecond),
	)

	var wg sync.WaitGroup

	// Run sensorSimulator worker
	valuesChan, err := sensorSimulator.Run(ctx, &wg)
	if err != nil {
		lg.Fatal("failed to run SensorSimulator", zap.Error(err))
	}

	serviceRunConfig := service.NewRunConfig(valuesChan, totalTimeoutRPC, limiter, &wg)
	if !serviceRunConfig.Valid() {
		lg.Fatal("invalid sensorService.Run configuration")
	}

	// Run send worker
	sensorService.Run(ctx, serviceRunConfig)

	select {
	case <-ctx.Done():
		lg.Debug("Parent context is done")
	case <-signalCh:
		lg.Info("Main routine received signal. Starting graceful shutdown")
		cancel() // cancel the context to signal other goroutines to stop
	}

	wg.Wait()

	lg.Debug("Main routine returned")
}
