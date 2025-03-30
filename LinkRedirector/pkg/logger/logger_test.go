package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	testCases := []struct {
		name       string
		level      string
		expectedLv zap.AtomicLevel
	}{
		{"Debug Level", "debug", zap.NewAtomicLevelAt(zap.DebugLevel)},
		{"Info Level", "info", zap.NewAtomicLevelAt(zap.InfoLevel)},
		{"Warn Level", "warn", zap.NewAtomicLevelAt(zap.WarnLevel)},
		{"Error Level", "error", zap.NewAtomicLevelAt(zap.ErrorLevel)},
		{"Default Level (Invalid Input)", "invalid", zap.NewAtomicLevelAt(zap.InfoLevel)},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := New(tc.level)

			assert.NoError(t, err)
			assert.NotNil(t, logger)

			assert.Equal(t, tc.expectedLv.Level(), logger.Level())
		})
	}
}

func TestNewLoggerBuildFailure(t *testing.T) {
	invalidConfig := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"/invalid/path"},
		ErrorOutputPaths: []string{"/invalid/path"},
	}

	_, err := invalidConfig.Build()
	assert.Error(t, err)
}
