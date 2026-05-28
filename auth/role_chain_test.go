package auth

import (
	"errors"
	"reflect"
	"testing"
)

// parentMap construye un parentOf de prueba a partir de un mapa
// hijo→padre. Un rol ausente del mapa se considera canónico (sin padre).
func parentMap(links map[string]string) func(string) (string, bool, error) {
	return func(id string) (string, bool, error) {
		p, ok := links[id]
		return p, ok, nil
	}
}

// TestResolveRoleChain_Canonical: un rol sin padre devuelve solo a sí mismo.
func TestResolveRoleChain_Canonical(t *testing.T) {
	got, err := ResolveRoleChain("teacher", parentMap(nil))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if want := []string{"teacher"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cadena incorrecta: got=%v want=%v", got, want)
	}
}

// TestResolveRoleChain_Depth1: alias → canónico (el caso real del seed).
func TestResolveRoleChain_Depth1(t *testing.T) {
	got, err := ResolveRoleChain("school_director", parentMap(map[string]string{
		"school_director": "school_admin",
	}))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if want := []string{"school_director", "school_admin"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cadena incorrecta: got=%v want=%v", got, want)
	}
}

// TestResolveRoleChain_DepthN: cadena de varios niveles, ordenada de
// más cercano a más lejano.
func TestResolveRoleChain_DepthN(t *testing.T) {
	got, err := ResolveRoleChain("a", parentMap(map[string]string{
		"a": "b",
		"b": "c",
		"c": "d",
	}))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if want := []string{"a", "b", "c", "d"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cadena incorrecta: got=%v want=%v", got, want)
	}
}

// TestResolveRoleChain_Cycle: un ciclo A→B→A corta sin error ni loop
// infinito, devolviendo los roles ya visitados una sola vez.
func TestResolveRoleChain_Cycle(t *testing.T) {
	got, err := ResolveRoleChain("a", parentMap(map[string]string{
		"a": "b",
		"b": "a",
	}))
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if want := []string{"a", "b"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cadena incorrecta: got=%v want=%v", got, want)
	}
}

// TestResolveRoleChain_EmptyParentStopsChain: ok=true con parent vacío
// equivale a "sin padre" (no agrega un eslabón vacío).
func TestResolveRoleChain_EmptyParentStopsChain(t *testing.T) {
	got, err := ResolveRoleChain("a", func(id string) (string, bool, error) {
		return "", true, nil // ok=true pero parent vacío
	})
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if want := []string{"a"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("cadena incorrecta: got=%v want=%v", got, want)
	}
}

// TestResolveRoleChain_PropagatesError: un error de parentOf se propaga
// y la cadena se descarta.
func TestResolveRoleChain_PropagatesError(t *testing.T) {
	sentinel := errors.New("fallo de BD")
	got, err := ResolveRoleChain("a", func(id string) (string, bool, error) {
		return "", false, sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("se esperaba el error centinela, got=%v", err)
	}
	if got != nil {
		t.Fatalf("se esperaba cadena nil ante error, got=%v", got)
	}
}
