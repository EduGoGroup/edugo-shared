# Auth

Servicios compartidos de autenticacion: hashing de password, JWT de acceso y refresh tokens.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/auth`
- Carpeta: `auth`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/auth
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Hashear passwords con bcrypt costo 12 y verificar hashes en login.
2. Generar access tokens con `JWTManager` y un `UserContext` activo con rol y permisos.
3. Validar access tokens exigiendo issuer valido y `ActiveContext` presente.
4. Generar minimal tokens para refresh y distinguirlos por `TokenUse=refresh`.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- Los access tokens requieren `ActiveContext`; los refresh usan un camino minimal separado.
- El limite de password es 72 bytes por la restriccion propia de bcrypt.
- El modulo tiene una bateria amplia de tests unitarios y benchmarks.
