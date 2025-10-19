package models

import (
	"strings"
)

type Program struct {
	Id          int    `gorm:"primaryKey"`
	Alias       string `gorm:"size:255;not null;unique"`
	Name        string `gorm:"size:255;not null"`
	Url         string `gorm:"not null;index"`
	Domain      string `gorm:"index;not null"` // Store as JSON string
	Description string `gorm:"default:null"`

	// // One-to-many relationships
	// // ImportJobs []ImportJob `gorm:"foreignKey:ProgramId"`
	// Endpoints []Endpoint  `gorm:"foreignKey:ProgramId"`
	// Requests  []MyRequest `gorm:"foreignKey:ProgramId"`

	// // Polymorphic relationships
	// Attachments []Attachment `gorm:"polymorphic:Reference;polymorphicValue:programs"`
	// Notes []Note `gorm:"polymorphic:Reference;polymorphicValue:programs"`
	// Images      []Image      `gorm:"polymorphic:Reference;polymorphicValue:programs"`
	// Taggables   []Taggable   `gorm:"polymorphic:Taggable;polymorphicValue:programs"`
}

type NewProgram struct {
	Alias       string  `json:"alias"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Domain      string  `json:"domain"`
	URL         string  `json:"url"`
}

type AllProgram struct {
	ID          int     `json:"id"`
	Alias       string  `json:"alias"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Domain      string  `json:"domain"`
	URL         string  `json:"url"`
}

func (p *Program) Text() string {
	return strings.Join([]string{p.Name,
		p.Alias,
		p.Domain,
		p.Url,
		p.Description,
	}, "\n")
}
