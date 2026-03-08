# Lifecycle

Manager de ciclo de vida para recursos de infraestructura con startup secuencial y cleanup LIFO.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/lifecycle`
- Carpeta: `lifecycle`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/lifecycle
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Registrar recursos con nombre, startup opcional y cleanup.
2. Ejecutar la fase de startup en orden de registro.
3. Detener el proceso si un startup falla y reportar el recurso causante.
4. Ejecutar cleanups en orden inverso acumulando errores sin abortar la limpieza total.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- El modulo es pequeno y estable, orientado a coordinacion in-process.
- Su valor principal es el orden de cleanup y el agregado de errores.
- Tiene tests unitarios sobre startup, cleanup y LIFO.
