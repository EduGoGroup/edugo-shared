package bootstrap

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	amqp "github.com/rabbitmq/amqp091-go"
)

// =============================================================================
// MESSAGE PUBLISHER IMPLEMENTATION
// =============================================================================

// defaultMessagePublisher es una implementación básica de MessagePublisher
type defaultMessagePublisher struct {
	channel *amqp.Channel
	factory RabbitMQFactory
}

// Publish publica un mensaje en una cola
func (p *defaultMessagePublisher) Publish(ctx context.Context, queueName string, body []byte) error {
	return p.PublishWithPriority(ctx, queueName, body, 0)
}

// PublishWithPriority publica un mensaje con prioridad
func (p *defaultMessagePublisher) PublishWithPriority(ctx context.Context, queueName string, body []byte, priority uint8) error {
	// Asegurar que la cola existe
	_, err := p.factory.DeclareQueue(p.channel, queueName)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Publicar mensaje
	err = p.channel.PublishWithContext(
		ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Priority:     priority,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

// Close cierra el publicador
func (p *defaultMessagePublisher) Close() error {
	if p.channel != nil {
		return p.channel.Close()
	}
	return nil
}

// =============================================================================
// STORAGE CLIENT IMPLEMENTATION
// =============================================================================

// defaultStorageClient es una implementación básica de StorageClient
type defaultStorageClient struct {
	client        *s3.Client
	presignClient interface{}
	bucket        string
}

// Upload sube un archivo al storage
func (c *defaultStorageClient) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	_, err := c.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object: %w", err)
	}

	// Retornar URL del objeto
	url := fmt.Sprintf("s3://%s/%s", c.bucket, key)
	return url, nil
}

// Download descarga un archivo del storage
func (c *defaultStorageClient) Download(ctx context.Context, key string) ([]byte, error) {
	result, err := c.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to download object: %w", err)
	}
	defer result.Body.Close()

	// Leer contenido completo
	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %w", err)
	}

	return data, nil
}

// Delete elimina un archivo del storage
func (c *defaultStorageClient) Delete(ctx context.Context, key string) error {
	_, err := c.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

// GetPresignedURL genera una URL pre-firmada para acceso temporal
func (c *defaultStorageClient) GetPresignedURL(ctx context.Context, key string, expirationMinutes int) (string, error) {
	// TODO: Implementar con presign client
	return "", fmt.Errorf("presigned URL not implemented yet")
}

// Exists verifica si un archivo existe
func (c *defaultStorageClient) Exists(ctx context.Context, key string) (bool, error) {
	_, err := c.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		// Si es error "not found", retornar false sin error
		return false, nil
	}

	return true, nil
}
