package textmatch

import (
	"errors"
	"reflect"
	"testing"
)

func TestRuleNamesOrderAndDefaults(t *testing.T) {
	wantNames := []string{RulePunctuation, RuleAccents, RuleNumberWords, RuleLeadingArticle, RuleGenderNumber}
	if got := RuleNames(); !reflect.DeepEqual(got, wantNames) {
		t.Fatalf("RuleNames() = %v, quiero el orden fijo %v", got, wantNames)
	}
	wantLiteral := []string{"punctuation", "accents", "number_words", "leading_article", "gender_number"}
	if !reflect.DeepEqual(wantNames, wantLiteral) {
		t.Fatalf("los nombres son contrato (answer_spec.rules): %v, quiero %v", wantNames, wantLiteral)
	}
	seen := map[string]bool{}
	for _, r := range Rules() {
		if seen[r.Name()] {
			t.Errorf("nombre repetido: %q", r.Name())
		}
		seen[r.Name()] = true
	}
	wantDefault := map[string]bool{
		RulePunctuation: true, RuleAccents: true, RuleNumberWords: true,
		RuleLeadingArticle: true, RuleGenderNumber: false,
	}
	for _, r := range Rules() {
		if r.DefaultEnabled() != wantDefault[r.Name()] {
			t.Errorf("%s nace %v, quiero %v", r.Name(), r.DefaultEnabled(), wantDefault[r.Name()])
		}
	}
}

func TestRulesReturnsACopy(t *testing.T) {
	rs := Rules()
	rs[0] = genderNumberRule{}
	if RuleNames()[0] != RulePunctuation {
		t.Fatal("modificar lo que devuelve Rules() cambia el registro")
	}
}

func TestLookupRule(t *testing.T) {
	for _, name := range RuleNames() {
		r, err := LookupRule(name)
		if err != nil || r.Name() != name {
			t.Errorf("LookupRule(%q) = %v, %v", name, r, err)
		}
	}
	_, err := LookupRule("plural")
	if !errors.Is(err, ErrUnknownRule) {
		t.Fatalf("una regla que no existe da %v, quiero ErrUnknownRule", err)
	}
}

func TestResolveRules(t *testing.T) {
	got, err := ResolveRules(nil)
	if err != nil {
		t.Fatalf("sin overrides: %v", err)
	}
	if got[RuleGenderNumber] || !got[RuleAccents] || len(got) != 5 {
		t.Errorf("valores por defecto = %v", got)
	}
	got, err = ResolveRules(map[string]bool{RuleGenderNumber: true, RuleAccents: false})
	if err != nil {
		t.Fatalf("con overrides: %v", err)
	}
	if !got[RuleGenderNumber] || got[RuleAccents] || !got[RulePunctuation] {
		t.Errorf("lo que la docente cambió no se aplica: %v", got)
	}
	if _, err := ResolveRules(map[string]bool{"disabled_rules": true}); !errors.Is(err, ErrUnknownRule) {
		t.Errorf("un nombre desconocido da %v, quiero ErrUnknownRule", err)
	}
}
