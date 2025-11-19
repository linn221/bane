package models

type Alias struct {
	Name          string `gorm:"primaryKey"`
	ReferenceId   int    `gorm:"index;not null"`
	ReferenceType string `gorm:"index;not null"`
}
