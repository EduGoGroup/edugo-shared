package auth

import (
	"reflect"
	"testing"
)

// TestMergeGrantChain_Union verifica que la unión de varios niveles de
// la cadena produce un único Grants con allow y deny combinados.
func TestMergeGrantChain_Union(t *testing.T) {
	parent := Grants{
		Allow: []string{"academic.units.*", "content.materials.*"},
		Deny:  []string{},
	}
	child := Grants{
		Allow: []string{"reports.read"},
		Deny:  []string{"*.delete"},
	}

	got := MergeGrantChain([]Grants{parent, child})
	want := Grants{
		Allow: []string{"academic.units.*", "content.materials.*", "reports.read"},
		Deny:  []string{"*.delete"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("union incorrecta:\n got=%+v\nwant=%+v", got, want)
	}
}

// TestMergeGrantChain_Dedup verifica que patterns repetidos entre
// niveles aparecen una sola vez (idempotencia de la unión).
func TestMergeGrantChain_Dedup(t *testing.T) {
	parent := Grants{
		Allow: []string{"academic.units.*", "reports.read"},
		Deny:  []string{"*.delete"},
	}
	child := Grants{
		Allow: []string{"reports.read", "menu.*"}, // reports.read duplicado
		Deny:  []string{"*.delete", "*.update"},   // *.delete duplicado
	}

	got := MergeGrantChain([]Grants{parent, child})
	want := Grants{
		Allow: []string{"academic.units.*", "reports.read", "menu.*"},
		Deny:  []string{"*.delete", "*.update"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("dedup incorrecto:\n got=%+v\nwant=%+v", got, want)
	}
}

// TestMergeGrantChain_Empty verifica que una cadena vacía o de un solo
// nivel se comporta como identidad y nunca devuelve nil slices.
func TestMergeGrantChain_Empty(t *testing.T) {
	got := MergeGrantChain(nil)
	if got.Allow == nil || got.Deny == nil {
		t.Fatalf("slices no inicializados: %+v", got)
	}
	if len(got.Allow) != 0 || len(got.Deny) != 0 {
		t.Fatalf("cadena vacía debería dar grants vacíos, got=%+v", got)
	}

	single := Grants{Allow: []string{"a.*"}, Deny: []string{"*.x"}}
	got = MergeGrantChain([]Grants{single})
	if !reflect.DeepEqual(got, single) {
		t.Fatalf("un solo nivel debería ser identidad:\n got=%+v\nwant=%+v", got, single)
	}
}
