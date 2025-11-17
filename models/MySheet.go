package models

import (
	"time"

	"gorm.io/gorm"
)

type MySheet struct {
	Id           int       `gorm:"primaryKey"`
	Title        string    `gorm:"not null"`
	Body         string    `gorm:"type:text"`
	Created      time.Time `gorm:"not null"`
	NextDate     time.Time `gorm:"index;not null"`
	PreviousDate time.Time `gorm:"index;not null"`
	Index        int       `gorm:"not null;default:0"`
	// Age is calculated, not stored
}

type MySheetInput struct {
	Title string  `json:"title"`
	Body  string  `json:"body"`
	Alias string  `json:"alias,omitempty"`
	Date  *MyDate `json:"date,omitempty"`

	UpdateNextDate bool
}

func (input *MySheetInput) Validate(db *gorm.DB, id int) error {
	// Add validation logic here if needed
	return nil
}

type MySheetFilter struct {
	Title        string  `json:"title,omitempty"`
	Search       string  `json:"search,omitempty"`
	NextDate     *MyDate `json:"nextDate,omitempty"`
	PreviousDate *MyDate `json:"previousDate,omitempty"`
}
