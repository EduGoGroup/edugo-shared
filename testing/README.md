# Testing

Infraestructura de testing basada en Testcontainers para PostgreSQL, MongoDB y RabbitMQ, expuesta principalmente via el package `containers`.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/testing`
- Carpeta: `testing`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/testing
```

El modulo se descarga como `testing`, pero la API consumible esta concentrada en el package `testing/containers`.

## Procesos documentados

1. Construir una configuracion fluida con `ConfigBuilder`.
2. Obtener un `Manager` singleton y crear solo los containers habilitados.
3. Exponer accesos a containers de PostgreSQL, MongoDB y RabbitMQ para integracion.
4. Limpiar estado entre tests con truncado, drop de colecciones o purge de colas.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- La documentacion historica del modulo se reemplaza por esta version centrada en procesos y arquitectura.
- El consumo real ocurre sobre el package `containers`, no sobre un package raiz `testing` con la misma densidad de API.
- El modulo tiene tests unitarios e integracion con Docker.
