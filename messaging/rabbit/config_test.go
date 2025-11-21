package rabbit

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, "amqp://guest:guest@localhost:5672/", config.URL)
	assert.Equal(t, "default_exchange", config.Exchange.Name)
	assert.Equal(t, "topic", config.Exchange.Type)
	assert.True(t, config.Exchange.Durable)
	assert.False(t, config.Exchange.AutoDelete)
	assert.Equal(t, "default_queue", config.Queue.Name)
	assert.True(t, config.Queue.Durable)
	assert.False(t, config.Queue.AutoDelete)
	assert.False(t, config.Queue.Exclusive)
	assert.NotNil(t, config.Queue.Args)
	assert.Equal(t, "default_consumer", config.Consumer.Name)
	assert.False(t, config.Consumer.AutoAck)
	assert.False(t, config.Consumer.Exclusive)
	assert.False(t, config.Consumer.NoLocal)
	assert.Equal(t, DefaultPrefetchCount, config.PrefetchCount)
}

func TestConfig_Fields(t *testing.T) {
	config := Config{
		URL: "amqp://user:pass@example.com:5672/test",
		Exchange: ExchangeConfig{
			Name:       "test_exchange",
			Type:       "direct",
			Durable:    false,
			AutoDelete: true,
		},
		Queue: QueueConfig{
			Name:       "test_queue",
			Durable:    false,
			AutoDelete: true,
			Exclusive:  true,
			Args: map[string]interface{}{
				"x-message-ttl": 60000,
			},
		},
		Consumer: ConsumerConfig{
			Name:          "test_consumer",
			AutoAck:       true,
			Exclusive:     true,
			NoLocal:       true,
			PrefetchCount: 10,
		},
		PrefetchCount: 20,
	}

	assert.Equal(t, "amqp://user:pass@example.com:5672/test", config.URL)
	assert.Equal(t, "test_exchange", config.Exchange.Name)
	assert.Equal(t, "direct", config.Exchange.Type)
	assert.False(t, config.Exchange.Durable)
	assert.True(t, config.Exchange.AutoDelete)
	assert.Equal(t, "test_queue", config.Queue.Name)
	assert.True(t, config.Queue.Exclusive)
	assert.Equal(t, 60000, config.Queue.Args["x-message-ttl"])
	assert.Equal(t, "test_consumer", config.Consumer.Name)
	assert.True(t, config.Consumer.AutoAck)
	assert.Equal(t, 10, config.Consumer.PrefetchCount)
	assert.Equal(t, 20, config.PrefetchCount)
}

func TestExchangeConfig_Types(t *testing.T) {
	exchangeTypes := []string{"direct", "topic", "fanout", "headers"}

	for _, exchangeType := range exchangeTypes {
		t.Run(exchangeType, func(t *testing.T) {
			config := ExchangeConfig{
				Name:       "test_" + exchangeType,
				Type:       exchangeType,
				Durable:    true,
				AutoDelete: false,
			}

			assert.Equal(t, "test_"+exchangeType, config.Name)
			assert.Equal(t, exchangeType, config.Type)
			assert.True(t, config.Durable)
			assert.False(t, config.AutoDelete)
		})
	}
}

func TestQueueConfig_WithArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     map[string]interface{}
		validate func(t *testing.T, args map[string]interface{})
	}{
		{
			name: "with TTL",
			args: map[string]interface{}{
				"x-message-ttl": 60000,
			},
			validate: func(t *testing.T, args map[string]interface{}) {
				assert.Equal(t, 60000, args["x-message-ttl"])
			},
		},
		{
			name: "with max length",
			args: map[string]interface{}{
				"x-max-length": 1000,
			},
			validate: func(t *testing.T, args map[string]interface{}) {
				assert.Equal(t, 1000, args["x-max-length"])
			},
		},
		{
			name: "with DLX",
			args: map[string]interface{}{
				"x-dead-letter-exchange":    "dlx_exchange",
				"x-dead-letter-routing-key": "dlx_routing",
			},
			validate: func(t *testing.T, args map[string]interface{}) {
				assert.Equal(t, "dlx_exchange", args["x-dead-letter-exchange"])
				assert.Equal(t, "dlx_routing", args["x-dead-letter-routing-key"])
			},
		},
		{
			name: "empty args",
			args: map[string]interface{}{},
			validate: func(t *testing.T, args map[string]interface{}) {
				assert.Empty(t, args)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := QueueConfig{
				Name:       "test_queue",
				Durable:    true,
				AutoDelete: false,
				Exclusive:  false,
				Args:       tt.args,
			}

			assert.NotNil(t, config.Args)
			tt.validate(t, config.Args)
		})
	}
}

func TestConsumerConfig_PrefetchCount(t *testing.T) {
	tests := []struct {
		name          string
		prefetchCount int
	}{
		{"default", 5},
		{"custom low", 1},
		{"custom high", 100},
		{"zero", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := ConsumerConfig{
				Name:          "test_consumer",
				AutoAck:       false,
				Exclusive:     false,
				NoLocal:       false,
				PrefetchCount: tt.prefetchCount,
			}

			assert.Equal(t, tt.prefetchCount, config.PrefetchCount)
		})
	}
}

func TestConfig_CustomPrefetchCount(t *testing.T) {
	tests := []struct {
		name          string
		prefetchCount int
	}{
		{"default", DefaultPrefetchCount},
		{"custom 1", 1},
		{"custom 10", 10},
		{"custom 100", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				URL:           "amqp://localhost:5672/",
				PrefetchCount: tt.prefetchCount,
			}

			assert.Equal(t, tt.prefetchCount, config.PrefetchCount)
		})
	}
}

func TestConsumerConfig_DLQConfiguration(t *testing.T) {
	dlqConfig := DLQConfig{
		Enabled:               true,
		MaxRetries:            5,
		RetryDelay:            10 * time.Second,
		DLXExchange:           "test.dlx",
		DLXRoutingKey:         "test.dlq",
		UseExponentialBackoff: true,
	}

	consumerConfig := ConsumerConfig{
		Name:    "test_consumer",
		AutoAck: false,
		DLQ:     dlqConfig,
	}

	assert.Equal(t, "test_consumer", consumerConfig.Name)
	assert.False(t, consumerConfig.AutoAck)
	assert.True(t, consumerConfig.DLQ.Enabled)
	assert.Equal(t, 5, consumerConfig.DLQ.MaxRetries)
	assert.Equal(t, 10*time.Second, consumerConfig.DLQ.RetryDelay)
	assert.Equal(t, "test.dlx", consumerConfig.DLQ.DLXExchange)
	assert.Equal(t, "test.dlq", consumerConfig.DLQ.DLXRoutingKey)
	assert.True(t, consumerConfig.DLQ.UseExponentialBackoff)
}
