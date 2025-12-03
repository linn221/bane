package config

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

func migrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Note{},
		&models.Endpoint{},
		&models.Vuln{},
		&models.VulnConnection{},
		&models.Job{},
		&models.Request{},
		&models.WordList{},
		&models.Word{},
		&models.Task{},
		&models.MySheet{},
		&models.MySheetLabel{},
		&models.Project{},
		&models.MyRequest{},
		&models.Alias{},
		// &models.Taggable{},
	)
	if err != nil {
		panic("Error migrating tables: " + err.Error())
	}
}
