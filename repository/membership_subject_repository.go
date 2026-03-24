package repository

import (
	"context"

	"github.com/EduGoGroup/edugo-infrastructure/postgres/entities"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type postgresMembershipSubjectRepository struct{ db *gorm.DB }

// NewPostgresMembershipSubjectRepository crea un nuevo repositorio de subjects de membresía.
func NewPostgresMembershipSubjectRepository(db *gorm.DB) MembershipSubjectRepository {
	return &postgresMembershipSubjectRepository{db: db}
}

func (r *postgresMembershipSubjectRepository) AssignSubjects(ctx context.Context, membershipID uuid.UUID, subjectIDs []uuid.UUID) error {
	if len(subjectIDs) == 0 {
		return nil
	}
	records := make([]entities.MembershipSubject, len(subjectIDs))
	for i, sid := range subjectIDs {
		records[i] = entities.MembershipSubject{
			MembershipID: membershipID,
			SubjectID:    sid,
		}
	}
	return r.db.WithContext(ctx).Create(&records).Error
}

func (r *postgresMembershipSubjectRepository) RemoveAllByMembership(ctx context.Context, membershipID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("membership_id = ?", membershipID).Delete(&entities.MembershipSubject{}).Error
}

func (r *postgresMembershipSubjectRepository) GetByMembership(ctx context.Context, membershipID uuid.UUID) ([]*entities.MembershipSubject, error) {
	var results []*entities.MembershipSubject
	if err := r.db.WithContext(ctx).Where("membership_id = ?", membershipID).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *postgresMembershipSubjectRepository) GetBySubject(ctx context.Context, subjectID uuid.UUID) ([]*entities.MembershipSubject, error) {
	var results []*entities.MembershipSubject
	if err := r.db.WithContext(ctx).Where("subject_id = ?", subjectID).Find(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
