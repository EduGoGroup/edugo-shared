package metrics

import "time"

// Business metric names.
const (
	MetricBusinessOpsTotal    = "business_operations_total"
	MetricBusinessOpsDuration = "business_operations_duration_seconds"
	MetricExportTotal         = "export_operations_total"
	MetricExportDuration      = "export_duration_seconds"
	MetricExportRows          = "export_rows_total"
	MetricAssessmentAttempts  = "assessment_attempts_total"
	MetricGradingTotal        = "grading_operations_total"
	MetricGradingDuration     = "grading_duration_seconds"
	MetricReviewTotal         = "review_operations_total"
	MetricReviewDuration      = "review_duration_seconds"
	MetricNotificationTotal   = "notification_operations_total"
)

// RecordBusinessOperation records a generic business operation.
// entity: "membership", "grade", "attendance", "assessment", "material", "sync", etc.
// operation: "create", "update", "delete", "list", "publish", "archive", etc.
func (m *Metrics) RecordBusinessOperation(entity, operation string, duration time.Duration, err error) {
	labels := map[string]string{
		"service":   m.service,
		"entity":    entity,
		"operation": operation,
		"status":    statusLabel(err),
	}
	m.recorder.CounterAdd(MetricBusinessOpsTotal, 1, labels)
	m.recorder.HistogramObserve(MetricBusinessOpsDuration, durationSeconds(duration), labels)
}

// RecordAssessmentAttempt records an assessment attempt action.
// action: "start", "save_answer", "submit", "view_result"
func (m *Metrics) RecordAssessmentAttempt(action string, duration time.Duration, err error) {
	labels := map[string]string{
		"service": m.service,
		"action":  action,
		"status":  statusLabel(err),
	}
	m.recorder.CounterAdd(MetricAssessmentAttempts, 1, labels)
	m.recorder.HistogramObserve(MetricBusinessOpsDuration, durationSeconds(duration), labels)
}

// RecordGrading records a grading operation for a question.
// questionType: "multiple_choice", "true_false", "short_answer", "open"
func (m *Metrics) RecordGrading(questionType string, duration time.Duration, err error) {
	labels := map[string]string{
		"service":       m.service,
		"question_type": questionType,
		"status":        statusLabel(err),
	}
	m.recorder.CounterAdd(MetricGradingTotal, 1, labels)
	m.recorder.HistogramObserve(MetricGradingDuration, durationSeconds(duration), labels)
}

// RecordReview records an assessment review operation.
// action: "submit", "request_revision", "approve"
func (m *Metrics) RecordReview(action string, duration time.Duration, err error) {
	labels := map[string]string{
		"service": m.service,
		"action":  action,
		"status":  statusLabel(err),
	}
	m.recorder.CounterAdd(MetricReviewTotal, 1, labels)
	m.recorder.HistogramObserve(MetricReviewDuration, durationSeconds(duration), labels)
}

// RecordNotification records a notification operation.
// channel: "push", "in_app", "email"
func (m *Metrics) RecordNotification(channel string, err error) {
	labels := map[string]string{
		"service": m.service,
		"channel": channel,
		"status":  statusLabel(err),
	}
	m.recorder.CounterAdd(MetricNotificationTotal, 1, labels)
}

// RecordExport records an export operation.
// format: "xlsx", "csv", "pdf", "markdown"
func (m *Metrics) RecordExport(format string, rows int, duration time.Duration, err error) {
	labels := map[string]string{
		"service": m.service,
		"format":  format,
		"status":  statusLabel(err),
	}
	m.recorder.CounterAdd(MetricExportTotal, 1, labels)
	m.recorder.HistogramObserve(MetricExportDuration, durationSeconds(duration), labels)
	m.recorder.CounterAdd(MetricExportRows, float64(rows), map[string]string{
		"service": m.service,
		"format":  format,
	})
}
