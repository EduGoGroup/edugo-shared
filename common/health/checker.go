// Package health define la interfaz comun de health check para recursos de infraestructura.
//
// Los modulos database/postgres, database/mongodb, cache/redis y messaging/rabbit
// implementan patrones de health check similares con timeout de 5 segundos.
// Este paquete centraliza el contrato para nuevos modulos.
package health

import (
	"context"
	"time"
)

// DefaultTimeout es el timeout por defecto para health checks.
// Coincide con el valor usado en database/postgres, database/mongodb y cache/redis.
const DefaultTimeout = 5 * time.Second

// Checker define la interfaz de health check para un recurso de infraestructura.
//
// Las implementaciones deben:
//   - Respetar la cancelacion del contexto
//   - Aplicar un timeout razonable (se recomienda [DefaultTimeout])
//   - Retornar nil si el recurso esta saludable, error en caso contrario
type Checker interface {
	// Check verifica el estado de salud del recurso.
	Check(ctx context.Context) error
}

// CheckerFunc es un adaptador para usar funciones como [Checker].
//
//	checker := health.CheckerFunc(func(ctx context.Context) error {
//	    return db.PingContext(ctx)
//	})
type CheckerFunc func(ctx context.Context) error

// Check implementa [Checker] delegando a la funcion subyacente.
func (f CheckerFunc) Check(ctx context.Context) error {
	return f(ctx)
}

// CheckWithTimeout ejecuta un [Checker] con el timeout especificado.
// Si timeout es 0, usa [DefaultTimeout].
func CheckWithTimeout(ctx context.Context, checker Checker, timeout time.Duration) error {
	if timeout == 0 {
		timeout = DefaultTimeout
	}
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return checker.Check(ctx)
}
