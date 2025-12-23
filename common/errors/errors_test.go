package errors_test

import (
	stderrors "errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-shared/common/errors"
)

func TestNew(t *testing.T) {
	err := errors.New(errors.ErrorCodeValidation, "validation failed")

	if err.Code != errors.ErrorCodeValidation {
		t.Errorf("Expected code %s, got %s", errors.ErrorCodeValidation, err.Code)
	}
	if err.Message != "validation failed" {
		t.Errorf("Expected message 'validation failed', got '%s'", err.Message)
	}
	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, err.StatusCode)
	}
}

func TestWrap(t *testing.T) {
	originalErr := fmt.Errorf("original error")
	err := errors.Wrap(originalErr, errors.ErrorCodeInternal, "wrapped error")

	if err.Code != errors.ErrorCodeInternal {
		t.Errorf("Expected code %s, got %s", errors.ErrorCodeInternal, err.Code)
	}
	if !stderrors.Is(err.Internal, originalErr) {
		t.Error("Internal error not set correctly")
	}
}

func TestAppError_Error(t *testing.T) {
	t.Run("without_internal_error", func(t *testing.T) {
		err := errors.New(errors.ErrorCodeNotFound, "resource not found")
		expected := "NOT_FOUND: resource not found"
		if err.Error() != expected {
			t.Errorf("Expected '%s', got '%s'", expected, err.Error())
		}
	})

	t.Run("with_internal_error", func(t *testing.T) {
		internal := fmt.Errorf("internal problem")
		err := errors.Wrap(internal, errors.ErrorCodeInternal, "server error")
		if !strings.Contains(err.Error(), "internal problem") {
			t.Error("Error message should include internal error")
		}
	})
}

func TestAppError_WithDetails(t *testing.T) {
	err := errors.New(errors.ErrorCodeValidation, "validation failed")
	err = err.WithDetails("field 'email' is invalid")

	if err.Details != "field 'email' is invalid" {
		t.Errorf("Details not set correctly: %s", err.Details)
	}
}

func TestAppError_WithField(t *testing.T) {
	err := errors.New(errors.ErrorCodeNotFound, "not found")
	err = err.WithField("resource", "user").WithField("id", 123)

	if err.Fields["resource"] != "user" {
		t.Error("Field 'resource' not set")
	}
	if err.Fields["id"] != 123 {
		t.Error("Field 'id' not set")
	}
}

func TestAppError_WithInternal(t *testing.T) {
	internal := fmt.Errorf("db connection failed")
	err := errors.New(errors.ErrorCodeDatabaseError, "database error")
	err = err.WithInternal(internal)

	if !stderrors.Is(err.Internal, internal) {
		t.Error("Internal error not set")
	}
}

func TestNewValidationError(t *testing.T) {
	err := errors.NewValidationError("invalid input")

	if err.Code != errors.ErrorCodeValidation {
		t.Error("Wrong error code")
	}
	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, err.StatusCode)
	}
}

func TestNewNotFoundError(t *testing.T) {
	err := errors.NewNotFoundError("user")

	if err.Code != errors.ErrorCodeNotFound {
		t.Error("Wrong error code")
	}
	if !strings.Contains(err.Message, "user") {
		t.Error("Message should contain resource name")
	}
	if err.Fields["resource"] != "user" {
		t.Error("Resource field not set")
	}
	if err.StatusCode != http.StatusNotFound {
		t.Error("Wrong status code")
	}
}

func TestNewAlreadyExistsError(t *testing.T) {
	err := errors.NewAlreadyExistsError("email")

	if err.Code != errors.ErrorCodeAlreadyExists {
		t.Error("Wrong error code")
	}
	if err.StatusCode != http.StatusConflict {
		t.Error("Wrong status code")
	}
}

func TestNewUnauthorizedError(t *testing.T) {
	t.Run("with_message", func(t *testing.T) {
		err := errors.NewUnauthorizedError("invalid credentials")
		if err.Message != "invalid credentials" {
			t.Error("Message not set correctly")
		}
	})

	t.Run("empty_message", func(t *testing.T) {
		err := errors.NewUnauthorizedError("")
		if err.Message != "unauthorized" {
			t.Error("Should use default message")
		}
	})

	err := errors.NewUnauthorizedError("test")
	if err.StatusCode != http.StatusUnauthorized {
		t.Error("Wrong status code")
	}
}

func TestNewForbiddenError(t *testing.T) {
	t.Run("with_message", func(t *testing.T) {
		err := errors.NewForbiddenError("insufficient permissions")
		if err.Message != "insufficient permissions" {
			t.Error("Message not set correctly")
		}
	})

	t.Run("empty_message", func(t *testing.T) {
		err := errors.NewForbiddenError("")
		if err.Message != "forbidden" {
			t.Error("Should use default message")
		}
	})

	err := errors.NewForbiddenError("test")
	if err.StatusCode != http.StatusForbidden {
		t.Error("Wrong status code")
	}
}

