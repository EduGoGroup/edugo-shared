package metrics

import (
	"strconv"
	"time"
)

// HTTP metric names.
const (
	MetricHTTPRequestsTotal   = "http_requests_total"
	MetricHTTPRequestDuration = "http_request_duration_seconds"
	MetricHTTPActiveRequests  = "http_active_requests"
)

// RecordHTTPRequest records an HTTP request with method, path template, status code, and duration.
func (m *Metrics) RecordHTTPRequest(method, path string, status int, duration time.Duration) {
	labels := map[string]string{
		"service": m.service,
		"method":  method,
		"path":    path,
		"status":  strconv.Itoa(status),
	}
	m.recorder.CounterAdd(MetricHTTPRequestsTotal, 1, labels)
	m.recorder.HistogramObserve(MetricHTTPRequestDuration, durationSeconds(duration), labels)
}

// SetHTTPActiveRequests sets the current number of active/in-flight HTTP requests.
func (m *Metrics) SetHTTPActiveRequests(count int) {
	m.recorder.GaugeSet(MetricHTTPActiveRequests, float64(count), map[string]string{
		"service": m.service,
	})
}
