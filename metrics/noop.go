package metrics

// NoopRecorder is a metrics recorder that does nothing.
// It is the default recorder used when no backend is configured.
// All methods are zero-cost no-ops safe for concurrent use.
type NoopRecorder struct{}

// CounterAdd es un no-op que satisface la interfaz Recorder.
func (n *NoopRecorder) CounterAdd(string, float64, map[string]string) {}

// HistogramObserve es un no-op que satisface la interfaz Recorder.
func (n *NoopRecorder) HistogramObserve(string, float64, map[string]string) {}

// GaugeSet es un no-op que satisface la interfaz Recorder.
func (n *NoopRecorder) GaugeSet(string, float64, map[string]string) {}
