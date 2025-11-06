package services

import (
	"errors"

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

func (s *aliasService) CreateAlias(tx *gorm.DB, referenceType string, referenceId int, alias string) error {
	if alias == "" {
		return nil
	}
	if err := validate.Validate(tx, validate.NewExistsRule(referenceType, referenceId, errors.New("referencenot found"), nil)); err != nil {
		return err
	}

	// Create new alias
	aliasRecord := models.Alias{
		Name:          alias,
		ReferenceId:   referenceId,
		ReferenceType: models.AliasReferenceType(referenceType),
	}
	return tx.Create(&aliasRecord).Error
}

func (s *aliasService) SetAlias(referenceType string, referenceId int, alias string) error {
	if alias == "" {
		return nil
	}

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
