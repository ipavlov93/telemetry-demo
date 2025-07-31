package logger

type Config struct {
	LogDestinations []LogOutput
}

// LogOutput is a single logging output config, e.g. stdout, file, or connection.
type LogOutput struct {
	Enabled  bool
	MinLevel string
	Type     logOutput

	File *FileLoggerConfig `json:"file,omitempty"`
}
