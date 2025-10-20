package services

import (
	"time"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

func GetNextDate(currentDate time.Time, index int) time.Time {
	var intervals = [...]int{1, 1, 5, 7, 14, 19, 30, 50, 60, 80, 100, 120}
	i := min(index, len(intervals)-1)
	return currentDate.AddDate(0, 0, intervals[i])
}
func GetTodayNotes(db *gorm.DB, currentDate time.Time) ([]*models.MemorySheet, error) {
	nextSheets, err := getMemorySheetsByNextDate(db, currentDate)
	if err != nil {
		return nil, err
	}
	tx := db.Begin()
	defer tx.Rollback()
	for _, nSheet := range nextSheets {
		id := nSheet.Id
		_, err := MemorySheetCrud.Update(tx, &models.NewMemorySheet{UpdateNextDate: true}, &id)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	currentSheets, err := getMemorySheetsByCurrentDate(tx, currentDate)
	if err != nil {
		return nil, err
	}
	return currentSheets, err
}

var MemorySheetCrud = GeneralCrud[models.NewMemorySheet, models.MemorySheet]{
	transform: func(input *models.NewMemorySheet) models.MemorySheet {
		result := models.MemorySheet{
			Value: input.Value,
		}
		result.CreateDate = utils.Today()
		result.CurrentDate = result.CreateDate
		result.NextDate = result.CurrentDate.AddDate(0, 0, 1)
		return result
	},
	updates: func(existing models.MemorySheet, input *models.NewMemorySheet) map[string]any {
		updates := map[string]any{}

		if input.UpdateNextDate {
			currentDate := existing.NextDate
			nextDate := GetNextDate(currentDate, existing.Index+1)
			// NextDate has moved for the note
			updates["CurrentDate"] = currentDate
			updates["NextDate"] = nextDate
			updates["Index"] = existing.Index + 1
		} else { // normal update coming from graphql
			if input.Value != "" {
				updates["Value"] = input.Value
			}
		}
		return updates
	},
	// ValidateWrite: func(db *gorm.DB, input models.NewMemorySheet, id int) error {
	// 	return input.Validate(db, id)
	// },
}

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
