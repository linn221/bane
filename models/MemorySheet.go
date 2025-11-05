package models

import (
	"time"

	"gorm.io/gorm"
)

type MemorySheet struct {
	Id          int       `gorm:"primaryKey"`
	Value       string    `gorm:"text;not null"`
	CreateDate  time.Time `gorm:"not null"`
	CurrentDate time.Time `gorm:"index;not null"`
	NextDate    time.Time `gorm:"index;not null"`
	Index       int       `gorm:"not null;default:0"`
	// Notes       []Note    `gorm:"foreignKey:MemorySheetId"`
}

type NewMemorySheet struct {
	Value string  `json:"value"`
	Alias string  `json:"alias,omitempty"`
	Date  *MyDate `json:"date,omitempty"`

	UpdateNextDate bool
}

func (input *NewMemorySheet) Validate(db *gorm.DB, id int) error {
	// Add validation logic here if needed
	return nil
}

// PatchMemorySheet represents partial updates for MemorySheet
type PatchMemorySheet struct {
	Value *string `json:"value,omitempty"`
	Alias *string `json:"alias,omitempty"`
	Date  *MyDate `json:"date,omitempty"`
}

func (input *PatchMemorySheet) Validate(db *gorm.DB, id int) error {
	// Add validation logic here if needed
	return nil
}
