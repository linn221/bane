package models

import (
	"strings"

	"github.com/linn221/bane/validate"
	"gorm.io/gorm"
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

// PatchProgram represents partial updates for Program
type PatchProgram struct {
	Name        *string `json:"name,omitempty"`
	Alias       *string `json:"alias,omitempty"`
	Description *string `json:"description,omitempty"`
	Domain      *string `json:"domain,omitempty"`
	URL         *string `json:"url,omitempty"`
}

func (input *PatchProgram) Validate(db *gorm.DB, id int) error {
	var rules []validate.Rule

	// Validate alias uniqueness if provided
	if input.Alias != nil && *input.Alias != "" {
		rules = append(rules, validate.NewUniqueRule("programs", "alias", *input.Alias, nil).Except(id).Say("duplicate alias for program"))
	}

	return validate.Validate(db, rules...)
}
