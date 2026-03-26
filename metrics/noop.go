package metrics

// NoopRecorder is a metrics recorder that does nothing.
// It is the default recorder used when no backend is configured.
// All methods are zero-cost no-ops safe for concurrent use.
type NoopRecorder struct{}

func (n *NoopRecorder) CounterAdd(string, float64, map[string]string)       {}
func (n *NoopRecorder) HistogramObserve(string, float64, map[string]string) {}
func (n *NoopRecorder) GaugeSet(string, float64, map[string]string)         {}
