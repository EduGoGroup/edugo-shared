# Audit PostgreSQL

Adaptador de auditoria que normaliza eventos y los persiste en `audit.events` usando GORM.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/audit/postgres`
- Carpeta: `audit/postgres`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/audit/postgres
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Recibir un `audit.AuditEvent` o extraerlo desde `gin.Context` por conveniencia.
2. Completar defaults de severidad, categoria y `serviceName` si el llamador no los define.
3. Transformar el evento al modelo `auditEventDB` con campos opcionales y JSON serializado.
4. Persistir el registro en la tabla `audit.events` con el contexto de la request.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- Es un adaptador especifico de almacenamiento, no el contrato general de auditoria.
- Asume que la tabla `audit.events` y sus serializers existen del lado de la base de datos.
- No se encontraron archivos `_test.go` propios en este modulo durante esta fase.
