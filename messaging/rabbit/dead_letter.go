package rabbit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"
	"unicode/utf8"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Cabeceras que el consumidor con cola muerta escribe en los mensajes que republica.
const (
	// HeaderRetryCount cuenta los reintentos ya hechos (0 en la primera entrega).
	HeaderRetryCount = "x-retry-count"
	// HeaderOriginalExchange es el exchange de la última entrega (vacío tras un reintento).
	HeaderOriginalExchange = "x-original-exchange"
	// HeaderOriginalRoutingKey es la clave de la última entrega (la cola, tras un reintento).
	HeaderOriginalRoutingKey = "x-original-routing-key"
	// HeaderFailedAt es el instante Unix (UTC) en que el mensaje se dio por muerto.
	HeaderFailedAt = "x-failed-at"
	// HeaderError lleva el texto del último error del handler, acotado a MaxErrorHeaderBytes.
	HeaderError = "x-error"
	// HeaderDeadReason dice por qué murió: DeadReasonPermanent o DeadReasonMaxRetries.
	HeaderDeadReason = "x-dead-reason"
	// HeaderFirstExchange es el exchange por el que el mensaje entró la primera vez.
	HeaderFirstExchange = "x-first-exchange"
	// HeaderFirstRoutingKey es la clave con la que el mensaje entró la primera vez.
	HeaderFirstRoutingKey = "x-first-routing-key"
)

const (
	// DeadReasonPermanent marca un mensaje muerto sin reintentos por un error permanente.
	DeadReasonPermanent = "permanent"
	// DeadReasonMaxRetries marca un mensaje muerto por agotar los reintentos.
	DeadReasonMaxRetries = "max_retries"

	// MaxErrorHeaderBytes es el tope, en bytes, del texto que viaja en HeaderError.
	MaxErrorHeaderBytes = 1024

	truncationMark = "...[truncated]"
	bodyHashLength = 12
)

// FailureEvent describe un fallo del handler para los ganchos OnRetry y OnDeadLetter.
type FailureEvent struct {
	Err        error         // Error que devolvió el handler
	Queue      string        // Cola de la que se consumió el mensaje
	Reason     string        // DeadReasonPermanent o DeadReasonMaxRetries; vacío en OnRetry
	BodyHash   string        // Primeros caracteres hex del SHA-256 del cuerpo
	Retry      int           // Reintentos ya hechos cuando falló (0 = primera entrega)
	MaxRetries int           // Máximo de reintentos configurado
	Backoff    time.Duration // Espera antes del reintento; 0 en OnDeadLetter
	BodySize   int           // Tamaño del cuerpo en bytes
}

// DeliveryInfo es lo que el handler puede saber de la entrega que está procesando.
type DeliveryInfo struct {
	Queue             string // Cola de la que se consumió el mensaje
	Retry             int    // Reintentos ya hechos (0 = primera entrega)
	MaxRetries        int    // Máximo de reintentos configurado
	DeadLetterEnabled bool   // Si el consumidor tiene cola muerta habilitada
}

// IsLastAttempt dice si un fallo de esta ejecución manda el mensaje a la cola muerta
// por agotar los reintentos. Sin cola muerta nunca hay última ejecución.
func (d DeliveryInfo) IsLastAttempt() bool {
	return d.DeadLetterEnabled && d.Retry >= d.MaxRetries
}

type deliveryInfoKey struct{}

func withDeliveryInfo(ctx context.Context, info DeliveryInfo) context.Context {
	return context.WithValue(ctx, deliveryInfoKey{}, info)
}

// DeliveryFromContext devuelve los datos de la entrega que ConsumeWithDLQ deja en el
// contexto del handler. El segundo valor es falso fuera de ese consumidor.
func DeliveryFromContext(ctx context.Context) (DeliveryInfo, bool) {
	info, ok := ctx.Value(deliveryInfoKey{}).(DeliveryInfo)
	return info, ok
}

type failureAction int

const (
	actionRequeue failureAction = iota
	actionRetry
	actionDeadLetter
)

// decideFailure decide qué se hace con un mensaje cuyo handler falló. Sin cola muerta
// el mensaje vuelve a la cola, sea cual sea la clase del error.
func decideFailure(dlqEnabled, permanent bool, retries, maxRetries int) (failureAction, string) {
	if !dlqEnabled {
		return actionRequeue, ""
	}
	if permanent {
		return actionDeadLetter, DeadReasonPermanent
	}
	if retries+1 > maxRetries {
		return actionDeadLetter, DeadReasonMaxRetries
	}
	return actionRetry, ""
}

// truncateError devuelve el texto del error acotado a limit bytes, cortado en frontera
// de runa y con una marca al final cuando hubo corte.
func truncateError(err error, limit int) string {
	if err == nil {
		return ""
	}
	text := err.Error()
	if len(text) <= limit {
		return text
	}
	cut := max(limit-len(truncationMark), 0)
	for cut > 0 && !utf8.RuneStart(text[cut]) {
		cut--
	}
	return text[:cut] + truncationMark
}

// withFirstRoute copia las cabeceras y guarda, si aún no están, el exchange y la clave
// por los que el mensaje entró. No muta las cabeceras recibidas.
func withFirstRoute(headers amqp.Table, exchange, routingKey string) amqp.Table {
	out := cloneHeaders(headers)
	if _, ok := out[HeaderFirstExchange]; !ok {
		out[HeaderFirstExchange] = exchange
	}
	if _, ok := out[HeaderFirstRoutingKey]; !ok {
		out[HeaderFirstRoutingKey] = routingKey
	}
	return out
}

// deadLetterHeaders arma las cabeceras del mensaje muerto sin mutar las de la entrega.
func deadLetterHeaders(msg *amqp.Delivery, reason string, cause error, now time.Time) amqp.Table {
	out := withFirstRoute(msg.Headers, msg.Exchange, msg.RoutingKey)
	out[HeaderOriginalExchange] = msg.Exchange
	out[HeaderOriginalRoutingKey] = msg.RoutingKey
	out[HeaderFailedAt] = now.UTC().Unix()
	out[HeaderRetryCount] = getRetryCount(msg.Headers)
	out[HeaderError] = truncateError(cause, MaxErrorHeaderBytes)
	out[HeaderDeadReason] = reason
	return out
}

// bodyHash devuelve un hash corto del cuerpo, para identificar un mensaje en el log.
func bodyHash(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])[:bodyHashLength]
}

func newFailureEvent(queue string, msg *amqp.Delivery, retries, maxRetries int, cause error) FailureEvent {
	return FailureEvent{
		Err:        cause,
		Queue:      queue,
		BodyHash:   bodyHash(msg.Body),
		Retry:      retries,
		MaxRetries: maxRetries,
		BodySize:   len(msg.Body),
	}
}
