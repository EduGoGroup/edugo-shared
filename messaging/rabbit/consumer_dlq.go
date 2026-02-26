package rabbit

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumerDLQ implementation of Consumer with Dead Letter Queue support
//
//nolint:revive // Name maintained for API compatibility
type ConsumerDLQ struct {
	conn   *Connection
	config ConsumerConfig
	// Inner consumer logic is shared with RabbitMQConsumer, but we implement specific
	// Consume method to handle DLQ setup
	base *RabbitMQConsumer
}

// NewConsumerDLQ creates a new Consumer with DLQ support
func NewConsumerDLQ(conn *Connection, config ConsumerConfig) Consumer {
	base := NewConsumer(conn, config).(*RabbitMQConsumer)
	return &ConsumerDLQ{
		conn:   conn,
		config: config,
		base:   base,
	}
}

// Consume starts consuming messages from a queue with DLQ support
func (c *ConsumerDLQ) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
	c.base.mu.Lock()
	if c.base.running {
		c.base.mu.Unlock()
		return fmt.Errorf("consumer already running")
	}
	c.base.running = true
	c.base.mu.Unlock()

	// Get channel from connection
	ch := c.conn.GetChannel()

	// Setup DLQ infrastructure if enabled
	if c.config.DLQ.Enabled {
		if err := c.setupDLQ(ch); err != nil {
			c.base.mu.Lock()
			c.base.running = false
			c.base.mu.Unlock()
			return fmt.Errorf("failed to setup DLQ: %w", err)
		}
	}

	// Start consuming
	msgs, err := ch.Consume(
		queueName,
		c.config.Name,
		c.config.AutoAck,
		c.config.Exclusive,
		c.config.NoLocal,
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		c.base.mu.Lock()
		c.base.running = false
		c.base.mu.Unlock()
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	c.base.wg.Add(1)
	go func() {
		defer c.base.wg.Done()
		defer func() {
			c.base.mu.Lock()
			c.base.running = false
			c.base.mu.Unlock()
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.base.stopCh:
				return
			case msg, ok := <-msgs:
				if !ok {
					select {
					case c.base.errChan <- fmt.Errorf("message channel closed unexpectedly"):
					default:
					}
					return
				}

				// Create delivery wrapper
				delivery := amqpDelivery{
					Body:        msg.Body,
					DeliveryTag: msg.DeliveryTag,
					Ack:         msg.Ack,
					Nack:        msg.Nack,
				}

				c.processMessage(ctx, ch, queueName, msg, delivery, handler)
			}
		}
	}()

	return nil
}

// setupDLQ configures the Dead Letter Exchange and Queue
// Updated to accept ChannelInterface
func (c *ConsumerDLQ) setupDLQ(ch ChannelInterface) error {
	// Declare DLX
	err := ch.ExchangeDeclare(
		c.config.DLQ.DLXExchange,
		"direct", // DLX is usually direct
		true,     // durable
		false,    // auto-delete
		false,    // internal
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLX: %w", err)
	}

	// Declare DLQ
	dlqName := c.config.DLQ.DLXRoutingKey

	_, err = ch.QueueDeclare(
		dlqName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to declare DLQ: %w", err)
	}

	// Bind DLQ to DLX
	err = ch.QueueBind(
		dlqName,
		c.config.DLQ.DLXRoutingKey,
		c.config.DLQ.DLXExchange,
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to bind DLQ: %w", err)
	}

	return nil
}

// processMessage handles message processing with DLQ logic
// Updated to accept ChannelInterface
func (c *ConsumerDLQ) processMessage(ctx context.Context, ch ChannelInterface, queueName string, msg amqp.Delivery, delivery amqpDelivery, handler MessageHandler) {
	// Process message
	err := handler(ctx, delivery.Body)

	if err == nil {
		if !c.config.AutoAck {
			_ = delivery.Ack(false)
		}
		return
	}

	// Error handling with DLQ
	if !c.config.AutoAck {
		// Get retry count from headers
		retries := 0
		if xDeath, ok := msg.Headers["x-death"].([]interface{}); ok && len(xDeath) > 0 {
			if deathMap, ok := xDeath[0].(amqp.Table); ok {
				if count, ok := deathMap["count"].(int64); ok {
					retries = int(count)
				}
			}
		}

		if retries >= c.config.DLQ.MaxRetries {
			// Max retries exceeded, send to DLQ (Reject with requeue=false)
			_ = delivery.Nack(false, false) // requeue = false -> DLX
		} else {
			// Retry mechanism
			// Ideally we would publish to a retry queue with TTL, but simple Nack(requeue=true)
			// might cause busy loop.
			// With DLQ config, we might want to implement delayed retry.
			// For this implementation, we'll assume basic Nack(requeue=false) sends to DLX immediately
			// if configured on the queue. If we want exponential backoff, we need a separate retry architecture.
			// Here we just Nack without requeue to send to DLX, assuming manual intervention or separate DLQ processor.
			_ = delivery.Nack(false, false)
		}
	}
}

// Delegate other methods to base
func (c *ConsumerDLQ) Wait() error { return c.base.Wait() }
func (c *ConsumerDLQ) Stop()       { c.base.Stop() }
func (c *ConsumerDLQ) Errors() <-chan error { return c.base.Errors() }
func (c *ConsumerDLQ) IsRunning() bool { return c.base.IsRunning() }
func (c *ConsumerDLQ) Close() error { return c.base.Close() }
