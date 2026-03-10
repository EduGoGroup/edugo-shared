# Database PostgreSQL

Capa de bajo nivel para sql.DB, GORM y wrappers de transaccion sobre PostgreSQL.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/database/postgres`
- Carpeta: `database/postgres`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/database/postgres
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Construir `Config` y un DSN con `search_path`, SSL y timeout.
2. Abrir `sql.DB`, configurar pool y verificar conectividad con `PingContext`.
3. Crear una conexion GORM equivalente con `ConnectGORM`.
4. Ejecutar health checks, cierre y lectura de stats del pool.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- El modulo ya documentaba testing de integracion, pero la nueva documentacion lo reubica dentro del esquema por fases.
- La API se mantiene en nivel de infraestructura y transacciones.
- Tiene tests unitarios e integracion con buena densidad sobre conexion y transacciones.
