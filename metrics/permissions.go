package metrics

import "time"

// Nombres de métricas de evaluación de permisos.
//
// Estas métricas permiten observar la salud del subsistema RBAC durante la
// migración legacy -> grants (feature flag RBAC_USE_GRANTS) y detectar drift
// entre ambos paths antes/después del cutover.
const (
	// MetricPermissionsLookupTotal cuenta las llamadas a GetUserPermissions
	// (resolución del set efectivo de permisos para un usuario). Se etiqueta
	// por "source" ∈ {"legacy","grants"} para validar paridad durante el
	// canary.
	MetricPermissionsLookupTotal = "permissions_lookup_total"

	// MetricPermissionsLookupDuration registra la duración (en segundos) de
	// GetUserPermissions. Útil para detectar regresiones de performance al
	// activar el path grants (que usa permission_matches).
	MetricPermissionsLookupDuration = "permissions_lookup_duration_seconds"

	// MetricPermissionsCheckTotal cuenta las evaluaciones de
	// RequirePermission (granted/denied), por permission. Sustituye al alias
	// histórico MetricAuthPermissionChecks.
	MetricPermissionsCheckTotal = "permissions_check_total"

	// MetricPermissionsDriftTotal contabiliza divergencias detectadas entre
	// el path legacy y el path grants para un mismo usuario. Mientras este
	// contador permanezca en 0, ambos paths son equivalentes.
	MetricPermissionsDriftTotal = "permissions_drift_total"
)

// RecordPermissionsLookup registra una resolución de permisos.
//
// source identifica el path: "legacy" (iam.role_permissions) o "grants"
// (iam.role_grants + iam.permission_matches). count es el tamaño del set
// resultante; hoy se ignora en la emisión de métricas (queda disponible para
// futuras extensiones — p.ej. gauge de tamaño promedio — sin romper la
// signature pública).
func (m *Metrics) RecordPermissionsLookup(source string, count int, duration time.Duration) {
	_ = count // reservado para futuras extensiones
	labels := map[string]string{
		"service": m.service,
		"source":  source,
	}
	m.recorder.CounterAdd(MetricPermissionsLookupTotal, 1, labels)
	m.recorder.HistogramObserve(MetricPermissionsLookupDuration, durationSeconds(duration), labels)
}

// RecordPermissionDrift incrementa el contador de drift detectado entre los
// paths legacy y grants. source identifica desde qué evaluador se detectó (por
// ejemplo "canary" cuando se ejecutan ambos en paralelo). userID se pasa para
// trazabilidad pero NO se emite como label (alta cardinalidad) — se registra
// vía logger en el sitio de detección.
func (m *Metrics) RecordPermissionDrift(userID, source string) {
	_ = userID // explícitamente fuera de labels: alta cardinalidad
	m.recorder.CounterAdd(MetricPermissionsDriftTotal, 1, map[string]string{
		"service": m.service,
		"source":  source,
	})
}
