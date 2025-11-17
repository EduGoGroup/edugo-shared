package evaluation

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Assessment representa un cuestionario generado o manual
type Assessment struct {
	ID               uuid.UUID  `json:"id" bson:"_id"`
	MaterialID       uuid.UUID  `json:"material_id" bson:"material_id"`                       // BREAKING CHANGE: int64 → uuid.UUID
	MongoDocID       string     `json:"mongo_doc_id,omitempty" bson:"mongo_doc_id,omitempty"` // Referencia a documento en MongoDB
	Title            string     `json:"title" bson:"title"`
	Description      string     `json:"description,omitempty" bson:"description,omitempty"`
	Type             string     `json:"type" bson:"type"`                   // "manual", "generated"
	Status           string     `json:"status" bson:"status"`               // "draft", "published", "archived"
	PassingScore     int        `json:"passing_score" bson:"passing_score"` // Porcentaje mínimo para aprobar (0-100)
	TotalQuestions   int        `json:"total_questions" bson:"total_questions"`
	TotalPoints      int        `json:"total_points" bson:"total_points"`
	MaxAttempts      *int       `json:"max_attempts,omitempty" bson:"max_attempts,omitempty"`             // Nullable - NULL = ilimitado
	TimeLimitMinutes *int       `json:"time_limit_minutes,omitempty" bson:"time_limit_minutes,omitempty"` // Nullable - NULL = sin límite
	CreatedBy        int64      `json:"created_by" bson:"created_by"`                                     // User ID
	CreatedAt        time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" bson:"updated_at"`
}

// Validate valida los campos del assessment
func (a *Assessment) Validate() error {
	if a.Title == "" {
		return errors.New("title is required")
	}
	if a.PassingScore < 0 || a.PassingScore > 100 {
		return errors.New("passing score must be between 0 and 100")
	}
	if a.MaxAttempts != nil && *a.MaxAttempts <= 0 {
		return errors.New("max attempts must be greater than 0")
	}
	if a.TimeLimitMinutes != nil && *a.TimeLimitMinutes <= 0 {
		return errors.New("time limit minutes must be greater than 0")
	}
	return nil
}

// IsPublished retorna si el assessment está publicado
func (a *Assessment) IsPublished() bool {
	return a.Status == "published"
}
