package models

import "gorm.io/gorm"

type Project struct {
	Id          int    `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
}

type ProjectInput struct {
	Name        string `json:"name"`
	Alias       string `json:"alias,omitempty"`
	Description string `json:"description,omitempty"`
}

func (input *ProjectInput) Validate(db *gorm.DB, id int) error {
	// Add validation logic here if needed
	return nil
}

type ProjectFilter struct {
	Name   string `json:"name,omitempty"`
	Search string `json:"search,omitempty"`
}
