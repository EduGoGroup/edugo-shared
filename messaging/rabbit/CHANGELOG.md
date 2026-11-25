# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/messaging/rabbit` deben registrarse aqui.

## [Unreleased]

### Fixed
- `ConsumeWithDLQ`: el mensaje que va a la cola muerta se publica **persistente** (`DeliveryMode` 2);
  antes salía transitorio y se perdía al reiniciar el broker (bug 0157, HU-IA.3).
- El error del handler ya no se pierde: viaja en la cabecera `x-error` del mensaje muerto, acotado a
  1024 bytes en frontera de runa y con marca de corte.
- `sendToDLQ` ya no muta las cabeceras de la entrega.

### Added
- `ConsumerConfig.IsPermanent func(error) bool`: si devuelve `true`, el mensaje va a la cola muerta sin
  reintentos (deuda 140). Con `nil`, la conducta es la de siempre. Solo se consulta con la cola muerta
  habilitada, y con el contexto cancelado el mensaje vuelve a la cola en vez de morir.
- `ConsumerConfig.OnRetry` y `ConsumerConfig.OnDeadLetter`, ganchos opcionales que reciben un
  `FailureEvent` (cola, reintento, máximo, espera, motivo, error, tamaño y hash corto del cuerpo).
- Cabeceras nuevas del mensaje muerto: `x-error`, `x-dead-reason` (`permanent` o `max_retries`),
  `x-first-exchange` y `x-first-routing-key` (el exchange y la clave por los que entró; se guardan en el
  primer fallo y sobreviven a los reintentos). `x-original-*`, `x-failed-at` y `x-retry-count` no cambian.
- `DeliveryFromContext(ctx)` devuelve `DeliveryInfo` (cola, reintento, máximo, `IsLastAttempt()`) dentro
  del handler de `ConsumeWithDLQ`. `MessageHandler` y `Consume` no cambian.
- Constantes exportadas con los nombres de cabecera (`Header*`) y los motivos (`DeadReason*`).
- Primeros tests de `ConsumeWithDLQ` contra un broker (`consumer_dlq_integration_test.go`).

### Changed
- `ConsumerConfig` lleva campos de función y deja de ser comparable con `==`.

## [0.900.1] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencia interna: `edugo-shared/testing` `v0.1.0` → `v0.900.0`.
- Dependencias directas: `github.com/rabbitmq/amqp091-go` `v1.10.0` → `v1.16.0`, `github.com/stretchr/testify` `v1.11.1` → `v1.12.1`.

## [0.900.0] - 2026-06-02

### Changed

- Reversionado a la serie `0.900.x` para escapar del fantasma de checksum del tag `v0.1.0`
  (su contenido cambió tras la limpieza de historial del repo, dejando go.sum de consumidores
  desincronizado; ver bug 0022 / ADR 0015). No hay cambios de código ni de API respecto al
  contenido actual de `main`: es una versión fresca con hash limpio para los consumidores
  (`edugo-api-learning`, `edugo-worker`).

## [0.100.0] - 2026-04-02

### Removed

- Deleted `publisher_unit_test.go` (100% covered by integration tests)

### Changed

- Added `//go:build integration` tag to `publisher_test.go`, `connection_test.go`, `consumer_test.go`

## [Unreleased]

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

### Added

- Baseline de documentacion de fase 1 con `README.md`, `docs/README.md` y organizacion local por modulo.
