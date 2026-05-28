package bootstrap

// LifecycleManager gestiona el ciclo de vida de recursos de infraestructura.
// Reemplaza el parametro `lifecycleManager any` con un contrato tipado.
type LifecycleManager interface {
	// RegisterCleanup registra una funcion de limpieza con nombre para shutdown.
	// Las funciones se ejecutan en orden LIFO durante el graceful shutdown.
	RegisterCleanup(name string, cleanup func() error)
}
