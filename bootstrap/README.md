# Bootstrap

Orquestador de inicializacion para logger, bases de datos, mensajeria y S3 mediante factories.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/bootstrap`
- Carpeta: `bootstrap`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/bootstrap
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Aplicar opciones de bootstrap y mezclar factories reales con mocks si corresponde.
2. Validar que existan factories para los recursos marcados como requeridos.
3. Inicializar primero el logger y luego PostgreSQL, MongoDB, RabbitMQ y S3.
4. Registrar recursos inicializados dentro de `Resources` y sus cleanups asociados.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- Es el modulo de mayor acoplamiento tecnico porque coordina recursos heterogeneos.
- El orden real observado es logger primero, luego bases de datos, mensajeria y finalmente storage.
- Tiene tests unitarios e integraciones para PostgreSQL, MongoDB, RabbitMQ y S3.
