# Features Faltantes en Módulos Existentes

## 🎯 Objetivo

Detallar las **features que faltan** en módulos que **ya existen** en edugo-shared, especificando qué agregar, cómo implementar, y tiempo estimado.

---

## 📦 messaging/rabbit/ (v0.5.0 → v0.6.0)

### Feature: Dead Letter Queue (DLQ) Support

#### Requerido Por
- ✅ **Worker** (CRÍTICO)
- ⚠️ api-mobile (opcional para MVP)

#### Estado Actual (v0.5.0)

**Archivo actual:** `consumer.go`

```go
// consumer.go (v0.5.0 - ACTUAL)
type Consumer struct {
    connection *amqp.Connection
    channel    *amqp.Channel
    queue      string
}

func (c *Consumer) Consume(handler func([]byte) error) error {
    msgs, err := c.channel.Consume(
        c.queue,
        "",    // consumer tag
        false, // auto-ack = false (manual)
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        return err
    }

    for msg := range msgs {
        if err := handler(msg.Body); err != nil {
            // ❌ PROBLEMA: Si falla, solo hace Nack con requeue=true
            // Esto reencola infinitamente si el mensaje siempre falla
            msg.Nack(false, true)
            continue
        }
        msg.Ack(false)
    }
    return nil
}
```

**Problema:**
1. No hay límite de reintentos
2. Mensajes con errores permanentes se reencolan infinitamente
3. No hay visibilidad de mensajes fallidos
4. Worker crashea si procesa mensaje malo repetidamente

---

#### Implementación Necesaria (v0.6.0)

**Archivo nuevo/modificado:** `dlq.go`

```go
// dlq.go (NUEVO en v0.6.0)
package rabbit

import (
    "fmt"
    "time"
)

// DLQConfig configura el Dead Letter Queue
type DLQConfig struct {
    Enabled          bool
    MaxRetries       int           // Default: 3
    RetryDelay       time.Duration // Default: 5s
    DLXExchange      string        // Dead Letter Exchange
    DLXRoutingKey    string        // Routing key para DLQ
    UseExponentialBackoff bool     // Default: true
}

// DefaultDLQConfig retorna configuración por defecto
func DefaultDLQConfig() DLQConfig {
    return DLQConfig{
        Enabled:               true,
        MaxRetries:            3,
        RetryDelay:            5 * time.Second,
        DLXExchange:           "dlx",
        DLXRoutingKey:         "dlq",
        UseExponentialBackoff: true,
    }
}

// calculateBackoff calcula el delay con exponential backoff
func (c *DLQConfig) calculateBackoff(attempt int) time.Duration {
    if !c.UseExponentialBackoff {
        return c.RetryDelay
    }
    // Exponential: 5s, 10s, 20s, 40s...
    return c.RetryDelay * time.Duration(1<<uint(attempt))
}
```

**Archivo modificado:** `consumer.go`

