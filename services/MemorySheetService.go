package services

import (
	"context"
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

type memorySheetService struct {
	GeneralCrud[models.MemorySheetInput, models.MemorySheet]
	db           *gorm.DB
	aliasService *aliasService
}

func newMemorySheetService(db *gorm.DB, aliasService *aliasService) *memorySheetService {
	return &memorySheetService{
		GeneralCrud: GeneralCrud[models.MemorySheetInput, models.MemorySheet]{
			transform: func(input *models.MemorySheetInput) models.MemorySheet {
				result := models.MemorySheet{
					Value: input.Value,
				}
				result.CreateDate = utils.Today()
				result.CurrentDate = result.CreateDate
				result.NextDate = result.CurrentDate.AddDate(0, 0, 1)
				return result
			},
			updates: func(existing models.MemorySheet, input *models.MemorySheetInput) map[string]any {
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
		},
		db:           db,
		aliasService: aliasService,
	}
}

func (mss *memorySheetService) Create(input *models.MemorySheetInput) (*models.MemorySheet, error) {
	var result *models.MemorySheet
	err := mss.db.Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = mss.GeneralCrud.Create(tx, input)
		if err != nil {
			return err
		}
		// Create alias (will be auto-generated if not provided)
		if err := mss.aliasService.CreateAlias(tx, "memory_sheets", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (mss *memorySheetService) Get(id *int) (*models.MemorySheet, error) {
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return mss.GeneralCrud.Get(mss.db, id)
}

func (mss *memorySheetService) Update(id *int, alias *string, input *models.MemorySheetInput) (*models.MemorySheet, error) {
	if id != nil {
		return mss.GeneralCrud.Update(mss.db, input, id)
	}
	if alias != nil {
		// Note: GetId doesn't have context, but we'll keep it for now since Update doesn't have context
		// TODO: Add context to Update method if needed
		memorySheetId, err := mss.aliasService.GetId(*alias)
		if err != nil {
			return nil, err
		}
		var result *models.MemorySheet
		err = mss.db.Transaction(func(tx *gorm.DB) error {
			result, err = mss.GeneralCrud.Update(tx, input, &memorySheetId)
			if err != nil {
				return err
			}
			// Create alias if provided
			if input.Alias != "" {
				if err := mss.aliasService.CreateAlias(tx, "memory_sheets", memorySheetId, input.Alias); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (mss *memorySheetService) Patch(id *int, alias *string, updates map[string]any) (*models.MemorySheet, error) {
	if id != nil {
		return mss.GeneralCrud.Patch(mss.db, updates, id)
	}
	if alias != nil {
		return mss.GeneralCrud.PatchByAlias(context.Background(), mss.db, mss.aliasService, updates, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (mss *memorySheetService) Delete(id *int, alias *string) (*models.MemorySheet, error) {
	if id != nil {
		return mss.GeneralCrud.Delete(mss.db, id)
	}
	if alias != nil {
		return mss.GeneralCrud.DeleteByAlias(context.Background(), mss.db, mss.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (mss *memorySheetService) GetTodayNotes(currentDate time.Time) ([]*models.MemorySheet, error) {
	nextSheets, err := getMemorySheetsByNextDate(mss.db, currentDate)
	if err != nil {
		return nil, err
	}
	tx := mss.db.Begin()
	defer tx.Rollback()
	for _, nSheet := range nextSheets {
		id := nSheet.Id
		_, err := mss.GeneralCrud.Update(tx, &models.MemorySheetInput{UpdateNextDate: true}, &id)
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
