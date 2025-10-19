package services

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

var TagCrud = GeneralCrud[models.NewTag, models.Tag]{
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
}
