package services

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type tagService struct {
	GeneralCrud[models.NewTag, models.Tag]
	db *gorm.DB
}

func newTagService(db *gorm.DB) *tagService {
	return &tagService{
		GeneralCrud: GeneralCrud[models.NewTag, models.Tag]{
			transform: func(input *models.NewTag) models.Tag {
				result := models.Tag{
					Name:        input.Name,
					Description: input.Description,
					Alias:       input.Alias,
					Priority:    input.Priority,
				}
				return result
			},
			updates: func(existing models.Tag, input *models.NewTag) map[string]any {
				return map[string]any{
					"Name":        input.Name,
					"Description": input.Description,
					"Alias":       input.Alias,
					"Priority":    input.Priority,
				}
			},
			validateWrite: func(db *gorm.DB, input *models.NewTag, id int) error {
				return input.Validate(db, id)
			},
		},
		db: db,
	}
}

func (ts *tagService) Create(input *models.NewTag) (*models.Tag, error) {
	return ts.GeneralCrud.Create(ts.db, input)
}

func (ts *tagService) Get(id *int, alias *string) (*models.Tag, error) {
	if id != nil {
		return ts.GeneralCrud.Get(ts.db, id)
	}
	if alias != nil {
		return ts.GeneralCrud.GetByAlias(ts.db, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ts *tagService) Update(id *int, alias *string, input *models.NewTag) (*models.Tag, error) {
	if id != nil {
		return ts.GeneralCrud.Update(ts.db, input, id)
	}
	if alias != nil {
		return ts.GeneralCrud.UpdateByAlias(ts.db, input, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ts *tagService) Delete(id *int, alias *string) (*models.Tag, error) {
	if id != nil {
		return ts.GeneralCrud.Delete(ts.db, id)
	}
	if alias != nil {
		return ts.GeneralCrud.DeleteByAlias(ts.db, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}
