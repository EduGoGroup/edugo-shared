# Audit

Contrato base para construir y despachar eventos de auditoria sin acoplar el almacenamiento.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/audit`
- Carpeta: `audit`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/audit
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Construir un `AuditEvent` con datos de actor, accion, recurso, request y metadata.
2. Enriquecer el evento con `AuditOption` como severidad, categoria, cambios, permisos o error.
3. Despachar el evento a traves de la interfaz `AuditLogger` sin fijar la implementacion concreta.
4. Usar `NoopAuditLogger` cuando se requiere un sink inerte para tests o escenarios locales.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- Este modulo documenta solo el contrato y el logger noop.
- La persistencia en PostgreSQL y la extraccion desde Gin se documentan aparte.
- El modulo tiene tests unitarios propios.
