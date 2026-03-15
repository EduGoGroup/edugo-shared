package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMembershipMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	require.NoError(t, err)

	return gormDB, mock
}

func TestMembershipRepository_Create(t *testing.T) {
	gormDB, mock := setupMembershipMockDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	membership := &entities.Membership{
		ID:       uuid.New(),
		UserID:   uuid.New(),
		SchoolID: uuid.New(),
		Role:     "student",
		IsActive: true,
	}

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "academic"."memberships"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnRows(sqlmock.NewRows([]string{"metadata"}).AddRow(nil))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "academic"."memberships"`)).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Create(context.Background(), membership)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMembershipRepository_FindByID(t *testing.T) {
	gormDB, mock := setupMembershipMockDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	membershipID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		want    *entities.Membership
		wantErr error
	}{
		{
			name: "Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "role"}).
					AddRow(membershipID, "student")
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE id = $1 ORDER BY "memberships"."id" LIMIT $2`)).
					WithArgs(membershipID, 1).
					WillReturnRows(rows)
			},
			want: &entities.Membership{
				ID:   membershipID,
				Role: "student",
			},
			wantErr: nil,
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE id = $1 ORDER BY "memberships"."id" LIMIT $2`)).
					WithArgs(membershipID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: ErrNotFound,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE id = $1 ORDER BY "memberships"."id" LIMIT $2`)).
					WithArgs(membershipID, 1).
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindByID(context.Background(), membershipID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, ErrNotFound) {
					assert.Equal(t, ErrNotFound, err)
				}
			} else {
				assert.NoError(t, err)
				if got != nil && tt.want != nil {
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.Role, got.Role)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMembershipRepository_FindByUser(t *testing.T) {
	gormDB, mock := setupMembershipMockDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	userID := uuid.New()

	tests := []struct {
		name    string
		filters ListFilters
		mock    func()
		wantLen int
		wantErr bool
	}{
		{
			name:    "No Filters",
			filters: ListFilters{},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE user_id = $1 AND is_active = true`)).
					WithArgs(userID).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "user_id"}).AddRow(uuid.New(), userID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE user_id = $1 AND is_active = true ORDER BY created_at DESC`)).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "DB Error",
			filters: ListFilters{},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE user_id = $1 AND is_active = true`)).
					WithArgs(userID).
					WillReturnError(errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			memberships, total, err := repo.FindByUser(context.Background(), userID, tt.filters)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, memberships, tt.wantLen)
				if tt.wantLen > 0 {
					assert.Equal(t, int64(1), total)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMembershipRepository_FindByUnit(t *testing.T) {
	gormDB, mock := setupMembershipMockDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	unitID := uuid.New()

	tests := []struct {
		name    string
		filters ListFilters
		mock    func()
		wantLen int
		wantErr bool
	}{
		{
			name:    "No Filters",
			filters: ListFilters{},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE academic_unit_id = $1 AND is_active = true`)).
					WithArgs(unitID).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "academic_unit_id"}).AddRow(uuid.New(), unitID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE academic_unit_id = $1 AND is_active = true ORDER BY created_at DESC`)).
					WithArgs(unitID).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:    "DB Error",
			filters: ListFilters{},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE academic_unit_id = $1 AND is_active = true`)).
					WithArgs(unitID).
					WillReturnError(errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			memberships, total, err := repo.FindByUnit(context.Background(), unitID, tt.filters)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, memberships, tt.wantLen)
				if tt.wantLen > 0 {
					assert.Equal(t, int64(1), total)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMembershipRepository_FindByUnitAndRole(t *testing.T) {
	gormDB, mock := setupMembershipMockDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	unitID := uuid.New()
	role := "student"

	tests := []struct {
		name       string
		activeOnly bool
		filters    ListFilters
		mock       func()
		wantLen    int
		wantErr    bool
	}{
		{
			name:       "Active Only",
			activeOnly: true,
			filters:    ListFilters{},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE (academic_unit_id = $1 AND role = $2) AND is_active = true`)).
					WithArgs(unitID, role).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "academic_unit_id", "role"}).AddRow(uuid.New(), unitID, role)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE (academic_unit_id = $1 AND role = $2) AND is_active = true ORDER BY created_at DESC`)).
					WithArgs(unitID, role).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:       "All (Active and Inactive)",
			activeOnly: false,
			filters:    ListFilters{},
			mock: func() {
				countRows := sqlmock.NewRows([]string{"count"}).AddRow(1)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE academic_unit_id = $1 AND role = $2`)).
					WithArgs(unitID, role).
					WillReturnRows(countRows)

				rows := sqlmock.NewRows([]string{"id", "academic_unit_id", "role"}).AddRow(uuid.New(), unitID, role)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE academic_unit_id = $1 AND role = $2 ORDER BY created_at DESC`)).
					WithArgs(unitID, role).
					WillReturnRows(rows)
			},
			wantLen: 1,
			wantErr: false,
		},
		{
			name:       "DB Error",
			activeOnly: false,
			filters:    ListFilters{},
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "academic"."memberships" WHERE academic_unit_id = $1 AND role = $2`)).
					WithArgs(unitID, role).
					WillReturnError(errors.New("db error"))
			},
			wantLen: 0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			memberships, total, err := repo.FindByUnitAndRole(context.Background(), unitID, role, tt.activeOnly, tt.filters)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, memberships, tt.wantLen)
				if tt.wantLen > 0 {
					assert.Equal(t, int64(1), total)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMembershipRepository_FindByUserAndSchool(t *testing.T) {
	gormDB, mock := setupMembershipMockDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	userID := uuid.New()
	schoolID := uuid.New()
	membershipID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		want    *entities.Membership
		wantErr error
	}{
		{
			name: "Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "user_id", "school_id"}).
					AddRow(membershipID, userID, schoolID)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE user_id = $1 AND school_id = $2 AND is_active = true ORDER BY "memberships"."id" LIMIT $3`)).
					WithArgs(userID, schoolID, 1).
					WillReturnRows(rows)
			},
			want: &entities.Membership{
				ID:       membershipID,
				UserID:   userID,
				SchoolID: schoolID,
			},
			wantErr: nil,
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE user_id = $1 AND school_id = $2 AND is_active = true ORDER BY "memberships"."id" LIMIT $3`)).
					WithArgs(userID, schoolID, 1).
					WillReturnError(gorm.ErrRecordNotFound)
			},
			want:    nil,
			wantErr: ErrNotFound,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "academic"."memberships" WHERE user_id = $1 AND school_id = $2 AND is_active = true ORDER BY "memberships"."id" LIMIT $3`)).
					WithArgs(userID, schoolID, 1).
					WillReturnError(errors.New("db error"))
			},
			want:    nil,
			wantErr: errors.New("db error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := repo.FindByUserAndSchool(context.Background(), userID, schoolID)
			if tt.wantErr != nil {
				assert.Error(t, err)
				if errors.Is(tt.wantErr, ErrNotFound) {
					assert.Equal(t, ErrNotFound, err)
				}
			} else {
				assert.NoError(t, err)
				if got != nil && tt.want != nil {
					assert.Equal(t, tt.want.ID, got.ID)
					assert.Equal(t, tt.want.UserID, got.UserID)
					assert.Equal(t, tt.want.SchoolID, got.SchoolID)
				}
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMembershipRepository_Update(t *testing.T) {
	gormDB, mock := setupMembershipMockDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	membership := &entities.Membership{
		ID:   uuid.New(),
		Role: "student",
	}

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."memberships" SET`)).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE "academic"."memberships" SET`)).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Update(context.Background(), membership)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMembershipRepository_Delete(t *testing.T) {
	gormDB, mock := setupMembershipMockDB(t)
	repo := NewPostgresMembershipRepository(gormDB)

	membershipID := uuid.New()

	tests := []struct {
		name    string
		mock    func()
		wantErr bool
	}{
		{
			name: "Success",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "academic"."memberships" WHERE id = $1`)).
					WithArgs(membershipID).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			wantErr: false,
		},
		{
			name: "DB Error",
			mock: func() {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "academic"."memberships" WHERE id = $1`)).
					WithArgs(membershipID).
					WillReturnError(errors.New("db error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			err := repo.Delete(context.Background(), membershipID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
