package auth

import (
	"encoding/json"
	"os"
	"testing"
)

// goldenMatcherCase es un caso del fixture cross-language para
// PermissionMatches. El fixture canónico vive en
// EduUI/edugo-ui-kmp/e2e-integration-plan/permissions-redesign-spec/
// fixtures/permission_matcher_golden.json y esta copia local
// (testdata/permission_matcher_golden.json) debe ser byte-idéntica.
type goldenMatcherCase struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Pattern  string `json:"pattern"`
	Request  string `json:"request"`
	Expected bool   `json:"expected"`
}

// goldenGrantsCase es un caso del fixture cross-language para
// EvaluateGrants.
type goldenGrantsCase struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Allow    []string `json:"allow"`
	Deny     []string `json:"deny"`
	Request  string   `json:"request"`
	Expected bool     `json:"expected"`
}

type goldenFile struct {
	Version      string              `json:"version"`
	MatcherCases []goldenMatcherCase `json:"matcher_cases"`
	GrantsCases  []goldenGrantsCase  `json:"grants_cases"`
}

func loadGolden(t *testing.T) goldenFile {
	t.Helper()
	raw, err := os.ReadFile("testdata/permission_matcher_golden.json")
	if err != nil {
		t.Fatalf("no se pudo leer el fixture golden: %v", err)
	}
	var gf goldenFile
	if err := json.Unmarshal(raw, &gf); err != nil {
		t.Fatalf("no se pudo deserializar el fixture golden: %v", err)
	}
	if len(gf.MatcherCases) == 0 {
		t.Fatalf("fixture golden sin matcher_cases")
	}
	if len(gf.GrantsCases) == 0 {
		t.Fatalf("fixture golden sin grants_cases")
	}
	return gf
}

// TestPermissionMatcherGolden valida que auth.PermissionMatches
// coincida 1:1 con el fixture cross-language.
func TestPermissionMatcherGolden(t *testing.T) {
	gf := loadGolden(t)
	for _, c := range gf.MatcherCases {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			got := PermissionMatches(c.Pattern, c.Request)
			if got != c.Expected {
				t.Errorf("[%s] %s: pattern=%q request=%q want=%v got=%v",
					c.ID, c.Name, c.Pattern, c.Request, c.Expected, got)
			}
		})
	}
}

// TestEvaluateGrantsGolden valida que auth.EvaluateGrants coincida
// 1:1 con el fixture cross-language.
func TestEvaluateGrantsGolden(t *testing.T) {
	gf := loadGolden(t)
	for _, c := range gf.GrantsCases {
		c := c
		t.Run(c.ID, func(t *testing.T) {
			g := Grants{Allow: c.Allow, Deny: c.Deny}
			got := EvaluateGrants(g, c.Request)
			if got != c.Expected {
				t.Errorf("[%s] %s: allow=%v deny=%v request=%q want=%v got=%v",
					c.ID, c.Name, c.Allow, c.Deny, c.Request, c.Expected, got)
			}
		})
	}
}
