package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestSlogAdapter(buf *bytes.Buffer) (*SlogAdapter, *bytes.Buffer) {
	if buf == nil {
		buf = &bytes.Buffer{}
	}
	handler := slog.NewJSONHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := slog.New(handler)
	adapter := NewSlogAdapter(l).(*SlogAdapter)
	return adapter, buf
}

func TestSlogAdapter_Info(t *testing.T) {
	adapter, buf := newTestSlogAdapter(nil)
	adapter.Info("hello", "key", "value")

	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "hello", entry["msg"])
	assert.Equal(t, "value", entry["key"])
}

func TestSlogAdapter_Debug(t *testing.T) {
	adapter, buf := newTestSlogAdapter(nil)
	adapter.Debug("debug msg", "foo", "bar")

	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "debug msg", entry["msg"])
}

func TestSlogAdapter_Warn(t *testing.T) {
	adapter, buf := newTestSlogAdapter(nil)
	adapter.Warn("warn msg")

	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "warn msg", entry["msg"])
}

func TestSlogAdapter_Error(t *testing.T) {
	adapter, buf := newTestSlogAdapter(nil)
	adapter.Error("error msg", "err", "something failed")

	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "error msg", entry["msg"])
	assert.Equal(t, "something failed", entry["err"])
}

func TestSlogAdapter_With(t *testing.T) {
	adapter, buf := newTestSlogAdapter(nil)
	enriched := adapter.With("user_id", "abc123")
	enriched.Info("with context")

	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "with context", entry["msg"])
	assert.Equal(t, "abc123", entry["user_id"])
}

func TestSlogAdapter_With_Chained(t *testing.T) {
	adapter, buf := newTestSlogAdapter(nil)
	enriched := adapter.With("request_id", "req1").With("user_id", "usr1")
	enriched.Info("chained")

	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	assert.Equal(t, "req1", entry["request_id"])
	assert.Equal(t, "usr1", entry["user_id"])
}

func TestSlogAdapter_With_DoesNotMutateOriginal(t *testing.T) {
	buf := &bytes.Buffer{}
	adapter, _ := newTestSlogAdapter(buf)

	_ = adapter.With("extra", "field")
	adapter.Info("original")

	var entry map[string]interface{}
	require.NoError(t, json.Unmarshal(buf.Bytes(), &entry))
	_, hasExtra := entry["extra"]
	assert.False(t, hasExtra, "original logger should not have extra field")
}

func TestSlogAdapter_Sync(t *testing.T) {
	adapter, _ := newTestSlogAdapter(nil)
	assert.NoError(t, adapter.Sync())
}

func TestSlogAdapter_SlogLogger(t *testing.T) {
	adapter, _ := newTestSlogAdapter(nil)
	sl := adapter.SlogLogger()
	assert.NotNil(t, sl)
}

func TestSlogAdapter_ImplementsLogger(t *testing.T) {
	adapter, _ := newTestSlogAdapter(nil)
	var _ Logger = adapter // compile-time check
	assert.NotNil(t, adapter)
}

func TestSlogAdapter_MultipleEntries(t *testing.T) {
	buf := &bytes.Buffer{}
	adapter, _ := newTestSlogAdapter(buf)

	adapter.Info("first")
	adapter.Info("second")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	assert.Len(t, lines, 2)
}
