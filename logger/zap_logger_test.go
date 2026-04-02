package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name   string
		level  string
		format string
	}{
		{"DebugJSON", "debug", "json"},
		{"InfoConsole", "info", "console"},
		{"WarnJSON", "warn", "json"},
		{"ErrorConsole", "error", "console"},
		{"FatalJSON", "fatal", "json"},
		{"DefaultInfo", "invalid", "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewZapLogger(tt.level, tt.format)
			assert.NotNil(t, l)
		})
	}
}

func TestZapLogger_JSONOutput(t *testing.T) {
	// Custom sink to capture output
	buffer := &bytes.Buffer{}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(buffer),
		zapcore.DebugLevel,
	)
	zapLog := zap.New(core).Sugar()
	logger := &zapLogger{logger: zapLog}

	logger.Info("test message", "key", "value")

	var output map[string]any
	err := json.Unmarshal(buffer.Bytes(), &output)
	assert.NoError(t, err)
	assert.Equal(t, "test message", output["msg"])
	assert.Equal(t, "info", output["level"])
	assert.Equal(t, "value", output["key"])
}

func TestZapLogger_Levels(t *testing.T) {
	// We can't easily mock stdout/stderr for the real constructor,
	// so we'll test the method calls on a mocked or buffered core if possible,
	// or just trust that NewZapLogger returns a working logger and we test the wrapper methods.

	// Testing wrapper methods with a buffer
	buffer := &bytes.Buffer{}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(buffer),
		zapcore.DebugLevel,
	)
	zapLog := zap.New(core).Sugar()
	logger := &zapLogger{logger: zapLog}

	logger.Debug("debug msg")
	assert.Contains(t, buffer.String(), "debug msg")
	buffer.Reset()

	logger.Info("info msg")
	assert.Contains(t, buffer.String(), "info msg")
	buffer.Reset()

	logger.Warn("warn msg")
	assert.Contains(t, buffer.String(), "warn msg")
	buffer.Reset()

	logger.Error("error msg")
	assert.Contains(t, buffer.String(), "error msg")
	buffer.Reset()

	// Fatal usually calls os.Exit, so we skip it or mock it if possible (zap allow hooking checkwrite)
	// For now we skip Fatal test to avoid crashing the test runner
}

func TestZapLogger_With(t *testing.T) {
	buffer := &bytes.Buffer{}
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(buffer),
		zapcore.DebugLevel,
	)
	zapLog := zap.New(core).Sugar()
	logger := &zapLogger{logger: zapLog}

	child := logger.With("context", "child")
	child.Info("child msg")

	var output map[string]any
	err := json.Unmarshal(buffer.Bytes(), &output)
	assert.NoError(t, err)
	assert.Equal(t, "child msg", output["msg"])
	assert.Equal(t, "child", output["context"])
}
