// Package screenconfig provee tipos y funciones para la configuracion de pantallas UI.
//
// # Tipos y Enums (types.go)
//
//   - [Pattern] — patrones de pantalla (login, form, list, dashboard, etc.)
//   - [ScreenType] — tipos de pantalla (list, detail, create, edit, etc.)
//   - [Platform] — plataformas soportadas (ios, android, mobile, desktop, web)
//
// # DTOs (dto.go)
//
//   - [ScreenTemplateDTO], [ScreenInstanceDTO], [CombinedScreenDTO]
//   - [ResourceScreenDTO], [NavigationItemDTO], [NavigationConfigDTO]
//
// # Menu y navegacion (menu_tree.go)
//
//   - [BuildMenuTree] — construye arbol jerarquico desde lista plana de [MenuNode]
//   - [MenuNode] — entrada generica de menu
//   - [MenuTreeItem] — nodo en el arbol construido
//
// # Permisos (permissions.go)
//
//   - [ExtractResourceKeys] — extrae claves de recursos desde permisos "resource:action"
//   - [HasPermission] — verifica existencia de un permiso
//
// # Overrides de plataforma (platform_overrides.go)
//
//   - [ApplyPlatformOverrides] — aplica overrides de zonas por plataforma con fallback
//
// # Slots (slots.go)
//
//   - [ResolveSlots] — reemplaza referencias "slot:xxx" con valores reales
//
// # Validacion (validation.go)
//
//   - [ValidatePattern], [ValidateScreenType], [ValidatePlatform]
//   - [ValidateTemplateDefinition] — valida estructura JSON de definitions
//   - [ResolvePlatformOverrideKey] — resuelve clave de plataforma con cadena de fallback
//   - [PlatformFallback] — mapa de fallback (ios->mobile, android->mobile)
package screenconfig
