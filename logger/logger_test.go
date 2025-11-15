package logger_test

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/EduGoGroup/edugo-shared/logger"
)

// Helper to capture logger output - creates logger AFTER redirecting stdout
func captureOutput(level, format string, f func(log logger.Logger)) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outCh := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outCh <- buf.String()
	}()

	// Create logger AFTER redirecting stdout
	log := logger.NewZapLogger(level, format)
	f(log)

	w.Close()
	os.Stdout = old
	return <-outCh
}

func TestNewZapLogger_DifferentLevels(t *testing.T) {
	levels := []string{"debug", "info", "warn", "error", "fatal", "invalid"}
	for _, level := range levels {
		t.Run(level, func(t *testing.T) {
			log := logger.NewZapLogger(level, "json")
			if log == nil {
				t.Error("NewZapLogger returned nil")
			}
		})
	}
}

func TestNewZapLogger_DifferentFormats(t *testing.T) {
	formats := []string{"json", "console"}
	for _, format := range formats {
		t.Run(format, func(t *testing.T) {
			log := logger.NewZapLogger("info", format)
			if log == nil {
				t.Error("NewZapLogger returned nil")
			}
		})
	}
}

func TestZapLogger_Debug(t *testing.T) {
	output := captureOutput("debug", "json", func(log logger.Logger) {
		log.Debug("debug message", "key", "value")
	})

	if !strings.Contains(output, "debug message") {
		t.Errorf("Debug message not found in: %s", output)
	}
}

func TestZapLogger_Info(t *testing.T) {
	output := captureOutput("info", "json", func(log logger.Logger) {
		log.Info("info message", "user_id", 123)
	})

	if !strings.Contains(output, "info message") {
		t.Errorf("Info message not found in: %s", output)
	}
}

func TestZapLogger_Warn(t *testing.T) {
	output := captureOutput("warn", "json", func(log logger.Logger) {
		log.Warn("warning message", "code", 404)
	})

	if !strings.Contains(output, "warning message") {
		t.Errorf("Warn message not found in: %s", output)
	}
}

func TestZapLogger_Error(t *testing.T) {
	output := captureOutput("error", "json", func(log logger.Logger) {
		log.Error("error message", "error_code", 500)
	})

	if !strings.Contains(output, "error message") {
		t.Errorf("Error message not found in: %s", output)
	}
}

func TestZapLogger_With(t *testing.T) {
	output := captureOutput("info", "json", func(log logger.Logger) {
		contextLog := log.With("request_id", "abc123")
		contextLog.Info("operation completed")
	})

	if !strings.Contains(output, "abc123") {
		t.Errorf("Context field not found in: %s", output)
	}
	if !strings.Contains(output, "operation completed") {
		t.Errorf("Message not found in: %s", output)
	}
}

func TestZapLogger_WithChaining(t *testing.T) {
	output := captureOutput("info", "json", func(log logger.Logger) {
		log.With("field1", "value1").With("field2", "value2").Info("chained")
	})

	if !strings.Contains(output, "value1") || !strings.Contains(output, "value2") {
		t.Errorf("Chained fields not found in: %s", output)
	}
}

func TestZapLogger_Sync(t *testing.T) {
	log := logger.NewZapLogger("info", "json")
	// Sync may return error on some systems for stdout, don't fail on error
	_ = log.Sync()
}

func TestZapLogger_JSONFormat(t *testing.T) {
	output := captureOutput("info", "json", func(log logger.Logger) {
		log.Info("json test", "key", "value")
	})

	// Should contain JSON-like structure
	if !strings.Contains(output, "\"message\"") || !strings.Contains(output, "\"level\"") {
		t.Errorf("JSON format not detected in: %s", output)
	}
}

func TestZapLogger_ConsoleFormat(t *testing.T) {
	output := captureOutput("info", "console", func(log logger.Logger) {
		log.Info("console test")
	})

	if !strings.Contains(output, "console test") {
		t.Errorf("Console message not found in: %s", output)
	}
}

func TestZapLogger_LevelFiltering(t *testing.T) {
	t.Run("info_filters_debug", func(t *testing.T) {
		output := captureOutput("info", "json", func(log logger.Logger) {
			log.Debug("should not appear")
			log.Info("should appear")
		})

		if strings.Contains(output, "should not appear") {
			t.Error("Debug message appeared when level is info")
		}
		if !strings.Contains(output, "should appear") {
			t.Errorf("Info message not found in: %s", output)
		}
	})

	t.Run("error_filters_lower_levels", func(t *testing.T) {
		output := captureOutput("error", "json", func(log logger.Logger) {
			log.Info("info should not appear")
			log.Warn("warn should not appear")
			log.Error("error should appear")
		})

		if strings.Contains(output, "info should not appear") {
			t.Error("Info appeared when level is error")
		}
		if strings.Contains(output, "warn should not appear") {
			t.Error("Warn appeared when level is error")
		}
		if !strings.Contains(output, "error should appear") {
			t.Errorf("Error message not found in: %s", output)
		}
	})
}

func TestZapLogger_MultipleFieldTypes(t *testing.T) {
	output := captureOutput("info", "json", func(log logger.Logger) {
		log.Info("multi-field",
			"string", "text",
			"int", 42,
			"bool", true,
			"float", 3.14,
		)
	})

	checks := []string{"text", "42", "true", "3.14"}
	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Field %s not found in: %s", check, output)
		}
	}
}

func TestZapLogger_AllLevelsOutput(t *testing.T) {
	output := captureOutput("debug", "json", func(log logger.Logger) {
		log.Debug("debug msg")
		log.Info("info msg")
		log.Warn("warn msg")
		log.Error("error msg")
	})

	messages := []string{"debug msg", "info msg", "warn msg", "error msg"}
	for _, msg := range messages {
		if !strings.Contains(output, msg) {
			t.Errorf("Message '%s' not found in output", msg)
		}
	}
}
