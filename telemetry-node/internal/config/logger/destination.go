package logger

import "fmt"

// logOutput enforces safe construction via NewLogOutput().
type logOutput string

const (
	LogOutputStdout logOutput = "stdout"
	LogOutputFile   logOutput = "file"
)

var validLogOutputs = map[logOutput]struct{}{
	LogOutputStdout: {},
	LogOutputFile:   {},
}

func (l logOutput) String() string {
	_, ok := validLogOutputs[l]
	if !ok {
		return fmt.Sprintf("logOutput(%s)", string(l))
	}
	return string(l)
}

// NewLogOutput constructor will return an error if input isn't supported log destination stream.
func NewLogOutput(s string) (logOutput, error) {
	destination := logOutput(s)
	if !destination.Valid() {
		return "", fmt.Errorf("not supported log destination stream: %q", s)
	}
	return destination, nil
}

func (l logOutput) Valid() bool {
	if _, ok := validLogOutputs[l]; !ok {
		return false
	}
	return true
}
