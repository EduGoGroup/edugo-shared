# Cache Redis

Modulo pequeño para conectar a Redis por URL y exponer una cache JSON con TTL y borrado por patron.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/cache/redis`
- Carpeta: `cache/redis`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/cache/redis
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Parsear una URL Redis o Rediss y validar conectividad con `PING`.
2. Crear un `CacheService` backed por `go-redis`.
3. Serializar payloads a JSON para `Set` con TTL.
4. Deserializar JSON en `Get` y borrar entradas por clave o por patron usando `SCAN`.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- El alcance actual es deliberadamente pequeno: conexion, JSON y borrado basico.
- No se observaron adapters de metrics, locking o cache warming.
- El modulo tiene tests unitarios propios.
