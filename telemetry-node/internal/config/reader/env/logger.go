package env

import (
	"fmt"
	"os"
	"time"

	"github.com/ipavlov93/telemetry-demo/pkg/env"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
)

const defaultLevel = "info"

// AppLoggerConfig parses environment variables, otherwise returns error.
// Note: it will set default log level if MinLevel isn't found.
func AppLoggerConfig(key string) (logger.ConfigMap, error) {
	destinations := logger.NewNopConfigMap(key)

	appLogDestination := os.Getenv("APP_LOGGER_DESTINATION")
	appLogOutput, err := logger.NewLogOutput(appLogDestination)
	if err != nil {
		return nil, err
	}

	if os.Getenv("APP_LOGGER_DESTINATION") == "" {
		return nil, nil
	}

	configuration := logger.Configuration{
		Enabled:     true,
		MinLevel:    env.EnvironmentVariable("APP_LOGGER_MIN_LOG_LEVEL", defaultLevel),
		Destination: appLogOutput,
	}

	fileLoggerCfg, err := ReadFileLoggerEnv(
		"APP_LOGGER_FILE_PATH",
		"APP_LOGGER_FILE_NAME",
		"APP_LOGGER_FILE_ROTATION_PERIOD",
		configuration,
	)

	if err != nil {
		configuration.File = fileLoggerCfg
		destinations[key] = configuration
	}

	return destinations, nil
}

// ReadFileLoggerEnv parses environment variables, otherwise returns error.
func ReadFileLoggerEnv(
	pathEnv string,
	nameEnv string,
	rotationPeriodEnv string,
	configuration logger.Configuration,
) (*logger.FileLoggerConfig, error) {
	if configuration.Destination != logger.LogOutputFile {
		return nil, nil
	}

	filePath, found := os.LookupEnv(pathEnv)
	if !found {
		return nil, fmt.Errorf("env variable %s not found", pathEnv)
	}
	fileName, found := os.LookupEnv(nameEnv)
	if !found {
		return nil, fmt.Errorf("env variable %s not found", nameEnv)
	}

	fileRotationPeriod, found := os.LookupEnv(rotationPeriodEnv)
	if !found {
		return nil, fmt.Errorf("env variable %s not found", fileRotationPeriod)
	}
	duration, err := time.ParseDuration(fileRotationPeriod)
	if err != nil {
		return nil, err
	}

	return &logger.FileLoggerConfig{
		FilePath:           filePath,
		FileName:           fileName,
		FileRotationPeriod: duration,
	}, nil
}

//func Configuration() logger.ConfigMap {
//	var destinations []logger.Configuration
//
//	destinations = append(destinations, loggerDestination(
//		"STDOUT_LOGGER_MIN_LOG_LEVEL",
//		string(logger.LogOutputStdout),
//		"info",
//	))
//
//	return logger.ConfigMap{
//		LogDestinations: destinations,
//	}
//}
//
//func loggerDestination(
//	key string,
//	dst string,
//	defaultLevel string,
//) logger.Configuration {
//	_, found := os.LookupEnv(key)
//	if !found {
//		return logger.Configuration{}
//	}
//
//	destination, _ := logger.NewLogOutput(dst)
//	//if destination, err := logger.NewLogOutput(dst); err != nil {		}
//
//	return logger.Configuration{
//		Enabled:  true,
//		Destination:     destination,
//		MinLevel: env.EnvironmentVariable(key, defaultLevel),
//	}
//}
