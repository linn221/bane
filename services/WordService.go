package services

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type wordService struct {
	db           *gorm.DB
	aliasService *aliasService
}

func newWordService(db *gorm.DB, aliasService *aliasService) *wordService {
	return &wordService{db: db, aliasService: aliasService}
}

func (ws *wordService) CreateWord(input *models.NewWord) (*models.Word, error) {
	wordCrud := GeneralCrud[models.NewWord, models.Word]{
		transform: func(input *models.NewWord) models.Word {
			result := models.Word{
				Word:        input.Word,
				WordType:    input.WordType,
				Description: input.Description,
			}
			return result
		},
		updates: func(existing models.Word, input *models.NewWord) map[string]any {
			return map[string]any{
				"Word":        input.Word,
				"WordType":    input.WordType,
				"Description": input.Description,
			}
		},
		validateWrite: func(db *gorm.DB, input *models.NewWord, id int) error {
			return input.Validate(db, id)
		},
	}
	result, err := wordCrud.Create(ws.db, input)
	if err != nil {
		return nil, err
	}
	// Set alias if provided
	if input.Alias != "" {
		if err := ws.aliasService.SetAlias(string(models.AliasReferenceTypeWord), result.Id, input.Alias); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (ws *wordService) CreateWordList(input *models.NewWordList) (*models.WordList, error) {
	wordListCrud := GeneralCrud[models.NewWordList, models.WordList]{
		transform: func(input *models.NewWordList) models.WordList {
			result := models.WordList{
				Name:        input.Name,
				Description: input.Description,
			}
			return result
		},
		updates: func(existing models.WordList, input *models.NewWordList) map[string]any {
			return map[string]any{
				"Name":        input.Name,
				"Description": input.Description,
			}
		},
		validateWrite: func(db *gorm.DB, input *models.NewWordList, id int) error {
			return input.Validate(db, id)
		},
	}
	result, err := wordListCrud.Create(ws.db, input)
	if err != nil {
		return nil, err
	}
	// Set alias if provided
	if input.Alias != "" {
		if err := ws.aliasService.SetAlias(string(models.AliasReferenceTypeWordList), result.Id, input.Alias); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (ws *wordService) GetWord(id *int, alias *string) (*models.Word, error) {
	wordCrud := GeneralCrud[models.NewWord, models.Word]{
		transform: func(input *models.NewWord) models.Word {
			result := models.Word{
				Word:        input.Word,
				WordType:    input.WordType,
				Description: input.Description,
			}
			return result
		},
		updates: func(existing models.Word, input *models.NewWord) map[string]any {
			return map[string]any{
				"Word":        input.Word,
				"WordType":    input.WordType,
				"Description": input.Description,
			}
		},
		validateWrite: func(db *gorm.DB, input *models.NewWord, id int) error {
			return input.Validate(db, id)
		},
	}
	if id != nil {
		return wordCrud.Get(ws.db, id)
	}
	if alias != nil {
		return wordCrud.GetByAlias(ws.db, ws.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ws *wordService) UpdateWord(id *int, alias *string, input *models.NewWord) (*models.Word, error) {
	wordCrud := GeneralCrud[models.NewWord, models.Word]{
		transform: func(input *models.NewWord) models.Word {
			result := models.Word{
				Word:        input.Word,
				WordType:    input.WordType,
				Description: input.Description,
			}
			return result
		},
		updates: func(existing models.Word, input *models.NewWord) map[string]any {
			return map[string]any{
				"Word":        input.Word,
				"WordType":    input.WordType,
				"Description": input.Description,
			}
		},
		validateWrite: func(db *gorm.DB, input *models.NewWord, id int) error {
			return input.Validate(db, id)
		},
	}
	var wordId int
	var result *models.Word
	var err error
	if id != nil {
		wordId = *id
		result, err = wordCrud.Update(ws.db, input, id)
	} else if alias != nil {
		wordId, err = ws.aliasService.GetId(*alias)
		if err != nil {
			return nil, err
		}
		idPtr := &wordId
		result, err = wordCrud.Update(ws.db, input, idPtr)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	// Set alias if provided
	if input.Alias != "" {
		if err := ws.aliasService.SetAlias(string(models.AliasReferenceTypeWord), wordId, input.Alias); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (ws *wordService) PatchWord(id *int, alias *string, updates map[string]any) (*models.Word, error) {
	if id != nil {
		wordCrud := GeneralCrud[models.NewWord, models.Word]{}
		return wordCrud.Patch(ws.db, updates, id)
	}
	if alias != nil {
		wordCrud := GeneralCrud[models.NewWord, models.Word]{}
		return wordCrud.PatchByAlias(ws.db, ws.aliasService, updates, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ws *wordService) DeleteWord(id *int, alias *string) (*models.Word, error) {
	wordCrud := GeneralCrud[models.NewWord, models.Word]{
		validateWrite: func(db *gorm.DB, input *models.NewWord, id int) error {
			return input.Validate(db, id)
		},
	}
	if id != nil {
		return wordCrud.Delete(ws.db, id)
	}
	if alias != nil {
		return wordCrud.DeleteByAlias(ws.db, ws.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ws *wordService) GetWordList(id *int, alias *string) (*models.WordList, error) {
	wordListCrud := GeneralCrud[models.NewWordList, models.WordList]{
		transform: func(input *models.NewWordList) models.WordList {
			result := models.WordList{
				Name:        input.Name,
				Description: input.Description,
			}
			return result
		},
		updates: func(existing models.WordList, input *models.NewWordList) map[string]any {
			return map[string]any{
				"Name":        input.Name,
				"Description": input.Description,
			}
		},
		validateWrite: func(db *gorm.DB, input *models.NewWordList, id int) error {
			return input.Validate(db, id)
		},
	}
	if id != nil {
		return wordListCrud.Get(ws.db, id)
	}
	if alias != nil {
		return wordListCrud.GetByAlias(ws.db, ws.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ws *wordService) UpdateWordList(id *int, alias *string, input *models.NewWordList) (*models.WordList, error) {
	wordListCrud := GeneralCrud[models.NewWordList, models.WordList]{
		transform: func(input *models.NewWordList) models.WordList {
			result := models.WordList{
				Name:        input.Name,
				Description: input.Description,
			}
			return result
		},
		updates: func(existing models.WordList, input *models.NewWordList) map[string]any {
			return map[string]any{
				"Name":        input.Name,
				"Description": input.Description,
			}
		},
		validateWrite: func(db *gorm.DB, input *models.NewWordList, id int) error {
			return input.Validate(db, id)
		},
	}
	var wordListId int
	var result *models.WordList
	var err error
	if id != nil {
		wordListId = *id
		result, err = wordListCrud.Update(ws.db, input, id)
	} else if alias != nil {
		wordListId, err = ws.aliasService.GetId(*alias)
		if err != nil {
			return nil, err
		}
		idPtr := &wordListId
		result, err = wordListCrud.Update(ws.db, input, idPtr)
	} else {
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		return nil, err
	}
	// Set alias if provided
	if input.Alias != "" {
		if err := ws.aliasService.SetAlias(string(models.AliasReferenceTypeWordList), wordListId, input.Alias); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (ws *wordService) PatchWordList(id *int, alias *string, updates map[string]any) (*models.WordList, error) {
	if id != nil {
		wordListCrud := GeneralCrud[models.NewWordList, models.WordList]{}
		return wordListCrud.Patch(ws.db, updates, id)
	}
	if alias != nil {
		wordListCrud := GeneralCrud[models.NewWordList, models.WordList]{}
		return wordListCrud.PatchByAlias(ws.db, ws.aliasService, updates, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ws *wordService) DeleteWordList(id *int, alias *string) (*models.WordList, error) {
	wordListCrud := GeneralCrud[models.NewWordList, models.WordList]{
		validateWrite: func(db *gorm.DB, input *models.NewWordList, id int) error {
			return input.Validate(db, id)
		},
	}
	if id != nil {
		return wordListCrud.Delete(ws.db, id)
	}
	if alias != nil {
		return wordListCrud.DeleteByAlias(ws.db, ws.aliasService, *alias)
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
