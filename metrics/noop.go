package metrics

// NoopRecorder es un recorder de métricas que no hace nada.
// Es el recorder por defecto cuando no se configura un backend.
// Todos los métodos son no-ops de cero costo, seguros para uso concurrente.
type NoopRecorder struct{}

// CounterAdd es un no-op que satisface la interfaz Recorder.
func (n *NoopRecorder) CounterAdd(string, float64, map[string]string) {}

// HistogramObserve es un no-op que satisface la interfaz Recorder.
func (n *NoopRecorder) HistogramObserve(string, float64, map[string]string) {}

// GaugeSet es un no-op que satisface la interfaz Recorder.
func (n *NoopRecorder) GaugeSet(string, float64, map[string]string) {}