func TestNewInternalError(t *testing.T) {
	originalErr := fmt.Errorf("panic occurred")

	t.Run("with_message", func(t *testing.T) {
		err := errors.NewInternalError("service crashed", originalErr)
		if err.Message != "service crashed" {
			t.Error("Message not set correctly")
		}
	})

	t.Run("empty_message", func(t *testing.T) {
		err := errors.NewInternalError("", originalErr)
		if err.Message != "internal server error" {
			t.Error("Should use default message")
		}
	})

	err := errors.NewInternalError("test", originalErr)
	if err.StatusCode != http.StatusInternalServerError {
		t.Error("Wrong status code")
	}
}

func TestNewDatabaseError(t *testing.T) {
	dbErr := fmt.Errorf("connection timeout")
	err := errors.NewDatabaseError("insert", dbErr)

	if !strings.Contains(err.Message, "insert") {
		t.Error("Message should contain operation")
	}
	if err.Fields["operation"] != "insert" {
		t.Error("Operation field not set")
	}
	if err.StatusCode != http.StatusInternalServerError {
		t.Error("Wrong status code")
	}
}

func TestNewBusinessRuleError(t *testing.T) {
	err := errors.NewBusinessRuleError("cannot delete active user")

	if err.Code != errors.ErrorCodeBusinessRule {
		t.Error("Wrong error code")
	}
	if err.StatusCode != http.StatusUnprocessableEntity {
		t.Error("Wrong status code")
	}
}

func TestNewConflictError(t *testing.T) {
	err := errors.NewConflictError("resource conflict")

	if err.Code != errors.ErrorCodeConflict {
		t.Error("Wrong error code")
	}
	if err.StatusCode != http.StatusConflict {
		t.Error("Wrong status code")
	}
}

func TestNewRateLimitError(t *testing.T) {
	err := errors.NewRateLimitError()

	if err.Code != errors.ErrorCodeRateLimit {
		t.Error("Wrong error code")
	}
	if err.StatusCode != http.StatusTooManyRequests {
		t.Error("Wrong status code")
	}
}

func TestGetDefaultStatusCode(t *testing.T) {
	tests := []struct {
		code       errors.ErrorCode
		statusCode int
	}{
		{errors.ErrorCodeValidation, http.StatusBadRequest},
		{errors.ErrorCodeInvalidInput, http.StatusBadRequest},
		{errors.ErrorCodeNotFound, http.StatusNotFound},
		{errors.ErrorCodeAlreadyExists, http.StatusConflict},
		{errors.ErrorCodeConflict, http.StatusConflict},
		{errors.ErrorCodeUnauthorized, http.StatusUnauthorized},
		{errors.ErrorCodeInvalidToken, http.StatusUnauthorized},
		{errors.ErrorCodeTokenExpired, http.StatusUnauthorized},
		{errors.ErrorCodeForbidden, http.StatusForbidden},
		{errors.ErrorCodeBusinessRule, http.StatusUnprocessableEntity},
		{errors.ErrorCodeInvalidState, http.StatusUnprocessableEntity},
		{errors.ErrorCodeRateLimit, http.StatusTooManyRequests},
		{errors.ErrorCodeTimeout, http.StatusRequestTimeout},
		{errors.ErrorCodeDatabaseError, http.StatusInternalServerError},
		{errors.ErrorCodeInternal, http.StatusInternalServerError},
		{errors.ErrorCodeExternalService, http.StatusInternalServerError},
		{errors.ErrorCode("UNKNOWN"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(string(tt.code), func(t *testing.T) {
			err := errors.New(tt.code, "test")
			if err.StatusCode != tt.statusCode {
				t.Errorf("Expected status %d for %s, got %d",
					tt.statusCode, tt.code, err.StatusCode)
			}
		})
	}
}

func TestIsAppError(t *testing.T) {
	appErr := errors.NewValidationError("test")
	stdErr := fmt.Errorf("standard error")

	if !errors.IsAppError(appErr) {
		t.Error("Should recognize AppError")
	}
	if errors.IsAppError(stdErr) {
		t.Error("Should not recognize standard error as AppError")
	}
}

func TestGetAppError(t *testing.T) {
	appErr := errors.NewValidationError("test")
	stdErr := fmt.Errorf("standard error")

	got, ok := errors.GetAppError(appErr)
	if !ok || got == nil {
		t.Error("Should convert to AppError")
	}

	got, ok = errors.GetAppError(stdErr)
	if ok {
		t.Error("Should not convert standard error")
	}
	if got != nil {
		t.Error("Should return nil for standard error")
	}
}

func TestAppError_Unwrap(t *testing.T) {
	internal := fmt.Errorf("internal error")
	err := errors.Wrap(internal, errors.ErrorCodeInternal, "wrapped")

	unwrapped := err.Unwrap()
	if !stderrors.Is(unwrapped, internal) {
		t.Error("Unwrap should return internal error")
	}

	errNoInternal := errors.New(errors.ErrorCodeValidation, "test")
	if errNoInternal.Unwrap() != nil {
		t.Error("Unwrap should return nil when no internal error")
	}
}
