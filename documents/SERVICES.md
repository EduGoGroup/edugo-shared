# Servicios Externos y Dependencias

## Resumen de Servicios Requeridos

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        SERVICIOS EXTERNOS REQUERIDOS                             │
├─────────────────────────────────────────────────────────────────────────────────┤
│                                                                                  │
│  ┌─────────────────┐  ┌─────────────────┐  ┌─────────────────┐                  │
│  │   PostgreSQL    │  │     MongoDB     │  │    RabbitMQ     │                  │
│  │     15.x+       │  │      7.0+       │  │     3.12+       │                  │
│  │   Puerto: 5432  │  │  Puerto: 27017  │  │  Puerto: 5672   │                  │
│  │   REQUERIDO     │  │   REQUERIDO     │  │   OPCIONAL      │                  │
│  └─────────────────┘  └─────────────────┘  └─────────────────┘                  │
│                                                                                  │
│  ┌─────────────────┐                                                            │
│  │     AWS S3      │                                                            │
│  │   o Compatible  │                                                            │
│  │   (MinIO, etc.) │                                                            │
│  │   OPCIONAL      │                                                            │
│  └─────────────────┘                                                            │
│                                                                                  │
└─────────────────────────────────────────────────────────────────────────────────┘
```

---

## 1. PostgreSQL

### Descripción
Base de datos relacional para datos transaccionales y estructurados.

### Versión Requerida
- **Mínimo:** PostgreSQL 15
- **Recomendado:** PostgreSQL 15-alpine (para containers)

### Configuración de Conexión

| Parámetro | Descripción | Valor por Defecto |
|-----------|-------------|-------------------|
| `host` | Hostname del servidor | `localhost` |
| `port` | Puerto de conexión | `5432` |
| `user` | Usuario de BD | - |
| `password` | Contraseña | - |
| `database` | Nombre de la BD | - |
| `ssl_mode` | Modo SSL | `disable` |
| `max_open_conns` | Conexiones máximas | `25` |
| `max_idle_conns` | Conexiones idle | `5` |
| `conn_max_lifetime` | Vida máxima de conexión | `5m` |

### Docker Compose
```yaml
services:
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: edugo_user
      POSTGRES_PASSWORD: edugo_pass
      POSTGRES_DB: edugo_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U edugo_user -d edugo_db"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  postgres_data:
```

### Health Check
```sql
SELECT 1;
```

### SSL Modes Soportados
| Mode | Descripción |
|------|-------------|
| `disable` | Sin SSL |
| `require` | SSL requerido, sin verificación de cert |
| `verify-ca` | Verificar CA del servidor |
| `verify-full` | Verificar CA y hostname |

---

## 2. MongoDB

### Descripción
Base de datos NoSQL para documentos flexibles y datos no estructurados.

### Versión Requerida
- **Mínimo:** MongoDB 7.0
- **Recomendado:** MongoDB 7.0 (imagen oficial)

### Configuración de Conexión

| Parámetro | Descripción | Valor por Defecto |
|-----------|-------------|-------------------|
| `uri` | URI de conexión completa | - |
| `database` | Nombre de la BD | - |
| `max_pool_size` | Tamaño máximo del pool | `100` |
| `min_pool_size` | Tamaño mínimo del pool | `10` |
| `connect_timeout` | Timeout de conexión | `30s` |

### Formato URI
```
mongodb://[username:password@]host:port/[database][?options]
```

### Ejemplos de URI
```bash
# Local sin auth
mongodb://localhost:27017

# Con autenticación
mongodb://user:password@localhost:27017/mydb?authSource=admin

# Replica Set
mongodb://host1:27017,host2:27017,host3:27017/mydb?replicaSet=rs0

# Atlas
mongodb+srv://user:password@cluster.mongodb.net/mydb
```

### Docker Compose
```yaml
services:
  mongodb:
    image: mongo:7.0
    environment:
      MONGO_INITDB_ROOT_USERNAME: root
      MONGO_INITDB_ROOT_PASSWORD: rootpass
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  mongodb_data:
```

### Health Check
```javascript
db.runCommand({ ping: 1 })
```

---

## 3. RabbitMQ

### Descripción
Message broker para comunicación asíncrona entre servicios.

### Versión Requerida
- **Mínimo:** RabbitMQ 3.12
- **Recomendado:** RabbitMQ 3.12-management-alpine

### Estado
**OPCIONAL** - El sistema puede funcionar sin RabbitMQ si no se requiere mensajería.

### Configuración de Conexión

| Parámetro | Descripción | Valor por Defecto |
|-----------|-------------|-------------------|
| `url` | AMQP URL completa | - |

### Formato URL
```
amqp://[username:password@]host:port/[vhost]
```

### Ejemplos de URL
```bash
# Local por defecto
amqp://guest:guest@localhost:5672/

# Con vhost
amqp://user:pass@localhost:5672/myvhost

