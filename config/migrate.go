package config

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Note{},
		&models.Endpoint{},
		&models.Job{},
		&models.Request{},
		&models.WordList{},
		&models.Word{},
		&models.Project{},
		&models.MyRequest{},
		&models.Alias{},
		// &models.Taggable{},
	)
	if err != nil {
		panic("Error migrating tables: " + err.Error())
	}
}