```go
// consumer.go (v0.6.0 - MODIFICADO)
package rabbit

import (
    "encoding/json"
    "fmt"
)

type ConsumerConfig struct {
    Queue            string
    ConsumerTag      string
    AutoAck          bool
    Exclusive        bool
    PrefetchCount    int
    DLQ              DLQConfig // ← NUEVO
}

type Consumer struct {
    connection *amqp.Connection
    channel    *amqp.Channel
    config     ConsumerConfig
}

// ConsumeWithDLQ consume mensajes con soporte para Dead Letter Queue
func (c *Consumer) ConsumeWithDLQ(handler func([]byte) error) error {
    // Configurar prefetch
    if err := c.channel.Qos(c.config.PrefetchCount, 0, false); err != nil {
        return fmt.Errorf("failed to set QoS: %w", err)
    }

    // Declarar DLX y DLQ si está habilitado
    if c.config.DLQ.Enabled {
        if err := c.setupDLQ(); err != nil {
            return fmt.Errorf("failed to setup DLQ: %w", err)
        }
    }

    msgs, err := c.channel.Consume(
        c.config.Queue,
        c.config.ConsumerTag,
        c.config.AutoAck,
        c.config.Exclusive,
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        return err
    }

    for msg := range msgs {
        // Obtener número de reintentos del header
        retries := getRetryCount(msg.Headers)

        // Procesar mensaje
        if err := handler(msg.Body); err != nil {
            logger.Error("Error processing message", map[string]interface{}{
                "error":   err.Error(),
                "retries": retries,
            })

            // Verificar si excedió reintentos
            if c.config.DLQ.Enabled && retries >= c.config.DLQ.MaxRetries {
                // Enviar a DLQ
                if err := c.sendToDLQ(msg); err != nil {
                    logger.Error("Failed to send to DLQ", map[string]interface{}{
                        "error": err.Error(),
                    })
                    msg.Nack(false, true) // Requeue como fallback
                } else {
                    msg.Ack(false) // Acknowledge porque ya está en DLQ
                }
            } else {
                // Reencolar con delay
                backoff := c.config.DLQ.calculateBackoff(retries)
                time.Sleep(backoff)

                // Incrementar contador de reintentos
                if msg.Headers == nil {
                    msg.Headers = amqp.Table{}
                }
                msg.Headers["x-retry-count"] = retries + 1

                msg.Nack(false, true) // Requeue
            }
        } else {
            // Procesado exitosamente
            msg.Ack(false)
        }
    }

    return nil
}

// setupDLQ configura el Dead Letter Exchange y Queue
func (c *Consumer) setupDLQ() error {
    // Declarar DLX (exchange para mensajes fallidos)
    if err := c.channel.ExchangeDeclare(
        c.config.DLQ.DLXExchange, // name
        "direct",                  // type
        true,                      // durable
        false,                     // auto-deleted
        false,                     // internal
        false,                     // no-wait
        nil,                       // arguments
    ); err != nil {
        return fmt.Errorf("failed to declare DLX: %w", err)
    }

    // Declarar DLQ (queue para mensajes fallidos)
    _, err := c.channel.QueueDeclare(
        c.config.DLQ.DLXRoutingKey, // name (usa routing key como nombre)
        true,                        // durable
        false,                       // delete when unused
        false,                       // exclusive
        false,                       // no-wait
        nil,                         // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare DLQ: %w", err)
    }

    // Bindear DLQ al DLX
    if err := c.channel.QueueBind(
        c.config.DLQ.DLXRoutingKey, // queue name
        c.config.DLQ.DLXRoutingKey, // routing key
        c.config.DLQ.DLXExchange,   // exchange
        false,                      // no-wait
        nil,                        // arguments
    ); err != nil {
        return fmt.Errorf("failed to bind DLQ: %w", err)
    }

    return nil
}

// sendToDLQ envía un mensaje al Dead Letter Queue
func (c *Consumer) sendToDLQ(msg amqp.Delivery) error {
    // Agregar metadata al mensaje
    headers := msg.Headers
    if headers == nil {
        headers = amqp.Table{}
    }
    headers["x-original-exchange"] = msg.Exchange
    headers["x-original-routing-key"] = msg.RoutingKey
    headers["x-failed-at"] = time.Now().Unix()
    headers["x-retry-count"] = getRetryCount(msg.Headers)

    // Publicar a DLX
    return c.channel.Publish(
        c.config.DLQ.DLXExchange,   // exchange
        c.config.DLQ.DLXRoutingKey, // routing key
        false,                      // mandatory
        false,                      // immediate
        amqp.Publishing{
            ContentType: msg.ContentType,
            Body:        msg.Body,
            Headers:     headers,
        },
    )
}

// getRetryCount extrae el número de reintentos del header
func getRetryCount(headers amqp.Table) int {
    if headers == nil {
        return 0
    }
    if count, ok := headers["x-retry-count"].(int); ok {
        return count
    }
    return 0
}
```

---

#### Tests Requeridos

**Archivo:** `dlq_test.go`

```go
package rabbit_test

import (
    "testing"
    "time"
    "github.com/EduGoGroup/edugo-shared/messaging/rabbit"
)

func TestDLQConfig_CalculateBackoff(t *testing.T) {
    config := rabbit.DLQConfig{
        RetryDelay:            5 * time.Second,
        UseExponentialBackoff: true,
    }

    tests := []struct {
        attempt int
        want    time.Duration
    }{
        {0, 5 * time.Second},  // 5s * 2^0 = 5s
        {1, 10 * time.Second}, // 5s * 2^1 = 10s
        {2, 20 * time.Second}, // 5s * 2^2 = 20s
        {3, 40 * time.Second}, // 5s * 2^3 = 40s
    }

    for _, tt := range tests {
        t.Run(fmt.Sprintf("attempt_%d", tt.attempt), func(t *testing.T) {
            got := config.calculateBackoff(tt.attempt)
            if got != tt.want {
                t.Errorf("calculateBackoff(%d) = %v, want %v", tt.attempt, got, tt.want)
            }
        })
    }
}

func TestDLQConfig_LinearBackoff(t *testing.T) {
    config := rabbit.DLQConfig{
        RetryDelay:            5 * time.Second,
        UseExponentialBackoff: false,
    }

    // Sin exponential backoff, siempre retorna RetryDelay
    for attempt := 0; attempt < 5; attempt++ {
        got := config.calculateBackoff(attempt)
        if got != 5*time.Second {
            t.Errorf("calculateBackoff(%d) = %v, want 5s", attempt, got)
        }
    }
}
```

**Archivo:** `consumer_dlq_test.go` (test de integración con Testcontainers)

