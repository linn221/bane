package services

import (
	"context"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type wordService struct {
	db           *gorm.DB
	aliasService *aliasService
}

func (ws *wordService) CreateWord(input *models.WordInput) (*models.Word, error) {
	word := models.Word{
		Word:        input.Word,
		WordType:    input.WordType,
		Description: input.Description,
	}
	
	var result *models.Word
	err := ws.db.Transaction(func(tx *gorm.DB) error {
		if err := input.Validate(tx, 0); err != nil {
			return err
		}
		err := tx.Create(&word).Error
		if err != nil {
			return err
		}
		result = &word
		// Create alias (will be auto-generated if not provided)
		if err := ws.aliasService.CreateAlias(tx, "words", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ws *wordService) CreateWordList(input *models.WordListInput) (*models.WordList, error) {
	wordList := models.WordList{
		Name:        input.Name,
		Description: input.Description,
	}
	
	var result *models.WordList
	err := ws.db.Transaction(func(tx *gorm.DB) error {
		if err := input.Validate(tx, 0); err != nil {
			return err
		}
		err := tx.Create(&wordList).Error
		if err != nil {
			return err
		}
		result = &wordList
		// Create alias (will be auto-generated if not provided)
		if err := ws.aliasService.CreateAlias(tx, "wordlists", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ws *wordService) GetWord(id *int, alias *string) (*models.Word, error) {
	if id != nil {
		var word models.Word
		err := ws.db.First(&word, *id).Error
		return &word, err
	}
	if alias != nil {
		wordId, err := ws.aliasService.GetReferenceId(context.Background(), *alias)
		if err != nil {
			return nil, err
		}
		var word models.Word
		err = ws.db.First(&word, wordId).Error
		return &word, err
	}
	return nil, gorm.ErrRecordNotFound
}

func (ws *wordService) GetWordList(id *int, alias *string) (*models.WordList, error) {
	if id != nil {
		var wordList models.WordList
		err := ws.db.First(&wordList, *id).Error
		return &wordList, err
	}
	if alias != nil {
		wordListId, err := ws.aliasService.GetReferenceId(context.Background(), *alias)
		if err != nil {
			return nil, err
		}
		var wordList models.WordList
		err = ws.db.First(&wordList, wordListId).Error
		return &wordList, err
	}
	return nil, gorm.ErrRecordNotFound
}

func (ws *wordService) AddWordsToWordList(wordListId int, words []string) error {
	// Validate if word list exists
	var wordList models.WordList
	if err := ws.db.First(&wordList, wordListId).Error; err != nil {
		return err
	}

	// Process each word
	for _, wordText := range words {
		if wordText == "" {
			continue // Skip empty words
		}

		// Check if word exists, if not create it
		var word models.Word
		err := ws.db.Where("word = ?", wordText).First(&word).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create new word with default type
				word = models.Word{
					Word:        wordText,
					WordType:    models.WordTypeFuzz, // Default to fuzz type
					Description: "",
				}
				if err := ws.db.Create(&word).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// Check if relationship already exists
		var count int64
		err = ws.db.Table("word_list_words").
			Where("word_id = ? AND word_list_id = ?", word.Id, wordListId).
			Count(&count).Error
		if err != nil {
			return err
		}

		// If relationship doesn't exist, create it
		if count == 0 {
			err = ws.db.Table("word_list_words").
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
