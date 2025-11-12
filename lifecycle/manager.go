package lifecycle

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/EduGoGroup/edugo-shared/logger"
	"go.uber.org/zap"
)

// Resource representa un recurso con startup y cleanup
type Resource struct {
	Name    string
	Startup func(ctx context.Context) error
	Cleanup func() error
}

// Manager gestiona el ciclo de vida de recursos de infraestructura
// Maneja startup y cleanup en orden LIFO (Last In, First Out)
type Manager struct {
	resources []Resource
	mu        sync.Mutex
	logger    logger.Logger
	startTime time.Time
}

// NewManager crea un nuevo lifecycle manager
func NewManager(log logger.Logger) *Manager {
	return &Manager{
		resources: make([]Resource, 0),
		logger:    log,
		startTime: time.Now(),
	}
}

// Register registra un recurso para gestión de ciclo de vida
// Los recursos se limpian en orden inverso al registro (LIFO)
func (m *Manager) Register(name string, startup func(ctx context.Context) error, cleanup func() error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.resources = append(m.resources, Resource{
		Name:    name,
		Startup: startup,
		Cleanup: cleanup,
	})

	m.logger.Debug("resource registered for lifecycle management",
		zap.String("resource", name),
		zap.Int("total_resources", len(m.resources)))
}

// RegisterSimple registra un recurso solo con cleanup (sin startup)
// Útil para recursos que ya están inicializados
func (m *Manager) RegisterSimple(name string, cleanup func() error) {
	m.Register(name, nil, cleanup)
}

// Startup ejecuta el startup de todos los recursos en orden de registro
func (m *Manager) Startup(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("starting lifecycle startup phase",
		zap.Int("total_resources", len(m.resources)))

	for i, resource := range m.resources {
		if resource.Startup == nil {
			m.logger.Debug("resource has no startup function, skipping",
				zap.String("resource", resource.Name))
			continue
		}

		m.logger.Debug("starting up resource",
			zap.String("resource", resource.Name),
			zap.Int("index", i+1),
			zap.Int("total", len(m.resources)))

		startTime := time.Now()

		if err := resource.Startup(ctx); err != nil {
			m.logger.Error("resource startup failed",
				zap.String("resource", resource.Name),
				zap.Error(err),
				zap.Duration("duration", time.Since(startTime)))
			return fmt.Errorf("failed to startup resource %s: %w", resource.Name, err)
		}

		m.logger.Debug("resource started successfully",
			zap.String("resource", resource.Name),
			zap.Duration("duration", time.Since(startTime)))
	}

	m.logger.Info("lifecycle startup phase completed",
		zap.Duration("total_duration", time.Since(m.startTime)))

	return nil
}

// Cleanup ejecuta cleanup de todos los recursos en orden inverso (LIFO)
// Continúa limpiando incluso si algunos fallan, acumulando errores
func (m *Manager) Cleanup() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logger.Info("starting lifecycle cleanup phase",
		zap.Int("total_resources", len(m.resources)))

	var errors []error
	cleanupStartTime := time.Now()

	// Cleanup en orden inverso (LIFO)
	for i := len(m.resources) - 1; i >= 0; i-- {
		resource := m.resources[i]

		if resource.Cleanup == nil {
			m.logger.Debug("resource has no cleanup function, skipping",
				zap.String("resource", resource.Name))
			continue
		}

		m.logger.Debug("cleaning up resource",
			zap.String("resource", resource.Name),
			zap.Int("index", i+1),
			zap.Int("total", len(m.resources)))

		startTime := time.Now()

		if err := resource.Cleanup(); err != nil {
			m.logger.Error("resource cleanup failed",
				zap.String("resource", resource.Name),
				zap.Error(err),
				zap.Duration("duration", time.Since(startTime)))
			errors = append(errors, fmt.Errorf("%s: %w", resource.Name, err))
		} else {
			m.logger.Debug("resource cleaned up successfully",
				zap.String("resource", resource.Name),
				zap.Duration("duration", time.Since(startTime)))
		}
	}

	if len(errors) > 0 {
		m.logger.Error("lifecycle cleanup completed with errors",
			zap.Int("error_count", len(errors)),
			zap.Duration("total_duration", time.Since(cleanupStartTime)))
		return fmt.Errorf("cleanup failed for %d resource(s): %v", len(errors), errors)
	}

	m.logger.Info("lifecycle cleanup phase completed successfully",
		zap.Duration("total_duration", time.Since(cleanupStartTime)))

	return nil
}

// Count retorna el número de recursos registrados
func (m *Manager) Count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.resources)
}

// Clear limpia la lista de recursos sin ejecutar cleanup
// Útil para testing
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.resources = make([]Resource, 0)
	m.logger.Debug("lifecycle manager cleared")
}
