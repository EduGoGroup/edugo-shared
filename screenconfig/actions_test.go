package screenconfig_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/EduGoGroup/edugo-shared/screenconfig"
)

// jsonRound serializa y deserializa con encoding/json para emular el
// shape real que vera ComposeActions en runtime: los []any/map[string]any
// que produce json.Unmarshal, no maps literales que hash distinto.
func jsonRound(t *testing.T, v any) map[string]any {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestComposeActions_LegacyTotalOverride(t *testing.T) {
	// Instancia legacy (announcement-form style): declara "actions"
	// sin actions_added/actions_removed. El composer debe devolver esa
	// lista tal cual (sort por order), ignorando los defaults del
	// template.
	tpl := jsonRound(t, map[string]any{
		"default_actions": []any{
			map[string]any{"id": "save_new", "scope": "form-submit", "order": 10},
			map[string]any{"id": "save", "scope": "form-submit", "order": 10},
			map[string]any{"id": "delete", "scope": "form-submit", "order": 20},
		},
	})
	slot := jsonRound(t, map[string]any{
		"actions": []any{
			map[string]any{"id": "save", "label": "Guardar", "permission": "academic.announcements.update", "order": 10},
			map[string]any{"id": "delete", "label": "Eliminar", "permission": "academic.announcements.delete", "order": 20},
		},
	})

	got, err := screenconfig.ComposeActions(tpl, slot, "academic.announcements.read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2 (legacy override)", len(got))
	}
	if got[0]["id"] != "save" || got[1]["id"] != "delete" {
		t.Fatalf("expected [save,delete], got [%v,%v]", got[0]["id"], got[1]["id"])
	}
}

func TestComposeActions_DefaultsOnly(t *testing.T) {
	// Sin actions/actions_added/actions_removed en la instancia: el
	// composer materializa los defaults del template y resuelve
	// $resource$.
	tpl := jsonRound(t, map[string]any{
		"default_actions": []any{
			map[string]any{
				"id":         "save_new",
				"scope":      "form-submit",
				"label":      "Guardar",
				"permission": "$resource$.create",
				"condition":  "create-only",
				"order":      10,
			},
			map[string]any{
				"id":         "delete",
				"scope":      "form-submit",
				"label":      "Eliminar",
				"permission": "$resource$.delete",
				"condition":  "edit-only",
				"order":      20,
			},
		},
	})
	slot := jsonRound(t, map[string]any{})

	got, err := screenconfig.ComposeActions(tpl, slot, "content.assessments.read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("len=%d, want 2", len(got))
	}
	if got[0]["permission"] != "content.assessments.create" {
		t.Errorf("save_new.permission = %v, want content.assessments.create", got[0]["permission"])
	}
	if got[1]["permission"] != "content.assessments.delete" {
		t.Errorf("delete.permission = %v, want content.assessments.delete", got[1]["permission"])
	}
}

func TestComposeActions_AddedOverridesDefaultByID(t *testing.T) {
	// added con id que colisiona con un default → el default se
	// reemplaza por el added (no se dupliza).
	tpl := jsonRound(t, map[string]any{
		"default_actions": []any{
			map[string]any{"id": "detail", "scope": "resource-toolbar", "label": "Detalle", "icon": "list", "order": 10},
		},
	})
	slot := jsonRound(t, map[string]any{
		"actions_added": []any{
			map[string]any{"id": "detail", "scope": "resource-toolbar", "label": "Preguntas", "icon": "help_outline", "event_id": "view-questions", "order": 15},
		},
	})

	got, err := screenconfig.ComposeActions(tpl, slot, "content.assessments.read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1 (override por id)", len(got))
	}
	if got[0]["label"] != "Preguntas" {
		t.Errorf("expected label=Preguntas (override), got %v", got[0]["label"])
	}
	if got[0]["icon"] != "help_outline" {
		t.Errorf("expected icon override, got %v", got[0]["icon"])
	}
}

func TestComposeActions_RemovedDropsDefault(t *testing.T) {
	// actions_removed: ["delete"] elimina el default por id.
	tpl := jsonRound(t, map[string]any{
		"default_actions": []any{
			map[string]any{"id": "save", "scope": "form-submit", "order": 10},
			map[string]any{"id": "delete", "scope": "form-submit", "order": 20},
		},
	})
	slot := jsonRound(t, map[string]any{
		"actions_removed": []any{"delete"},
	})

	got, err := screenconfig.ComposeActions(tpl, slot, "content.materials.read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len=%d, want 1 (delete removido)", len(got))
	}
	if got[0]["id"] != "save" {
		t.Fatalf("expected [save], got %v", got[0]["id"])
	}
}

