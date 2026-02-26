package repository

import (
	"context"
	"errors"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresSchoolRepository struct{ db *gorm.DB }

// NewPostgresSchoolRepository crea una nueva instancia del repositorio de escuelas con PostgreSQL.
func NewPostgresSchoolRepository(db *gorm.DB) SchoolRepository {
	return &postgresSchoolRepository{db: db}
}

func (r *postgresSchoolRepository) Create(ctx context.Context, school *entities.School) error {
	return r.db.WithContext(ctx).Create(school).Error
}

func (r *postgresSchoolRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error) {
	var s entities.School
	if err := r.db.WithContext(ctx).First(&s, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresSchoolRepository) FindByCode(ctx context.Context, code string) (*entities.School, error) {
	var s entities.School
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&s).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &s, nil
}

func (r *postgresSchoolRepository) Update(ctx context.Context, school *entities.School) error {
	return r.db.WithContext(ctx).Save(school).Error
}

func (r *postgresSchoolRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.School{}, "id = ?", id).Error
}

func (r *postgresSchoolRepository) List(ctx context.Context, filters ListFilters) ([]*entities.School, error) {
	query := r.db.WithContext(ctx).Model(&entities.School{})
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	query = filters.ApplySearch(query)
	query = query.Order("created_at DESC")
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}
	var schools []*entities.School
	err := query.Find(&schools).Error
	return schools, err
}

func (r *postgresSchoolRepository) ExistsByCode(ctx context.Context, code string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.School{}).Where("code = ?", code).Count(&count).Error
	return count > 0, err
}
