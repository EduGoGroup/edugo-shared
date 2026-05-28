package metrics

import (
	"testing"
	"time"
)

// TestRecordPermissionsLookup_NoOp asegura que con el recorder default (NoOp)
// la llamada no panic ni efectos colaterales.
func TestRecordPermissionsLookup_NoOp(t *testing.T) {
	m := New("identity")
	m.RecordPermissionsLookup("legacy", 12, 5*time.Millisecond)
	m.RecordPermissionsLookup("grants", 0, time.Microsecond)
}

func TestRecordPermissionDrift_NoOp(t *testing.T) {
	m := New("identity")
	m.RecordPermissionDrift("user-123", "canary")
}

// TestRecordPermissionsLookup_EmitsCounterAndHistogram verifica que se emitan
// el counter total + histograma de duración + histograma de tamaño con los
// labels esperados.
func TestRecordPermissionsLookup_EmitsCounterAndHistogram(t *testing.T) {
	rec := &spyRecorder{}
	m := New("identity", rec)

	m.RecordPermissionsLookup("grants", 7, 250*time.Millisecond)

	// Counter total
	if len(rec.counters) != 1 {
		t.Fatalf("expected 1 counter, got %d", len(rec.counters))
	}
	if rec.counters[0].name != MetricPermissionsLookupTotal {
		t.Errorf("expected counter %s, got %s", MetricPermissionsLookupTotal, rec.counters[0].name)
	}
	if rec.counters[0].labels["service"] != "identity" {
		t.Errorf("expected service label 'identity', got %q", rec.counters[0].labels["service"])
	}
	if rec.counters[0].labels["source"] != "grants" {
		t.Errorf("expected source label 'grants', got %q", rec.counters[0].labels["source"])
	}
	if rec.counters[0].value != 1 {
		t.Errorf("expected counter value 1, got %f", rec.counters[0].value)
	}

	// Un histograma: duración (en segundos).
	if len(rec.histograms) != 1 {
		t.Fatalf("expected 1 histogram (duration), got %d", len(rec.histograms))
	}
	h := rec.histograms[0]
	if h.name != MetricPermissionsLookupDuration {
		t.Errorf("expected histogram %s, got %s", MetricPermissionsLookupDuration, h.name)
	}
	if h.value < 0.2 || h.value > 0.3 {
		t.Errorf("expected duration ~0.25s, got %f", h.value)
	}
	if h.labels["source"] != "grants" {
		t.Errorf("duration: expected source 'grants', got %q", h.labels["source"])
	}
	if h.labels["service"] != "identity" {
		t.Errorf("duration: expected service 'identity', got %q", h.labels["service"])
	}
}

func TestRecordPermissionsLookup_LegacySource(t *testing.T) {
	rec := &spyRecorder{}
	m := New("identity", rec)

	m.RecordPermissionsLookup("legacy", 3, time.Millisecond)

	if rec.counters[0].labels["source"] != "legacy" {
		t.Errorf("expected source 'legacy', got %q", rec.counters[0].labels["source"])
	}
}

// TestRecordPermissionDrift_Emits verifica labels y que userID NO aparezca
// como label (alta cardinalidad).
func TestRecordPermissionDrift_Emits(t *testing.T) {
	rec := &spyRecorder{}
	m := New("identity", rec)

	m.RecordPermissionDrift("user-abc", "canary")

	if len(rec.counters) != 1 {
		t.Fatalf("expected 1 counter, got %d", len(rec.counters))
	}
	c := rec.counters[0]
	if c.name != MetricPermissionsDriftTotal {
		t.Errorf("expected metric %s, got %s", MetricPermissionsDriftTotal, c.name)
	}
	if c.labels["service"] != "identity" {
		t.Errorf("expected service 'identity', got %q", c.labels["service"])
	}
	if c.labels["source"] != "canary" {
		t.Errorf("expected source 'canary', got %q", c.labels["source"])
	}
	if _, hasUser := c.labels["user_id"]; hasUser {
		t.Error("user_id MUST NOT be emitted as a label (cardinality)")
	}
}
