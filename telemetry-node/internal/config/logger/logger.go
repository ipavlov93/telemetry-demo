package logger

type ConfigMap map[string]Configuration // map[name]Configuration

func NewNopConfigMap(key string) ConfigMap {
	return map[string]Configuration{
		key: NewNopConfiguration(),
	}
}

// Configuration is a single logger config
type Configuration struct {
	Enabled     bool
	MinLevel    string
	Destination logOutput

	File *FileLoggerConfig `json:"file,omitempty"`
}

func NewNopConfiguration() Configuration {
	return Configuration{
		Enabled: false,
	}
}
