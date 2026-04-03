package bootstrap

import "testing"

func TestDefaultGORMOptions(t *testing.T) {
	opts := DefaultGORMOptions()

	if !opts.SimpleProtocol {
		t.Error("DefaultGORMOptions().SimpleProtocol should be true")
	}
	if opts.PrepareStmt {
		t.Error("DefaultGORMOptions().PrepareStmt should be false")
	}
	if opts.Logger != nil {
		t.Error("DefaultGORMOptions().Logger should be nil")
	}
}

func TestApplyGORMOptions_NoOpts(t *testing.T) {
	opts := ApplyGORMOptions()
	defaults := DefaultGORMOptions()

	if opts.SimpleProtocol != defaults.SimpleProtocol {
		t.Error("ApplyGORMOptions() without opts should match defaults")
	}
	if opts.PrepareStmt != defaults.PrepareStmt {
		t.Error("ApplyGORMOptions() without opts should match defaults")
	}
}

func TestWithSimpleProtocol(t *testing.T) {
	opts := ApplyGORMOptions(WithSimpleProtocol(false))
	if opts.SimpleProtocol {
		t.Error("WithSimpleProtocol(false) should set SimpleProtocol to false")
	}

	opts = ApplyGORMOptions(WithSimpleProtocol(true))
	if !opts.SimpleProtocol {
		t.Error("WithSimpleProtocol(true) should set SimpleProtocol to true")
	}
}

func TestWithPrepareStmt(t *testing.T) {
	opts := ApplyGORMOptions(WithPrepareStmt(true))
	if !opts.PrepareStmt {
		t.Error("WithPrepareStmt(true) should set PrepareStmt to true")
	}

	opts = ApplyGORMOptions(WithPrepareStmt(false))
	if opts.PrepareStmt {
		t.Error("WithPrepareStmt(false) should set PrepareStmt to false")
	}
}

func TestWithGORMLogger(t *testing.T) {
	logger := "fake-logger"
	opts := ApplyGORMOptions(WithGORMLogger(logger))
	if opts.Logger != logger {
		t.Error("WithGORMLogger should set Logger")
	}
}

func TestApplyGORMOptions_Multiple(t *testing.T) {
	opts := ApplyGORMOptions(
		WithSimpleProtocol(false),
		WithPrepareStmt(true),
		WithGORMLogger("my-logger"),
	)

	if opts.SimpleProtocol {
		t.Error("SimpleProtocol should be false")
	}
	if !opts.PrepareStmt {
		t.Error("PrepareStmt should be true")
	}
	if opts.Logger != "my-logger" {
		t.Error("Logger should be set")
	}
}

func TestApplyGORMOptions_LastWins(t *testing.T) {
	opts := ApplyGORMOptions(
		WithSimpleProtocol(true),
		WithSimpleProtocol(false),
	)

	if opts.SimpleProtocol {
		t.Error("Last option should win: SimpleProtocol should be false")
	}
}
