# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/crypto/envelope` se registran aquí.

## [Unreleased]

## [0.900.0] - 2026-11-14

### Changed
- Directiva `go 1.27.4` (plan 054, homologación de Go y dependencias).
- Dependencias directas: `github.com/stretchr/testify` `v1.11.1` → `v1.12.1`, `golang.org/x/crypto` `v0.51.0` → `v0.59.0`.
- 1 dependencias indirectas suben con `go mod tidy` (ver `go.mod`).
- El módulo entra a la banda `v0.900.x` (ADR 0010); el tag anterior `v0.1.x` queda como procedencia.

## [0.1.0] - 2026-07-15

### Added

- Módulo nuevo `crypto/envelope` (plan 025 F0.1): promoción del envelope crypto que el spike F0.2 validó
  inline en `edugo-api-messaging`. Dos capas:
  - **Simétrica (AES-256-GCM):** `Envelope` con `NewEnvelope`, `Seal`, `Open`, `Overhead`; DEK de 32 bytes,
    nonce aleatorio de 12 bytes prefijado, overhead fijo de 28 bytes, autenticación por tag GCM.
    Constantes `DEKSize`/`Overhead`; errores `ErrKeySize`/`ErrBlobTooShort`.
  - **Asimétrica (sellado X25519 vía NaCl box anónimo):** `GenerateKeyPair`, `SealFor`, `OpenWith` para el
    modelo zero-knowledge (ADR 0029). Constantes `PublicKeySize`/`PrivateKeySize`; errores
    `ErrPublicKeySize`/`ErrPrivateKeySize`/`ErrOpenFailed`.
  - Tests de vectores: round-trip de ambas capas, fallo con clave equivocada, fallo ante manipulación,
    no-determinismo del sellado anónimo, tamaños de clave inválidos, y un caso cruzado del flujo completo
    sellar-DEK → abrir → descifrar.
