# Diagramas de Flujo y Procesos

## 1. Flujo de Autenticación JWT

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           FLUJO DE AUTENTICACIÓN                                 │
└─────────────────────────────────────────────────────────────────────────────────┘

                                    LOGIN
                                    ═════
┌──────────┐                                              ┌──────────────────┐
│  Client  │                                              │   Auth Service   │
└────┬─────┘                                              └────────┬─────────┘
     │                                                             │
     │  POST /auth/login                                           │
     │  { email, password }                                        │
     │ ──────────────────────────────────────────────────────────► │
     │                                                             │
     │                                         ┌───────────────────┴───────────────────┐
     │                                         │  1. Buscar usuario por email          │
     │                                         │  2. Verificar password con bcrypt      │
     │                                         │  3. Generar JWT token                  │
     │                                         │     - Claims: user_id, email, role     │
     │                                         │     - Firmado con HS256               │
     │                                         │  4. Generar refresh token              │
     │                                         └───────────────────┬───────────────────┘
     │                                                             │
     │  { access_token, refresh_token, expires_in }                │
     │ ◄────────────────────────────────────────────────────────── │
     │                                                             │


                              ACCESO A RECURSO PROTEGIDO
                              ══════════════════════════
┌──────────┐                 ┌─────────────────┐              ┌──────────────────┐
│  Client  │                 │  JWT Middleware │              │     Handler      │
└────┬─────┘                 └────────┬────────┘              └────────┬─────────┘
     │                                │                                │
     │  GET /api/resource             │                                │
     │  Authorization: Bearer {token} │                                │
     │ ─────────────────────────────► │                                │
     │                                │                                │
     │              ┌─────────────────┴─────────────────┐              │
     │              │  1. Extraer token del header      │              │
     │              │  2. Validar firma con secret      │              │
     │              │  3. Verificar expiración          │              │
     │              │  4. Extraer claims                │              │
     │              │  5. Guardar en context            │              │
     │              └─────────────────┬─────────────────┘              │
     │                                │                                │
     │                                │  c.Set("user_id", claims.UserID)
     │                                │  c.Set("email", claims.Email)  │
     │                                │  c.Set("role", claims.Role)    │
     │                                │ ─────────────────────────────► │
     │                                │                                │
     │                                │                   ┌────────────┴────────────┐
     │                                │                   │  Procesar request con   │
     │                                │                   │  información del usuario │
     │                                │                   └────────────┬────────────┘
     │                                │                                │
     │  { response }                  │                                │
     │ ◄───────────────────────────────────────────────────────────── │
     │                                │                                │


                                 REFRESH TOKEN
                                 ═════════════
┌──────────┐                                              ┌──────────────────┐
│  Client  │                                              │   Auth Service   │
└────┬─────┘                                              └────────┬─────────┘
     │                                                             │
     │  POST /auth/refresh                                         │
     │  { refresh_token }                                          │
     │ ──────────────────────────────────────────────────────────► │
     │                                                             │
     │                                         ┌───────────────────┴───────────────────┐
     │                                         │  1. Validar refresh token             │
     │                                         │  2. Verificar no revocado             │
     │                                         │  3. Generar nuevo access token        │
     │                                         │  4. (Opcional) Rotar refresh token    │
     │                                         └───────────────────┬───────────────────┘
     │                                                             │
     │  { access_token, refresh_token, expires_in }                │
     │ ◄────────────────────────────────────────────────────────── │
     │                                                             │
