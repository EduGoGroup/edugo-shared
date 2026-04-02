# Changelog

Todos los cambios relevantes de `github.com/EduGoGroup/edugo-shared/screenconfig` se registran aquí.

## [0.100.0] - 2026-04-02

### Added

- **ValidateTemplateDefinition**: Validación completa de definiciones de template (ID, screen type, plataforma, JSON, overrides).
- **ValidatePattern**: Validación de patrón de ID con regex.
- **ValidateScreenType**: Validación de tipo de screen permitido.
- **ValidatePlatform**: Validación de plataforma permitida.
- **ApplyPlatformOverrides**: Aplicación de overrides específicos de plataforma con fallback automático.
- **PlatformFallback**: Mapa de fallback para plataformas (ios/android → mobile).
- **ResolvePlatformOverrideKey**: Resolución de clave de override con fallback.
- **ResolveSlots**: Resolución de placeholders `slot:*` en definiciones JSON.
- **BuildMenuTree**: Construcción de árbol jerárquico desde nodos planos.
- **MenuNode**: Estructura de nodo de menú con relaciones padre-hijo.
- **FilterMenuByPermissions**: Filtrado de árbol de menú basado en permisos del usuario.
- **ExtractResourceKeys**: Extracción de claves de recurso desde ResourceKey (formato resource:action).
- **ScreenTemplateDTO**: DTO para definición de template con validación y overrides.
- **ScreenInstanceDTO**: DTO para instancia de pantalla renderizada.
- Suite completa de tests unitarios cobriendo validación, overrides, slots, árbol de menú y permisos.
- Documentación técnica detallada en docs/README.md con componentes, flujos comunes y arquitectura.
- Makefile con targets: build, test, test-race, check, lint, fmt, vet, tidy, deps, release.

### Design Notes

- Declarativo, no persistente: módulo transforma y valida, no almacena datos.
- JSON agnóstico: usa `json.RawMessage` para no acoplarse a estructuras de UI rígidas.
- Funciones puras: todas las transformaciones son determinísticas sin side effects.
- Validaciones exhaustivas: patrones, tipos, plataformas, JSON structure.
- Fallback de plataforma: ios/android retroaceden a mobile automáticamente.
- Sin lógica de negocio: proporciona transformaciones genéricas sin reglas específicas del dominio.

## [0.1.0] - 2026-03-26

### Added

- Baseline de documentación de fase 1 con `README.md` y `docs/README.md`.
