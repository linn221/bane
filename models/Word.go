package models

import (
	"strings"
	"time"

	"github.com/linn221/bane/validate"
	"gorm.io/gorm"
)

// Word represents a word in a wordlist
type Word struct {
	Id          int       `gorm:"primaryKey"`
	Word        string    `gorm:"size:255;not null;index"`
	WordType    WordType  `gorm:"not null;index"`
	Description string    `gorm:"text;default:null"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

type NewWord struct {
	Word        string   `json:"word"`
	WordType    WordType `json:"wordType"`
	Description string   `json:"description,omitempty"`
}

func (input *NewWord) Validate(db *gorm.DB, id int) error {
	err := validate.Validate(db,
		validate.NewUniqueRule("words", "word", input.Word, nil).Except(id).Say("duplicate word"),
	)
	return err
}

func (word Word) Text() string {
	return strings.Join([]string{word.Word, word.Description}, "\n")
}
