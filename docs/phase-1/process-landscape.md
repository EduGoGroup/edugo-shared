# Mapa de procesos

Esta vista no reemplaza la documentacion por modulo. Su objetivo es mostrar como se distribuyen los procesos dominantes dentro del repositorio sin repetir el detalle fino.

## 1. Carga de configuracion y bootstrap

1. `config` carga archivo y entorno.
2. `bootstrap` valida factories y decide que recursos deben arrancar.
3. `logger`, `database/*`, `messaging/rabbit` y S3 se inicializan en orden.
4. `lifecycle` puede recibir cleanups y ordenar el shutdown.

## 2. Pipeline HTTP autenticado

1. `middleware/gin` valida JWT con `auth`.
2. El middleware guarda claims y contexto de usuario en `gin.Context`.
3. Las reglas de permisos se aplican antes de llegar al handler.
4. `AuditMiddleware` registra las mutaciones usando `audit` o `audit/postgres` segun el adapter disponible.

## 3. Persistencia relacional y no relacional

1. `database/postgres` y `database/mongodb` resuelven conectividad, pools y health checks.
2. `repository` se monta sobre GORM y entidades externas para CRUD y listados seguros.
3. `database/postgres` encapsula transacciones SQL y GORM.

## 4. Mensajeria asincrona

1. `messaging/events` define payloads de dominio serializables.
2. `messaging/rabbit` publica y consume esos mensajes via RabbitMQ.
3. La politica de retries/DLQ vive en `messaging/rabbit`, no en el modulo de eventos.

## 5. Configuracion dinamica de pantallas

1. `screenconfig` valida estructuras JSON de templates.
2. Aplica overrides por plataforma y resuelve slots de runtime.
3. Construye arboles de menu y utilidades de permisos para clientes consumidores.

## 6. Testing de integracion

1. `testing/containers` crea containers reutilizables por backend.
2. Los modulos de datos, mensajeria y bootstrap los usan en tests de integracion.
3. El cleanup puede hacerse sin destruir cada container entre test y test.

## Punto transversal

- `common` provee errores, validator, UUIDs, enums y helpers de entorno que aparecen repartidos por varios flujos.
