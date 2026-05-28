# Messaging Events

Define eventos de dominio serializables, hoy centrados en `MaterialUploadedEvent`.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/messaging/events`
- Carpeta: `messaging/events`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/messaging/events
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Validar campos requeridos del evento y del payload.
2. Estampar `Timestamp` automaticamente al construir el evento.
3. Serializar la estructura resultante a JSON para su transporte.
4. Resolver campos de compatibilidad legacy como `GetS3Key`, `GetMaterialID` y `GetAuthorID`.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- Hoy el modulo expone un evento de dominio principal; no se observaron otros schemas aun.
- El transporte o delivery queda en otros modulos.
- Tiene README historico y tests unitarios propios, ahora absorbidos por la nueva estructura.
