package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestNewLogrusLogger(t *testing.T) {
	l := logrus.New()
	logger := NewLogrusLogger(l)
	assert.NotNil(t, logger)
}

func TestLogrusLogger_Levels(t *testing.T) {
	buffer := &bytes.Buffer{}
	l := logrus.New()
	l.SetOutput(buffer)
	l.SetLevel(logrus.DebugLevel)
	l.SetFormatter(&logrus.JSONFormatter{})

	logger := NewLogrusLogger(l)

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
}

func TestLogrusLogger_With(t *testing.T) {
	buffer := &bytes.Buffer{}
	l := logrus.New()
	l.SetOutput(buffer)
	l.SetFormatter(&logrus.JSONFormatter{})
	logger := NewLogrusLogger(l)

	child := logger.With("context", "child")
	child.Info("child msg")

	var output map[string]any
	err := json.Unmarshal(buffer.Bytes(), &output)
	assert.NoError(t, err)
	assert.Equal(t, "child msg", output["msg"])
	assert.Equal(t, "child", output["context"])
}

func TestConvertToLogrusFields(t *testing.T) {
	// Private function test
	fields := convertToLogrusFields("key1", "value1", "key2", 123)
	assert.Equal(t, "value1", fields["key1"])
	assert.Equal(t, 123, fields["key2"])

	// Test odd number of arguments
	fields = convertToLogrusFields("key1", "value1", "ignored")
	assert.Equal(t, "value1", fields["key1"])
	assert.NotContains(t, fields, "ignored")

	// Test non-string key
	fields = convertToLogrusFields(123, "value1", "key2", "value2")
	assert.NotContains(t, fields, 123)
	assert.Equal(t, "value2", fields["key2"])
}
