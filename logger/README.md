# Logger

Abstraccion de logging estructurado con implementaciones en Zap y Logrus.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/logger`
- Carpeta: `logger`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/logger
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Construir un logger Zap con nivel y formato configurables.
2. Adaptar una instancia Logrus a la interfaz comun del repositorio.
3. Agregar contexto con `With(fields...)` y reutilizarlo aguas abajo.
4. Emitir logs estructurados y sincronizar buffers cuando el backend lo necesita.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- El modulo no define politicas de logging; solo la abstraccion y sus implementaciones.
- Zap soporta formato `json` o `console`; Logrus se adapta a la interfaz comun.
- Tiene tests unitarios sobre niveles, campos y formatos.