func TestComposeActions_MixedActionsAndAdded_Error(t *testing.T) {
	tpl := jsonRound(t, map[string]any{
		"default_actions": []any{
			map[string]any{"id": "save", "order": 10},
		},
	})
	slot := jsonRound(t, map[string]any{
		"actions":       []any{map[string]any{"id": "save", "order": 10}},
		"actions_added": []any{map[string]any{"id": "detail", "order": 15}},
	})

	_, err := screenconfig.ComposeActions(tpl, slot, "content.assessments.read")
	if !errors.Is(err, screenconfig.ErrActionsMixedWithAddedRemoved) {
		t.Fatalf("expected ErrActionsMixedWithAddedRemoved, got %v", err)
	}
}

func TestComposeActions_PlaceholderExpansion(t *testing.T) {
	tpl := jsonRound(t, map[string]any{
		"default_actions": []any{
			map[string]any{
				"id":         "save_new",
				"permission": "$resource$.create",
				"order":      10,
			},
			map[string]any{
				"id":         "save",
				"permission": "$resource$.update",
				"order":      10,
			},
			map[string]any{
				"id":         "delete",
				"permission": "$resource$.delete",
				"order":      20,
			},
		},
	})
	slot := jsonRound(t, map[string]any{})

	got, err := screenconfig.ComposeActions(tpl, slot, "academic.announcements.read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 3 {
		t.Fatalf("len=%d, want 3", len(got))
	}
	expected := map[string]string{
		"save_new": "academic.announcements.create",
		"save":     "academic.announcements.update",
		"delete":   "academic.announcements.delete",
	}
	for _, a := range got {
		id := a["id"].(string)
		if a["permission"] != expected[id] {
			t.Errorf("%s.permission = %v, want %v", id, a["permission"], expected[id])
		}
	}
}

func TestComposeActions_SortByOrderStable(t *testing.T) {
	// orden mezclado; ties por orden de aparicion (defaults antes que added).
	tpl := jsonRound(t, map[string]any{
		"default_actions": []any{
			map[string]any{"id": "save", "order": 10},
			map[string]any{"id": "delete", "order": 30},
		},
	})
	slot := jsonRound(t, map[string]any{
		"actions_added": []any{
			map[string]any{"id": "publish", "order": 20},
			map[string]any{"id": "archive", "order": 20}, // tie con publish → stable: publish primero
			map[string]any{"id": "detail", "order": 10},  // tie con save → stable: save primero (defaults antes que added)
		},
	})

	got, err := screenconfig.ComposeActions(tpl, slot, "content.assessments.read")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	gotIDs := make([]string, len(got))
	for i, a := range got {
		gotIDs[i] = a["id"].(string)
	}
	want := []string{"save", "detail", "publish", "archive", "delete"}
	if len(gotIDs) != len(want) {
		t.Fatalf("len=%d, want %d. Got: %v", len(gotIDs), len(want), gotIDs)
	}
	for i := range want {
		if gotIDs[i] != want[i] {
			t.Errorf("at %d: got %s, want %s. Full: %v", i, gotIDs[i], want[i], gotIDs)
		}
	}
}

func TestComposeActions_AddedOnlyNoDefaults(t *testing.T) {
	// Template sin default_actions: el composer devuelve solo los added.
	// Edge case: alguien sembro un template viejo sin defaults; la
	// instancia agrega los suyos via actions_added.
	tpl := jsonRound(t, map[string]any{})
	slot := jsonRound(t, map[string]any{
		"actions_added": []any{
			map[string]any{"id": "create", "order": 10},
		},
	})

	got, err := screenconfig.ComposeActions(tpl, slot, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0]["id"] != "create" {
		t.Fatalf("expected [create], got %v", got)
	}
}

func TestComposeActions_EmptyEverything(t *testing.T) {
	// Sin defaults ni slot data: lista vacia, sin error.
	got, err := screenconfig.ComposeActions(map[string]any{}, map[string]any{}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty, got %v", got)
	}
}

// --- Tests para el wrapper ComposeActionsForResolve ---

// mustMarshal serializa v a json.RawMessage y aborta el test en error. Util
// para construir slotDataRaw/templateDef sin escribir strings JSON a mano.
func mustMarshal(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return json.RawMessage(b)
}

