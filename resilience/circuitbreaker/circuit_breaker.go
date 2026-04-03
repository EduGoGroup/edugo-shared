package circuitbreaker

import (
	"context"
	"errors"
	"sync"
	"time"
)

// State representa el estado del circuit breaker
type State int

const (
	StateClosed   State = iota // Permite todas las peticiones
	StateOpen                  // Rechaza todas las peticiones
	StateHalfOpen              // Permite peticiones limitadas para probar
)

// String convierte el estado a string para logging
func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

var (
	// ErrCircuitOpen se retorna cuando el circuit breaker esta abierto
	ErrCircuitOpen = errors.New("circuit breaker is open")
	// ErrTooManyRequests se retorna cuando hay demasiadas peticiones en half-open
	ErrTooManyRequests = errors.New("too many requests")
)

// MetricsHook permite inyectar metricas opcionales.
type MetricsHook interface {
	SetState(name string, state State)
	RecordTransition(name string, from, to State)
}

// Config configuracion del circuit breaker
type Config struct {
	Name              string        // Nombre del circuit breaker para metricas
	MaxFailures       uint32        // Numero maximo de fallos antes de abrir
	Timeout           time.Duration // Tiempo antes de pasar de open a half-open
	MaxRequests       uint32        // Maximo de peticiones en half-open
	SuccessThreshold  uint32        // Exitos necesarios en half-open para cerrar
	FailureRateWindow time.Duration // Ventana de tiempo para calcular tasa de fallos
	Metrics           MetricsHook   // Hook opcional de metricas
}

// DefaultConfig retorna una configuracion por defecto
func DefaultConfig(name string) Config {
	return Config{
		Name:              name,
		MaxFailures:       5,
		Timeout:           60 * time.Second,
		MaxRequests:       1,
		SuccessThreshold:  2,
		FailureRateWindow: 30 * time.Second,
	}
}

// CircuitBreaker implementa el patron circuit breaker
type CircuitBreaker struct {
	config Config

	mu              sync.RWMutex
	state           State
	failures        uint32
	successes       uint32
	requests        uint32
	lastStateChange time.Time
	lastFailure     time.Time
	metrics         MetricsHook
}

// New crea un nuevo circuit breaker
func New(config Config) *CircuitBreaker {
	cb := &CircuitBreaker{
		config:          config,
		state:           StateClosed,
		lastStateChange: time.Now(),
		metrics:         config.Metrics,
	}

	// Registrar estado inicial en metricas
	if cb.metrics != nil {
		cb.metrics.SetState(config.Name, StateClosed)
	}

	return cb
}

// WithMetrics configura el hook de metricas y retorna el circuit breaker para encadenar.
func (cb *CircuitBreaker) WithMetrics(hook MetricsHook) *CircuitBreaker {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.metrics = hook
	if hook != nil {
		hook.SetState(cb.config.Name, cb.state)
	}
	return cb
}

// Execute ejecuta la funcion proporcionada con proteccion del circuit breaker
func (cb *CircuitBreaker) Execute(ctx context.Context, fn func(context.Context) error) error {
	// Verificar si podemos ejecutar
	if err := cb.beforeRequest(); err != nil {
		return err
	}

	// Ejecutar la funcion
	err := fn(ctx)

	// Registrar el resultado
	cb.afterRequest(err)

	return err
}

// beforeRequest verifica si se puede hacer la peticion
func (cb *CircuitBreaker) beforeRequest() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	switch cb.state {
	case StateClosed:
		// Permitir la peticion
		return nil

	case StateOpen:
		// Verificar si ha pasado el timeout
		if now.Sub(cb.lastStateChange) >= cb.config.Timeout {
			cb.setState(StateHalfOpen, now)
			cb.requests++
			return nil
		}
		return ErrCircuitOpen

	case StateHalfOpen:
		// Limitar el numero de peticiones de forma atomica
		// Verificar e incrementar en el mismo paso para evitar race conditions
		if cb.requests < cb.config.MaxRequests {
			cb.requests++
			return nil
		}
		return ErrTooManyRequests

	default:
		return ErrCircuitOpen
	}
}

// afterRequest registra el resultado de la peticion
func (cb *CircuitBreaker) afterRequest(err error) {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	now := time.Now()

	if err != nil {
		// Peticion fallida
		cb.failures++
		cb.lastFailure = now

		switch cb.state {
		case StateClosed:
			// Verificar si debemos abrir el circuit
			if cb.failures >= cb.config.MaxFailures {
				cb.setState(StateOpen, now)
			}

		case StateHalfOpen:
			// Un fallo en half-open vuelve a abrir el circuit
			cb.setState(StateOpen, now)
		}
	} else {
		// Peticion exitosa
		cb.successes++

		switch cb.state {
		case StateClosed:
			// Resetear el contador de fallos si llevamos tiempo sin fallos
			if now.Sub(cb.lastFailure) >= cb.config.FailureRateWindow {
				cb.failures = 0
			}

		case StateHalfOpen:
			// Verificar si debemos cerrar el circuit
			if cb.successes >= cb.config.SuccessThreshold {
				cb.setState(StateClosed, now)
			}
		}
	}
}

// setState cambia el estado del circuit breaker
func (cb *CircuitBreaker) setState(newState State, now time.Time) {
	if cb.state == newState {
		return
	}

	oldState := cb.state
	cb.state = newState
	cb.lastStateChange = now

	// Resetear contadores segun el nuevo estado
	switch newState {
	case StateClosed:
		cb.failures = 0
		cb.successes = 0
		cb.requests = 0
	case StateOpen:
		cb.requests = 0
		cb.successes = 0
	case StateHalfOpen:
		cb.requests = 0
		cb.successes = 0
	}

	// Registrar transicion en metricas
	if cb.metrics != nil {
		cb.metrics.SetState(cb.config.Name, newState)
		cb.metrics.RecordTransition(cb.config.Name, oldState, newState)
	}
}

// State retorna el estado actual del circuit breaker
func (cb *CircuitBreaker) State() State {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.state
}

// Failures retorna el numero de fallos actuales
func (cb *CircuitBreaker) Failures() uint32 {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.failures
}

// Successes retorna el numero de exitos actuales
func (cb *CircuitBreaker) Successes() uint32 {
	cb.mu.RLock()
	defer cb.mu.RUnlock()
	return cb.successes
}
