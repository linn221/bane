package models

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

type ProjectFilter struct {
	Name   string `json:"name,omitempty"`
	Search string `json:"search,omitempty"`
}
