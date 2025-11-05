package models

import (
	"strings"
	"time"

	"github.com/linn221/bane/validate"
	"gorm.io/gorm"
)

// WordList represents a collection of words
type WordList struct {
	Id          int       `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null;index"`
	Description string    `gorm:"text;default:null"`
	Words       []Word    `gorm:"many2many:word_list_words"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type NewWordList struct {
	Name        string `json:"name"`
	Alias       string `json:"alias,omitempty"`
	Description string `json:"description,omitempty"`
}

type AllWordList struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

func (input *NewWordList) Validate(db *gorm.DB, id int) error {
	var rules []validate.Rule

	// Validate name uniqueness
	rules = append(rules, validate.NewUniqueRule("word_lists", "name", input.Name, nil).Except(id).Say("duplicate wordlist name"))

	// Validate alias uniqueness if provided
	if input.Alias != "" {
		rules = append(rules, validate.NewUniqueRule("word_lists", "alias", input.Alias, nil).Except(id).Say("duplicate alias for wordlist"))
	}

	return validate.Validate(db, rules...)
}

func (wordList WordList) Text() string {
	return strings.Join([]string{wordList.Name, wordList.Description}, "\n")
}

// PatchWordList represents partial updates for WordList
type PatchWordList struct {
	Name        *string `json:"name,omitempty"`
	Alias       *string `json:"alias,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (input *PatchWordList) Validate(db *gorm.DB, id int) error {
	var rules []validate.Rule

	// Validate name uniqueness if provided
	if input.Name != nil && *input.Name != "" {
		rules = append(rules, validate.NewUniqueRule("word_lists", "name", *input.Name, nil).Except(id).Say("duplicate name"))
	}

	// Validate alias uniqueness if provided
	if input.Alias != nil && *input.Alias != "" {
		rules = append(rules, validate.NewUniqueRule("word_lists", "alias", *input.Alias, nil).Except(id).Say("duplicate alias for wordlist"))
	}

	return validate.Validate(db, rules...)
}
