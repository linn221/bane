package services

import (
	"context"

	"gorm.io/gorm"
)

func firstById[T any](db *gorm.DB, id any) (*T, error) {
	var v T
	if err := db.First(&v, id).Error; err != nil {
		return nil, err
	}

	return &v, nil
}

func first[T any](ctx context.Context, db *gorm.DB, aliasService *aliasService, alias string) (*T, error) {
	var v T
	// Use AliasService to get the ID
	id, err := aliasService.GetReferenceId(ctx, alias)
	if err != nil {
		return nil, err
	}

	// Get the record by ID - GORM will infer the table from the type
	if err := db.WithContext(ctx).First(&v, id).Error; err != nil {
		return nil, err
	}

	return &v, nil
}

func getIdByAlias[T any](ctx context.Context, db *gorm.DB, aliasService *aliasService, alias string) (int, error) {
	return aliasService.GetReferenceId(ctx, alias)
}
