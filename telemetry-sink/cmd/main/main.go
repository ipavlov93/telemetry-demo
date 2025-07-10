package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	sensorapi "github.com/ipavlov93/telemetry-demo/pkg/grpc/generated/v1/sensor_service"
	"github.com/ipavlov93/telemetry-demo/pkg/logger"
	"github.com/ipavlov93/telemetry-demo/telemetry-sink/internal/config"
	"github.com/ipavlov93/telemetry-demo/telemetry-sink/internal/processor"
	"github.com/ipavlov93/telemetry-demo/telemetry-sink/internal/writer"
	"github.com/ipavlov93/telemetry-demo/telemetry-sink/server"
	"github.com/ipavlov93/telemetry-demo/telemetry-sink/server/interceptor"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

const (
	defaultMinLogLevel = zapcore.InfoLevel
)

func main() {
	// Parse env/config or set default values
	appConfig := config.LoadConfigEnv()

	// set minimal log level to log
	logLevel := logger.ParseLevel(appConfig.LoggerMinLogLevel, defaultMinLogLevel)

	lg := logger.New(os.Stdout, logLevel)
	defer lg.Sync()

	// create tcp socket
	tcpListener, err := net.Listen("tcp", appConfig.GrpcServerSocket)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Initialize the token bucket rate limiter.
	limiter := rate.NewLimiter(rate.Limit(appConfig.RatePerSecond), appConfig.RatePerSecond)

	grpcServer := grpc.NewServer(
		// Important: insecure connection is allowed to use only for local development

		// registers custom rate limiter
		grpc.UnaryInterceptor(interceptor.ByteRateLimiterInterceptor(limiter, lg)),
	)

	// create server implementation
	srv := server.NewServer(lg)

	// register gRPC server
	sensorapi.RegisterSensorServiceServer(grpcServer, srv)

	lg.Info("Server configured to drop incoming messages that exceed the allowed bandwidth, bytes/second",
		zap.Float64("limit", float64(limiter.Limit())),
		zap.Float64("burst", float64(limiter.Burst())),
	)

	// Create a context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create a channel to listen for OS signals
	signalCh := make(chan os.Signal, 1)

	// Notify on SIGINT (Ctrl+C) and SIGTERM (docker stop, kill)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	// create IntervalSensor instance
	sensorDataProcessor, err := processor.NewBufferedProcessor(appConfig.BufferFlushInterval, appConfig.BufferSize)
	if err != nil {
		lg.Error("failed to run BufferedProcessor", zap.Error(err))
		os.Exit(1)
	}

	// create JsonWriter instance
	jsonWriter := writer.NewJsonWriter("", appConfig.FilePath, lg)

	// waitGroup is used for synchronisation
	var wg sync.WaitGroup

	// Run processor workers
	sensorDataProcessor.Run(ctx, srv.Out(), &wg)
	jsonWriter.Run(ctx, sensorDataProcessor.Out(), &wg)

	// Start serving gRPC requests blocks
	// routine ignores parent context
	go func(wg *sync.WaitGroup) {
		defer wg.Done()

		lg.Info("Server is configured listening socket", zap.String("socket", tcpListener.Addr().String()))
		err := grpcServer.Serve(tcpListener)
		if err != nil {
			lg.Error("gRPC Server Serve", zap.Error(err))
		}
	}(&wg)
	wg.Add(1)

	// Block until a signal is received or context is cancelled
	select {
	case <-signalCh:
		lg.Info("Main routine received signal. Starting graceful shutdown")
		cancel() // cancel the context to signal other goroutines to stop

		grpcServer.GracefulStop() // stopping the Server
	}

	// Wait for other goroutines to stop
	wg.Wait()

	lg.Info("Main routine returned")
}