```

---

## 2. Flujo de Bootstrap de Aplicación

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         FLUJO DE INICIALIZACIÓN                                  │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────────┐
│    main()   │
└──────┬──────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│  1. LOAD CONFIGURATION                                                        │
│  ═══════════════════                                                          │
│                                                                               │
│  config.yaml ──► config.Load() ──► BaseConfig                                │
│                                    ├── Environment: "local"                   │
│                                    ├── Server: { port: 8080, ... }           │
│                                    ├── Database: { host, port, ... }         │
│                                    └── MongoDB: { uri, database, ... }       │
└──────────────────────────────────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│  2. CREATE LIFECYCLE MANAGER                                                  │
│  ═══════════════════════════                                                  │
│                                                                               │
│  lifecycle.NewManager(logger) ──► Manager { resources: [] }                  │
│                                                                               │
└──────────────────────────────────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│  3. BOOTSTRAP RESOURCES                                                       │
│  ═════════════════════                                                        │
│                                                                               │
│  ┌────────────────────────────────────────────────────────────────────────┐  │
│  │  Step 1: Initialize Logger                                              │  │
│  │  ─────────────────────────                                              │  │
│  │  LoggerFactory.CreateLogger() ──► *logrus.Logger                        │  │
│  │  ✓ Log: "Starting application bootstrap..."                             │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                              │                                                │
│                              ▼                                                │
│  ┌────────────────────────────────────────────────────────────────────────┐  │
│  │  Step 2: Initialize PostgreSQL                                          │  │
│  │  ────────────────────────────                                           │  │
│  │  PostgreSQLFactory.CreateConnection(config) ──► *gorm.DB                │  │
│  │  lifecycleManager.RegisterSimple("postgresql", cleanup)                 │  │
│  │  ✓ Log: "PostgreSQL connection established"                             │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                              │                                                │
│                              ▼                                                │
│  ┌────────────────────────────────────────────────────────────────────────┐  │
│  │  Step 3: Initialize MongoDB                                             │  │
│  │  ──────────────────────────                                             │  │
│  │  MongoDBFactory.CreateConnection(config) ──► *mongo.Client              │  │
│  │  MongoDBFactory.GetDatabase(client, dbName) ──► *mongo.Database         │  │
│  │  lifecycleManager.RegisterSimple("mongodb", cleanup)                    │  │
│  │  ✓ Log: "MongoDB connection established"                                │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                              │                                                │
│                              ▼                                                │
│  ┌────────────────────────────────────────────────────────────────────────┐  │
│  │  Step 4: Initialize RabbitMQ (if required)                              │  │
│  │  ─────────────────────────────────────────                              │  │
│  │  RabbitMQFactory.CreateConnection(config) ──► *amqp.Connection          │  │
│  │  RabbitMQFactory.CreateChannel(conn) ──► *amqp.Channel                  │  │
│  │  lifecycleManager.RegisterSimple("rabbitmq", cleanup)                   │  │
│  │  ✓ Log: "RabbitMQ connection established"                               │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                              │                                                │
│                              ▼                                                │
│  ┌────────────────────────────────────────────────────────────────────────┐  │
│  │  Step 5: Initialize S3 (if required)                                    │  │
│  │  ───────────────────────────────────                                    │  │
│  │  S3Factory.CreateClient(config) ──► *s3.Client                          │  │
│  │  S3Factory.CreatePresignClient(client) ──► PresignClient                │  │
│  │  ✓ Log: "S3 client initialized"                                         │  │
│  └────────────────────────────────────────────────────────────────────────┘  │
│                                                                               │
└──────────────────────────────────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│  4. HEALTH CHECKS                                                             │
│  ═══════════════                                                              │
│                                                                               │
│  ┌─────────────────────┐    ┌─────────────────────┐                          │
│  │ PostgreSQL: SELECT 1│    │ MongoDB: Ping()     │                          │
│  │      ✓ PASS         │    │      ✓ PASS         │                          │
│  └─────────────────────┘    └─────────────────────┘                          │
│                                                                               │
│  ✓ Log: "All health checks passed"                                           │
└──────────────────────────────────────────────────────────────────────────────┘
       │
       ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│  5. RETURN RESOURCES                                                          │
│  ═══════════════════                                                          │
│                                                                               │
│  &Resources{                                                                  │
│      Logger:           *logrus.Logger,                                        │
│      PostgreSQL:       *gorm.DB,                                              │
│      MongoDB:          *mongo.Client,                                         │
│      MongoDatabase:    *mongo.Database,                                       │
│      MessagePublisher: MessagePublisher,                                      │
│      StorageClient:    StorageClient,                                         │
│  }                                                                            │
│                                                                               │
│  ✓ Log: "Application bootstrap completed successfully"                        │
└──────────────────────────────────────────────────────────────────────────────┘
```

---

## 3. Flujo de Mensajería (RabbitMQ)

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                           FLUJO DE PUBLICACIÓN                                   │
└─────────────────────────────────────────────────────────────────────────────────┘

