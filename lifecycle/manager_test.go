package lifecycle

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/logger"
)

func TestNewManager(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	if mgr == nil {
		t.Fatal("NewManager() returned nil")
	}

	if mgr.Count() != 0 {
		t.Errorf("Count() = %d, want 0", mgr.Count())
	}
}

func TestManager_Register(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	cleanup := func() error {
		return nil
	}

	mgr.Register("test-resource", nil, cleanup)

	if mgr.Count() != 1 {
		t.Errorf("Count() = %d, want 1", mgr.Count())
	}
}

func TestManager_RegisterSimple(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	called := false
	cleanup := func() error {
		called = true
		return nil
	}

	mgr.RegisterSimple("test-resource", cleanup)

	if mgr.Count() != 1 {
		t.Errorf("Count() = %d, want 1", mgr.Count())
	}

	// Ejecutar cleanup para verificar que funciona
	err := mgr.Cleanup()
	if err != nil {
		t.Errorf("Cleanup() error = %v, want nil", err)
	}

	if !called {
		t.Error("cleanup function was not called")
	}
}

func TestManager_Startup_Success(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	startupCalled := false
	startup := func(ctx context.Context) error {
		startupCalled = true
		return nil
	}

	mgr.Register("test-resource", startup, nil)

	ctx := context.Background()
	err := mgr.Startup(ctx)

	if err != nil {
		t.Errorf("Startup() error = %v, want nil", err)
	}

	if !startupCalled {
		t.Error("startup function was not called")
	}
}

func TestManager_Startup_Error(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	expectedErr := errors.New("startup failed")
	startup := func(ctx context.Context) error {
		return expectedErr
	}

	mgr.Register("failing-resource", startup, nil)

	ctx := context.Background()
	err := mgr.Startup(ctx)

	if err == nil {
		t.Error("Startup() error = nil, want error")
	}
}

func TestManager_Cleanup_Success(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	cleanupCalled := false
	cleanup := func() error {
		cleanupCalled = true
		return nil
	}

	mgr.Register("test-resource", nil, cleanup)

	err := mgr.Cleanup()

	if err != nil {
		t.Errorf("Cleanup() error = %v, want nil", err)
	}

	if !cleanupCalled {
		t.Error("cleanup function was not called")
	}
}

func TestManager_Cleanup_MultipleResources_LIFO(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	var order []string

	cleanup1 := func() error {
		order = append(order, "resource1")
		return nil
	}

	cleanup2 := func() error {
		order = append(order, "resource2")
		return nil
	}

	cleanup3 := func() error {
		order = append(order, "resource3")
		return nil
	}

	mgr.Register("resource1", nil, cleanup1)
	mgr.Register("resource2", nil, cleanup2)
	mgr.Register("resource3", nil, cleanup3)

	err := mgr.Cleanup()

	if err != nil {
		t.Errorf("Cleanup() error = %v, want nil", err)
	}

	// Verificar orden LIFO (Last In, First Out)
	expectedOrder := []string{"resource3", "resource2", "resource1"}
	if len(order) != len(expectedOrder) {
		t.Errorf("cleanup order length = %d, want %d", len(order), len(expectedOrder))
	}

	for i, name := range order {
		if name != expectedOrder[i] {
			t.Errorf("cleanup order[%d] = %s, want %s", i, name, expectedOrder[i])
		}
	}
}

func TestManager_Cleanup_WithErrors(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	cleanup1 := func() error {
		return nil
	}

	cleanup2 := func() error {
		return errors.New("cleanup2 failed")
	}

	cleanup3 := func() error {
		return nil
	}

	mgr.Register("resource1", nil, cleanup1)
	mgr.Register("resource2", nil, cleanup2)
	mgr.Register("resource3", nil, cleanup3)

	err := mgr.Cleanup()

	// Debe retornar error pero haber intentado limpiar todos los recursos
	if err == nil {
		t.Error("Cleanup() error = nil, want error")
	}
}

func TestManager_Clear(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	mgr.Register("resource1", nil, func() error { return nil })
	mgr.Register("resource2", nil, func() error { return nil })

	if mgr.Count() != 2 {
		t.Errorf("Count() = %d, want 2", mgr.Count())
	}

	mgr.Clear()

	if mgr.Count() != 0 {
		t.Errorf("Count() after Clear() = %d, want 0", mgr.Count())
	}
}

func TestManager_Startup_WithContext(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	mgr := NewManager(log)

	var receivedCtx context.Context
	startup := func(ctx context.Context) error {
		receivedCtx = ctx
		return nil
	}

	mgr.Register("test-resource", startup, nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := mgr.Startup(ctx)

	if err != nil {
		t.Errorf("Startup() error = %v, want nil", err)
	}

	if receivedCtx == nil {
		t.Error("context was not passed to startup function")
	}
}
