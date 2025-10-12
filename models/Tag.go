package models

import "time"

// Tag represents a tag that can be applied to various resources
type Tag struct {
	Id          UInt      `gorm:"primaryKey"`
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
