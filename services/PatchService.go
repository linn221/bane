package services

import (
	"context"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

func PatchModel(ctx context.Context, db *gorm.DB, aliasService *aliasService, aliasName string, patches models.PatchInput) (bool, error) {
	alias, err := aliasService.getAliasByName(ctx, aliasName)
	if err != nil {
		return false, err
	}

	modelStruct := toModelStruct(alias.ReferenceType)

	updates := map[string]any{}
	for _, v := range patches.Values {
		updates[v.Key] = v.Value
	}
	for _, v := range patches.ValuesInt {
		updates[v.Key] = v.Value
	}

	err = db.WithContext(ctx).Model(&modelStruct).Where("id = ?", alias.ReferenceId).Updates(updates).Error
	if err != nil {
		return false, err
	}

	return true, nil
}
