package models

import (
	"strings"
	"time"

	"github.com/linn221/bane/validate"
	"gorm.io/gorm"
)

// Word represents a word in a wordlist
type Word struct {
	Id          int        `gorm:"primaryKey"`
	Word        string     `gorm:"size:255;not null;index"`
	WordType    WordType   `gorm:"not null;index"`
	Description string     `gorm:"text;default:null"`
	WordLists   []WordList `gorm:"many2many:word_list_words"`
	CreatedAt   time.Time  `gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime"`
}

type WordInput struct {
	Word        string   `json:"word"`
	Alias       string   `json:"alias,omitempty"`
	WordType    WordType `json:"wordType"`
	Description string   `json:"description,omitempty"`
}

func (input *WordInput) Validate(db *gorm.DB, id int) error {
	var rules []validate.Rule

	// Validate word uniqueness
	rules = append(rules, validate.NewUniqueRule("words", "word", input.Word, nil).Except(id).Say("duplicate word"))

	// Validate alias uniqueness if provided
	if input.Alias != "" {
		rules = append(rules, validate.NewUniqueRule("words", "alias", input.Alias, nil).Except(id).Say("duplicate alias for word"))
	}

	return validate.Validate(db, rules...)
}

func (word Word) Text() string {
	return strings.Join([]string{word.Word, word.Description}, "\n")
}

// PatchWord represents partial updates for Word
type PatchWord struct {
	Word        *string   `json:"word,omitempty"`
	Alias       *string   `json:"alias,omitempty"`
	WordType    *WordType `json:"wordType,omitempty"`
	Description *string   `json:"description,omitempty"`
}

func (input *PatchWord) Validate(db *gorm.DB, id int) error {
	var rules []validate.Rule

	// Validate word uniqueness if provided
	if input.Word != nil && *input.Word != "" {
		rules = append(rules, validate.NewUniqueRule("words", "word", *input.Word, nil).Except(id).Say("duplicate word"))
	}

	// Validate alias uniqueness if provided
	if input.Alias != nil && *input.Alias != "" {
		rules = append(rules, validate.NewUniqueRule("words", "alias", *input.Alias, nil).Except(id).Say("duplicate alias for word"))
	}

	return validate.Validate(db, rules...)
}
