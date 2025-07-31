package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogOutput_Valid(t *testing.T) {
	tests := []struct {
		name        string
		destination logOutput
		want        bool // valid
	}{
		{
			name:        "should return false",
			destination: "",
			want:        false,
		},
		{
			name:        "should return false",
			destination: "not_supported_destination",
			want:        false,
		},
		{
			name:        "should return true",
			destination: LogOutputFile,
			want:        true,
		},
		{
			name:        "should return true",
			destination: LogOutputStdout,
			want:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.destination.Valid())
		})
	}
}