```go
package rabbit_test

import (
    "context"
    "errors"
    "testing"
    "time"
    "github.com/EduGoGroup/edugo-shared/messaging/rabbit"
    "github.com/EduGoGroup/edugo-shared/testing/containers"
)

func TestConsumer_DLQ_Integration(t *testing.T) {
    // Setup RabbitMQ con Testcontainers
    ctx := context.Background()
    rabbitContainer, err := containers.NewRabbitMQContainer(ctx)
    if err != nil {
        t.Fatalf("failed to start RabbitMQ: %v", err)
    }
    defer rabbitContainer.Terminate(ctx)

    url := rabbitContainer.ConnectionString()

    // Crear consumer con DLQ
    config := rabbit.ConsumerConfig{
        Queue:         "test-queue",
        PrefetchCount: 1,
        DLQ: rabbit.DLQConfig{
            Enabled:       true,
            MaxRetries:    3,
            DLXExchange:   "dlx",
            DLXRoutingKey: "dlq",
        },
    }

    consumer, err := rabbit.NewConsumer(url, config)
    if err != nil {
        t.Fatalf("failed to create consumer: %v", err)
    }
    defer consumer.Close()

    // Publicar mensaje de prueba
    publisher, _ := rabbit.NewPublisher(url)
    publisher.Publish("", "test-queue", []byte("test message"))

    // Handler que siempre falla
    failCount := 0
    handler := func(body []byte) error {
        failCount++
        return errors.New("simulated error")
    }

    // Consumir con timeout
    go consumer.ConsumeWithDLQ(handler)

    time.Sleep(20 * time.Second) // Esperar reintentos

    // Verificar que falló exactamente 3 veces
    if failCount != 3 {
        t.Errorf("expected 3 failures, got %d", failCount)
    }

    // Verificar que el mensaje está en DLQ
    dlqCount := rabbit.GetQueueMessageCount(url, "dlq")
    if dlqCount != 1 {
        t.Errorf("expected 1 message in DLQ, got %d", dlqCount)
    }
}
```

---

#### Versión Objetivo

**Tag:** `messaging/rabbit/v0.6.0`

**Changelog:**
```markdown
# Changelog - messaging/rabbit

## [0.6.0] - 2025-11-XX

### Added
- DLQ (Dead Letter Queue) support
- DLQConfig struct with configurable retries and backoff
- ConsumeWithDLQ method with automatic retry logic
- Exponential backoff support
- sendToDLQ helper para manejar mensajes fallidos
- setupDLQ helper para configurar DLX y DLQ
- Tests de integración con Testcontainers
- Documentación de uso de DLQ

### Changed
- ConsumerConfig ahora incluye DLQ field
- Consumer.Consume ahora soporta headers para retry count

### Fixed
- Mensajes con errores permanentes ya no se reencolan infinitamente
```

---

#### Tiempo Estimado

**Implementación:** 2-3 horas
- DLQConfig y helpers: 1 hora
- ConsumeWithDLQ modificado: 1 hora
- setupDLQ y sendToDLQ: 1 hora

**Tests:** 1-2 horas
- Tests unitarios de backoff: 30 min
- Tests de integración: 1-1.5 horas

**Total:** 3-5 horas

---

#### Impacto si NO se Implementa

- 🔴 Worker crashea con mensajes malos
- 🔴 Mensajes se reencolan infinitamente
- 🔴 No hay visibilidad de mensajes fallidos
- 🔴 Imposible depurar errores de procesamiento

---

## 🔐 auth/ (v0.5.0 → v0.6.0)

### Feature: Refresh Token Support

#### Requerido Por
- ✅ api-mobile (UX mejorada)
- ✅ api-admin (UX mejorada)

#### Estado Actual (v0.5.0)

⚠️ **VERIFICAR:** No se pudo confirmar si refresh tokens ya está implementado debido a `go mod tidy` requerido

**Supuesto:** NO existe (implementar)

**Archivo actual:** `jwt.go`

```go
// jwt.go (v0.5.0 - SUPUESTO ACTUAL)
type JWTManager struct {
    secretKey string
    issuer    string
}

func (j *JWTManager) GenerateToken(userID uuid.UUID, email string, role enum.SystemRole, expiration time.Duration) (string, error) {
    claims := Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(expiration).Unix(),
            Issuer:    j.issuer,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.secretKey))
}
```

---

#### Implementación Necesaria (v0.6.0)

**Archivo modificado:** `jwt.go`

