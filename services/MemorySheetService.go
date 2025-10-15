package services

import (
	"time"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
)

func GetNextDate(currentDate time.Time, index int) time.Time {
	var intervals = [...]int{1, 1, 5, 7, 14, 19, 30, 50, 60, 80, 100, 120}
	i := min(index, len(intervals)-1)
	return currentDate.AddDate(0, 0, intervals[i])
}

var MemorySheetCrud = GeneralCrud[models.NewMemorySheet, models.MemorySheet]{
	Transform: func(input models.NewMemorySheet) models.MemorySheet {
		result := models.MemorySheet{
			Value: input.Value,
		}
		if input.Date != nil {
			result.CreateDate = input.Date.Time
		} else {
			result.CreateDate = utils.Today()
		}
		return result
	},
	Updates: func(existing models.MemorySheet, input models.NewMemorySheet) map[string]any {
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

			if input.Date != nil {
				updates["CreateDate"] = input.Date.Time
			}
		}
		return updates
	},
	// ValidateWrite: func(db *gorm.DB, input models.NewMemorySheet, id int) error {
	// 	return input.Validate(db, id)
	// },
}
