# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/messaging/events` deben registrarse aqui.

## [Unreleased]

## [0.900.1] - 2026-06-11

### Added

- Campo `academic_unit_id` (omitempty) en `AssessmentAttemptRecordedPayload` y `AssessmentReviewedPayload`: lleva la unidad académica activa del emisor para que el worker propague el tenant al push, habilitando el deep-link multi-tenant (plan 020 N5, F4.6). Un evento viejo sin el campo deserializa a vacío sin romper compatibilidad.

## [0.900.0] - 2026-06-07

### Changed

- **Breaking** en `AssessmentAssignedPayload`: redefinido para apuntar a una sesión de materia (N4 F3, ADR 0019). Se elimina el targeting genérico `target_type` (`"student"`/`"unit"`) + `target_id` y se reemplaza por `subject_offering_id` (los destinatarios son los alumnos inscritos en la oferta); rename `assigned_by_id → assigned_by_membership_id` (re-llaveo a `academic.memberships`); se agrega `due_date` opcional (`*time.Time`, omitempty). El constructor `NewAssessmentAssignedEvent` ahora exige `SubjectOfferingID` y `AssignedByMembershipID`. Los consumidores deben migrar al nuevo contrato.
- **Breaking** en `AssessmentAttemptRecordedPayload`: rename `student_id → student_membership_id` (re-llaveo de `auth.users.id` a `academic.memberships`, ADR 0019) y `total_points → max_score`. Los consumidores del payload (worker, learning) deben migrar al nuevo contrato.

### Added

- Campos `subject_id` (permite al worker resolver la oferta de la materia cross-schema) y `status` (gate `completed`/`pending_review`) en `AssessmentAttemptRecordedPayload`.

## [0.1.0] - 2026-05-28

### Added
- Reinicio de la versión del módulo a `v0.1.0` (borrón y cuenta nueva).
- Conservación del código de producción estable del módulo.

## [v0.52.0] - 2026-03-28

### Added

- Add optional TeacherID and Title fields to AssessmentAttemptRecordedPayload
- Add optional StudentID and Title fields to AssessmentReviewedPayload
- Add optional Title field to AssessmentAssignedPayload
- All fields use omitempty for backward compatibility
- Added serialization and backward-compatibility tests

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
