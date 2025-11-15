package bootstrap

import (
	"testing"
)

func TestDefaultBootstrapOptions(t *testing.T) {
	opts := DefaultBootstrapOptions()

	if len(opts.RequiredResources) != 1 || opts.RequiredResources[0] != "logger" {
		t.Errorf("RequiredResources = %v, want [logger]", opts.RequiredResources)
	}

	if len(opts.OptionalResources) != 0 {
		t.Errorf("OptionalResources = %v, want []", opts.OptionalResources)
	}

	if opts.SkipHealthCheck {
		t.Error("SkipHealthCheck should be false by default")
	}

	if opts.MockFactories != nil {
		t.Error("MockFactories should be nil by default")
	}

	if !opts.StopOnFirstError {
		t.Error("StopOnFirstError should be true by default")
	}
}

func TestWithRequiredResources(t *testing.T) {
	opts := &BootstrapOptions{}
	WithRequiredResources("postgres", "mongodb", "rabbitmq")(opts)

	if len(opts.RequiredResources) != 3 {
		t.Errorf("Expected 3 required resources, got %d", len(opts.RequiredResources))
	}

	expected := []string{"postgres", "mongodb", "rabbitmq"}
	for i, res := range expected {
		if opts.RequiredResources[i] != res {
			t.Errorf("RequiredResources[%d] = %s, want %s", i, opts.RequiredResources[i], res)
		}
	}
}

func TestWithOptionalResources(t *testing.T) {
	opts := &BootstrapOptions{}
	WithOptionalResources("s3", "redis")(opts)

	if len(opts.OptionalResources) != 2 {
		t.Errorf("Expected 2 optional resources, got %d", len(opts.OptionalResources))
	}

	expected := []string{"s3", "redis"}
	for i, res := range expected {
		if opts.OptionalResources[i] != res {
			t.Errorf("OptionalResources[%d] = %s, want %s", i, opts.OptionalResources[i], res)
		}
	}
}

func TestWithSkipHealthCheck(t *testing.T) {
	opts := &BootstrapOptions{
		SkipHealthCheck: false,
	}

	WithSkipHealthCheck()(opts)

	if !opts.SkipHealthCheck {
		t.Error("SkipHealthCheck should be true after applying WithSkipHealthCheck")
	}
}

func TestWithMockFactories(t *testing.T) {
	mockFactories := &MockFactories{
		Logger: &mockLoggerFactory{},
	}

	opts := &BootstrapOptions{}
	WithMockFactories(mockFactories)(opts)

	if opts.MockFactories == nil {
		t.Fatal("MockFactories should not be nil")
	}

	if opts.MockFactories.Logger == nil {
		t.Error("MockFactories.Logger should not be nil")
	}
}

func TestWithStopOnFirstError(t *testing.T) {
	t.Run("set_to_true", func(t *testing.T) {
		opts := &BootstrapOptions{
			StopOnFirstError: false,
		}

		WithStopOnFirstError(true)(opts)

		if !opts.StopOnFirstError {
			t.Error("StopOnFirstError should be true")
		}
	})

	t.Run("set_to_false", func(t *testing.T) {
		opts := &BootstrapOptions{
			StopOnFirstError: true,
		}

		WithStopOnFirstError(false)(opts)

		if opts.StopOnFirstError {
			t.Error("StopOnFirstError should be false")
		}
	})
}

func TestApplyOptions(t *testing.T) {
	opts := DefaultBootstrapOptions()

	ApplyOptions(opts,
		WithRequiredResources("postgres", "mongodb"),
		WithOptionalResources("s3"),
		WithSkipHealthCheck(),
		WithStopOnFirstError(false),
	)

	if len(opts.RequiredResources) != 2 {
		t.Errorf("Expected 2 required resources, got %d", len(opts.RequiredResources))
	}

	if len(opts.OptionalResources) != 1 {
		t.Errorf("Expected 1 optional resource, got %d", len(opts.OptionalResources))
	}

	if !opts.SkipHealthCheck {
		t.Error("SkipHealthCheck should be true")
	}

	if opts.StopOnFirstError {
		t.Error("StopOnFirstError should be false")
	}
}

func TestApplyOptions_Empty(t *testing.T) {
	opts := DefaultBootstrapOptions()
	original := *opts

	ApplyOptions(opts)

	// Should remain unchanged
	if len(opts.RequiredResources) != len(original.RequiredResources) {
		t.Error("RequiredResources should not change with empty options")
	}
}

func TestApplyOptions_MultipleApplications(t *testing.T) {
	opts := &BootstrapOptions{}

	// Apply first set of options
	ApplyOptions(opts,
		WithRequiredResources("postgres"),
		WithStopOnFirstError(true),
	)

	// Apply second set of options (should override)
	ApplyOptions(opts,
		WithRequiredResources("mongodb", "rabbitmq"),
		WithStopOnFirstError(false),
	)

	// Last applied options should win
	if len(opts.RequiredResources) != 2 {
		t.Errorf("Expected 2 required resources, got %d", len(opts.RequiredResources))
	}

	if opts.StopOnFirstError {
		t.Error("StopOnFirstError should be false (last applied)")
	}
}

func TestBootstrapOption_Chaining(t *testing.T) {
	opts := &BootstrapOptions{}

	// Test that options can be chained
	option1 := WithRequiredResources("postgres")
	option2 := WithOptionalResources("s3")
	option3 := WithSkipHealthCheck()

	option1(opts)
	option2(opts)
	option3(opts)

	if len(opts.RequiredResources) != 1 || opts.RequiredResources[0] != "postgres" {
		t.Error("Required resources not set correctly")
	}

	if len(opts.OptionalResources) != 1 || opts.OptionalResources[0] != "s3" {
		t.Error("Optional resources not set correctly")
	}

	if !opts.SkipHealthCheck {
		t.Error("SkipHealthCheck not set correctly")
	}
}

func TestMockFactories_AllFields(t *testing.T) {
	mocks := &MockFactories{
		Logger:     &mockLoggerFactory{},
		PostgreSQL: &mockPostgreSQLFactory{},
		MongoDB:    &mockMongoDBFactory{},
		RabbitMQ:   &mockRabbitMQFactory{},
		S3:         &mockS3Factory{},
	}

	if mocks.Logger == nil {
		t.Error("Logger factory should not be nil")
	}
	if mocks.PostgreSQL == nil {
		t.Error("PostgreSQL factory should not be nil")
	}
	if mocks.MongoDB == nil {
		t.Error("MongoDB factory should not be nil")
	}
	if mocks.RabbitMQ == nil {
		t.Error("RabbitMQ factory should not be nil")
	}
	if mocks.S3 == nil {
		t.Error("S3 factory should not be nil")
	}
}
