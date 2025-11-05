package services

import (
	"errors"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type aliasService struct {
	db *gorm.DB
}

func newAliasService(db *gorm.DB) *aliasService {
	return &aliasService{db: db}
}

func (s *aliasService) SetAlias(referenceType string, referenceId int, alias string) error {
	if alias == "" {
		return nil
	}

	// Delete existing alias for this reference if it exists
	s.db.Where("reference_type = ? AND reference_id = ?", referenceType, referenceId).Delete(&models.Alias{})

	// Create new alias
	aliasRecord := models.Alias{
		Name:          alias,
		ReferenceId:   referenceId,
		ReferenceType: models.AliasReferenceType(referenceType),
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
