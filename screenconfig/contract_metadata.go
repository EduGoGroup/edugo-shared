package screenconfig

import "strings"

// ContractMetadata describe el "contrato" publico que el composer expone al
// frontend SDUI para que un GenericListContract / GenericFormContract pueda
// instanciarse sin codigo Kotlin especifico. Vive en shared para que cualquier
// API (platform, academic, learning, identity) reuse la misma proyeccion.
//
// El bloque se construye a partir de los campos opcionales que el
// screen_instance declara en su slot_data JSONB. La fuente de verdad de
// estos campos es la convencion descrita en
// tech-debt/sdui-refactor-spec/phase-3-data-driven-screens/design.md.
type ContractMetadata struct {
	APIPrefix     string         `json:"apiPrefix"`
	BasePath      string         `json:"basePath"`
	Resource      string         `json:"resource"`
	FormScreenKey *string        `json:"formScreenKey,omitempty"`
	ListScreenKey *string        `json:"listScreenKey,omitempty"`
	ParentIDParam *string        `json:"parentIdParam,omitempty"`
	Transforms    map[string]any `json:"transforms"`
}

// ExtractContractMetadata proyecta el contrato publico desde el slot_data del
// screen_instance. Devuelve nil cuando la metadata es insuficiente
// (sin api_prefix). El composer debe omitir el bloque "contract" del payload
// en ese caso.
//
// Reglas de derivacion (ver F3-REQ-1.2):
//   - apiPrefix      <- slot_data["api_prefix"]; si vacio/ausente => nil.
//   - resource       <- slot_data["resource"]; si vacio, derivado de
//     requiredPermission parseado como "prefix.resource.action" tomando el
//     segmento intermedio.
//   - basePath       <- slot_data["api_base_path"]; si vacio y resource != ""
//     usa la convencion "/api/v1/{resource}".
//   - formScreenKey, listScreenKey, parentIdParam <- campos string opcionales
//     del slot_data; nil cuando estan vacios/ausentes.
//   - transforms     <- slot_data["transforms"] casteado a map[string]any.
//     Si no es map, se devuelve un mapa vacio (no nil) para serializar "{}".
func ExtractContractMetadata(slotData map[string]any, requiredPermission string) *ContractMetadata {
	apiPrefix, _ := slotData["api_prefix"].(string)
	if apiPrefix == "" {
		return nil
	}

	resource, _ := slotData["resource"].(string)
	if resource == "" {
		resource = parseResourceFromPermission(requiredPermission)
	}

	basePath, _ := slotData["api_base_path"].(string)
	if basePath == "" && resource != "" {
		basePath = "/api/v1/" + resource
	}

	transforms, ok := slotData["transforms"].(map[string]any)
	if !ok {
		transforms = map[string]any{}
	}

	return &ContractMetadata{
		APIPrefix:     apiPrefix,
		BasePath:      basePath,
		Resource:      resource,
		FormScreenKey: optStringFromSlot(slotData, "form_screen_key"),
		ListScreenKey: optStringFromSlot(slotData, "list_screen_key"),
		ParentIDParam: optStringFromSlot(slotData, "parent_id_param"),
		Transforms:    transforms,
	}
}

// optStringFromSlot devuelve un *string apuntando al valor de la clave si
// existe y es un string no vacio; en caso contrario devuelve nil para que
// el campo se serialice con omitempty.
func optStringFromSlot(m map[string]any, k string) *string {
	s, ok := m[k].(string)
	if !ok || s == "" {
		return nil
	}
	return &s
}

// parseResourceFromPermission deriva el nombre del recurso a partir de una
// permission con forma "prefix.resource.action" (p. ej. "platform.colors.read"
// => "colors"). Si la permission no tiene la forma esperada, devuelve "".
//
// TODO(fase-4): esta funcion duplica deliberadamente la logica de
// edugo-api-platform/internal/core/usecase/screen_instance/compose.go:
// resourcePrefixFromPermission. Cuando la Fase 4 consolide los helpers SDUI
// en edugo-shared/screenconfig (ADR-5), este parser debe unificarse con el
// equivalente compartido y eliminar la duplicacion.
func parseResourceFromPermission(perm string) string {
	if perm == "" {
		return ""
	}
	_, rest, ok := strings.Cut(perm, ".")
	if !ok {
		return ""
	}
	resource, _, ok := strings.Cut(rest, ".")
	if !ok {
		return ""
	}
	return resource
}