# Con parámetros
amqp://user:pass@localhost:5672/myvhost?heartbeat=30
```

### Docker Compose
```yaml
services:
  rabbitmq:
    image: rabbitmq:3.12-management-alpine
    environment:
      RABBITMQ_DEFAULT_USER: edugo_user
      RABBITMQ_DEFAULT_PASS: edugo_pass
    ports:
      - "5672:5672"   # AMQP
      - "15672:15672" # Management UI
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    healthcheck:
      test: rabbitmq-diagnostics -q ping
      interval: 30s
      timeout: 10s
      retries: 5

volumes:
  rabbitmq_data:
```

### Puertos
| Puerto | Protocolo | Uso |
|--------|-----------|-----|
| 5672 | AMQP | Conexión de aplicaciones |
| 15672 | HTTP | Management UI |
| 15692 | HTTP | Prometheus metrics |

### Management UI
Acceso en `http://localhost:15672` con las credenciales configuradas.

---

## 4. AWS S3 (o Compatible)

### Descripción
Servicio de almacenamiento de objetos para archivos y assets.

### Estado
**OPCIONAL** - El sistema puede funcionar sin S3 si no se requiere almacenamiento de archivos.

### Configuración

| Parámetro | Descripción | Requerido |
|-----------|-------------|-----------|
| `bucket` | Nombre del bucket | Sí |
| `region` | Región AWS | Sí |
| `access_key_id` | AWS Access Key | Sí |
| `secret_access_key` | AWS Secret Key | Sí |

### Alternativas Compatibles

#### MinIO (Self-hosted)
```yaml
services:
  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minioadmin
      MINIO_ROOT_PASSWORD: minioadmin
    ports:
      - "9000:9000"  # API
      - "9001:9001"  # Console
    volumes:
      - minio_data:/data

volumes:
  minio_data:
```

#### LocalStack (para testing)
```yaml
services:
  localstack:
    image: localstack/localstack:latest
    environment:
      SERVICES: s3
      DEFAULT_REGION: us-east-1
    ports:
      - "4566:4566"
```

---

## Docker Compose Completo (Desarrollo)

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15-alpine
    container_name: edugo-postgres
    environment:
      POSTGRES_USER: edugo_user
      POSTGRES_PASSWORD: edugo_pass
      POSTGRES_DB: edugo_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U edugo_user -d edugo_db"]
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - edugo-network

  mongodb:
    image: mongo:7.0
    container_name: edugo-mongodb
    ports:
      - "27017:27017"
    volumes:
      - mongodb_data:/data/db
    healthcheck:
      test: echo 'db.runCommand("ping").ok' | mongosh localhost:27017/test --quiet
      interval: 10s
      timeout: 5s
      retries: 5
    networks:
      - edugo-network

  rabbitmq:
    image: rabbitmq:3.12-management-alpine
    container_name: edugo-rabbitmq
    environment:
      RABBITMQ_DEFAULT_USER: edugo_user
      RABBITMQ_DEFAULT_PASS: edugo_pass
    ports:
      - "5672:5672"
      - "15672:15672"
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    healthcheck:
      test: rabbitmq-diagnostics -q ping
      interval: 30s
      timeout: 10s
      retries: 5
    networks:
      - edugo-network

  minio:
    image: minio/minio:latest
    container_name: edugo-minio
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: edugo_user
      MINIO_ROOT_PASSWORD: edugo_pass
    ports:
      - "9000:9000"
      - "9001:9001"
    volumes:
      - minio_data:/data
    networks:
      - edugo-network

volumes:
  postgres_data:
  mongodb_data:
  rabbitmq_data:
  minio_data:

networks:
  edugo-network:
    driver: bridge
```

---

## Comandos Útiles

### Iniciar Todos los Servicios
```bash
docker-compose up -d
```

### Verificar Estado
```bash
docker-compose ps
docker-compose logs -f
```

### Detener Servicios
```bash
docker-compose down
```

### Limpiar Datos (⚠️ Destructivo)
```bash
docker-compose down -v
```

---

## Variables de Entorno Recomendadas

```bash
# PostgreSQL
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=edugo_user
DATABASE_PASSWORD=edugo_pass
DATABASE_NAME=edugo_db
DATABASE_SSL_MODE=disable

# MongoDB
MONGODB_URI=mongodb://localhost:27017
MONGODB_DATABASE=edugo_db

# RabbitMQ
RABBITMQ_URL=amqp://edugo_user:edugo_pass@localhost:5672/

# S3 / MinIO
S3_BUCKET=edugo-bucket
S3_REGION=us-east-1
S3_ACCESS_KEY_ID=edugo_user
S3_SECRET_ACCESS_KEY=edugo_pass
S3_ENDPOINT=http://localhost:9000  # Solo para MinIO/LocalStack
```

---

## Verificación de Conectividad

### PostgreSQL
```bash
psql -h localhost -U edugo_user -d edugo_db -c "SELECT 1;"
```

### MongoDB
```bash
mongosh "mongodb://localhost:27017/edugo_db" --eval "db.runCommand({ping:1})"
```

### RabbitMQ
```bash
curl -u edugo_user:edugo_pass http://localhost:15672/api/overview
```

### MinIO
```bash
curl http://localhost:9000/minio/health/live
```
