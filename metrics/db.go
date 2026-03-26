package metrics

import "time"

// Database metric names.
const (
	MetricDBQueriesTotal    = "db_queries_total"
	MetricDBQueryDuration   = "db_query_duration_seconds"
	MetricDBConnectionsOpen = "db_connections_open"
)

// RecordDBQuery records a database query with type (postgres/mongodb), operation (select/insert/update/delete),
// table name, duration, and error status.
func (m *Metrics) RecordDBQuery(dbType, operation, table string, duration time.Duration, err error) {
	labels := map[string]string{
		"service":   m.service,
		"db_type":   dbType,
		"operation": operation,
		"table":     table,
		"status":    statusLabel(err),
	}
	m.recorder.CounterAdd(MetricDBQueriesTotal, 1, labels)
	m.recorder.HistogramObserve(MetricDBQueryDuration, durationSeconds(duration), labels)
}

// SetDBConnectionsOpen sets the current number of open database connections.
func (m *Metrics) SetDBConnectionsOpen(dbType string, count int) {
	m.recorder.GaugeSet(MetricDBConnectionsOpen, float64(count), map[string]string{
		"service": m.service,
		"db_type": dbType,
	})
}
