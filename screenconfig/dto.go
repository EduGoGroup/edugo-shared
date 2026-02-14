package screenconfig

import (
	"encoding/json"
	"time"
)

type ScreenTemplateDTO struct {
	ID          string          `json:"id"`
	Pattern     Pattern         `json:"pattern"`
	Name        string          `json:"name"`
	Description string          `json:"description,omitempty"`
	Version     int             `json:"version"`
	Definition  json.RawMessage `json:"definition"`
	IsActive    bool            `json:"is_active"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type ScreenInstanceDTO struct {
	ID                 string          `json:"id"`
	ScreenKey          string          `json:"screen_key"`
	TemplateID         string          `json:"template_id"`
	Name               string          `json:"name"`
	Description        string          `json:"description,omitempty"`
	SlotData           json.RawMessage `json:"slot_data"`
	Actions            json.RawMessage `json:"actions"`
	DataEndpoint       string          `json:"data_endpoint,omitempty"`
	DataConfig         json.RawMessage `json:"data_config,omitempty"`
	Scope              string          `json:"scope"`
	RequiredPermission string          `json:"required_permission,omitempty"`
	IsActive           bool            `json:"is_active"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
}

type CombinedScreenDTO struct {
	ScreenID        string          `json:"screenId"`
	ScreenKey       string          `json:"screenKey"`
	ScreenName      string          `json:"screenName"`
	Pattern         Pattern         `json:"pattern"`
	Version         int             `json:"version"`
	Template        json.RawMessage `json:"template"`
	DataEndpoint    string          `json:"dataEndpoint,omitempty"`
	DataConfig      json.RawMessage `json:"dataConfig,omitempty"`
	Actions         json.RawMessage `json:"actions"`
	UserPreferences json.RawMessage `json:"userPreferences,omitempty"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type ResourceScreenDTO struct {
	ResourceID  string `json:"resource_id"`
	ResourceKey string `json:"resource_key"`
	ScreenKey   string `json:"screen_key"`
	ScreenType  string `json:"screen_type"`
	IsDefault   bool   `json:"is_default"`
}

type ActionDefinitionDTO struct {
	ID            string          `json:"id"`
	Trigger       string          `json:"trigger"`
	TriggerSlotID string          `json:"triggerSlotId,omitempty"`
	Type          ActionType      `json:"type"`
	Config        json.RawMessage `json:"config"`
}
