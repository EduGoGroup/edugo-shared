# Middleware Gin

Middleware HTTP para validar JWT, poblar contexto, verificar permisos y registrar auditoria.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/middleware/gin`
- Carpeta: `middleware/gin`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/middleware/gin
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Validar el header `Authorization: Bearer ...` y el JWT usando `auth.JWTManager`.
2. Poblar el `gin.Context` con `user_id`, `email`, `role` y `jwt_claims`.
3. Aplicar middlewares de permiso simple, any-of o all-of segun permisos del contexto activo.
4. Registrar eventos de auditoria para metodos mutantes y derivar recurso e identificador desde la ruta.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- Este modulo asume que los claims validos incluyen `ActiveContext` con permisos.
- El middleware de auditoria ignora requests de solo lectura.
- Tiene tests unitarios sobre auth, permisos, contexto y auditoria.
