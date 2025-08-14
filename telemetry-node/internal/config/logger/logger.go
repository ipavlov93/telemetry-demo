// Package logger contains log configuration related types.
package logger

type ConfigMap map[string]Configuration // map[name]Configuration

func NewNopConfigMap(key string) ConfigMap {
	return map[string]Configuration{
		key: NewNopConfiguration(),
	}
}

// Configuration is a single logger config
type Configuration struct {
	Enabled     bool      `json:"enabled"`
	MinLevel    string    `json:"minLevel"`
	Destination logOutput `json:"destination"`

	File *FileLoggerConfig `json:"file,omitempty"`
}

func NewNopConfiguration() Configuration {
	return Configuration{
		Enabled: false,
	}
}
