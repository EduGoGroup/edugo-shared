package containers

import (
	"context"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestRabbitMQContainer_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test en modo short")
	}

	ctx := context.Background()
	cfg := &RabbitConfig{
		Image:    "rabbitmq:3.12-alpine",
		Username: "test_user",
		Password: "test_pass",
	}

	container, err := createRabbitMQ(ctx, cfg)
	if err != nil {
		t.Fatalf("Error creando container: %v", err)
	}
	defer container.Terminate(ctx)

	t.Run("ConnectionString", func(t *testing.T) {
		connStr, err := container.ConnectionString(ctx)
		if err != nil {
			t.Errorf("Error obteniendo connection string: %v", err)
		}
		if connStr == "" {
			t.Error("Connection string está vacío")
		}
	})

	t.Run("Connection_Access", func(t *testing.T) {
		conn := container.Connection()
		if conn == nil {
			t.Fatal("Connection no debería ser nil")
		}

		if conn.IsClosed() {
			t.Error("Connection debería estar abierta")
		}
	})

	t.Run("Channel_Creation", func(t *testing.T) {
		ch, err := container.Channel()
		if err != nil {
			t.Fatalf("Error creando canal: %v", err)
		}
		defer ch.Close()

		if ch == nil {
			t.Error("Channel no debería ser nil")
		}
	})

	t.Run("Queue_Declaration_And_Publish", func(t *testing.T) {
		ch, err := container.Channel()
		if err != nil {
			t.Fatalf("Error creando canal: %v", err)
		}
		defer ch.Close()

		// Declarar cola
		queueName := "test_queue"
		queue, err := ch.QueueDeclare(
			queueName,
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			t.Fatalf("Error declarando cola: %v", err)
		}

		// Publicar mensaje
		err = ch.Publish(
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte("test message"),
			},
		)
		if err != nil {
			t.Errorf("Error publicando mensaje: %v", err)
		}

		// Verificar que el mensaje está en la cola
		queue, err = ch.QueueInspect(queueName)
		if err != nil {
			t.Fatalf("Error inspeccionando cola: %v", err)
		}
		if queue.Messages != 1 {
			t.Errorf("Esperado 1 mensaje en cola, obtenido %d", queue.Messages)
		}
	})

	t.Run("PurgeQueue", func(t *testing.T) {
		ch, err := container.Channel()
		if err != nil {
			t.Fatalf("Error creando canal: %v", err)
		}
		defer ch.Close()

		// Declarar cola y publicar mensajes
		queueName := "test_purge_queue"
		queue, _ := ch.QueueDeclare(queueName, false, false, false, false, nil)

		for i := 0; i < 5; i++ {
			ch.Publish("", queue.Name, false, false, amqp.Publishing{
				Body: []byte("message"),
			})
		}

		// Verificar que hay mensajes
		queue, _ = ch.QueueInspect(queueName)
		if queue.Messages != 5 {
			t.Errorf("Esperado 5 mensajes antes de purge, obtenido %d", queue.Messages)
		}

		// Purgar cola usando el método del container
		err = container.PurgeQueue(queueName)
		if err != nil {
			t.Errorf("Error purgando cola: %v", err)
		}

		// Verificar que la cola está vacía
		queue, _ = ch.QueueInspect(queueName)
		if queue.Messages != 0 {
			t.Errorf("Esperado 0 mensajes después de purge, obtenido %d", queue.Messages)
		}
	})

	t.Run("DeleteQueue", func(t *testing.T) {
		ch, err := container.Channel()
		if err != nil {
			t.Fatalf("Error creando canal: %v", err)
		}
		defer ch.Close()

		// Declarar cola
		queueName := "test_delete_queue"
		ch.QueueDeclare(queueName, false, false, false, false, nil)

		// Eliminar cola usando el método del container
		err = container.DeleteQueue(queueName)
		if err != nil {
			t.Errorf("Error eliminando cola: %v", err)
		}

		// Verificar que la cola no existe (debería dar error al inspeccionar)
		_, err = ch.QueueInspect(queueName)
		if err == nil {
			t.Error("QueueInspect debería dar error para cola eliminada")
		}
	})
}

func TestCreateRabbitMQ_NilConfig(t *testing.T) {
	ctx := context.Background()
	_, err := createRabbitMQ(ctx, nil)
	if err == nil {
		t.Error("createRabbitMQ con config nil debería dar error")
	}
}