┌──────────────────┐          ┌──────────────────┐          ┌──────────────────┐
│    Service A     │          │    RabbitMQ      │          │    Service B     │
│   (Publisher)    │          │     Broker       │          │   (Consumer)     │
└────────┬─────────┘          └────────┬─────────┘          └────────┬─────────┘
         │                             │                             │
         │  1. Publish(ctx, exchange,  │                             │
         │     routingKey, message)    │                             │
         │ ──────────────────────────► │                             │
         │                             │                             │
         │          ┌──────────────────┴──────────────────┐          │
         │          │  - Serialize to JSON                │          │
         │          │  - Set ContentType: application/json│          │
         │          │  - Set DeliveryMode: Persistent     │          │
         │          │  - Route to queue by routing key    │          │
         │          └──────────────────┬──────────────────┘          │
         │                             │                             │
         │                             │  2. Deliver to queue        │
         │                             │ ──────────────────────────► │
         │                             │                             │
         │                             │          ┌──────────────────┴──────────────────┐
         │                             │          │  3. MessageHandler(ctx, body)       │
         │                             │          │     - Deserialize JSON              │
         │                             │          │     - Process message               │
         │                             │          │     - ACK or NACK                   │
         │                             │          └──────────────────┬──────────────────┘
         │                             │                             │
         │                             │  4. Acknowledgment          │
         │                             │ ◄────────────────────────── │
         │                             │                             │


┌─────────────────────────────────────────────────────────────────────────────────┐
│                         FLUJO DE DEAD LETTER QUEUE                               │
└─────────────────────────────────────────────────────────────────────────────────┘

┌──────────────────┐          ┌──────────────────┐          ┌──────────────────┐
│     Message      │          │   Main Queue     │          │      DLQ         │
└────────┬─────────┘          └────────┬─────────┘          └────────┬─────────┘
         │                             │                             │
         │  1. Message arrives         │                             │
         │ ──────────────────────────► │                             │
         │                             │                             │
         │          ┌──────────────────┴──────────────────┐          │
         │          │  2. Consumer processes message      │          │
         │          │     ❌ Error occurs                 │          │
         │          │     ❌ Max retries exceeded         │          │
         │          └──────────────────┬──────────────────┘          │
         │                             │                             │
         │                             │  3. NACK (no requeue)       │
         │                             │  Message to DLQ             │
         │                             │ ──────────────────────────► │
         │                             │                             │
         │                             │          ┌──────────────────┴──────────────────┐
         │                             │          │  4. DLQ Consumer                    │
         │                             │          │     - Log error                     │
         │                             │          │     - Alert if needed               │
         │                             │          │     - Manual intervention           │
         │                             │          └─────────────────────────────────────┘
```

---

## 4. Flujo de Lifecycle Cleanup

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                         FLUJO DE CLEANUP (LIFO)                                  │
└─────────────────────────────────────────────────────────────────────────────────┘

                    REGISTRO (FIFO)                    CLEANUP (LIFO)
                    ═══════════════                    ═════════════

    ┌─────────────────────────────┐           ┌─────────────────────────────┐
    │  Register("postgresql", ...)│           │  Cleanup "s3"               │
    │         [Index: 0]          │           │  ✓ S3 client closed         │
    └──────────────┬──────────────┘           └──────────────┬──────────────┘
                   │                                         │
                   ▼                                         ▼
    ┌─────────────────────────────┐           ┌─────────────────────────────┐
    │  Register("mongodb", ...)   │           │  Cleanup "rabbitmq"         │
    │         [Index: 1]          │           │  ✓ Channel closed           │
    └──────────────┬──────────────┘           │  ✓ Connection closed        │
                   │                          └──────────────┬──────────────┘
                   ▼                                         │
    ┌─────────────────────────────┐                         ▼
    │  Register("rabbitmq", ...)  │           ┌─────────────────────────────┐
    │         [Index: 2]          │           │  Cleanup "mongodb"          │
    └──────────────┬──────────────┘           │  ✓ Client disconnected      │
                   │                          └──────────────┬──────────────┘
                   ▼                                         │
    ┌─────────────────────────────┐                         ▼
    │  Register("s3", ...)        │           ┌─────────────────────────────┐
    │         [Index: 3]          │           │  Cleanup "postgresql"       │
    └─────────────────────────────┘           │  ✓ Connection pool closed   │
                                              └─────────────────────────────┘

    Resources list after register:             Cleanup order:
    [postgresql, mongodb, rabbitmq, s3]        [s3, rabbitmq, mongodb, postgresql]
```

---

## 5. Flujo de Manejo de Errores

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                          FLUJO DE MANEJO DE ERRORES                              │
└─────────────────────────────────────────────────────────────────────────────────┘

