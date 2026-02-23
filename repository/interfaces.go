package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
)

// ListFilters represents common filters for listing entities
type ListFilters struct {
	IsActive *bool
	Limit    int
	Offset   int
}

// UserRepository defines persistence operations for User
type UserRepository interface {
	Create(ctx context.Context, user *entities.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
	FindByEmail(ctx context.Context, email string) (*entities.User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *entities.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters ListFilters) ([]*entities.User, error)
}

// MembershipRepository defines persistence operations for Membership
type MembershipRepository interface {
	Create(ctx context.Context, membership *entities.Membership) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.Membership, error)
	FindByUser(ctx context.Context, userID uuid.UUID) ([]*entities.Membership, error)
	FindByUnit(ctx context.Context, unitID uuid.UUID) ([]*entities.Membership, error)
	FindByUnitAndRole(ctx context.Context, unitID uuid.UUID, role string, activeOnly bool) ([]*entities.Membership, error)
	FindByUserAndSchool(ctx context.Context, userID, schoolID uuid.UUID) (*entities.Membership, error)
	Update(ctx context.Context, membership *entities.Membership) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// SchoolRepository defines persistence operations for School
type SchoolRepository interface {
	Create(ctx context.Context, school *entities.School) error
	FindByID(ctx context.Context, id uuid.UUID) (*entities.School, error)
	FindByCode(ctx context.Context, code string) (*entities.School, error)
	Update(ctx context.Context, school *entities.School) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, filters ListFilters) ([]*entities.School, error)
	ExistsByCode(ctx context.Context, code string) (bool, error)
}
