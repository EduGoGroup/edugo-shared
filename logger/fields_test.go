package logger

import (
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWithRequestID(t *testing.T) {
	attr := WithRequestID("req-123")
	assert.Equal(t, FieldRequestID, attr.Key)
	assert.Equal(t, "req-123", attr.Value.String())
}

func TestWithUserID(t *testing.T) {
	attr := WithUserID("user-456")
	assert.Equal(t, FieldUserID, attr.Key)
	assert.Equal(t, "user-456", attr.Value.String())
}

func TestWithCorrelationID(t *testing.T) {
	attr := WithCorrelationID("corr-789")
	assert.Equal(t, FieldCorrelationID, attr.Key)
	assert.Equal(t, "corr-789", attr.Value.String())
}

func TestWithError(t *testing.T) {
	attr := WithError(errors.New("something failed"))
	assert.Equal(t, FieldError, attr.Key)
	assert.Equal(t, "something failed", attr.Value.String())
}

func TestWithDuration(t *testing.T) {
	attr := WithDuration(1500 * time.Millisecond)
	assert.Equal(t, FieldDuration, attr.Key)
	assert.Equal(t, int64(1500), attr.Value.Int64())
}

func TestWithComponent(t *testing.T) {
	attr := WithComponent("auth_service")
	assert.Equal(t, FieldComponent, attr.Key)
	assert.Equal(t, "auth_service", attr.Value.String())
}

func TestWithSchoolID(t *testing.T) {
	attr := WithSchoolID("school-abc")
	assert.Equal(t, FieldSchoolID, attr.Key)
	assert.Equal(t, "school-abc", attr.Value.String())
}

func TestWithRole(t *testing.T) {
	attr := WithRole("teacher")
	assert.Equal(t, FieldRole, attr.Key)
	assert.Equal(t, "teacher", attr.Value.String())
}

func TestWithIP(t *testing.T) {
	attr := WithIP("192.168.1.1")
	assert.Equal(t, FieldIP, attr.Key)
	assert.Equal(t, "192.168.1.1", attr.Value.String())
}

// Benchmarks

func BenchmarkNewSlogProvider(b *testing.B) {
	cfg := SlogConfig{
		Level:   "info",
		Format:  "json",
		Service: "bench-service",
		Env:     "test",
		Version: "1.0.0",
	}
	b.ResetTimer()
	for b.Loop() {
		_ = NewSlogProvider(cfg)
	}
}

func BenchmarkSlogAdapter_Info(b *testing.B) {
	l := NewSlogProvider(SlogConfig{Level: "error", Format: "json"})
	adapter := NewSlogAdapter(l)
	b.ResetTimer()
	for b.Loop() {
		adapter.Info("benchmark message", "key", "value")
	}
}

func BenchmarkSlogAdapter_With(b *testing.B) {
	l := NewSlogProvider(SlogConfig{Level: "info", Format: "json"})
	adapter := NewSlogAdapter(l)
	b.ResetTimer()
	for b.Loop() {
		_ = adapter.With("request_id", "abc123")
	}
}

func BenchmarkFromContext(b *testing.B) {
	l := NewSlogProvider(SlogConfig{Level: "info", Format: "json"})
	ctx := NewContext(nil, l)
	b.ResetTimer()
	for b.Loop() {
		_ = FromContext(ctx)
	}
}

func BenchmarkTypedHelpers(b *testing.B) {
	for b.Loop() {
		_ = WithRequestID("req-123")
		_ = WithUserID("user-456")
		_ = WithDuration(100 * time.Millisecond)
	}
}

func BenchmarkSlogDirect_Info(b *testing.B) {
	l := NewSlogProvider(SlogConfig{Level: "error", Format: "json"})
	b.ResetTimer()
	for b.Loop() {
		l.Info("benchmark message", slog.String("key", "value"))
	}
}
