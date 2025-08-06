package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config"
	simulatorfactory "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/domain/simulator/factory"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/infra/logger/factory/writer"
	zapfactory "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/infra/logger/factory/zap"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/service"
	servicefactory "github.com/ipavlov93/telemetry-demo/telemetry-node/internal/service/factory"
	retryfactory "github.com/ipavlov93/telemetry-demo/telemetry-node/pkg/retry"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	appConfig := config.NewAppConfig()

	// Create a context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	lg, err := zapfactory.NewLogger(
		appConfig.AppLoggerConfig,
		writer.NewLogWriter,
		zapfactory.NewEncoder(),
	)
	if err != nil {
		lg = zap.NewNop()
	}

	defer lg.Sync()

	sensorSimulator, err := simulatorfactory.NewRandomValueSimulator(
		appConfig.SensorName,
		appConfig.SensorValueRatePerSecond,
		lg,
	)
	if err != nil {
		lg.Fatal("failed to initialize SensorSimulator", zap.Error(err))
	}

	// Important: insecure is allowed to use only for local development
	grpcClientOptions := retryfactory.NewInsecure(
		retryfactory.NewRetryUnaryInterceptor(
			appConfig.GrpcClientMaxRetryAttempts,
			lg,
		)...,
	)

	clientConn, err := grpc.NewClient(appConfig.GrpcServerSocket, grpcClientOptions...)
	if err != nil {
		lg.Fatal("failed to dial gRPC server", zap.Error(err))
	}
	defer clientConn.Close()

	sensorService, err := servicefactory.NewSensorServiceRPS(
		appConfig.RequestRatePerSecond,
		clientConn,
		lg,
	)
	if err != nil {
		lg.Fatal("failed to initialize sensor service", zap.Error(err))
	}

	lg.Info("Client configured to send requests to server socket",
		zap.String("socket", appConfig.GrpcServerSocket),
		zap.Uint("max_attempts", appConfig.GrpcClientMaxRetryAttempts),
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

	// Important: if totalTimeoutRPC <= perRetryTimeout then retryStrategy will never run
	totalTimeoutRPC := retryfactory.NewTotalTimeout(appConfig.GrpcClientMaxRetryAttempts)
	serviceRunConfig := service.NewRunConfig(valuesChan, totalTimeoutRPC, &wg)
	if !serviceRunConfig.Valid() {
		lg.Fatal("invalid configuration for sensorService.Run")
	}

	// Run sensor service worker
	err = sensorService.Run(ctx, serviceRunConfig)
	if err != nil {
		lg.Fatal("failed to Run SensorService", zap.Error(err))
	}

	select {
	case <-ctx.Done():
	case <-signalCh:
		lg.Info("Main routine received signal. Starting graceful shutdown")
		cancel() // cancel the context to signal other goroutines to stop
	}

	wg.Wait()
}
