package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type aliasService struct {
	db *gorm.DB
}

func (s *aliasService) getAliasByReference(db *gorm.DB, referenceType string, referenceId int) (*models.Alias, error) {
	var alias models.Alias
	if err := db.Where("reference_type = ? AND reference_id = ?", referenceType, referenceId).First(&alias).Error; err != nil {
		return nil, err
	}
	return &alias, nil
}

func (s *aliasService) getAliasByName(ctx context.Context, aliasStr string) (*models.Alias, error) {
	var alias models.Alias
	if err := s.db.WithContext(ctx).Where("name = ?", aliasStr).First(&alias).Error; err != nil {
		return nil, err
	}
	return &alias, nil
}

func (s *aliasService) CreateAlias(tx *gorm.DB, referenceType string, referenceId int, aliasStr string) error {
	// If alias is empty, generate it automatically using referenceTypePrefix + referenceId
	if aliasStr == "" {
		aliasStr = fmt.Sprintf("%s%d", referenceType, referenceId)
	}
	//2d validate if alias name is unique

	// Create new alias
	aliasRecord := models.Alias{
		Name:          aliasStr,
		ReferenceId:   referenceId,
		ReferenceType: referenceType,
	}
	return tx.Create(&aliasRecord).Error
}

func (s *aliasService) DestroyReference(ctx context.Context, aliasStr string) (bool, error) {

	aliasRecord, err := s.getAliasByName(ctx, aliasStr)
	if err != nil {
		return false, err
	}
	tx := s.db.WithContext(ctx).Begin()
	defer tx.Rollback()
	sql := fmt.Sprintf("DELETE FROM %s WHERE id = ?", aliasRecord.ReferenceType)
	if err := tx.Exec(sql, aliasRecord.ReferenceId).Error; err != nil {
		return false, err
	}
	if err := tx.Delete(&aliasRecord).Error; err != nil {
		return false, err
	}
	if err := tx.Commit().Error; err != nil {
		return false, err
	}
	return true, nil
}

func (s *aliasService) ScopeReference(ctx context.Context, db *gorm.DB, aliasName string) (*gorm.DB, error) {
	aliasRecord, err := s.getAliasByName(ctx, aliasName)
	if err != nil {
		return nil, err
	}
	return db.Table(aliasRecord.ReferenceType).Where("id = ?", aliasRecord.ReferenceId), nil
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
		return 0, err
	}
	return aliasRecord.ReferenceId, nil
}

func (s *aliasService) GetReferenceId(ctx context.Context, alias string) (int, error) {
	arecord, err := s.getAliasByName(ctx, alias)
	if err != nil {
		return 0, err
	}
	return arecord.ReferenceId, nil
}

func (s *aliasService) GetIdAndType(ctx context.Context, alias string) (int, string, error) {
	aliasRecord, err := s.getAliasByName(ctx, alias)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, "", gorm.ErrRecordNotFound
		}
		return 0, "", err
	}
	return aliasRecord.ReferenceId, aliasRecord.ReferenceType, nil
}

func GetRecordByAlias[T any](db *gorm.DB, refType string, alias string) (*T, error) {
	var a models.Alias
	if err := db.Where("name = ?", alias).First(&a).Error; err != nil {
		return nil, err
	}
	var v T
	if err := db.Where("id = ?", a.ReferenceId).First(&v).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func GetRecordByAliasOrId[T any](db *gorm.DB, refType string, alias *string, id *int) (*T, error) {
	var v T
	if alias != nil {
		var a models.Alias
		if err := db.Where("name = ?", alias).First(&a).Error; err != nil {
			return nil, err
		}
		if err := db.Where("id = ?", a.ReferenceId).First(&v).Error; err != nil {
			return nil, err
		}
	}
	if id != nil {
		if err := db.First(&v, *id).Error; err != nil {
			return nil, err
		}
	}
	return &v, nil
}
