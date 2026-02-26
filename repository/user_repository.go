package repository

import (
	"context"
	"errors"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresUserRepository struct{ db *gorm.DB }

// NewPostgresUserRepository crea una nueva instancia del repositorio de usuarios con PostgreSQL.
func NewPostgresUserRepository(db *gorm.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

func (r *postgresUserRepository) Create(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *postgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	var u entities.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresUserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var u entities.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *postgresUserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *postgresUserRepository) Update(ctx context.Context, user *entities.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *postgresUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&entities.User{}, "id = ?", id).Error
}

func (r *postgresUserRepository) List(ctx context.Context, filters ListFilters) ([]*entities.User, error) {
	query := r.db.WithContext(ctx).Model(&entities.User{})
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
	var users []*entities.User
	err := query.Find(&users).Error
	return users, err
}
