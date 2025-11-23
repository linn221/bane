package services

import (
	"context"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type noteService struct {
	db *gorm.DB
}

func (s *noteService) Create(ctx context.Context, input *models.NoteInput) (*models.Note, error) {
	today := utils.Today()
	note := models.Note{
		ReferenceType: input.ReferenceType,
		ReferenceID:   input.ReferenceId,
		Value:         input.Value,
		NoteDate:      models.MyDate{Time: today},
	}
	err := s.db.WithContext(ctx).Create(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}

func (s *noteService) List(ctx context.Context, filter *models.NoteFilter) ([]*models.Note, error) {
	dbctx := s.db.WithContext(ctx).Model(&models.Note{})
	if filter != nil {
		if !filter.NoteDate.IsZero() {
			dbctx.Where("note_date = ?", filter.NoteDate)
		}

		if filter.ReferenceType != "" {
			if filter.ReferenceID > 0 {
				dbctx.Where("reference_type = ? AND reference_id = ?", filter.ReferenceType, filter.ReferenceID)
			} else {
				dbctx.Where("reference_type = ?", filter.ReferenceType)
			}
		}

		if filter.Search != "" {
			dbctx.Where("value LIKE ?", "%"+filter.Search+"%")
		}

	}

	var results []*models.Note
	err := dbctx.Find(&results).Error
	return results, err
}

func (s *noteService) Delete(ctx context.Context, id *int) (*models.Note, error) {
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	var note models.Note
	err := s.db.WithContext(ctx).First(&note, *id).Error
	if err != nil {
		return nil, err
	}
	err = s.db.WithContext(ctx).Delete(&note).Error
	if err != nil {
		return nil, err
	}
	return &note, nil
}
