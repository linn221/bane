package models

type Program struct {
	Id          int    `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null"`
	Url         string `gorm:"not null;index"`
	Description string `gorm:"default:null"`
	Domain      string `gorm:"type:index;not null"` // Store as JSON string

	// // One-to-many relationships
	// // ImportJobs []ImportJob `gorm:"foreignKey:ProgramId"`
	// Endpoints []Endpoint  `gorm:"foreignKey:ProgramId"`
	// Requests  []MyRequest `gorm:"foreignKey:ProgramId"`

	// // Polymorphic relationships
	// Attachments []Attachment `gorm:"polymorphic:Reference;polymorphicValue:programs"`
	// Notes       []Note       `gorm:"polymorphic:Reference;polymorphicValue:programs"`
	// Images      []Image      `gorm:"polymorphic:Reference;polymorphicValue:programs"`
	// Taggables   []Taggable   `gorm:"polymorphic:Taggable;polymorphicValue:programs"`
}
