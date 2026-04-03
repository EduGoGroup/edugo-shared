package bootstrap

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrMissingFactory_Error(t *testing.T) {
	err := ErrMissingFactory{Resource: "postgresql"}
	got := err.Error()
	want := "missing required factory: postgresql"
	if got != want {
		t.Errorf("ErrMissingFactory.Error() = %q, want %q", got, want)
	}
}

func TestErrConnectionFailed_Error(t *testing.T) {
	cause := fmt.Errorf("connection refused")
	err := ErrConnectionFailed{Resource: "mongodb", Err: cause}
	got := err.Error()
	want := "bootstrap/mongodb: connection failed: connection refused"
	if got != want {
		t.Errorf("ErrConnectionFailed.Error() = %q, want %q", got, want)
	}
}

func TestErrConnectionFailed_Unwrap(t *testing.T) {
	cause := fmt.Errorf("timeout")
	err := ErrConnectionFailed{Resource: "rabbitmq", Err: cause}

	if !errors.Is(err, cause) {
		t.Error("ErrConnectionFailed.Unwrap() should return the wrapped error")
	}
}

func TestErrConnectionFailed_ErrorsAs(t *testing.T) {
	cause := fmt.Errorf("dial tcp: connection refused")
	err := ErrConnectionFailed{Resource: "postgresql", Err: cause}

	var connErr ErrConnectionFailed
	if !errors.As(err, &connErr) {
		t.Error("errors.As should match ErrConnectionFailed")
	}
	if connErr.Resource != "postgresql" {
		t.Errorf("Resource = %q, want %q", connErr.Resource, "postgresql")
	}
}
