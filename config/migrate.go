package config

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(&models.Tag{})
	if err != nil {
		panic("Error migrating tables: " + err.Error())
	}
}
