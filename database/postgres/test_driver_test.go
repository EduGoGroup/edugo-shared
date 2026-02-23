package postgres

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

type fakeSQLBehavior struct {
	beginErr    error
	commitErr   error
	rollbackErr error
	pingErr     error
	closeErr    error

	beginCalls    int32
	commitCalls   int32
	rollbackCalls int32
	pingCalls     int32
	closeCalls    int32
}

type fakeSQLDriver struct{}

type fakeSQLConn struct {
	behavior *fakeSQLBehavior
}

type fakeSQLTx struct {
	behavior *fakeSQLBehavior
}

var (
	fakeDriverName = registerFakeSQLDriver()
	fakeCasesMu    sync.Mutex
	fakeCases      = map[string]*fakeSQLBehavior{}
	fakeCaseID     int64
)

func registerFakeSQLDriver() string {
	name := "edugo_postgres_fake_driver"
	sql.Register(name, &fakeSQLDriver{})
	return name
}

func openFakeDB(t *testing.T, behavior *fakeSQLBehavior) *sql.DB {
	t.Helper()

	caseID := fmt.Sprintf("case-%d", atomic.AddInt64(&fakeCaseID, 1))

	fakeCasesMu.Lock()
	fakeCases[caseID] = behavior
	fakeCasesMu.Unlock()

	db, err := sql.Open(fakeDriverName, caseID)
	if err != nil {
		t.Fatalf("opening fake db: %v", err)
	}

	t.Cleanup(func() {
		if closeErr := db.Close(); closeErr != nil {
			t.Errorf("closing fake db: %v", closeErr)
		}
		fakeCasesMu.Lock()
		delete(fakeCases, caseID)
		fakeCasesMu.Unlock()
	})

	return db
}

func (d *fakeSQLDriver) Open(name string) (driver.Conn, error) {
	fakeCasesMu.Lock()
	behavior := fakeCases[name]
	fakeCasesMu.Unlock()

	if behavior == nil {
		return nil, fmt.Errorf("fake behavior not found for case %q", name)
	}

	return &fakeSQLConn{behavior: behavior}, nil
}

func (c *fakeSQLConn) Prepare(string) (driver.Stmt, error) {
	return nil, fmt.Errorf("prepare not implemented in fake driver")
}

func (c *fakeSQLConn) Close() error {
	atomic.AddInt32(&c.behavior.closeCalls, 1)
	return c.behavior.closeErr
}

func (c *fakeSQLConn) Begin() (driver.Tx, error) {
	return c.BeginTx(context.Background(), driver.TxOptions{})
}

func (c *fakeSQLConn) BeginTx(context.Context, driver.TxOptions) (driver.Tx, error) {
	atomic.AddInt32(&c.behavior.beginCalls, 1)
	if c.behavior.beginErr != nil {
		return nil, c.behavior.beginErr
	}
	return &fakeSQLTx{behavior: c.behavior}, nil
}

func (c *fakeSQLConn) Ping(context.Context) error {
	atomic.AddInt32(&c.behavior.pingCalls, 1)
	return c.behavior.pingErr
}

func (tx *fakeSQLTx) Commit() error {
	atomic.AddInt32(&tx.behavior.commitCalls, 1)
	return tx.behavior.commitErr
}

func (tx *fakeSQLTx) Rollback() error {
	atomic.AddInt32(&tx.behavior.rollbackCalls, 1)
	return tx.behavior.rollbackErr
}
