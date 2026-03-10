# Repository

Adaptadores GORM para CRUD y listados seguros sobre entidades externas de PostgreSQL.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/repository`
- Carpeta: `repository`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/repository
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Construir queries base por entidad (`User`, `School`, `Membership`) con contexto.
2. Aplicar filtros de actividad segun el tipo de repositorio.
3. Aplicar busqueda segura con `ILIKE` y escaping de patrones via `ListFilters`.
4. Aplicar paginacion y retornar el total antes de materializar la lista.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- La dependencia mas fuerte del modulo es hacia entidades externas, no hacia modulos internos de `edugo-shared`.
- Los tests encontrados cubren `ListFilters` y sus reglas de seguridad, no CRUD real sobre base de datos.
- Requerira especial atencion en fase 3 por estar fuera de la validacion raiz.
