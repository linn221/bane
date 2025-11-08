package services

import (
	"errors"
	"fmt"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/validate"
	"gorm.io/gorm"
)

type aliasService struct {
	db *gorm.DB
}

func newAliasService(db *gorm.DB) *aliasService {
	return &aliasService{db: db}
}

// getReferenceTypePrefix returns the prefix for a given reference type
// Examples: endpoints -> "end", programs -> "prog", words -> "word"
func getReferenceTypePrefix(referenceType models.AliasReferenceType) string {
	switch referenceType {
	case models.AliasReferenceTypeEndpoint:
		return "end"
	case models.AliasReferenceTypeProgram:
		return "prog"
	case models.AliasReferenceTypeWord:
		return "word"
	case models.AliasReferenceTypeWordList:
		return "wl"
	case models.AliasReferenceTypeVuln:
		return "vuln"
	case models.AliasReferenceTypeTag:
		return "tag"
	case models.AliasReferenceTypeMemorySheet:
		return "ms"
	default:
		// Fallback: use first 3 characters of the reference type
		refTypeStr := string(referenceType)
		if len(refTypeStr) >= 3 {
			return refTypeStr[:3]
		}
		return refTypeStr
	}
}

func (s *aliasService) CreateAlias(tx *gorm.DB, referenceType string, referenceId int, alias string) error {
	refType := models.AliasReferenceType(referenceType)

	// If alias is empty, generate it automatically using referenceTypePrefix + referenceId
	if alias == "" {
		prefix := getReferenceTypePrefix(refType)
		alias = fmt.Sprintf("%s%d", prefix, referenceId)
	}

	if err := validate.Validate(tx, validate.NewExistsRule(referenceType, referenceId, errors.New("referencenot found"), nil)); err != nil {
		return err
	}

	// Create new alias
	aliasRecord := models.Alias{
		Name:          alias,
		ReferenceId:   referenceId,
		ReferenceType: refType,
	}
	return tx.Create(&aliasRecord).Error
}

func (s *aliasService) SetAlias(referenceType string, referenceId int, alias string) error {
	refType := models.AliasReferenceType(referenceType)

	// If alias is empty, generate it automatically using referenceTypePrefix + referenceId
	if alias == "" {
		prefix := getReferenceTypePrefix(refType)
		alias = fmt.Sprintf("%s%d", prefix, referenceId)
	}

	// Create new alias
	aliasRecord := models.Alias{
		Name:          alias,
		ReferenceId:   referenceId,
		ReferenceType: refType,
	}

	return s.db.Create(&aliasRecord).Error
}

func (s *aliasService) GetId(alias string) (int, error) {
	var aliasRecord models.Alias
	if err := s.db.Where("name = ?", alias).First(&aliasRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, gorm.ErrRecordNotFound
		}
		return 0, err
	}
	return aliasRecord.ReferenceId, nil
}

func (s *aliasService) GetIdAndType(alias string) (int, string, error) {
	var aliasRecord models.Alias
	if err := s.db.Where("name = ?", alias).First(&aliasRecord).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, "", gorm.ErrRecordNotFound
		}
		return 0, "", err
	}
	return aliasRecord.ReferenceId, string(aliasRecord.ReferenceType), nil
}