┌───────────────┐      ┌───────────────┐      ┌───────────────┐      ┌──────────┐
│    Handler    │      │    Service    │      │  Repository   │      │    DB    │
└───────┬───────┘      └───────┬───────┘      └───────┬───────┘      └────┬─────┘
        │                      │                      │                   │
        │  CreateUser(req)     │                      │                   │
        │ ───────────────────► │                      │                   │
        │                      │                      │                   │
        │                      │  repo.Create(user)   │                   │
        │                      │ ───────────────────► │                   │
        │                      │                      │                   │
        │                      │                      │  INSERT INTO...   │
        │                      │                      │ ────────────────► │
        │                      │                      │                   │
        │                      │                      │  ❌ Duplicate Key │
        │                      │                      │ ◄──────────────── │
        │                      │                      │                   │
        │                      │   ┌──────────────────┴──────────────────┐
        │                      │   │  Wrap error:                        │
        │                      │   │  errors.NewAlreadyExistsError("user")│
        │                      │   │    .WithField("email", email)       │
        │                      │   └──────────────────┬──────────────────┘
        │                      │                      │
        │                      │  AppError{           │
        │                      │    Code: ALREADY_EXISTS
        │                      │    StatusCode: 409   │
        │                      │  }                   │
        │                      │ ◄─────────────────── │
        │                      │                      │
        │   ┌──────────────────┴──────────────────┐   │
        │   │  Check error type                   │   │
        │   │  Log error if internal              │   │
        │   │  Return appropriate response        │   │
        │   └──────────────────┬──────────────────┘   │
        │                      │                      │
        │  HTTP 409 Conflict   │                      │
        │  {                   │                      │
        │    "error": "user already exists",          │
        │    "code": "ALREADY_EXISTS",                │
        │    "fields": { "email": "..." }             │
        │  }                   │                      │
        │ ◄─────────────────── │                      │
        │                      │                      │
```

---

## 6. Flujo de Testing con Containers

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        FLUJO DE TESTING CON CONTAINERS                           │
└─────────────────────────────────────────────────────────────────────────────────┘

┌──────────────────────────────────────────────────────────────────────────────────┐
│  TestMain(m *testing.M)                                                          │
│  ═══════════════════════                                                         │
│                                                                                  │
│  ┌────────────────────────────────────────────────────────────────────────────┐ │
│  │  1. Create Config                                                          │ │
│  │     config := containers.NewConfig().                                      │ │
│  │         WithPostgreSQL(nil).                                               │ │
│  │         WithMongoDB(nil).                                                  │ │
│  │         Build()                                                            │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                           │
│                                      ▼                                           │
│  ┌────────────────────────────────────────────────────────────────────────────┐ │
│  │  2. Get Manager (Singleton)                                                │ │
│  │     manager, err := containers.GetManager(nil, config)                     │ │
│  │                                                                            │ │
│  │     ┌─────────────────────────────────────────────────────────────────┐   │ │
│  │     │  First call:                                                    │   │ │
│  │     │  - Start PostgreSQL container (~5s)                             │   │ │
│  │     │  - Start MongoDB container (~5s)                                │   │ │
│  │     │  - Wait for readiness                                           │   │ │
│  │     │  - Store in singleton                                           │   │ │
│  │     └─────────────────────────────────────────────────────────────────┘   │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                           │
│                                      ▼                                           │
│  ┌────────────────────────────────────────────────────────────────────────────┐ │
│  │  3. Run Tests                                                              │ │
│  │     os.Exit(m.Run())                                                       │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
│                                      │                                           │
│                                      ▼                                           │
│  ┌────────────────────────────────────────────────────────────────────────────┐ │
│  │  4. Cleanup (defer)                                                        │ │
│  │     manager.Cleanup(ctx)                                                   │ │
│  │     - Stop PostgreSQL container                                            │ │
│  │     - Stop MongoDB container                                               │ │
│  │     - Remove containers                                                    │ │
│  └────────────────────────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────────┘


┌──────────────────────────────────────────────────────────────────────────────────┐
│  Individual Test                                                                 │
│  ═══════════════                                                                 │
│                                                                                  │
│  func TestUserRepository(t *testing.T) {                                         │
│      // Get same manager (instant - already started)                             │
│      manager, _ := containers.GetManager(t, nil)                                 │
│                                                                                  │
│      // Get database                                                             │
│      db := manager.PostgreSQL().DB()                                             │
│                                                                                  │
│      // Setup test data                                                          │
│      db.Exec("INSERT INTO users ...")                                            │
│                                                                                  │
│      // Run test assertions                                                      │
│      ...                                                                         │
│                                                                                  │
│      // Cleanup for next test                                                    │
│      manager.CleanPostgreSQL(ctx, "users")                                       │
│  }                                                                               │
└──────────────────────────────────────────────────────────────────────────────────┘
```
