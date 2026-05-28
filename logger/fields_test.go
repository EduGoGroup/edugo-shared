package logger

import (
	"context"
	"log/slog"
	"testing"
)

// Benchmarks

func BenchmarkNewSlogProvider(b *testing.B) {
	cfg := SlogConfig{
		Level:   "info",
		Format:  "json",
		Service: "bench-service",
		Env:     "test",
		Version: "1.0.0",
	}
	for b.Loop() {
		_ = NewSlogProvider(cfg)
	}
}

func BenchmarkSlogAdapter_Info(b *testing.B) {
	l := NewSlogProvider(SlogConfig{Level: "error", Format: "json"})
	adapter := NewSlogAdapter(l)
	for b.Loop() {
		adapter.Info("benchmark message", "key", "value")
	}
}

func BenchmarkSlogAdapter_With(b *testing.B) {
	l := NewSlogProvider(SlogConfig{Level: "info", Format: "json"})
	adapter := NewSlogAdapter(l)
	for b.Loop() {
		_ = adapter.With("request_id", "abc123")
	}
}

func BenchmarkFromContext(b *testing.B) {
	l := NewSlogProvider(SlogConfig{Level: "info", Format: "json"})
	ctx := NewContext(context.Background(), NewSlogAdapter(l))
	for b.Loop() {
		_ = FromContext(ctx)
	}
}

func BenchmarkSlogDirect_Info(b *testing.B) {
	l := NewSlogProvider(SlogConfig{Level: "error", Format: "json"})
	for b.Loop() {
		l.Info("benchmark message", slog.String("key", "value"))
	}
}
