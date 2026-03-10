# ScreenConfig

Utilidades para validar templates, resolver slots, aplicar overrides por plataforma y construir arboles de menu.

## Alcance

- Modulo Go: `github.com/EduGoGroup/edugo-shared/screenconfig`
- Carpeta: `screenconfig`
- Fase documental actual: `fase 1`, solo con evidencia local del repositorio

## Instalacion

```bash
go get github.com/EduGoGroup/edugo-shared/screenconfig
```

El modulo se puede versionar y consumir de forma independiente gracias a su `go.mod` propio.

## Procesos documentados

1. Validar patrones, screen types, plataformas y definiciones JSON de templates.
2. Aplicar overrides por plataforma con fallback `ios/android -> mobile`.
3. Resolver placeholders `slot:*` dentro de definiciones JSON.
4. Construir arboles de menu jerarquicos a partir de nodos planos.

## Navegacion

- [Documentacion del modulo](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operacion local

- `make build`
- `make test`
- `make check`
- Revisar `docs/README.md` para notas especificas de tests e integracion

## Notas actuales

- El foco del modulo es declarativo: transforma y valida, pero no persiste.
- El fallback de plataforma esta explicitamente codificado en `PlatformFallback`.
- Tiene tests unitarios sobre menu, permisos, overrides, slots y validacion.
