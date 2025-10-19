package services

import (
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

func (s *GeneralCrud[I, M]) Update(tx *gorm.DB, input *I, id int) (*M, error) {

	var existing M
	if err := tx.First(&existing, id).Error; err != nil {
		return nil, err
	}
	if err := s.validateWrite(tx, input, id); err != nil {
		return nil, err
	}
	err := tx.Model(&existing).Updates(s.updates(existing, input)).Error
	return &existing, err
}

func (s *GeneralCrud[I, M]) UpdateByAlias(tx *gorm.DB, input *I, alias string) (*M, error) {
	id, err := getIdByAlias[M](tx, alias)
	if err != nil {
		return nil, err
	}
	return s.Update(tx, input, id)
}

func (s *GeneralCrud[I, M]) DeleteByAlias(tx *gorm.DB, alias string) (*M, error) {
	id, err := getIdByAlias[M](tx, alias)
	if err != nil {
		return nil, err
	}
	return s.Delete(tx, id)
}

func (s *GeneralCrud[I, M]) Delete(tx *gorm.DB, id int) (*M, error) {

	var existing M
	if err := tx.First(&existing, id).Error; err != nil {
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

func (s *GeneralCrud[I, M]) Get(db *gorm.DB, id int) (*M, error) {
	var result M
	err := db.First(&result, id).Error
	return &result, err
}

func (s *GeneralCrud[I, M]) GetByAlias(db *gorm.DB, alias string) (*M, error) {

	var result M
	err := db.Where("alias = ?", alias).First(&result).Error
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
