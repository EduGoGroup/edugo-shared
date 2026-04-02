# ScreenConfig — Documentación técnica

Utilidades para validar templates, resolver slots, aplicar overrides por plataforma y construir árboles de menú.

## Propósito

Proporcionar funciones puras de transformación y validación para:
- Validar integridad de definiciones JSON de templates
- Aplicar overrides específicos de plataforma (iOS, Android, Web, etc.)
- Resolver placeholders (slots) en definiciones JSON
- Construir árboles de menú jerárquicos desde nodos planos
- Filtrar menús basado en permisos del usuario

## Componentes principales

### ValidateTemplateDefinition — Validación completa

Valida que una definición de template sea estructuralmente correcta.

**Función:**
```go
func ValidateTemplateDefinition(template *ScreenTemplateDTO) error
```

**Valida:**
- Patrón de ID (formato válido)
- Tipo de screen (detail, list, form, etc.)
- Plataforma (ios, android, web, mobile, etc.)
- Estructura JSON de JSONDef (JSON válido)
- Estructura JSON de overrides por plataforma

**Retorna:**
- `nil` si la validación es exitosa
- `error` si algún campo es inválido

**Ejemplo:**
```go
template := &ScreenTemplateDTO{
    ID:         "user-detail",
    Name:       "User Detail Screen",
    ScreenType: "detail",
    Platform:   "mobile",
    JSONDef:    json.RawMessage(`{"title": "User", "fields": []}`),
}

if err := ValidateTemplateDefinition(template); err != nil {
    log.Printf("Validación fallida: %v", err)
}
```

### ApplyPlatformOverrides — Aplicar overrides por plataforma

Aplica overrides específicos de plataforma con fallback automático (ios/android → mobile).

**Función:**
```go
func ApplyPlatformOverrides(template *ScreenTemplateDTO, platform string) *ScreenTemplateDTO
```

**Comportamiento:**
- Si existe override para plataforma exacta: usar ese override
- Si no existe: intentar fallback según PlatformFallback
- Si tampoco: usar JSONDef original

**Fallback chart:**
- `ios` → `mobile` → `JSONDef`
- `android` → `mobile` → `JSONDef`
- `web` → `JSONDef`
- `mobile` → `JSONDef`

**Ejemplo:**
```go
template := &ScreenTemplateDTO{
    JSONDef: json.RawMessage(`{"color": "blue"}`),
    Overrides: map[string]json.RawMessage{
        "ios":    json.RawMessage(`{"color": "red"}`),
        "mobile": json.RawMessage(`{"color": "green"}`),
    },
}

// Para iOS: retorna override iOS
ios := ApplyPlatformOverrides(template, "ios")

// Para Web: retorna JSONDef original (sin override)
web := ApplyPlatformOverrides(template, "web")
```

### ResolveSlots — Resolver placeholders

Reemplaza placeholders `slot:*` en definiciones JSON con valores de un diccionario.

**Función:**
```go
func ResolveSlots(jsonDef json.RawMessage, slots map[string]json.RawMessage) json.RawMessage
```

**Comportamiento:**
- Busca strings que coincidan con patrón `slot:*`
- Reemplaza con valor del diccionario `slots`
- Si no existe key: mantiene placeholder original
- Funciona en profundidad (nested objects)

**Ejemplo:**
```go
jsonDef := json.RawMessage(`{
    "header": "slot:header_content",
    "body": {
        "main": "slot:main_content",
        "sidebar": "slot:sidebar"
    }
}`)

slots := map[string]json.RawMessage{
    "header_content": json.RawMessage(`{"title": "Dashboard"}`),
    "main_content":   json.RawMessage(`{"type": "widget"}`),
    // Note: sidebar no definido, se mantiene placeholder
}

resolved := ResolveSlots(jsonDef, slots)
// Resultado: header y main_content reemplazados, sidebar mantiene "slot:sidebar"
```

### BuildMenuTree — Construir árbol jerárquico

Construye árbol jerárquico de menú desde lista plana de nodos.

**Función:**
```go
func BuildMenuTree(nodes []MenuNode) (*MenuNode, error)
```

**Estructura MenuNode:**
```go
type MenuNode struct {
    ID          string      // Identificador único
    ParentID    string      // ID del nodo padre (vacío para raíz)
    Label       string      // Texto visible
    ResourceKey string      // Clave de permiso (resource:action)
    Children    []*MenuNode // Nodos hijos
}
```

**Comportamiento:**
- Encuentra raíz (ParentID vacío)
- Organiza nodos según relaciones padre-hijo
- Retorna árbol con estructura jerárquica
- Error si hay ciclos o referencias inválidas

**Ejemplo:**
```go
nodes := []MenuNode{
    {ID: "root", ParentID: "", Label: "Main Menu"},
    {ID: "products", ParentID: "root", Label: "Products", ResourceKey: "products:list"},
    {ID: "users", ParentID: "root", Label: "Users", ResourceKey: "users:list"},
    {ID: "user-detail", ParentID: "users", Label: "Detail", ResourceKey: "users:view"},
}

tree, err := BuildMenuTree(nodes)
// tree.Children contiene productos y usuarios
// tree.Children[1].Children contiene user-detail
```

### FilterMenuByPermissions — Filtrar por permisos

Filtra árbol de menú manteniendo solo nodos con permisos.

**Función:**
```go
func FilterMenuByPermissions(tree *MenuNode, permissions []string) *MenuNode
```

