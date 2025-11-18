package services

import (
	"errors"
	"fmt"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type aliasService struct {
	db *gorm.DB
}

func newAliasService(db *gorm.DB) *aliasService {
	return &aliasService{db: db}
}

func (s *aliasService) getAlias(db *gorm.DB, referenceType string, referenceId int) (*models.Alias, error) {
	var alias models.Alias
	if err := db.Where("reference_type = ? AND reference_id = ?", referenceType, referenceId).First(&alias).Error; err != nil {
		return nil, err
	}
	return &alias, nil
}

func (s *aliasService) CreateAlias(tx *gorm.DB, referenceType string, referenceId int, aliasStr string) error {
	refType := models.AliasReferenceType(referenceType)

	// If alias is empty, generate it automatically using referenceTypePrefix + referenceId
	if aliasStr == "" {
		aliasStr = fmt.Sprintf("%s%d", referenceType, referenceId)
	}
	//2d validate if alias name is unique

	// Create new alias
	aliasRecord := models.Alias{
		Name:          aliasStr,
		ReferenceId:   referenceId,
		ReferenceType: refType,
	}
	return tx.Create(&aliasRecord).Error
}

func (s *aliasService) DestroyReference(tx *gorm.DB, alias string) error {

	panic("todo")
}

func (s *aliasService) CatchReference(db *gorm.DB, alias string) *gorm.DB {
	panic("todo")
}

// func (s *aliasService) SetAlias(referenceType string, referenceId int, alias string) error {
// 	refType := models.AliasReferenceType(referenceType)

// 	// If alias is empty, generate it automatically using referenceTypePrefix + referenceId
// 	if alias == "" {
// 		prefix := getReferenceTypePrefix(refType)
// 		alias = fmt.Sprintf("%s%d", prefix, referenceId)
// 	}

// 	// Create new alias
// 	aliasRecord := models.Alias{
// 		Name:          alias,
// 		ReferenceId:   referenceId,
// 		ReferenceType: refType,
// 	}

// 	return s.db.Create(&aliasRecord).Error
// }

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
