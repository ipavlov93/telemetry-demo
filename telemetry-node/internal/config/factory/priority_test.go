package factory_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/app"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/factory"
	"github.com/ipavlov93/telemetry-demo/telemetry-node/internal/config/logger"
	"github.com/stretchr/testify/require"
)

const configFileName = "config.test.json"

func TestNewAppConfig(t *testing.T) {
	defaultConfig := &app.Config{
		SensorName:                 "default-sensor-name",
		RequestRatePerSecond:       1,
		SensorValueRatePerSecond:   5,
		GrpcServerSocket:           "localhost:8000",
		GrpcClientMaxRetryAttempts: 1,
		AppLoggerConfig: logger.ConfigMap{
			"APP_LOGGER": logger.Configuration{
				Enabled:     true,
				MinLevel:    "info",
				Destination: "stdout"},
		},
	}

	tests := []struct {
		name          string
		configContent any
		expected      *app.Config
		expectErr     bool
		missingFile   bool
	}{
		{
			name:      "valid JSON file and env vars",
			expectErr: false,
			expected: &app.Config{
				SensorName:                 "sensor-name-123",
				RequestRatePerSecond:       11,
				SensorValueRatePerSecond:   12,
				GrpcServerSocket:           "localhost:13",
				GrpcClientMaxRetryAttempts: 14,
				AppLoggerConfig: logger.ConfigMap{
					"APP_LOGGER": logger.Configuration{
						Enabled:     false,
						MinLevel:    "error",
						Destination: logger.LogOutputFile,
						File: &logger.FileLoggerConfig{
							FilePath:           "../../",
							FileName:           "appLogFile.log",
							FileRotationPeriod: "1m",
						},
					},
				},
			},
			configContent: `
			{
			  "sensorName": "sensor-name-123",
			  "requestRatePerSecond": 11,
			  "sensorValueRatePerSecond": 12,
			  "grpcServerSocket": "localhost:13",
			  "grpcClientMaxRetryAttempts": 14,
			  "appLoggerConfig": {
			    "APP_LOGGER": {
			      "enabled": false,
			      "minLevel": "error",
			      "destination": "file",
			      "file": {
			        "filePath": "../../",
			        "fileName": "appLogFile.log",
			        "fileRotationPeriod": "1m"
			      }
			    }
			  }
			}
			`,
		},
		{
			name:        "missing config file",
			missingFile: true,
			expectErr:   false,
			expected:    defaultConfig,
		},
		{
			name:          "invalid JSON file",
			configContent: "{ this is invalid json ",
			expectErr:     false,
			expected:      defaultConfig,
		},
	}

	origWd, err := os.Getwd()
	require.NoError(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			err = os.Chdir(tempDir)
			require.NoError(t, err)

			if !tt.missingFile {
				f, err := os.Create(filepath.Join(tempDir, configFileName))
				require.NoError(t, err)

				switch v := tt.configContent.(type) {
				case string:
					_, err = f.WriteString(v)
					require.NoError(t, err)
				default:
					enc := json.NewEncoder(f)
					err = enc.Encode(v)
					require.NoError(t, err)
				}

				_ = f.Close()
			}

			cfg, err := factory.NewAppConfig(configFileName)

			if tt.expectErr {
				require.Error(t, err)
				require.Nil(t, cfg)
			} else {
				require.NoError(t, err)
				require.NotNil(t, cfg)

				require.EqualValues(t, tt.expected, cfg)
				//reflect.DeepEqual(t, tt.expected, cfg)
			}

			_ = os.Chdir(origWd)
		})
	}
}
