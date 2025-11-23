package models

type Alias struct {
	Id            int    `gorm:"primaryKey"`
	Name          string `gorm:"unique;not null"`
	ReferenceId   int    `gorm:"index;not null"`
	ReferenceType string `gorm:"index;not null"`
}
