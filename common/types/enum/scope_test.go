package enum

import "testing"

func TestScope_String(t *testing.T) {
	if got := ScopeNotificationsDispatch.String(); got != "notifications.dispatch" {
		t.Errorf("Scope.String() = %q, want %q", got, "notifications.dispatch")
	}
}

func TestScope_IsValid(t *testing.T) {
	if !ScopeNotificationsDispatch.IsValid() {
		t.Errorf("ScopeNotificationsDispatch debería ser válido")
	}
	if Scope("unknown.scope").IsValid() {
		t.Errorf("un scope desconocido no debería ser válido")
	}
}

// TestAllScopes_MapIntegrity verifica que cada Scope declarado esté en el
// catálogo cerrado AllScopes (mismo contrato que AllPermissions).
func TestAllScopes_MapIntegrity(t *testing.T) {
	declared := []Scope{
		ScopeNotificationsDispatch,
	}
	for _, s := range declared {
		if !AllScopes[s] {
			t.Errorf("scope %q declarado pero ausente en AllScopes", s)
		}
	}
	if len(AllScopes) != len(declared) {
		t.Errorf("AllScopes tiene %d entradas, se declararon %d scopes", len(AllScopes), len(declared))
	}
}
