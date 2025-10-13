package services

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

var WordCrud = GeneralCrud[models.NewWord, models.Word]{
	Transform: func(input models.NewWord) models.Word {
		result := models.Word{
			Word:        input.Word,
			WordType:    input.WordType,
			Description: input.Description,
		}
		return result
	},
	Updates: func(input models.NewWord) map[string]any {
		return map[string]any{
			"Word":        input.Word,
			"WordType":    input.WordType,
			"Description": input.Description,
		}
	},
	ValidateWrite: func(db *gorm.DB, input models.NewWord, id int) error {
		return input.Validate(db, id)
	},
}

var WordListCrud = GeneralCrud[models.NewWordList, models.WordList]{
	Transform: func(input models.NewWordList) models.WordList {
		result := models.WordList{
			Name:        input.Name,
			Description: input.Description,
		}
		return result
	},
	Updates: func(input models.NewWordList) map[string]any {
		return map[string]any{
			"Name":        input.Name,
			"Description": input.Description,
		}
	},
	ValidateWrite: func(db *gorm.DB, input models.NewWordList, id int) error {
		return input.Validate(db, id)
	},
}
