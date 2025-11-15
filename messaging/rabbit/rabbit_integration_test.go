package rabbit_test

import (
	"context"
	"testing"
	"time"

	"github.com/EduGoGroup/edugo-shared/messaging/rabbit"
	"github.com/EduGoGroup/edugo-shared/testing/containers"
	amqp "github.com/rabbitmq/amqp091-go"
)

func TestConnect_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := containers.NewConfig().
		WithRabbitMQ(nil).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	rabbitContainer := manager.RabbitMQ()
	if rabbitContainer == nil {
		t.Fatal("RabbitMQ container no disponible")
	}

	// Obtener connection string
	ctx := context.Background()
	connStr, err := rabbitContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo connection string: %v", err)
	}

	// Test de conexión
	conn, err := rabbit.Connect(connStr)
	if err != nil {
		t.Fatalf("Error conectando a RabbitMQ: %v", err)
	}

	if conn == nil {
		t.Fatal("Connection es nil")
	}

	// Verificar que podemos obtener el canal
	ch := conn.GetChannel()
	if ch == nil {
		t.Fatal("Channel es nil")
	}

	// Verificar que la conexión no está cerrada
	if conn.IsClosed() {
		t.Error("Connection está cerrada inmediatamente después de conectar")
	}
}

func TestPublisher_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := containers.NewConfig().
		WithRabbitMQ(nil).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	rabbitContainer := manager.RabbitMQ()
	ctx := context.Background()
	connStr, err := rabbitContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo connection string: %v", err)
	}

	conn, err := rabbit.Connect(connStr)
	if err != nil {
		t.Fatalf("Error conectando: %v", err)
	}

	publisher := rabbit.NewPublisher(conn)
	if publisher == nil {
		t.Fatal("Publisher es nil")
	}

	// Declarar cola temporal para el test
	queueName := "test_publisher_queue"
	queueCfg := rabbit.QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
		Exclusive:  false,
	}

	_, err = conn.DeclareQueue(queueCfg)
	if err != nil {
		t.Fatalf("Error declarando cola: %v", err)
	}

	// Publicar mensaje simple
	testMsg := map[string]string{"test": "message"}
	err = publisher.Publish(ctx, "", queueName, testMsg)
	if err != nil {
		t.Errorf("Error publicando mensaje: %v", err)
	}

	// Publicar con prioridad
	err = publisher.PublishWithPriority(ctx, "", queueName, testMsg, 5)
	if err != nil {
		t.Errorf("Error publicando mensaje con prioridad: %v", err)
	}
}

func TestConsumer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	config := containers.NewConfig().
		WithRabbitMQ(nil).
		Build()

	manager, err := containers.GetManager(t, config)
	if err != nil {
		t.Fatalf("Error creando manager: %v", err)
	}

	rabbitContainer := manager.RabbitMQ()
	ctx := context.Background()
	connStr, err := rabbitContainer.ConnectionString(ctx)
	if err != nil {
		t.Fatalf("Error obteniendo connection string: %v", err)
	}

	conn, err := rabbit.Connect(connStr)
	if err != nil {
		t.Fatalf("Error conectando: %v", err)
	}

	// Crear cola y publicar mensaje de prueba
	queueName := "test_consumer_queue"
	queueCfg := rabbit.QueueConfig{
		Name:       queueName,
		Durable:    false,
		AutoDelete: true,
	}

	_, err = conn.DeclareQueue(queueCfg)
	if err != nil {
		t.Fatalf("Error declarando cola: %v", err)
	}

	// Publicar mensaje usando el canal directo
	ch := conn.GetChannel()
	err = ch.PublishWithContext(ctx, "", queueName, false, false, amqp.Publishing{
		Body: []byte(`{"test":"message for consumer"}`),
	})
	if err != nil {
		t.Fatalf("Error publicando mensaje: %v", err)
	}

	// Crear consumer
	consumerCfg := rabbit.ConsumerConfig{
		Name:          "test-consumer",
		AutoAck:       true,
		PrefetchCount: 10,
	}

	consumer := rabbit.NewConsumer(conn, consumerCfg)
	if consumer == nil {
		t.Fatal("Consumer es nil")
	}

	// Consumir con timeout
	consumeCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	messageReceived := false
	handler := func(ctx context.Context, body []byte) error {
		if len(body) > 0 {
			messageReceived = true
			cancel() // Cancelar después de recibir mensaje
		}
		return nil
	}

	err = consumer.Consume(consumeCtx, queueName, handler)
	if err != nil {
		t.Fatalf("Error iniciando consumer: %v", err)
	}

	// Esperar un poco para recibir el mensaje
	time.Sleep(500 * time.Millisecond)

	if !messageReceived {
		t.Error("No se recibió el mensaje esperado")
	}
}
