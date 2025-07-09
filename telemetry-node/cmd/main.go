package main

import (
	"context"
	"math"
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
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/service"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

const defaultMinLogLevel = zapcore.InfoLevel

func main() {
	appConfig := config.LoadConfigEnv()

	// set minimal log level to log
	logLevel := logger.ParseLevel(appConfig.LoggerMinLogLevel, defaultMinLogLevel)

	lg := logger.New(os.Stdout, logLevel)
	defer lg.Sync()

	// actualRPS is "ceilled" and converted to int.
	// More details are in the Client-side Rate Limiter Setup section.
	actualRPS := int(math.Ceil(float64(appConfig.RequestRatePerSecond)))

	// create IntervalSensor instance
	sensor, err := domain.NewRateSensor(
		// func generates pseudo random values
		func() int64 {
			return rand.Int64N(int64(2 << 16))
		},
		appConfig.RequestRatePerSecond,
		appConfig.SensorName,
		lg,
	)
	if err != nil {
		lg.Error("can't create IntervalSensor using NewRateSensor()", zap.Error(err))
		os.Exit(1)
	}

	// Create a context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for OS signals
	signalCh := make(chan os.Signal, 1)

	// Notify on SIGINT (Ctrl+C) and SIGTERM (docker stop, kill)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	perRetryTimeout := 100 * time.Millisecond //  Timeout for each individual retry attempt

	// Important: if totalTimeoutPerRPCall <= perRetryTimeout then retryStrategy will never run
	totalTimeoutPerRPCall := perRetryTimeout * time.Duration(appConfig.GrpcClientMaxRetryAttempts)

	grpcClientOptions := []grpc.DialOption{
		// Important: insecure is allowed to use only for local development
		grpc.WithTransportCredentials(insecure.NewCredentials()),

		// Important: there won't be retry attempts if server is totally unavailable due network issues
		// add your own implementation of retry strategy with transport connection reestablishment
		grpc.WithUnaryInterceptor(
			retry.UnaryClientInterceptor(
				retry.WithCodes(codes.Unavailable, codes.ResourceExhausted, codes.DeadlineExceeded), // gRPC codes that should be retried
				retry.WithMax(appConfig.GrpcClientMaxRetryAttempts),                                 // Maximum retry attempts
				retry.WithPerRetryTimeout(perRetryTimeout),                                          // Timeout for each individual retry attempt
				retry.WithOnRetryCallback(func(ctx context.Context, attempt uint, err error) {
					lg.Debug("", zap.Uint("retry_attempt", attempt), zap.Error(err)) // log retry attempts
				}),
			)),
	}

	clientConn, err := grpc.NewClient(appConfig.GrpcServerSocket, grpcClientOptions...)
	if err != nil {
		lg.Error("failed to dial gRPC server", zap.Error(err))
		os.Exit(1)
	}
	defer clientConn.Close()

	sensorClient := sensorapi.NewSensorServiceClient(clientConn)

	sensorService := service.NewSensorService(sensorClient, lg)

	// --- Client-side Rate Limiter Setup ---
	// rate.NewLimiter(r Limit, b int)
	// r: rate limit, int (tokens per second)
	// b: burst, float64 (max tokens that can be accumulated, allowing for short bursts above the rate)
	// Setting burst = r allows for some flexibility if requests come in slightly faster than rate but average out.
	//
	// rps replaced with actualRPS to make rate = burst.
	limiter := rate.NewLimiter(rate.Limit(actualRPS), actualRPS)

	lg.Info("Client configured to send requests to server socket",
		zap.String("socket", appConfig.GrpcServerSocket),
		zap.Uint("max_attempts", appConfig.GrpcClientMaxRetryAttempts),
		zap.Duration("per_retry_timeout_second", perRetryTimeout),
	)
	lg.Info("Client configured to send requests with RPS",
		zap.Float64("limit", float64(limiter.Limit())),
		zap.Float64("burst", float64(limiter.Burst())),
	)

	// waitGroup is used for synchronisation
	var wg sync.WaitGroup

	// Run sensor worker
	valuesChan, err := sensor.Run(ctx, &wg)
	if err != nil {
		lg.Error("failed to run IntervalSensor", zap.Error(err))
		os.Exit(1)
	}

	//wg.Add(1) // go build -race found negative wg counter

	// Run send worker
	go func(rpcCallTimeout time.Duration, wg *sync.WaitGroup) {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				lg.Debug("send worker received context done, returning")
				return
			default:
				// Block until the rate limiter allows sending the next request
				// This is the core of maintaining constant RPS
				err = limiter.Wait(context.Background())
				if err != nil {
					// This error typically means the context passed to Wait was cancelled.
					lg.Error("Rate limiter wait interrupted", zap.Error(err))
					return
				}
				sensorService.SendSensorValues(ctx, rpcCallTimeout, valuesChan)
			}
		}
	}(totalTimeoutPerRPCall, &wg)
	wg.Add(1)

	// Block until a signal is received or context is cancelled
	select {
	case <-signalCh:
		lg.Info("Main routine received signal. Starting graceful shutdown")
		cancel() // cancel the context to signal other goroutines to stop
	}

	// Wait for other goroutines to stop
	wg.Wait()

	lg.Info("Main routine returned")
}
