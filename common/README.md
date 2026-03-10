# Common

Modulo base del repositorio: centraliza primitives reutilizables y subpaquetes de bajo acoplamiento.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/common`
- Carpeta: `common`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/common
```

El modulo se descarga como `common`, pero el consumo real ocurre via subpaquetes como `common/errors`, `common/config`, `common/validator`, `common/types` y `common/types/enum`.

## Procesos documentados

1. Resolver variables de entorno y detectar ambiente con `common/config`.
2. Construir `AppError` tipados y mapearlos a status HTTP en `common/errors`.
3. Acumular errores de validacion y helpers de formato en `common/validator`.
4. Generar, parsear y serializar UUIDs en `common/types`.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- Consumir `common` implica importar subpaquetes concretos; no existe un package raiz unico para toda la API.
- Aqui estan varios contratos que otros modulos consideran fundacionales.
- Los tests cubren errores, validator, UUID y enums.
