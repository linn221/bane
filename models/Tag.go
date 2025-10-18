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

	Taggables []Taggable `gorm:"foreignKey:TagID"`
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

// TaggableType represents the type of taggable resource
type TaggableType string

const (
	TaggableTypePrograms  TaggableType = "programs"
	TaggableTypeEndpoints TaggableType = "endpoints"
	TaggableTypeRequests  TaggableType = "requests"
	TaggableTypeVulns     TaggableType = "vulns"
	TaggableTypeNotes     TaggableType = "notes"
)

// Taggable represents the many-to-many polymorphic relationship between tags and resources
type Taggable struct {
	ID           int       `gorm:"primaryKey"`
	TagID        int       `gorm:"column:tag_id;not null;index"`
	TaggableType string    `gorm:"column:taggable_type;size:20;not null;index"` // "programs", "endpoints", "requests", "vulns", "notes"
	TaggableID   int       `gorm:"column:taggable_id;not null;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`

	// Relationships
	Tag Tag `gorm:"foreignKey:TagID"`
}
