# ScreenConfig

Utilidades para validar templates, resolver slots, aplicar overrides por plataforma y construir árboles de menú.

## Instalación

```bash
go get github.com/EduGoGroup/edugo-shared/screenconfig
```

El módulo se versionan y consume de forma independiente gracias a su `go.mod` propio.

## Quick Start

### Validar definición de template

```go
import (
    "github.com/EduGoGroup/edugo-shared/screenconfig"
)

templateDef := &screenconfig.ScreenTemplateDTO{
    ID:         "template-123",
    Name:       "UserProfile",
    ScreenType: "detail",
    Platform:   "ios",
    JSONDef:    json.RawMessage(`{"title": "User", "sections": []}`),
}

// Validar estructura completa
if err := screenconfig.ValidateTemplateDefinition(templateDef); err != nil {
    log.Printf("Template inválido: %v", err)
}
```

### Aplicar overrides por plataforma

```go
// Estructura base
template := &screenconfig.ScreenTemplateDTO{
    Platform:   "mobile", // Plataforma genérica
    JSONDef:    json.RawMessage(`{"color": "blue"}`),
    Overrides: map[string]json.RawMessage{
        "ios":     json.RawMessage(`{"color": "red"}`),
        "android": json.RawMessage(`{"color": "green"}`),
    },
}

// Aplicar override para plataforma específica
overridden := screenconfig.ApplyPlatformOverrides(template, "ios")
// overridden.JSONDef contiene el override de iOS
```

### Resolver placeholders de slots

```go
// Definición con placeholders
jsonDef := json.RawMessage(`{
    "title": "Dashboard",
    "content": "slot:dashboard_content",
    "sidebar": "slot:navigation"
}`)

// Resolver slots con diccionario
slots := map[string]json.RawMessage{
    "dashboard_content": json.RawMessage(`{"type": "widget", "name": "Analytics"}`),
    "navigation":       json.RawMessage(`{"type": "menu", "items": []}`),
}

resolved := screenconfig.ResolveSlots(jsonDef, slots)
// resolved contiene todos los placeholders reemplazados
```

### Construir árbol de menú jerárquico

```go
// Nodos de menú planos
nodes := []screenconfig.MenuNode{
    {ID: "root", ParentID: "", Label: "Menu"},
    {ID: "home", ParentID: "root", Label: "Home", ResourceKey: "screen:view"},
    {ID: "users", ParentID: "root", Label: "Users", ResourceKey: "users:list"},
    {ID: "settings", ParentID: "users", Label: "Settings", ResourceKey: "users:edit"},
}

// Construir árbol jerárquico
tree, err := screenconfig.BuildMenuTree(nodes)
if err != nil {
    log.Printf("Error constructing menu: %v", err)
}

// Filtrar por permisos
permissions := []string{"screen:view", "users:list"}
filtered := screenconfig.FilterMenuByPermissions(tree, permissions)
// filtered contiene solo nodos con permisos
```

## Componentes principales

- **ScreenTemplateDTO**: Definición de template con validación y overrides
- **MenuNode**: Nodo de menú con relaciones jerárquicas
- **ValidateTemplateDefinition**: Validación completa de estructura de template
- **ApplyPlatformOverrides**: Aplicación de overrides específicos de plataforma
- **ResolveSlots**: Resolución de placeholders slot:* en definiciones JSON
- **BuildMenuTree**: Construcción de árbol jerárquico a partir de nodos planos
- **FilterMenuByPermissions**: Filtrado de menú basado en permisos

## Documentación

- [Documentación técnica](docs/README.md)
- [Changelog](CHANGELOG.md)

## Operación local

```bash
make build    # Compilar módulo
make test     # Ejecutar tests
make test-race # Tests con race detector
make check    # Validar (fmt, vet, lint, test)
```

## Notas de diseño

- **Declarativo, no persistente**: Módulo transforma y valida, no persiste datos
- **JSON agnóstico**: Usa `json.RawMessage` para no acoplarse a estructuras de UI rígidas
- **Validaciones exhaustivas**: Patrones, tipos de screen, plataformas, definiciones JSON
- **Fallback de plataforma**: ios/android retroaceden a mobile
- **Funciones puras**: Todas las transformaciones son determinísticas sin side effects
