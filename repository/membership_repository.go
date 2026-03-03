package repository

import (
	"context"
	"errors"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresMembershipRepository struct{ db *gorm.DB }

// NewPostgresMembershipRepository crea una nueva instancia del repositorio de membresías con PostgreSQL.
func NewPostgresMembershipRepository(db *gorm.DB) MembershipRepository {
	return &postgresMembershipRepository{db: db}
}

func (r *postgresMembershipRepository) Create(ctx context.Context, m *entities.Membership) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *postgresMembershipRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error) {
	var m entities.Membership
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *postgresMembershipRepository) FindByUser(ctx context.Context, userID uuid.UUID, filters ListFilters) ([]*entities.Membership, int64, error) {
	baseQuery := r.db.WithContext(ctx).Model(&entities.Membership{}).Where("user_id = ? AND is_active = true", userID)
	baseQuery = filters.ApplySearch(baseQuery)

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := baseQuery.Order("created_at DESC")
	query = filters.ApplyPagination(query)
	var memberships []*entities.Membership
	if err := query.Find(&memberships).Error; err != nil {
		return nil, 0, err
	}
	return memberships, total, nil
}

func (r *postgresMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID, filters ListFilters) ([]*entities.Membership, int64, error) {
	baseQuery := r.db.WithContext(ctx).Model(&entities.Membership{}).Where("academic_unit_id = ? AND is_active = true", unitID)
	baseQuery = filters.ApplySearch(baseQuery)

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := baseQuery.Order("created_at DESC")
	query = filters.ApplyPagination(query)
	var memberships []*entities.Membership
	if err := query.Find(&memberships).Error; err != nil {
		return nil, 0, err
	}
	return memberships, total, nil
}

func (r *postgresMembershipRepository) FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool, filters ListFilters) ([]*entities.Membership, int64, error) {
	baseQuery := r.db.WithContext(ctx).Model(&entities.Membership{}).Where("academic_unit_id = ? AND role = ?", unitID, role)
	if activeOnly {
		baseQuery = baseQuery.Where("is_active = true")
	}
	baseQuery = filters.ApplySearch(baseQuery)

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	query := baseQuery.Order("created_at DESC")
	query = filters.ApplyPagination(query)
	var memberships []*entities.Membership
	if err := query.Find(&memberships).Error; err != nil {
		return nil, 0, err
	}
	return memberships, total, nil
}

func (r *postgresMembershipRepository) FindByUserAndSchool(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error) {
	var m entities.Membership
	if err := r.db.WithContext(ctx).Where("user_id = ? AND school_id = ? AND is_active = true", userID, schoolID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &m, nil
}

func (r *postgresMembershipRepository) Update(ctx context.Context, m *entities.Membership) error {
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *postgresMembershipRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.Membership{}, "id = ?", id).Error
}
