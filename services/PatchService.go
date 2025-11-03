package services

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

func PatchModel[T any](db *gorm.DB, alias string, patches models.PatchInput) (*T, error) {
	result, err := first[T](db, alias)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{}
	for _, v := range patches.Values {
		updates[v.Key] = v.Value
	}
	for _, v := range patches.ValuesInt {
		updates[v.Key] = v.Value
	}
	err = db.Model(&result).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
