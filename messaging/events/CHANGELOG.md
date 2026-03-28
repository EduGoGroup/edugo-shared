# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/messaging/events` deben registrarse aqui.

## [Unreleased]

## [v0.51.0] - 2026-03-28

### Added

- `AssessmentPublishedEvent` — CloudEvent emitido cuando un profesor publica una evaluación. Payload: `AssessmentPublishedEvent`.
- `AssessmentAssignedEvent` — CloudEvent emitido cuando una evaluación es asignada a un estudiante o unidad académica. Payload: `AssessmentAssignedPayload`.
- `AssessmentAttemptRecordedEvent` — CloudEvent emitido cuando un estudiante finaliza un intento de evaluación. Payload: `AssessmentAttemptRecordedPayload`.
- `AssessmentReviewedEvent` — CloudEvent emitido cuando un profesor completa la revisión manual de un intento. Payload: `AssessmentReviewedPayload`.

### Fixed

- Error strings en constructores lowercaseados para cumplir con la convención Go (ST1005): `"title no puede estar vacío"`, `"status no puede estar vacío"`.

## [v0.50.1] - 2026-02-25

### Changed

- Actualización de dependencias a versiones más recientes.

### Added

- Baseline de documentacion de fase 1 con `README.md`, `docs/README.md` y organizacion local por modulo.
