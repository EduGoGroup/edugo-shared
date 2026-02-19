package screenconfig

import (
	"encoding/json"
	"fmt"
)

var validPatterns = map[Pattern]bool{
	PatternLogin:        true,
	PatternForm:         true,
	PatternList:         true,
	PatternDashboard:    true,
	PatternSettings:     true,
	PatternDetail:       true,
	PatternSearch:       true,
	PatternProfile:      true,
	PatternModal:        true,
	PatternNotification: true,
	PatternOnboarding:   true,
	PatternEmptyState:   true,
}

var validActionTypes = map[ActionType]bool{
	ActionNavigate:     true,
	ActionNavigateBack: true,
	ActionAPICall:      true,
	ActionSubmitForm:   true,
	ActionRefresh:      true,
	ActionConfirm:      true,
	ActionLogout:       true,
	ActionCustom:       true,
}

var validScreenTypes = map[ScreenType]bool{
	ScreenTypeList:      true,
	ScreenTypeDetail:    true,
	ScreenTypeCreate:    true,
	ScreenTypeEdit:      true,
	ScreenTypeDashboard: true,
	ScreenTypeSettings:  true,
}

// ValidatePattern valida que el string sea un Pattern valido
func ValidatePattern(p string) error {
	if !validPatterns[Pattern(p)] {
		return fmt.Errorf("invalid pattern: %q", p)
	}
	return nil
}

// ValidateActionType valida que el string sea un ActionType valido
func ValidateActionType(a string) error {
	if !validActionTypes[ActionType(a)] {
		return fmt.Errorf("invalid action type: %q", a)
	}
	return nil
}

// ValidateScreenType valida que el string sea un ScreenType valido
func ValidateScreenType(st string) error {
	if !validScreenTypes[ScreenType(st)] {
		return fmt.Errorf("invalid screen type: %q", st)
	}
	return nil
}

// templateDefinition es la estructura interna para validar definitions
type templateDefinition struct {
	Zones []zone `json:"zones"`
}

type zone struct {
	ID    string `json:"id"`
	Type  string `json:"type"`
	Slots []slot `json:"slots"`
}

type slot struct {
	ID          string `json:"id"`
	ControlType string `json:"controlType"`
}

// ValidateTemplateDefinition valida que el JSON de definition tenga la estructura correcta
func ValidateTemplateDefinition(definition []byte) error {
	var def templateDefinition
	if err := json.Unmarshal(definition, &def); err != nil {
		return fmt.Errorf("invalid template definition JSON: %w", err)
	}

	if len(def.Zones) == 0 {
		return fmt.Errorf("template definition must have at least one zone")
	}

	for i, z := range def.Zones {
		if z.ID == "" {
			return fmt.Errorf("zone at index %d is missing 'id'", i)
		}
		if z.Type == "" {
			return fmt.Errorf("zone %q is missing 'type'", z.ID)
		}
		for j, s := range z.Slots {
			if s.ID == "" {
				return fmt.Errorf("slot at index %d in zone %q is missing 'id'", j, z.ID)
			}
			if s.ControlType == "" {
				return fmt.Errorf("slot %q in zone %q is missing 'controlType'", s.ID, z.ID)
			}
		}
	}

	return nil
}
