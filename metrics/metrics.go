// Package metrics provee una fachada para el registro de métricas de aplicación.
// Incluye un NoopRecorder por defecto. Para habilitar métricas reales,
// pasa un recorder de Prometheus, Datadog u OpenTelemetry a New().
//
// Uso:
//
//	m := metrics.New("my-service")                    // NoOp (por defecto)
//	m := metrics.New("my-service", prometheusRecorder) // Métricas reales
//
// Los métodos de Metrics delegan al Recorder proporcionado.
// La seguridad de concurrencia y ausencia de panics depende de la implementación del Recorder.
// NoopRecorder (el default) es seguro para uso concurrente y nunca hace panic.
package metrics

import "time"

// Recorder es la interfaz que los backends de métricas deben implementar.
// Implementaciones: NoopRecorder (incluido), futuras: PrometheusRecorder, DatadogRecorder, OTelRecorder.
type Recorder interface {
	// CounterAdd incrementa un contador por el valor dado.
	CounterAdd(name string, value float64, labels map[string]string)
	// HistogramObserve registra un valor en un histograma/distribución.
	HistogramObserve(name string, value float64, labels map[string]string)
	// GaugeSet establece un gauge al valor dado.
	GaugeSet(name string, value float64, labels map[string]string)
}

// Metrics es el punto de entrada central para registrar métricas en el ecosistema EduGo.
// Crea una instancia por servicio y pásala a los componentes que necesiten instrumentación.
type Metrics struct {
	recorder Recorder
	service  string
}

// New crea una instancia de Metrics para el servicio dado.
// Si no se proporciona un recorder, se usa NoopRecorder (cero overhead).
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

// Service retorna el nombre del servicio para esta instancia de Metrics.
func (m *Metrics) Service() string {
	return m.service
}

// Recorder retorna el recorder subyacente para casos de uso avanzados.
func (m *Metrics) Recorder() Recorder {
	return m.recorder
}

// statusLabel retorna "success" o "error" según el error.
func statusLabel(err error) string {
	if err != nil {
		return "error"
	}
	return "success"
}

// durationSeconds convierte un time.Duration a segundos como float64.
func durationSeconds(d time.Duration) float64 {
	return d.Seconds()
}
