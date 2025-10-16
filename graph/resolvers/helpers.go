package resolvers

import (
	"time"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

func getMemorySheetsByNextDate(db *gorm.DB, nextDate time.Time) ([]*models.MemorySheet, error) {
	var memorySheets []*models.MemorySheet
	err := db.Where("next_date = ?", nextDate).Find(&memorySheets).Error
	return memorySheets, err
}
func getMemorySheetsByCurrentDate(db *gorm.DB, currentDate time.Time) ([]*models.MemorySheet, error) {
	var memorySheets []*models.MemorySheet
	err := db.Where("current_date = ?", currentDate).Find(&memorySheets).Error
	return memorySheets, err
}