func TestComposeActionsForResolve_HappyPath(t *testing.T) {
	// slot con actions_added; el wrapper compone, elimina los campos
	// compositivos y devuelve un map ya resuelto + JSON encoded.
	tplDef := mustMarshal(t, map[string]any{
		"default_actions": []any{
			map[string]any{"id": "save", "scope": "form-submit", "order": 10},
		},
	})
	slotDataRaw := mustMarshal(t, map[string]any{
		"actions_added": []any{
			map[string]any{"id": "detail", "label": "Detalle", "order": 20},
		},
		"actions_removed": []any{},
		"some_other":      "preserved",
	})

	encoded, slotMap := screenconfig.ComposeActionsForResolve(slotDataRaw, tplDef, "content.assessments.read")

	if slotMap == nil {
		t.Fatalf("slotMap is nil; expected resolved map")
	}
	// El map ya no contiene actions_added/actions_removed.
	if _, ok := slotMap["actions_added"]; ok {
		t.Errorf("slotMap should not contain actions_added")
	}
	if _, ok := slotMap["actions_removed"]; ok {
		t.Errorf("slotMap should not contain actions_removed")
	}
	if _, ok := slotMap["actions"]; !ok {
		t.Errorf("slotMap should contain actions")
	}
	if slotMap["some_other"] != "preserved" {
		t.Errorf("expected some_other preserved, got %v", slotMap["some_other"])
	}

	// El JSON encoded refleja lo mismo.
	var decoded map[string]any
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("encoded JSON invalid: %v", err)
	}
	if _, ok := decoded["actions"]; !ok {
		t.Errorf("encoded JSON missing actions key")
	}
	if _, ok := decoded["actions_added"]; ok {
		t.Errorf("encoded JSON should not contain actions_added")
	}
	if _, ok := decoded["actions_removed"]; ok {
		t.Errorf("encoded JSON should not contain actions_removed")
	}
	actions, ok := decoded["actions"].([]any)
	if !ok {
		t.Fatalf("actions is not a []any, got %T", decoded["actions"])
	}
	if len(actions) != 2 {
		t.Fatalf("expected 2 actions composed (default + added), got %d", len(actions))
	}
}

func TestComposeActionsForResolve_EmptySlotData(t *testing.T) {
	// slotDataRaw vacio (nil y len 0) → devuelve input intacto, map nil,
	// sin panic.
	tplDef := mustMarshal(t, map[string]any{
		"default_actions": []any{
			map[string]any{"id": "save", "order": 10},
		},
	})

	// Caso nil.
	encoded, slotMap := screenconfig.ComposeActionsForResolve(nil, tplDef, "content.x.read")
	if encoded != nil {
		t.Errorf("expected nil encoded, got %v", encoded)
	}
	if slotMap != nil {
		t.Errorf("expected nil slotMap, got %v", slotMap)
	}

	// Caso len==0.
	empty := json.RawMessage{}
	encoded, slotMap = screenconfig.ComposeActionsForResolve(empty, tplDef, "content.x.read")
	if len(encoded) != 0 {
		t.Errorf("expected empty encoded, got %v", encoded)
	}
	if slotMap != nil {
		t.Errorf("expected nil slotMap for empty input, got %v", slotMap)
	}
}

func TestComposeActionsForResolve_InvalidJSON(t *testing.T) {
	// slotDataRaw con bytes no-JSON → devuelve input intacto, map nil,
	// sin panic.
	tplDef := mustMarshal(t, map[string]any{
		"default_actions": []any{
			map[string]any{"id": "save", "order": 10},
		},
	})
	bad := json.RawMessage([]byte("{not valid json"))

	encoded, slotMap := screenconfig.ComposeActionsForResolve(bad, tplDef, "content.x.read")
	if string(encoded) != string(bad) {
		t.Errorf("expected encoded == input on invalid JSON, got %s", string(encoded))
	}
	if slotMap != nil {
		t.Errorf("expected nil slotMap on invalid JSON, got %v", slotMap)
	}
}

func TestComposeActionsForResolve_ComposerError_LiteralMixed(t *testing.T) {
	// slot con "actions" literal + "actions_added" → ComposeActions
	// retorna ErrActionsMixedWithAddedRemoved. El wrapper debe devolver
	// el slotDataRaw original (no encodea slot) y el map ya parseado
	// (no nil).
	tplDef := mustMarshal(t, map[string]any{
		"default_actions": []any{
			map[string]any{"id": "save", "order": 10},
		},
	})
	slotDataRaw := mustMarshal(t, map[string]any{
		"actions":       []any{map[string]any{"id": "save", "order": 10}},
		"actions_added": []any{map[string]any{"id": "detail", "order": 15}},
	})

	encoded, slotMap := screenconfig.ComposeActionsForResolve(slotDataRaw, tplDef, "content.assessments.read")

	if string(encoded) != string(slotDataRaw) {
		t.Errorf("expected encoded == slotDataRaw on composer error, got %s vs %s", string(encoded), string(slotDataRaw))
	}
	if slotMap == nil {
		t.Fatalf("expected non-nil slotMap (parsed before composer error)")
	}
	// El map preserva ambas claves (no se limpiaron porque hubo error).
	if _, ok := slotMap["actions"]; !ok {
		t.Errorf("slotMap should still contain actions (composer error path)")
	}
	if _, ok := slotMap["actions_added"]; !ok {
		t.Errorf("slotMap should still contain actions_added (composer error path)")
	}
}
