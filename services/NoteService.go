package services

import (
	"context"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type noteService struct {
	GeneralCrud[models.NoteInput, models.Note]
	db *gorm.DB
}

func (s *noteService) Create(ctx context.Context, input *models.NoteInput) (*models.Note, error) {
	return s.GeneralCrud.Create(s.db.WithContext(ctx), input)
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
	return s.GeneralCrud.Delete(s.db.WithContext(ctx), id)
}
