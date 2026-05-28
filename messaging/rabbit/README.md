# Messaging RabbitMQ

Runtime RabbitMQ para publicar y consumir mensajes JSON con soporte opcional de DLQ y retries.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/messaging/rabbit`
- Carpeta: `messaging/rabbit`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/messaging/rabbit
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Conectar a RabbitMQ y abrir un canal reutilizable.
2. Declarar exchanges, queues, bindings y QoS/prefetch segun configuracion.
3. Publicar mensajes serializados a JSON con prioridad y contexto.
4. Consumir mensajes, ejecutar handlers y hacer ack o nack segun resultado.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- El modulo cubre tanto happy path como reintentos y DLQ.
- La salud de la conexion se valida mediante un exchange temporal para evitar race conditions.
- Tiene tests unitarios e integracion bastante amplios.
