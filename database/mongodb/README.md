# Database MongoDB

Modulo de bajo nivel para abrir, validar y cerrar conexiones MongoDB.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/database/mongodb`
- Carpeta: `database/mongodb`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/database/mongodb
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Construir `Config` con defaults de URI, database, timeout y tamanos de pool.
2. Configurar `ClientOptions` con URI, pool y timeouts.
3. Conectar el cliente y validar conectividad contra `readpref.Primary()`.
4. Exponer acceso a una database concreta con `GetDatabase`.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- El modulo trabaja a nivel de cliente y database, no de repositorios o DAOs.
- La validacion de salud usa `Ping` al primary.
- Existe cobertura combinada de tests unitarios e integracion.
