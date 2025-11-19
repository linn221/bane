package services

import (
	"context"

	"gorm.io/gorm"
)

// service code with no run time dependencies
type GeneralCrud[I any, M any] struct {
	transform      func(input *I) M
	updates        func(existing M, input *I) map[string]any
	validateWrite  func(db *gorm.DB, input *I, id int) error
	validateDelete func(db *gorm.DB, existing M) error
}

func (s *GeneralCrud[I, M]) Create(db *gorm.DB, input *I) (*M, error) {
	if s.validateWrite != nil {
		if err := s.validateWrite(db, input, 0); err != nil {
			return nil, err
		}
	}
	m := s.transform(input)
	err := db.Create(&m).Error
	return &m, err
}

func (s *GeneralCrud[I, M]) Update(tx *gorm.DB, input *I, id *int) (*M, error) {

	var existing M
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	if err := tx.First(&existing, *id).Error; err != nil {
		return nil, err
	}
	if err := s.validateWrite(tx, input, *id); err != nil {
		return nil, err
	}
	err := tx.Model(&existing).Updates(s.updates(existing, input)).Error
	return &existing, err
}

func (s *GeneralCrud[I, M]) UpdateByAlias(ctx context.Context, tx *gorm.DB, aliasService *aliasService, input *I, alias string) (*M, error) {
	id, err := getIdByAlias[M](ctx, tx, aliasService, alias)
	if err != nil {
		return nil, err
	}
	return s.Update(tx, input, &id)
}

func (s *GeneralCrud[I, M]) DeleteByAlias(ctx context.Context, tx *gorm.DB, aliasService *aliasService, alias string) (*M, error) {
	id, err := getIdByAlias[M](ctx, tx, aliasService, alias)
	if err != nil {
		return nil, err
	}
	return s.Delete(tx, &id)
}

func (s *GeneralCrud[I, M]) Delete(tx *gorm.DB, id *int) (*M, error) {

	var existing M
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	if err := tx.First(&existing, *id).Error; err != nil {
		return nil, err
	}
	if vld := s.validateDelete; vld != nil {
		err := vld(tx, existing)
		if err != nil {
			return nil, err
		}
	}
	err := tx.Delete(&existing).Error
	return &existing, err
}

func (s *GeneralCrud[I, M]) Get(db *gorm.DB, id *int) (*M, error) {
	var result M
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	err := db.First(&result, *id).Error
	return &result, err
}

func (s *GeneralCrud[I, M]) GetByAlias(ctx context.Context, db *gorm.DB, aliasService *aliasService, alias string) (*M, error) {
	var result M
	id, _, err := aliasService.GetIdAndType(ctx, alias)
	if err != nil {
		return nil, err
	}
	err = db.WithContext(ctx).First(&result, id).Error
	return &result, err
}

// Patch updates only the fields provided in the updates map
func (s *GeneralCrud[I, M]) Patch(tx *gorm.DB, updates map[string]any, id *int) (*M, error) {
	var existing M
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	if err := tx.First(&existing, *id).Error; err != nil {
		return nil, err
	}

	// If no fields to update, return existing
	if len(updates) == 0 {
		return &existing, nil
	}

	// Note: Patch does not call validateWrite (as requested)
	err := tx.Model(&existing).Updates(updates).Error
	return &existing, err
}

// PatchByAlias updates only the fields provided in the updates map using alias
func (s *GeneralCrud[I, M]) PatchByAlias(ctx context.Context, tx *gorm.DB, aliasService *aliasService, updates map[string]any, alias string) (*M, error) {
	id, err := getIdByAlias[M](ctx, tx, aliasService, alias)
	if err != nil {
		return nil, err
	}
	return s.Patch(tx, updates, &id)
}

// func (s *GeneralCrud[I, M]) List(db *gorm.DB, filter func(*gorm.DB) *gorm.DB) ([]*M, error) {
// 	var v M
// 	dbCtx := db.WithContext(ctx).Model(&v)
// 	if filter != nil {
// 		dbCtx = filter(dbCtx)
// 	}
// 	var results []*M
// 	if err := dbCtx.Find(&results).Error; err != nil {
// 		return nil, err
// 	}

// 	return results, nil
// }
