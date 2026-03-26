// Package metrics provides a facade for application metrics recording.
// It ships with a NoopRecorder by default. To enable real metrics,
// pass a Prometheus, Datadog, or OpenTelemetry recorder to New().
//
// Usage:
//
//	m := metrics.New("my-service")                    // NoOp (default)
//	m := metrics.New("my-service", prometheusRecorder) // Real metrics
//
// All methods are safe for concurrent use and never panic.
package metrics

import "time"

// Recorder is the interface that metric backends must implement.
// Implementations: NoopRecorder (built-in), future: PrometheusRecorder, DatadogRecorder, OTelRecorder.
type Recorder interface {
	// CounterAdd increments a counter metric by the given value.
	CounterAdd(name string, value float64, labels map[string]string)
	// HistogramObserve records a value in a histogram/distribution metric.
	HistogramObserve(name string, value float64, labels map[string]string)
	// GaugeSet sets a gauge metric to the given value.
	GaugeSet(name string, value float64, labels map[string]string)
}

// Metrics is the central entry point for recording metrics across the EduGo ecosystem.
// Create one instance per service and pass it to components that need instrumentation.
type Metrics struct {
	recorder Recorder
	service  string
}

// New creates a Metrics instance for the given service.
// If no recorder is provided, a NoopRecorder is used (zero overhead).
func New(service string, recorder ...Recorder) *Metrics {
	var r Recorder = &NoopRecorder{}
	if len(recorder) > 0 && recorder[0] != nil {
		r = recorder[0]
	}
	return &Metrics{
		recorder: r,
		service:  service,
	}
}

// Service returns the service name this Metrics instance was created for.
func (m *Metrics) Service() string {
	return m.service
}

// Recorder returns the underlying recorder for advanced use cases.
func (m *Metrics) Recorder() Recorder {
	return m.recorder
}

// statusLabel returns "success" or "error" based on the error.
func statusLabel(err error) string {
	if err != nil {
		return "error"
	}
	return "success"
}

// durationSeconds converts a time.Duration to seconds as float64.
func durationSeconds(d time.Duration) float64 {
	return d.Seconds()
}
