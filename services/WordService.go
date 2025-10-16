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
	Updates: func(existing models.Word, input models.NewWord) map[string]any {
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
	Updates: func(existing models.WordList, input models.NewWordList) map[string]any {
		return map[string]any{
			"Name":        input.Name,
			"Description": input.Description,
		}
	},
	ValidateWrite: func(db *gorm.DB, input models.NewWordList, id int) error {
		return input.Validate(db, id)
	},
}

func AddWordsToWordList(db *gorm.DB, wordListId int, words []string) error {
	// Validate if word list exists
	var wordList models.WordList
	if err := db.First(&wordList, wordListId).Error; err != nil {
		return err
	}

	// Process each word
	for _, wordText := range words {
		if wordText == "" {
			continue // Skip empty words
		}

		// Check if word exists, if not create it
		var word models.Word
		err := db.Where("word = ?", wordText).First(&word).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new word with default type
				word = models.Word{
					Word:        wordText,
					WordType:    models.WordTypeFuzz, // Default to fuzz type
					Description: "",
				}
				if err := db.Create(&word).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Check if relationship already exists
		var count int64
		err = db.Table("word_list_words").
			Where("word_id = ? AND word_list_id = ?", word.Id, wordListId).
			Count(&count).Error
		if err != nil {
			return err
		}

		// If relationship doesn't exist, create it
		if count == 0 {
			err = db.Table("word_list_words").
				Create(map[string]interface{}{
					"word_id":      word.Id,
					"word_list_id": wordListId,
				}).Error
			if err != nil {
				return err
			}
		}
	}

	return nil
}
