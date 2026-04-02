// Package retry provee utilidades de backoff exponencial reutilizables.
//
// Extraido del patron comun en messaging/rabbit (connection.go y dlq.go).
// Disponible para cualquier modulo que necesite reintentos con backoff.
package retry

import "time"

// NextBackoff duplica el delay actual sin exceder maxDelay.
//
// Uso tipico en un loop de reconexion:
//
//	delay := initialDelay
//	for {
//	    if err := connect(); err != nil {
//	        time.Sleep(delay)
//	        delay = retry.NextBackoff(delay, maxDelay)
//	        continue
//	    }
//	    break
//	}
func NextBackoff(current, maxDelay time.Duration) time.Duration {
	next := current * 2
	if next > maxDelay {
		return maxDelay
	}
	return next
}

// ExponentialDelay calcula el delay para un intento dado con backoff exponencial.
//
// Formula: baseDelay * 2^attempt
// El attempt se limita a 30 para evitar overflow (~5.7 anos con base 5s).
//
// Con baseDelay=5s: 5s, 10s, 20s, 40s, 80s...
func ExponentialDelay(baseDelay time.Duration, attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}
	if attempt > 30 {
		attempt = 30
	}
	return baseDelay * time.Duration(1<<uint(attempt)) //nolint:gosec // attempt validado arriba
}
