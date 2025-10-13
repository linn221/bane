package models

import (
	"strings"
	"time"

	"github.com/linn221/bane/validate"
	"gorm.io/gorm"
)

// Tag represents a tag that can be applied to various resources
type Tag struct {
	Id          int       `gorm:"primaryKey"`
	Name        string    `gorm:"size:100;not null;index"`
	Description string    `gorm:"text;default:null"`
	Alias       string    `gorm:"index;default:null"`
	Priority    int       `gorm:"default:0"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Many-to-many polymorphic relationships through Taggable
	// Taggables []Taggable `gorm:"foreignKey:TagID"`
}

type NewTag struct {
	Name        string `json:"name"`
	Alias       string `json:"alias,omitempty"`
	Description string `json:"description,omitempty"`
	Priority    int    `json:"priority,omitempty"`
}

func (input *NewTag) Validate(db *gorm.DB, id int) error {
	err := validate.Validate(db,
		validate.NewUniqueRule("tags", "name", input.Name, nil).Except(id).Say("duplicate name"),
		validate.NewUniqueRule("tags", "alias", input.Alias, nil).Except(id).Say("duplicate alias for tag"),
	)
	return err
}

func (tag Tag) Text() string {
	return strings.Join([]string{tag.Name, tag.Alias, tag.Description}, "\n")
}
