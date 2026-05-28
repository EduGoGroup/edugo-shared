package screenconfig

import (
	"encoding/json"
	"fmt"
	"maps"
	"sort"
)

// ComposeActions resuelve la lista canonica de actions de una pantalla
// a partir del template + slot_data de la instancia.
//
// Pipeline:
//
//	defaults (templateDef.default_actions, con $resource$ resuelto)
//	  - actions_removed (por id, de slotData)
//	  + actions_added   (de slotData; overridea default si colisiona por id)
//	  -> ordenar por "order" ascendente (sort stable: ties por orden de aparicion)
//	  -> []map[string]any final
//
// Retrocompat: si slotData declara la lista legacy "actions" SIN
// added/removed, esa lista es override total (no aplica defaults). Si
// ademas hay actions_added/actions_removed, devuelve
// ErrActionsMixedWithAddedRemoved.
//
// $resource$ placeholder: cualquier valor string en un default que
// contenga "$resource$" se sustituye por el prefijo extraido de
// requiredPermission. Para "content.assessments.read" el prefijo es
// "content.assessments"; el placeholder solo aplica a la lista de
// defaults del template, no a entries de actions_added (que son
// especificas de la instancia y se declaran con el permission ya
// resuelto).
//
// El composer es funcion pura: no toca slotData; el caller es quien
// decide donde escribir el resultado y como limpiar los campos
// compositivos.
func ComposeActions(templateDef map[string]any, slotData map[string]any, requiredPermission string) ([]map[string]any, error) {
	legacyActions, hasLegacyActions := readActionList(slotData, "actions")
	added, hasAdded := readActionList(slotData, "actions_added")
	removed, hasRemoved := readStringList(slotData, "actions_removed")

	// Retrocompat: actions legacy sin added/removed = override total.
	if hasLegacyActions && !hasAdded && !hasRemoved {
		return sortStableByOrder(legacyActions), nil
	}

	// Mezcla: error explicito.
	if hasLegacyActions && (hasAdded || hasRemoved) {
		return nil, ErrActionsMixedWithAddedRemoved
	}

	// Defaults del template (puede ser nil si el template no declara).
	defaults, _ := readActionList(templateDef, "default_actions")
	resourcePrefix := ResourcePrefixFromPermission(requiredPermission)
	defaults = ExpandResourcePlaceholders(defaults, resourcePrefix)

	// Indice de removed.
	removedSet := make(map[string]struct{}, len(removed))
	for _, r := range removed {
		removedSet[r] = struct{}{}
	}

	// Indice de added por id (para detectar overrides).
	addedByID := make(map[string]map[string]any, len(added))
	for _, a := range added {
		if id, ok := a["id"].(string); ok && id != "" {
			addedByID[id] = a
		}
	}

	// Construir lista final: defaults filtrados por removed y por
	// override-de-added; luego concatenar added en su orden de declaracion.
	out := make([]map[string]any, 0, len(defaults)+len(added))
	for _, def := range defaults {
		id, _ := def["id"].(string)
		if _, isRemoved := removedSet[id]; isRemoved {
			continue
		}
		if _, isOverridden := addedByID[id]; isOverridden {
			// Se materializara desde el bloque added.
			continue
		}
		out = append(out, def)
	}
	for _, a := range added {
		out = append(out, a)
	}

	return sortStableByOrder(out), nil
}

// ComposeActionsForResolve aplica el composer SDUI a slot_data. Tolerante:
// si el slot_data o template no son JSON validos, retorna slot_data sin
// modificar (pantallas legacy siguen funcionando).
//
// Devuelve tambien el slot_data como map[string]any (ya compuesto, sin
// actions_added/actions_removed) para que callers puedan extraer metadata
// adicional sin re-unmarshalear. El map puede ser nil si el slot_data no
// es JSON valido o esta vacio; los callers deben tolerarlo (ExtractContractMetadata
// retorna nil ante map nil porque "api_prefix" no esta presente).
func ComposeActionsForResolve(slotDataRaw, templateDef json.RawMessage, requiredPerm string) (json.RawMessage, map[string]any) {
	if len(slotDataRaw) == 0 {
		return slotDataRaw, nil
	}
	var slot map[string]any
	if err := json.Unmarshal(slotDataRaw, &slot); err != nil || slot == nil {
		return slotDataRaw, nil
	}
	var tplDef map[string]any
	if len(templateDef) > 0 {
		_ = json.Unmarshal(templateDef, &tplDef)
	}
	actions, err := ComposeActions(tplDef, slot, requiredPerm)
	if err != nil {
		return slotDataRaw, slot
	}
	slot["actions"] = actions
	delete(slot, "actions_added")
	delete(slot, "actions_removed")
	encoded, err := json.Marshal(slot)
	if err != nil {
		return slotDataRaw, slot
	}
	return encoded, slot
}

// readActionList lee una lista de actions ([]map[string]any) bajo la
// clave dada. Retorna (lista, true) si la clave existe y el valor es un
// slice; (nil, false) en cualquier otro caso (clave ausente, nil,
// type-mismatch). El composer asume JSON ya decodificado, donde los
// arrays vienen como []any.
func readActionList(m map[string]any, key string) ([]map[string]any, bool) {
	if m == nil {
		return nil, false
	}
	raw, ok := m[key]
	if !ok || raw == nil {
		return nil, false
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, false
	}
	out := make([]map[string]any, 0, len(arr))
	for _, item := range arr {
		if obj, ok := item.(map[string]any); ok {
			// Copia superficial para no mutar el input del template.
			cp := make(map[string]any, len(obj))
			maps.Copy(cp, obj)
			out = append(out, cp)
		}
	}
	return out, true
}

// readStringList lee una lista de strings bajo la clave dada. Mismo
// contrato que readActionList: distingue ausente de vacia.
func readStringList(m map[string]any, key string) ([]string, bool) {
	if m == nil {
		return nil, false
	}
	raw, ok := m[key]
	if !ok || raw == nil {
		return nil, false
	}
	arr, ok := raw.([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok {
			out = append(out, s)
		}
	}
	return out, true
}

// sortStableByOrder ordena las actions ascendentemente por "order"
// (numero entero / float). Sort stable garantiza que los ties respetan
// el orden de aparicion (defaults primero, added despues). Actions sin
// "order" toman 0 — el seed deberia declararlo siempre, pero ser
// defensivo evita panics en pantallas viejas.
func sortStableByOrder(actions []map[string]any) []map[string]any {
	type indexed struct {
		idx   int
		order int
		val   map[string]any
	}
	rows := make([]indexed, len(actions))
	for i, a := range actions {
		rows[i] = indexed{idx: i, order: orderOf(a), val: a}
	}
	sort.SliceStable(rows, func(i, j int) bool {
		return rows[i].order < rows[j].order
	})
	out := make([]map[string]any, len(rows))
	for i, r := range rows {
		out[i] = r.val
	}
	return out
}

// orderOf extrae el campo "order" como int. Acepta float64 (JSON
// numbers vienen como float tras encoding/json.Unmarshal), int y int64
// por defensividad ante callers que construyan el map a mano en tests.
// Default 0 cuando ausente o de tipo no numerico.
func orderOf(a map[string]any) int {
	raw, ok := a["order"]
	if !ok || raw == nil {
		return 0
	}
	switch v := raw.(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	case int32:
		return int(v)
	default:
		// Tipo inesperado — log silencioso (composer es pure function);
		// no panic. El test catches assertion lo destapa.
		_ = fmt.Sprint(v)
		return 0
	}
}
