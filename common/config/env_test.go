package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetEnv(t *testing.T) {
	key := "TEST_GET_ENV"

	// Default value
	assert.Equal(t, "default", GetEnv(key, "default"))

	// Value set
	require.NoError(t, os.Setenv(key, "set_value"))
	defer func() { require.NoError(t, os.Unsetenv(key)) }()
	assert.Equal(t, "set_value", GetEnv(key, "default"))
}

func TestGetEnvRequired(t *testing.T) {
	key := "TEST_GET_ENV_REQUIRED"

	// Panic when not set
	assert.Panics(t, func() { GetEnvRequired(key) })

	// Return value when set
	require.NoError(t, os.Setenv(key, "required_value"))
	defer func() { require.NoError(t, os.Unsetenv(key)) }()
	assert.Equal(t, "required_value", GetEnvRequired(key))
}

func TestGetEnvInt(t *testing.T) {
	key := "TEST_GET_ENV_INT"

	// Default value when not set
	assert.Equal(t, 42, GetEnvInt(key, 42))

	// Default value when not parsable
	require.NoError(t, os.Setenv(key, "not_an_int"))
	assert.Equal(t, 42, GetEnvInt(key, 42))

	// Valid int
	require.NoError(t, os.Setenv(key, "100"))
	defer func() { require.NoError(t, os.Unsetenv(key)) }()
	assert.Equal(t, 100, GetEnvInt(key, 42))
}

func TestGetEnvBool(t *testing.T) {
	key := "TEST_GET_ENV_BOOL"

	// Default value when not set
	assert.Equal(t, true, GetEnvBool(key, true))

	// Default value when not parsable
	require.NoError(t, os.Setenv(key, "not_a_bool"))
	assert.Equal(t, true, GetEnvBool(key, true))

	// Valid bool (true)
	require.NoError(t, os.Setenv(key, "true"))
	assert.Equal(t, true, GetEnvBool(key, false))

	// Valid bool (false)
	require.NoError(t, os.Setenv(key, "false"))
	defer func() { require.NoError(t, os.Unsetenv(key)) }()
	assert.Equal(t, false, GetEnvBool(key, true))
}

func TestGetEnvDuration(t *testing.T) {
	key := "TEST_GET_ENV_DURATION"
	defaultDuration := 5 * time.Second

	// Default value when not set
	assert.Equal(t, defaultDuration, GetEnvDuration(key, defaultDuration))

	// Default value when not parsable
	require.NoError(t, os.Setenv(key, "not_a_duration"))
	assert.Equal(t, defaultDuration, GetEnvDuration(key, defaultDuration))

	// Valid duration
	require.NoError(t, os.Setenv(key, "10s"))
	defer func() { require.NoError(t, os.Unsetenv(key)) }()
	assert.Equal(t, 10*time.Second, GetEnvDuration(key, defaultDuration))
}

func TestMustGetEnv(t *testing.T) {
	key := "TEST_MUST_GET_ENV"

	// Panic when not set
	assert.Panics(t, func() { MustGetEnv(key) })

	// Return value when set
	require.NoError(t, os.Setenv(key, "required_value"))
	defer func() { require.NoError(t, os.Unsetenv(key)) }()
	assert.Equal(t, "required_value", MustGetEnv(key))
}

func TestSetEnvAndUnsetEnv(t *testing.T) {
	key := "TEST_SET_UNSET_ENV"
	value := "some_value"

	// Set
	err := SetEnv(key, value)
	assert.NoError(t, err)
	assert.Equal(t, value, os.Getenv(key))

	// Unset
	err = UnsetEnv(key)
	assert.NoError(t, err)
	assert.Equal(t, "", os.Getenv(key))
}

func TestLookupEnv(t *testing.T) {
	key := "TEST_LOOKUP_ENV"

	// Not set
	val, exists := LookupEnv(key)
	assert.False(t, exists)
	assert.Equal(t, "", val)

	// Set
	require.NoError(t, os.Setenv(key, "lookup_value"))
	defer func() { require.NoError(t, os.Unsetenv(key)) }()
	val, exists = LookupEnv(key)
	assert.True(t, exists)
	assert.Equal(t, "lookup_value", val)
}

func TestGetEnvironment(t *testing.T) {
	key := "APP_ENV"

	// Default
	assert.Equal(t, "development", GetEnvironment())

	// Custom
	require.NoError(t, os.Setenv(key, "custom_env"))
	defer func() { require.NoError(t, os.Unsetenv(key)) }()
	assert.Equal(t, "custom_env", GetEnvironment())
}

func TestIsDevelopment(t *testing.T) {
	key := "APP_ENV"
	defer func() { require.NoError(t, os.Unsetenv(key)) }()

	require.NoError(t, os.Setenv(key, "development"))
	assert.True(t, IsDevelopment())

	require.NoError(t, os.Setenv(key, "dev"))
	assert.True(t, IsDevelopment())

	require.NoError(t, os.Setenv(key, "local"))
	assert.True(t, IsDevelopment())

	require.NoError(t, os.Setenv(key, "production"))
	assert.False(t, IsDevelopment())
}

func TestIsProduction(t *testing.T) {
	key := "APP_ENV"
	defer func() { require.NoError(t, os.Unsetenv(key)) }()

	require.NoError(t, os.Setenv(key, "production"))
	assert.True(t, IsProduction())

	require.NoError(t, os.Setenv(key, "prod"))
	assert.True(t, IsProduction())

	require.NoError(t, os.Setenv(key, "development"))
	assert.False(t, IsProduction())
}

func TestIsStaging(t *testing.T) {
	key := "APP_ENV"
	defer func() { require.NoError(t, os.Unsetenv(key)) }()

	require.NoError(t, os.Setenv(key, "staging"))
	assert.True(t, IsStaging())

	require.NoError(t, os.Setenv(key, "stage"))
	assert.True(t, IsStaging())

	require.NoError(t, os.Setenv(key, "qa"))
	assert.True(t, IsStaging())

	require.NoError(t, os.Setenv(key, "production"))
	assert.False(t, IsStaging())
}
