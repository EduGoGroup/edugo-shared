package metrics

import (
	"errors"
	"testing"
	"time"
)

func TestNew_DefaultsToNoop(t *testing.T) {
	m := New("test-service")
	if m == nil {
		t.Fatal("expected non-nil Metrics")
	}
	if m.Service() != "test-service" {
		t.Errorf("expected service 'test-service', got '%s'", m.Service())
	}
	if _, ok := m.Recorder().(*NoopRecorder); !ok {
		t.Error("expected NoopRecorder as default")
	}
}

func TestNew_WithCustomRecorder(t *testing.T) {
	rec := &spyRecorder{}
	m := New("svc", rec)
	if m.Recorder() != rec {
		t.Error("expected custom recorder")
	}
}

func TestNew_NilRecorderFallsBackToNoop(t *testing.T) {
	m := New("svc", nil)
	if _, ok := m.Recorder().(*NoopRecorder); !ok {
		t.Error("expected NoopRecorder when nil passed")
	}
}

func TestNoop_NeverPanics(t *testing.T) {
	m := New("test")
	// All methods should be no-ops without panic
	m.RecordHTTPRequest("GET", "/test", 200, time.Second)
	m.RecordDBQuery("postgres", "select", "users", time.Millisecond, nil)
	m.RecordDBQuery("mongodb", "insert", "docs", time.Millisecond, errors.New("fail"))
	m.RecordLogin(true, time.Second)
	m.RecordLogin(false, time.Second)
	m.RecordTokenRefresh(true, time.Second)
	m.RecordRateLimitHit("login")
	m.RecordPermissionCheck("users:read", true)
	m.RecordPermissionCheck("users:write", false)
	m.RecordBusinessOperation("membership", "create", time.Second, nil)
	m.RecordBusinessOperation("grade", "update", time.Second, errors.New("fail"))
	m.RecordAssessmentAttempt("start", time.Second, nil)
	m.RecordAssessmentAttempt("submit", time.Second, errors.New("timeout"))
	m.RecordGrading("multiple_choice", time.Millisecond, nil)
	m.RecordExport("xlsx", 100, time.Second, nil)
	m.RecordMessageProcessed("material_uploaded", time.Second, nil)
	m.RecordMessageRetry("material_uploaded", 3)
	m.SetMessagesInQueue("events", 42)
	m.SetCircuitBreakerState("openai", "closed")
	m.SetCircuitBreakerState("openai", "half_open")
	m.SetCircuitBreakerState("openai", "open")
	m.SetHTTPActiveRequests(5)
	m.SetDBConnectionsOpen("postgres", 10)
}

func TestSpyRecorder_CapturesMetrics(t *testing.T) {
	rec := &spyRecorder{}
	m := New("test", rec)

	m.RecordLogin(true, 500*time.Millisecond)

	if len(rec.counters) == 0 {
		t.Fatal("expected counter to be recorded")
	}
	if rec.counters[0].name != MetricAuthLoginsTotal {
		t.Errorf("expected metric name %s, got %s", MetricAuthLoginsTotal, rec.counters[0].name)
	}
	if rec.counters[0].labels["status"] != "success" {
		t.Errorf("expected status 'success', got '%s'", rec.counters[0].labels["status"])
	}

	if len(rec.histograms) == 0 {
		t.Fatal("expected histogram to be recorded")
	}
	if rec.histograms[0].value < 0.4 || rec.histograms[0].value > 0.6 {
		t.Errorf("expected duration ~0.5s, got %f", rec.histograms[0].value)
	}
}

func TestStatusLabel(t *testing.T) {
	if statusLabel(nil) != "success" {
		t.Error("nil error should be 'success'")
	}
	if statusLabel(errors.New("fail")) != "error" {
		t.Error("non-nil error should be 'error'")
	}
}

// spyRecorder captures metric calls for assertions.
type spyRecorder struct {
	counters   []spyMetric
	histograms []spyMetric
	gauges     []spyMetric
}

type spyMetric struct {
	name   string
	value  float64
	labels map[string]string
}

func (s *spyRecorder) CounterAdd(name string, value float64, labels map[string]string) {
	s.counters = append(s.counters, spyMetric{name, value, labels})
}

func (s *spyRecorder) HistogramObserve(name string, value float64, labels map[string]string) {
	s.histograms = append(s.histograms, spyMetric{name, value, labels})
}

func (s *spyRecorder) GaugeSet(name string, value float64, labels map[string]string) {
	s.gauges = append(s.gauges, spyMetric{name, value, labels})
}
