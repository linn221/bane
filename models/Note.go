package models

import "strings"

// Note represents a note linked to a program, endpoint, or request
type Note struct {
	Id            int    `gorm:"primaryKey"`
	ReferenceType string `gorm:"size:20;not null;index"` // "programs", "endpoints", "requests"
	ReferenceID   int    `gorm:"not null;index"`         // Changed from ReferenceId to ReferenceID for GORM polymorphic
	Value         string `gorm:"type:text;not null"`
	NoteDate      MyDate `gorm:"index; not null"`

	// Polymorphic relationships
	Taggables []Taggable `gorm:"polymorphic:Taggable;polymorphicValue:notes"`
}

type NoteInput struct {
	Value         string
	RId           int
	ReferenceId   int
	ReferenceType string
}

type NoteFilter struct {
	RID           int    `json:"rId,omitempty"`
	ReferenceID   int    `json:"referenceId,omitempty"`
	ReferenceType string `json:"referenceType,omitempty"`
	NoteDate      MyDate `json:"noteDate,omitempty"`
	Search        string `json:"search,omitempty"`
}

func (n *Note) Text() string {
	return strings.Join([]string{n.Value, n.NoteDate.String()}, "\n")
}
