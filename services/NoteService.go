package services

import (
	"github.com/linn221/bane/app"
	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type noteService struct {
	GeneralCrud[models.NewNote, models.Note]
}

var NoteService = noteService{
	GeneralCrud: GeneralCrud[models.NewNote, models.Note]{
		transform: func(input *models.NewNote) models.Note {
			today := utils.Today()
			return models.Note{
				ReferenceType: input.ReferenceType,
				ReferenceID:   input.ReferenceId,
				Value:         input.Value,
				NoteDate:      models.MyDate{Time: today},
			}
		},
	},
}

func (ns *noteService) CreateNote(app *app.App, db *gorm.DB, input *models.NewNote) (*models.Note, error) {
	if input.RId > 0 {
		input.ReferenceId, input.ReferenceType = app.Deducer.ReadRId(input.RId)
	}

	return ns.Create(db, input)
}

func (ns *noteService) ListNotes(app *app.App, db *gorm.DB, filter *models.NoteFilter) ([]*models.Note, error) {
	dbctx := db.Model(&models.Note{})
	if filter != nil {
		if !filter.NoteDate.IsZero() {
			dbctx.Where("note_date = ?", filter.NoteDate)
		}

		if filter.RID > 0 {
			filter.ReferenceID, filter.ReferenceType = app.Deducer.ReadRId(filter.RID)
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

func (ns *noteService) Delete(db *gorm.DB, id *int) (*models.Note, error) {
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return ns.GeneralCrud.Delete(db, id)
}
