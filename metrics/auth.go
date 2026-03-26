package metrics

import "time"

// Auth metric names.
const (
	MetricAuthLoginsTotal       = "auth_logins_total"
	MetricAuthLoginDuration     = "auth_login_duration_seconds"
	MetricAuthTokenRefreshTotal = "auth_token_refresh_total" //nolint:gosec // no es una credencial, es un nombre de métrica
	MetricAuthRateLimitHits     = "auth_rate_limit_hits_total"
	MetricAuthPermissionChecks  = "auth_permission_checks_total"
)

// RecordLogin records a login attempt with success/failure status and duration.
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

// RecordTokenRefresh records a token refresh attempt.
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
}

// RecordRateLimitHit records when a rate limit is triggered for a resource.
func (m *Metrics) RecordRateLimitHit(resource string) {
	m.recorder.CounterAdd(MetricAuthRateLimitHits, 1, map[string]string{
		"service":  m.service,
		"resource": resource,
	})
}

// RecordPermissionCheck records a permission check result.
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
