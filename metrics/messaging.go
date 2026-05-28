package metrics

import (
	"strconv"
	"time"
)

// Messaging metric names.
const (
	MetricMessagesProcessedTotal = "messages_processed_total"
	MetricMessageDuration        = "message_processing_duration_seconds"
	MetricMessageRetriesTotal    = "message_retries_total"
	MetricMessagesInQueue        = "messages_in_queue"
	MetricCircuitBreakerState    = "circuit_breaker_state"
)

// RecordMessageProcessed records a processed message from a queue.
// eventType: "material_uploaded", "material_deleted", "assessment_attempt", etc.
func (m *Metrics) RecordMessageProcessed(eventType string, duration time.Duration, err error) {
	labels := map[string]string{
		"service":    m.service,
		"event_type": eventType,
		"status":     statusLabel(err),
	}
	m.recorder.CounterAdd(MetricMessagesProcessedTotal, 1, labels)
	m.recorder.HistogramObserve(MetricMessageDuration, durationSeconds(duration), labels)
}

// RecordMessageRetry records a message retry attempt.
func (m *Metrics) RecordMessageRetry(eventType string, attempt int) {
	m.recorder.CounterAdd(MetricMessageRetriesTotal, 1, map[string]string{
		"service":    m.service,
		"event_type": eventType,
		"attempt":    strconv.Itoa(attempt),
	})
}

// SetMessagesInQueue sets the current number of messages waiting in the queue.
func (m *Metrics) SetMessagesInQueue(queueName string, count int) {
	m.recorder.GaugeSet(MetricMessagesInQueue, float64(count), map[string]string{
		"service": m.service,
		"queue":   queueName,
	})
}

// SetCircuitBreakerState records the current state of a circuit breaker.
// state: "closed", "half_open", "open"
func (m *Metrics) SetCircuitBreakerState(targetService, state string) {
	stateValue := 0.0
	switch state {
	case "closed":
		stateValue = 0
	case "half_open":
		stateValue = 1
	case "open":
		stateValue = 2
	}
	m.recorder.GaugeSet(MetricCircuitBreakerState, stateValue, map[string]string{
		"service":        m.service,
		"target_service": targetService,
		"state":          state,
	})
}
