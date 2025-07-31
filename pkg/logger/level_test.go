package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestParseLevelOrDefault(t *testing.T) {
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
			got := ParseLevelOrDefault(tt.givenLevel, tt.fallback)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name       string
		givenLevel string
		want       zapcore.Level
		wantErr    bool
	}{
		{
			name:       "should return error",
			givenLevel: "",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "should return error",
			givenLevel: "123456",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "should return error",
			givenLevel: "-123456",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "should return error",
			givenLevel: "not_supported_level",
			want:       0,
			wantErr:    true,
		},
		{
			name:       "should set given level",
			givenLevel: "debug",
			want:       zapcore.DebugLevel,
			wantErr:    false,
		},
		{
			name:       "should set given level",
			givenLevel: "info",
			want:       zapcore.InfoLevel,
			wantErr:    false,
		},
		{
			name:       "should set given level",
			givenLevel: "DEBUG",
			want:       zapcore.DebugLevel,
			wantErr:    false,
		},
		{
			name:       "should set given level",
			givenLevel: "warn",
			want:       zapcore.WarnLevel,
			wantErr:    false,
		},
		{
			name:       "should set given level",
			givenLevel: "error",
			want:       zapcore.ErrorLevel,
			wantErr:    false,
		},
		{
			name:       "should set given level",
			givenLevel: "fatal",
			want:       zapcore.FatalLevel,
			wantErr:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseLevel(tt.givenLevel)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