```go
// jwt.go (v0.6.0 - MODIFICADO)
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

// TokenPair contiene access y refresh tokens
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
    ExpiresIn    int64  `json:"expires_in"` // Segundos hasta expiración del access token
    TokenType    string `json:"token_type"` // "Bearer"
}

// RefreshClaims son claims específicos para refresh tokens
type RefreshClaims struct {
    UserID uuid.UUID `json:"user_id"`
    Email  string    `json:"email"`
    Role   string    `json:"role"`
    jwt.RegisteredClaims
}

// GenerateTokenPair genera access y refresh tokens
func (j *JWTManager) GenerateTokenPair(userID uuid.UUID, email string, role enum.SystemRole) (*TokenPair, error) {
    // Access token (corta duración: 15 minutos)
    accessExpiration := 15 * time.Minute
    accessToken, err := j.GenerateToken(userID, email, role, accessExpiration)
    if err != nil {
        return nil, fmt.Errorf("failed to generate access token: %w", err)
    }

    // Refresh token (larga duración: 7 días)
    refreshExpiration := 7 * 24 * time.Hour
    refreshClaims := RefreshClaims{
        UserID: userID,
        Email:  email,
        Role:   string(role),
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshExpiration)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    j.issuer,
        },
    }

    refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    refreshToken, err := refreshTokenObj.SignedString([]byte(j.secretKey))
    if err != nil {
        return nil, fmt.Errorf("failed to generate refresh token: %w", err)
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        ExpiresIn:    int64(accessExpiration.Seconds()),
        TokenType:    "Bearer",
    }, nil
}

// RefreshAccessToken valida refresh token y genera nuevo access token
func (j *JWTManager) RefreshAccessToken(refreshToken string) (string, error) {
    // Parsear y validar refresh token
    token, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(j.secretKey), nil
    })

    if err != nil {
        return "", fmt.Errorf("invalid refresh token: %w", err)
    }

    claims, ok := token.Claims.(*RefreshClaims)
    if !ok || !token.Valid {
        return "", errors.New("invalid refresh token claims")
    }

    // Generar nuevo access token
    accessToken, err := j.GenerateToken(
        claims.UserID,
        claims.Email,
        enum.SystemRole(claims.Role),
        15*time.Minute,
    )
    if err != nil {
        return "", fmt.Errorf("failed to generate new access token: %w", err)
    }

    return accessToken, nil
}

// ValidateRefreshToken valida un refresh token
func (j *JWTManager) ValidateRefreshToken(refreshToken string) (*RefreshClaims, error) {
    token, err := jwt.ParseWithClaims(refreshToken, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(j.secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*RefreshClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid token")
    }

    return claims, nil
}
```

---

#### Tests Requeridos

```go
// refresh_token_test.go
func TestJWTManager_GenerateTokenPair(t *testing.T) {
    manager := NewJWTManager("secret", "test")
    userID := uuid.New()

    pair, err := manager.GenerateTokenPair(userID, "test@example.com", enum.SystemRoleStudent)

    if err != nil {
        t.Fatalf("GenerateTokenPair failed: %v", err)
    }

    if pair.AccessToken == "" {
        t.Error("access token is empty")
    }
    if pair.RefreshToken == "" {
        t.Error("refresh token is empty")
    }
    if pair.ExpiresIn != 900 { // 15 min = 900 seg
        t.Errorf("expected ExpiresIn=900, got %d", pair.ExpiresIn)
    }
}

func TestJWTManager_RefreshAccessToken(t *testing.T) {
    manager := NewJWTManager("secret", "test")
    userID := uuid.New()

    // Generar token pair inicial
    pair, _ := manager.GenerateTokenPair(userID, "test@example.com", enum.SystemRoleStudent)

    // Usar refresh token para obtener nuevo access token
    newAccessToken, err := manager.RefreshAccessToken(pair.RefreshToken)

    if err != nil {
        t.Fatalf("RefreshAccessToken failed: %v", err)
    }

    if newAccessToken == "" {
        t.Error("new access token is empty")
    }

    // Validar nuevo access token
    claims, err := manager.ValidateToken(newAccessToken)
    if err != nil {
        t.Fatalf("new access token is invalid: %v", err)
    }

    if claims.UserID != userID {
        t.Errorf("expected userID=%v, got %v", userID, claims.UserID)
    }
}
```

---

#### Versión Objetivo

**Tag:** `auth/v0.6.0`

**Tiempo:** 2-3 horas

---

## 📋 Resumen de Features Faltantes

| Módulo | Feature | Prioridad | Tiempo Est. | Versión Objetivo |
|--------|---------|-----------|-------------|------------------|
| messaging/rabbit | DLQ Support | P0 | 3-5 horas | v0.6.0 |
| auth | Refresh Tokens | P1 | 2-3 horas | v0.6.0 |

**Total features:** 2  
**Total tiempo:** 5-8 horas (~1 día)

---

**Documento generado:** 15 de Noviembre, 2025  
**Próximo documento:** `05-PLAN_SPRINTS.md`
