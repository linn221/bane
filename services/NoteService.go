package services

import (
	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type noteService struct {
	GeneralCrud[models.NoteInput, models.Note]
	db      *gorm.DB
	deducer Deducer
}

func newNoteService(db *gorm.DB, deducer Deducer) *noteService {
	return &noteService{
		GeneralCrud: GeneralCrud[models.NoteInput, models.Note]{
			transform: func(input *models.NoteInput) models.Note {
				today := utils.Today()
				return models.Note{
					ReferenceType: input.ReferenceType,
					ReferenceID:   input.ReferenceId,
					Value:         input.Value,
					NoteDate:      models.MyDate{Time: today},
				}
			},
		},
		db:      db,
		deducer: deducer,
	}
}

func (ns *noteService) Create(input *models.NoteInput) (*models.Note, error) {
	if input.RId > 0 {
		input.ReferenceId, input.ReferenceType = ns.deducer.ReadRId(input.RId)
	}

	return ns.GeneralCrud.Create(ns.db, input)
}

func (ns *noteService) List(filter *models.NoteFilter) ([]*models.Note, error) {
	dbctx := ns.db.Model(&models.Note{})
	if filter != nil {
		if !filter.NoteDate.IsZero() {
			dbctx.Where("note_date = ?", filter.NoteDate)
		}

		if filter.RID > 0 {
			filter.ReferenceID, filter.ReferenceType = ns.deducer.ReadRId(filter.RID)
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

func (ns *noteService) Delete(id *int) (*models.Note, error) {
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return ns.GeneralCrud.Delete(ns.db, id)
}
