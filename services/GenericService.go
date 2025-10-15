package services

import (
	"gorm.io/gorm"
)

// service code with no run time dependencies
type GeneralCrud[I any, M any] struct {
	Transform      func(input I) M
	Updates        func(existing M, input I) map[string]any
	ValidateWrite  func(db *gorm.DB, input I, id int) error
	ValidateDelete func(db *gorm.DB, existing M) error
}

func (s *GeneralCrud[I, M]) Create(db *gorm.DB, input I) (*M, error) {
	if s.ValidateWrite != nil {
		if err := s.ValidateWrite(db, input, 0); err != nil {
			return nil, err
		}
	}
	m := s.Transform(input)
	err := db.Create(&m).Error
	return &m, err
}

func (s *GeneralCrud[I, M]) Update(tx *gorm.DB, input I, id int) (*M, error) {

	var existing M
	if err := tx.First(&existing, id).Error; err != nil {
		return nil, err
	}
	if err := s.ValidateWrite(tx, input, id); err != nil {
		return nil, err
	}
	err := tx.Model(&existing).Updates(s.Updates(existing, input)).Error
	return &existing, err
}

func (s *GeneralCrud[I, M]) Delete(tx *gorm.DB, id int) (*M, error) {

	var existing M
	if err := tx.First(&existing, id).Error; err != nil {
		return nil, err
	}
	if vld := s.ValidateDelete; vld != nil {
		err := vld(tx, existing)
		if err != nil {
			return nil, err
		}
	}
	err := tx.Delete(&existing).Error
	return &existing, err
}

func (s *GeneralCrud[I, M]) Get(db *gorm.DB, id int) (*M, error) {
	var result M
	err := db.First(&result, id).Error
	return &result, err
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
