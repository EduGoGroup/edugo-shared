package postgres

import (
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestConnect_ReturnsErrorWhenPingFails(t *testing.T) {
	cfg := DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = 1
	cfg.ConnectTimeout = 1 * time.Second

	db, err := Connect(&cfg)
	if err == nil {
		if db != nil {
			if closeErr := db.Close(); closeErr != nil {
				t.Errorf("closing db after unexpected success: %v", closeErr)
			}
		}
		t.Fatal("expected ping error, got nil")
	}
	if db != nil {
		t.Fatal("expected nil db on ping failure")
	}
	if !strings.Contains(err.Error(), "failed to ping database") {
		t.Fatalf("expected ping failure message, got: %v", err)
	}
}

func TestHealthCheck_Success(t *testing.T) {
	behavior := &fakeSQLBehavior{}
	db := openFakeDB(t, behavior)

	if err := HealthCheck(db); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}

	if calls := atomic.LoadInt32(&behavior.pingCalls); calls == 0 {
		t.Fatal("expected PingContext to be called")
	}
}

func TestHealthCheck_ReturnsWrappedError(t *testing.T) {
	pingErr := errors.New("ping failed")
	behavior := &fakeSQLBehavior{pingErr: pingErr}
	db := openFakeDB(t, behavior)

	err := HealthCheck(db)
	if err == nil {
		t.Fatal("expected health check error")
	}
	if !errors.Is(err, pingErr) {
		t.Fatalf("expected wrapped ping error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "database health check failed") {
		t.Fatalf("expected health check prefix, got: %v", err)
	}
}

func TestGetStats_ReturnsConfiguredMaxOpenConnections(t *testing.T) {
	db := openFakeDB(t, &fakeSQLBehavior{})
	db.SetMaxOpenConns(7)

	stats := GetStats(db)
	if stats.MaxOpenConnections != 7 {
		t.Fatalf("expected MaxOpenConnections=7, got: %d", stats.MaxOpenConnections)
	}
}

func TestClose_NonNilDB(t *testing.T) {
	behavior := &fakeSQLBehavior{}
	db := openFakeDB(t, behavior)

	if err := db.Ping(); err != nil {
		t.Fatalf("expected ping to succeed, got: %v", err)
	}

	if err := Close(db); err != nil {
		t.Fatalf("expected close success, got: %v", err)
	}

	if calls := atomic.LoadInt32(&behavior.closeCalls); calls == 0 {
		t.Fatal("expected underlying connection close to be called")
	}
}
