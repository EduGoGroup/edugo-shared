package postgres

import (
	"context"
	"errors"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/EduGoGroup/edugo-shared/repository"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresMembershipRepository struct{ db *gorm.DB }

func NewPostgresMembershipRepository(db *gorm.DB) repository.MembershipRepository {
	return &postgresMembershipRepository{db: db}
}

func (r *postgresMembershipRepository) Create(ctx context.Context, m *entities.Membership) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *postgresMembershipRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error) {
	var m entities.Membership
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (r *postgresMembershipRepository) FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error) {
	var memberships []*entities.Membership
	err := r.db.WithContext(ctx).Where("user_id = ? AND is_active = true", userID).Order("created_at DESC").Find(&memberships).Error
	return memberships, err
}

func (r *postgresMembershipRepository) FindByUnit(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error) {
	var memberships []*entities.Membership
	err := r.db.WithContext(ctx).Where("academic_unit_id = ? AND is_active = true", unitID).Order("created_at DESC").Find(&memberships).Error
	return memberships, err
}

func (r *postgresMembershipRepository) FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool) ([]*entities.Membership, error) {
	query := r.db.WithContext(ctx).Where("academic_unit_id = ? AND role = ?", unitID, role)
	if activeOnly {
		query = query.Where("is_active = true")
	}
	var memberships []*entities.Membership
	err := query.Find(&memberships).Error
	return memberships, err
}

func (r *postgresMembershipRepository) FindByUserAndSchool(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error) {
	var m entities.Membership
	if err := r.db.WithContext(ctx).Where("user_id = ? AND school_id = ? AND is_active = true", userID, schoolID).First(&m).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
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
