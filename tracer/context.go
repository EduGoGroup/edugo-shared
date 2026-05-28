package tracer

import (
	"context"

	"go.opentelemetry.io/otel/trace"
)

// TraceIDFromContext extracts the trace ID from the context.
// Returns empty string if no active span.
func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	if sc.HasTraceID() {
		return sc.TraceID().String()
	}
	return ""
}

// SpanIDFromContext extracts the span ID from the context.
// Returns empty string if no active span.
func SpanIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	sc := span.SpanContext()
	if sc.HasSpanID() {
		return sc.SpanID().String()
	}
	return ""
}
