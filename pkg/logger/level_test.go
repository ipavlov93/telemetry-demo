package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name       string
		givenLevel string
		fallback   zapcore.Level
		want       zapcore.Level
	}{
		{
			name:       "should set default level",
			givenLevel: "",
			fallback:   zapcore.DebugLevel,
			want:       zapcore.DebugLevel,
		},
		{
			name:       "should set given level",
			givenLevel: "debug",
			fallback:   zapcore.InfoLevel,
			want:       zapcore.DebugLevel,
		},
		{
			name:       "should set given level",
			givenLevel: "info",
			fallback:   zapcore.DebugLevel,
			want:       zapcore.InfoLevel,
		},
		{
			name:       "should set given level",
			givenLevel: "DEBUG",
			fallback:   zapcore.InfoLevel,
			want:       zapcore.DebugLevel,
		},
		{
			name:       "should set given level",
			givenLevel: "warn",
			fallback:   zapcore.InfoLevel,
			want:       zapcore.WarnLevel,
		},
		{
			name:       "should set given level",
			givenLevel: "error",
			fallback:   zapcore.InfoLevel,
			want:       zapcore.ErrorLevel,
		},
		{
			name:       "should set given level",
			givenLevel: "fatal",
			fallback:   zapcore.InfoLevel,
			want:       zapcore.FatalLevel,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLevel(tt.givenLevel, tt.fallback)
			assert.Equal(t, tt.want, got)
		})
	}
}
