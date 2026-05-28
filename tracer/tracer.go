package tracer

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	Endpoint       string // OTLP endpoint; empty = stdout exporter
	Insecure       bool
	Enabled        bool // false = noop tracer (zero overhead)
}

type Shutdown func(context.Context) error

// Init initializes the OpenTelemetry tracer and logger providers and registers
// them as the global providers. Spans go to OTLP/gRPC (or stdout when Endpoint
// is empty); slog records routed via NewSlogHandler are exported to the same
// OTLP/gRPC endpoint as OTel LogRecords.
//
// Returns a Shutdown function that must be called on application exit; it
// flushes and closes both providers. If cfg.Enabled is false, returns a noop
// shutdown (zero overhead).
func Init(ctx context.Context, cfg Config) (Shutdown, error) {
	if !cfg.Enabled {
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			semconv.DeploymentEnvironmentKey.String(cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	var spanExporter sdktrace.SpanExporter
	if cfg.Endpoint != "" {
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(cfg.Endpoint),
		}
		if cfg.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		spanExporter, err = otlptracegrpc.New(ctx, opts...)
	} else {
		spanExporter, err = stdouttrace.New(stdouttrace.WithPrettyPrint())
	}
	if err != nil {
		return nil, fmt.Errorf("creating span exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(spanExporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// LoggerProvider: only when an OTLP endpoint is configured. If Endpoint is
	// empty (stdout traces) we skip the log exporter — the stdout span exporter
	// already provides full visibility for local debugging.
	var lp *sdklog.LoggerProvider
	if cfg.Endpoint != "" {
		logOpts := []otlploggrpc.Option{
			otlploggrpc.WithEndpoint(cfg.Endpoint),
		}
		if cfg.Insecure {
			logOpts = append(logOpts, otlploggrpc.WithInsecure())
		}
		logExporter, err := otlploggrpc.New(ctx, logOpts...)
		if err != nil {
			return nil, fmt.Errorf("creating log exporter: %w", err)
		}
		lp = sdklog.NewLoggerProvider(
			sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
			sdklog.WithResource(res),
		)
		global.SetLoggerProvider(lp)
	}

	return func(ctx context.Context) error {
		var errs []error
		if err := tp.Shutdown(ctx); err != nil {
			errs = append(errs, fmt.Errorf("tracer provider shutdown: %w", err))
		}
		if lp != nil {
			if err := lp.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Errorf("logger provider shutdown: %w", err))
			}
		}
		return errors.Join(errs...)
	}, nil
}

// Tracer returns a named tracer from the global provider.
// Usage: tracer.Tracer("identity").Start(ctx, "operation-name")
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}
