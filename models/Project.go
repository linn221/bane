package models

type Project struct {
	Id          int    `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string `gorm:"type:text"`
	Url         string `gorm:"type:text;default:null"`
}

type ProjectInput struct {
	Name        string `json:"name"`
	Alias       string `json:"alias,omitempty"`
	Description string `json:"description,omitempty"`
	Url         string `json:"url,omitempty"`
}

type ProjectFilter struct {
	Name   string `json:"name,omitempty"`
	Search string `json:"search,omitempty"`
}
