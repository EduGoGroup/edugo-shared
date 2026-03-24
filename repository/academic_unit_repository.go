package repository

import (
	"context"
	"errors"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresAcademicUnitRepository struct{ db *gorm.DB }

// NewPostgresAcademicUnitRepository crea un repositorio de solo lectura para unidades académicas.
func NewPostgresAcademicUnitRepository(db *gorm.DB) AcademicUnitRepository {
	return &postgresAcademicUnitRepository{db: db}
}

func (r *postgresAcademicUnitRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.AcademicUnit, error) {
	var unit entities.AcademicUnit
	if err := r.db.WithContext(ctx).First(&unit, "id = ? AND is_active = true", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &unit, nil
}

func (r *postgresAcademicUnitRepository) FindBySchoolID(ctx context.Context, schoolID uuid.UUID, filters ListFilters) ([]*entities.AcademicUnit, int64, error) {
	buildBase := func() *gorm.DB {
		q := r.db.WithContext(ctx).Model(&entities.AcademicUnit{}).Where("school_id = ? AND is_active = true", schoolID)
		q = filters.ApplySearch(q)
		return q
	}

	var total int64
	if err := buildBase().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var units []*entities.AcademicUnit
	query := buildBase().Order("name ASC")
	query = filters.ApplyPagination(query)
	if err := query.Find(&units).Error; err != nil {
		return nil, 0, err
	}
	return units, total, nil
}
