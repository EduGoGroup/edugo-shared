package metrics

import "time"

// Nombres de métricas de autenticación.
const (
	MetricAuthLoginsTotal          = "auth_logins_total"
	MetricAuthLoginDuration        = "auth_login_duration_seconds"
	MetricAuthTokenRefreshTotal    = "auth_token_refresh_total"    //nolint:gosec // no es una credencial, es un nombre de métrica
	MetricAuthTokenRefreshDuration = "auth_token_refresh_duration" //nolint:gosec // no es una credencial, es un nombre de métrica
	MetricAuthRateLimitHits        = "auth_rate_limit_hits_total"
	MetricAuthPermissionChecks     = "auth_permission_checks_total"
)

// RecordLogin registra un intento de login con estado exitoso/fallido y duración.
func (m *Metrics) RecordLogin(success bool, duration time.Duration) {
	status := "success"
	if !success {
		status = "failure"
	}
	labels := map[string]string{
		"service": m.service,
		"status":  status,
	}
	m.recorder.CounterAdd(MetricAuthLoginsTotal, 1, labels)
	m.recorder.HistogramObserve(MetricAuthLoginDuration, durationSeconds(duration), labels)
}

// RecordTokenRefresh registra un intento de refresh de token con estado y duración.
func (m *Metrics) RecordTokenRefresh(success bool, duration time.Duration) {
	status := "success"
	if !success {
		status = "failure"
	}
	labels := map[string]string{
		"service": m.service,
		"status":  status,
	}
	m.recorder.CounterAdd(MetricAuthTokenRefreshTotal, 1, labels)
	m.recorder.HistogramObserve(MetricAuthTokenRefreshDuration, durationSeconds(duration), labels)
}

// RecordRateLimitHit registra cuando se activa un rate limit para un recurso.
func (m *Metrics) RecordRateLimitHit(resource string) {
	m.recorder.CounterAdd(MetricAuthRateLimitHits, 1, map[string]string{
		"service":  m.service,
		"resource": resource,
	})
}

// RecordPermissionCheck registra el resultado de una verificación de permisos.
func (m *Metrics) RecordPermissionCheck(permission string, granted bool) {
	result := "granted"
	if !granted {
		result = "denied"
	}
	m.recorder.CounterAdd(MetricAuthPermissionChecks, 1, map[string]string{
		"service":    m.service,
		"permission": permission,
		"result":     result,
	})
}
