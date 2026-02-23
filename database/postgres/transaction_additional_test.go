package postgres

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
)

func TestWithTransaction_ReturnsBeginError(t *testing.T) {
	beginErr := errors.New("begin failed")
	db := openFakeDB(t, &fakeSQLBehavior{beginErr: beginErr})

	err := WithTransaction(context.Background(), db, func(*sql.Tx) error { return nil })
	if err == nil {
		t.Fatal("expected begin error")
	}
	if !errors.Is(err, beginErr) {
		t.Fatalf("expected wrapped begin error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "failed to begin transaction") {
		t.Fatalf("expected begin failure prefix, got: %v", err)
	}
}

func TestWithTransaction_RollsBackOnCallbackError(t *testing.T) {
	behavior := &fakeSQLBehavior{}
	db := openFakeDB(t, behavior)
	callbackErr := errors.New("callback failed")

	err := WithTransaction(context.Background(), db, func(*sql.Tx) error {
		return callbackErr
	})
	if err == nil {
		t.Fatal("expected callback error")
	}
	if !errors.Is(err, callbackErr) {
		t.Fatalf("expected callback error to be returned, got: %v", err)
	}
	if calls := atomic.LoadInt32(&behavior.rollbackCalls); calls != 1 {
		t.Fatalf("expected rollback to be called once, got: %d", calls)
	}
	if calls := atomic.LoadInt32(&behavior.commitCalls); calls != 0 {
		t.Fatalf("expected commit to not be called, got: %d", calls)
	}
}

func TestWithTransaction_ReturnsJoinedErrorWhenRollbackFails(t *testing.T) {
	rollbackErr := errors.New("rollback failed")
	behavior := &fakeSQLBehavior{rollbackErr: rollbackErr}
	db := openFakeDB(t, behavior)
	callbackErr := errors.New("callback failed")

	err := WithTransaction(context.Background(), db, func(*sql.Tx) error {
		return callbackErr
	})
	if err == nil {
		t.Fatal("expected joined error")
	}
	if !errors.Is(err, callbackErr) {
		t.Fatalf("expected callback error to be wrapped, got: %v", err)
	}
	if !errors.Is(err, rollbackErr) {
		t.Fatalf("expected rollback error to be wrapped, got: %v", err)
	}
}

func TestWithTransaction_ReturnsCommitError(t *testing.T) {
	commitErr := errors.New("commit failed")
	behavior := &fakeSQLBehavior{commitErr: commitErr}
	db := openFakeDB(t, behavior)

	err := WithTransaction(context.Background(), db, func(*sql.Tx) error { return nil })
	if err == nil {
		t.Fatal("expected commit error")
	}
	if !errors.Is(err, commitErr) {
		t.Fatalf("expected wrapped commit error, got: %v", err)
	}
	if !strings.Contains(err.Error(), "failed to commit transaction") {
		t.Fatalf("expected commit failure prefix, got: %v", err)
	}
}

func TestWithTransaction_CommitsOnSuccess(t *testing.T) {
	behavior := &fakeSQLBehavior{}
	db := openFakeDB(t, behavior)
	called := false

	err := WithTransaction(context.Background(), db, func(*sql.Tx) error {
		called = true
		return nil
	})
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}
	if !called {
		t.Fatal("expected callback to be called")
	}
	if calls := atomic.LoadInt32(&behavior.commitCalls); calls != 1 {
		t.Fatalf("expected commit once, got: %d", calls)
	}
}

func TestWithTransaction_RollsBackAndRepanics(t *testing.T) {
	behavior := &fakeSQLBehavior{}
	db := openFakeDB(t, behavior)

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected panic to be propagated")
		}
		if recovered != "boom" {
			t.Fatalf("unexpected panic value: %v", recovered)
		}
		if calls := atomic.LoadInt32(&behavior.rollbackCalls); calls != 1 {
			t.Fatalf("expected rollback on panic, got: %d", calls)
		}
	}()

	if err := WithTransaction(context.Background(), db, func(*sql.Tx) error {
		panic("boom")
	}); err != nil {
		t.Fatalf("unexpected error (panic path should re-panic): %v", err)
	}
}

func TestWithTransactionIsolation_ReturnsBeginError(t *testing.T) {
	beginErr := errors.New("begin failed")
	db := openFakeDB(t, &fakeSQLBehavior{beginErr: beginErr})

	err := WithTransactionIsolation(
		context.Background(),
		db,
		sql.LevelSerializable,
		func(*sql.Tx) error { return nil },
	)
	if err == nil {
		t.Fatal("expected begin error")
	}
	if !errors.Is(err, beginErr) {
		t.Fatalf("expected wrapped begin error, got: %v", err)
	}
}

func TestWithTransactionIsolation_RollsBackOnCallbackError(t *testing.T) {
	behavior := &fakeSQLBehavior{}
	db := openFakeDB(t, behavior)
	callbackErr := errors.New("isolation callback failed")

	err := WithTransactionIsolation(
		context.Background(),
		db,
		sql.LevelReadCommitted,
		func(*sql.Tx) error { return callbackErr },
	)
	if err == nil {
		t.Fatal("expected callback error")
	}
	if !errors.Is(err, callbackErr) {
		t.Fatalf("expected callback error to be returned, got: %v", err)
	}
	if calls := atomic.LoadInt32(&behavior.rollbackCalls); calls != 1 {
		t.Fatalf("expected rollback once, got: %d", calls)
	}
}

func TestWithTransactionIsolation_ReturnsCommitError(t *testing.T) {
	commitErr := errors.New("commit failed")
	behavior := &fakeSQLBehavior{commitErr: commitErr}
	db := openFakeDB(t, behavior)

	err := WithTransactionIsolation(
		context.Background(),
		db,
		sql.LevelReadCommitted,
		func(*sql.Tx) error { return nil },
	)
	if err == nil {
		t.Fatal("expected commit error")
	}
	if !errors.Is(err, commitErr) {
		t.Fatalf("expected wrapped commit error, got: %v", err)
	}
}

func TestWithTransactionIsolation_RollsBackAndRepanics(t *testing.T) {
	behavior := &fakeSQLBehavior{}
	db := openFakeDB(t, behavior)

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected panic to be propagated")
		}
		if recovered != "boom isolation" {
			t.Fatalf("unexpected panic value: %v", recovered)
		}
		if calls := atomic.LoadInt32(&behavior.rollbackCalls); calls != 1 {
			t.Fatalf("expected rollback on panic, got: %d", calls)
		}
	}()

	if err := WithTransactionIsolation(
		context.Background(),
		db,
		sql.LevelSerializable,
		func(*sql.Tx) error {
			panic("boom isolation")
		},
	); err != nil {
		t.Fatalf("unexpected error (panic path should re-panic): %v", err)
	}
}