**Comportamiento:**
- Extrae recurso:acción de ResourceKey
- Compara contra lista de permisos
- Mantiene nodo si tiene permiso O si algún hijo tiene permisos
- Retorna árbol filtrado

**Ejemplo:**
```go
userPermissions := []string{
    "products:list",
    "products:view",
    "users:list", // Pero NO "users:edit"
}

filtered := FilterMenuByPermissions(tree, userPermissions)
// tree contiene solo nodos de productos y users:list
// Omite users:edit y sus hijos
```

## Flujos comunes

### 1. Validar y aplicar overrides a template

```go
func prepareTemplate(templateDef *ScreenTemplateDTO, userPlatform string) (*ScreenTemplateDTO, error) {
    // Validar estructura
    if err := ValidateTemplateDefinition(templateDef); err != nil {
        return nil, fmt.Errorf("template validation failed: %w", err)
    }

    // Aplicar override para plataforma del usuario
    overridden := ApplyPlatformOverrides(templateDef, userPlatform)

    return overridden, nil
}
```

### 2. Resolver slots en configuración de pantalla

```go
func buildScreenConfig(template *ScreenTemplateDTO, slotDefinitions map[string]json.RawMessage) (json.RawMessage, error) {
    // Validar
    if err := ValidateTemplateDefinition(template); err != nil {
        return nil, err
    }

    // Resolver todos los slots
    resolved := ResolveSlots(template.JSONDef, slotDefinitions)

    return resolved, nil
}
```

### 3. Construir menú con permisos

```go
func buildUserMenu(userPermissions []string) (*MenuNode, error) {
    // Nodos de menú base (ej: desde base de datos o configuración)
    allMenuNodes := getMenuNodes() // Implementar según necesidad

    // Construir árbol
    tree, err := BuildMenuTree(allMenuNodes)
    if err != nil {
        return nil, err
    }

    // Filtrar por permisos del usuario
    userMenu := FilterMenuByPermissions(tree, userPermissions)

    return userMenu, nil
}
```

### 4. Flujo completo: Template → Slots → Platform override

```go
func renderScreen(
    templateID string,
    userPlatform string,
    userPermissions []string,
) (json.RawMessage, error) {
    // 1. Cargar template
    template := loadTemplate(templateID)

    // 2. Validar
    if err := ValidateTemplateDefinition(template); err != nil {
        return nil, err
    }

    // 3. Resolver slots
    slotDefs := loadSlotDefinitions(template.SlotKeys())
    resolved := ResolveSlots(template.JSONDef, slotDefs)

    // 4. Aplicar override de plataforma
    template.JSONDef = resolved
    final := ApplyPlatformOverrides(template, userPlatform)

    // 5. Filtrar menú por permisos (si aplica)
    if menuTree := extractMenuFromConfig(final.JSONDef); menuTree != nil {
        filtered := FilterMenuByPermissions(menuTree, userPermissions)
        updateMenuInConfig(final.JSONDef, filtered)
    }

    return final.JSONDef, nil
}
```

## Arquitectura

Flujo de transformación de templates:

```
1. Cargar template JSON
   ↓
2. ValidateTemplateDefinition
   ├─ Validar patrón de ID
   ├─ Validar screen type
   ├─ Validar plataforma
   ├─ Validar JSON
   └─ Validar overrides
   ↓
3. ResolveSlots (opcional)
   ├─ Buscar placeholders slot:*
   ├─ Reemplazar con diccionario
   └─ Mantener placeholder si no existe
   ↓
4. ApplyPlatformOverrides
   ├─ Buscar override para plataforma
   ├─ Fallback a mobile si aplica
   └─ Usar JSONDef si no existe override
   ↓
5. Renderizar en cliente
```

Flujo de construcción de menú:

```
1. Obtener nodos planos
   ↓
2. BuildMenuTree
   ├─ Encontrar raíz (ParentID vacío)
   ├─ Organizar relaciones padre-hijo
   └─ Validar integridad
   ↓
3. ExtractResourceKeys (de ResourceKey)
   ↓
4. FilterMenuByPermissions
   ├─ Validar permisos
   ├─ Mantener si tiene permiso
   └─ Mantener si hijo tiene permiso
   ↓
5. Retornar menú filtrado
```

## Dependencias

- **Internas**: Ninguna (módulo autónomo)
- **Externas**: `encoding/json` (stdlib)

## Testing

Suite de tests completa:

- Validación de patrones, screen types, plataformas
- Validación de definiciones JSON
- Aplicación de overrides por plataforma
- Fallback de plataforma (ios/android → mobile)
- Resolución de slots
- Construcción de árbol de menú
- Filtrado por permisos
- Casos edge: slots no existentes, ciclos en árbol, permisos parciales

Ejecutar:
```bash
make test          # Tests básicos
make test-race     # Tests con race detector
make check         # Tests + linting + format
```

## Notas de diseño

- **Declarativo, no persistente**: Módulo transforma y valida, no almacena datos
- **JSON agnóstico**: Usa `json.RawMessage` para no acoplarse a estructuras de UI rígidas
- **Funciones puras**: Todas las transformaciones son determinísticas sin side effects
- **Validaciones exhaustivas**: Patrones de ID, tipos de screen, plataformas, JSON structure
- **Fallback explícito**: ios/android retroaceden a mobile según PlatformFallback
- **Sem lógica de negocio**: Módulo proporciona transformaciones genéricas sin reglas específicas del dominio
