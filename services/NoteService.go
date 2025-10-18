package services

import (
	"github.com/linn221/bane/app"
	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type noteCrud struct {
	GeneralCrud[models.NewNote, models.Note]
}

var NoteService = noteCrud{
	GeneralCrud: GeneralCrud[models.NewNote, models.Note]{
		Transform: func(input models.NewNote) models.Note {
			today := utils.Today()
			return models.Note{
				ReferenceType: input.ReferenceType,
				ReferenceID:   input.ReferenceId,
				Value:         input.Value,
				NoteDate:      today,
			}
		},
	},
}

func (ns *noteCrud) CreateNote(app *app.App, db *gorm.DB, input *models.NewNote) (*models.Note, error) {
	if input.RId > 0 {
		input.ReferenceId, input.ReferenceType = app.Deducer.ReadRId(input.RId)
	}

	return ns.Create(db, *input)
}
