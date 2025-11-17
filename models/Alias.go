package models

type Alias struct {
	Name          string             `gorm:"primaryKey"`
	ReferenceId   int                `gorm:"index;not null"`
	ReferenceType AliasReferenceType `gorm:"index;not null"`
}

type AliasReferenceType string

const (
	AliasReferenceTypeProgram     AliasReferenceType = "programs"
	AliasReferenceTypeWord        AliasReferenceType = "words"
	AliasReferenceTypeWordList    AliasReferenceType = "word_lists"
	AliasReferenceTypeEndpoint    AliasReferenceType = "endpoints"
	AliasReferenceTypeVuln        AliasReferenceType = "vulns"
	AliasReferenceTypeTag         AliasReferenceType = "tags"
	AliasReferenceTypeMemorySheet AliasReferenceType = "memory_sheets"
	AliasReferenceTypeMySheet     AliasReferenceType = "my_sheets"
	AliasReferenceTypeProject     AliasReferenceType = "projects"
	AliasReferenceTypeTodo        AliasReferenceType = "todos"
)
