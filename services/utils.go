package services

import "gorm.io/gorm"

func firstById[T any](db *gorm.DB, id any) (*T, error) {
	var v T
	if err := db.First(&v, id).Error; err != nil {
		return nil, err
	}

	return &v, nil
}

func first[T any](db *gorm.DB, alias string) (*T, error) {
	var v T
	if err := db.Where("alias = ?", alias).First(&v).Error; err != nil {
		return nil, err
	}

	return &v, nil
}

func getIdByAlias[T any](db *gorm.DB, alias string) (int, error) {
	var v T
	var id int
	err := db.Model(&v).Where("alias = ?", alias).Select("id").Scan(&id).Error
	return id, err
}
