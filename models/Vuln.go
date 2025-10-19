package models

type Vuln struct {
	Id          int    `gorm:"primaryKey"`
	Name        string `gorm:"size:255;not null"`
	Alias       string `gorm:"index;default:null"`
	Description string `gorm:"default:null"`
	// Notes relationship is handled through polymorphic pattern in Note model
	// Notes can be linked via ReferenceType="vulns" and ReferenceID=Id
}

type NewVuln struct {
	Name        string `json:"name"`
	Alias       string `json:"alias,omitempty"`
	Description string `json:"description,omitempty"`
}

type VulnReferenceType string

const (
	VulnReferenceTypeProgram    VulnReferenceType = "programs"
	VulnReferenceTypeEndpoint   VulnReferenceType = "endpoints"
	VulnReferenceTypeRequest    VulnReferenceType = "requests"
	VulnReferenceTypeNote       VulnReferenceType = "notes"
	VulnReferenceTypeAttachment VulnReferenceType = "attachments"
	VulnReferenceTypeVuln       VulnReferenceType = "vulns"
)

type VulnConnection struct {
	VulnId        int               `gorm:"not null;index"`
	ReferenceId   int               `gorm:"not null;index"`
	ReferenceType VulnReferenceType `gorm:"not null;index"`
}
