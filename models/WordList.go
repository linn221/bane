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
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type NewWordList struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func (input *NewWordList) Validate(db *gorm.DB, id int) error {
	err := validate.Validate(db,
		validate.NewUniqueRule("word_lists", "name", input.Name, nil).Except(id).Say("duplicate wordlist name"),
	)
	return err
}

func (wordList WordList) Text() string {
	return strings.Join([]string{wordList.Name, wordList.Description}, "\n")
}
