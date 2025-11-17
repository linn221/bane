package services

import (
	"time"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type mySheetService struct {
	GeneralCrud[models.MySheetInput, models.MySheet]
	db           *gorm.DB
	aliasService *aliasService
}

func newMySheetService(db *gorm.DB, aliasService *aliasService) *mySheetService {
	return &mySheetService{
		GeneralCrud: GeneralCrud[models.MySheetInput, models.MySheet]{
			transform: func(input *models.MySheetInput) models.MySheet {
				result := models.MySheet{
					Title: input.Title,
					Body:  input.Body,
				}
				today := utils.Today()
				if input.Date != nil {
					today = input.Date.Time
				}
				result.Created = today
				result.NextDate = today.AddDate(0, 0, 1)
				result.PreviousDate = time.Time{}
				return result
			},
			updates: func(existing models.MySheet, input *models.MySheetInput) map[string]any {
				updates := map[string]any{}

				if input.UpdateNextDate {
					currentDate := existing.NextDate
					nextDate := GetNextDate(currentDate, existing.Index+1)
					// NextDate has moved for the sheet
					updates["PreviousDate"] = currentDate
					updates["NextDate"] = nextDate
					updates["Index"] = existing.Index + 1
				} else { // normal update coming from graphql
					if input.Title != "" {
						updates["Title"] = input.Title
					}
					if input.Body != "" {
						updates["Body"] = input.Body
					}
				}
				return updates
			},
		},
		db:           db,
		aliasService: aliasService,
	}
}

func (mss *mySheetService) Create(input *models.MySheetInput) (*models.MySheet, error) {
	result, err := mss.GeneralCrud.Create(mss.db, input)
	if err != nil {
		return nil, err
	}
	// Set alias (will be auto-generated if not provided)
	if err := mss.aliasService.SetAlias(string(models.AliasReferenceTypeMySheet), result.Id, input.Alias); err != nil {
		return nil, err
	}
	return result, nil
}

func (mss *mySheetService) Get(id *int, alias *string) (*models.MySheet, error) {
	if id != nil {
		return mss.GeneralCrud.Get(mss.db, id)
	}
	if alias != nil {
		return mss.GeneralCrud.GetByAlias(mss.db, mss.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (mss *mySheetService) Update(id *int, alias *string, input *models.MySheetInput) (*models.MySheet, error) {
	if id != nil {
		return mss.GeneralCrud.Update(mss.db, input, id)
	}
	if alias != nil {
		mySheetId, err := mss.aliasService.GetId(*alias)
		if err != nil {
			return nil, err
		}
		result, err := mss.GeneralCrud.Update(mss.db, input, &mySheetId)
		if err != nil {
			return nil, err
		}
		// Set alias if provided
		if input.Alias != "" {
			if err := mss.aliasService.SetAlias(string(models.AliasReferenceTypeMySheet), mySheetId, input.Alias); err != nil {
				return nil, err
			}
		}
		return result, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (mss *mySheetService) Delete(id *int, alias *string) (*models.MySheet, error) {
	if id != nil {
		return mss.GeneralCrud.Delete(mss.db, id)
	}
	if alias != nil {
		return mss.GeneralCrud.DeleteByAlias(mss.db, mss.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (mss *mySheetService) List(filter *models.MySheetFilter) ([]*models.MySheet, error) {
	dbctx := mss.db.Model(&models.MySheet{})
	if filter != nil {
		if filter.Title != "" {
			dbctx = dbctx.Where("title LIKE ?", "%"+filter.Title+"%")
		}
		if filter.Search != "" {
			dbctx = dbctx.Where("title LIKE ? OR body LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
		}
		if filter.NextDate != nil && !filter.NextDate.Time.IsZero() {
			dbctx = dbctx.Where("next_date = ?", filter.NextDate.Time)
		}
		if filter.PreviousDate != nil && !filter.PreviousDate.Time.IsZero() {
			dbctx = dbctx.Where("previous_date = ?", filter.PreviousDate.Time)
		}
	}
	var results []*models.MySheet
	err := dbctx.Find(&results).Error
	return results, err
}

func (mss *mySheetService) GetTodaySheets(currentDate time.Time) ([]*models.MySheet, error) {
	nextSheets, err := getMySheetsByNextDate(mss.db, currentDate)
	if err != nil {
		return nil, err
	}
	tx := mss.db.Begin()
	defer tx.Rollback()
	for _, nSheet := range nextSheets {
		id := nSheet.Id
		_, err := mss.GeneralCrud.Update(tx, &models.MySheetInput{UpdateNextDate: true}, &id)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	currentSheets, err := getMySheetsByPreviousDate(tx, currentDate)
	if err != nil {
		return nil, err
	}
	return currentSheets, err
}

func getMySheetsByNextDate(db *gorm.DB, nextDate time.Time) ([]*models.MySheet, error) {
	var mySheets []*models.MySheet
	err := db.Where("next_date = ?", nextDate).Find(&mySheets).Error
	return mySheets, err
}

func getMySheetsByPreviousDate(db *gorm.DB, previousDate time.Time) ([]*models.MySheet, error) {
	var mySheets []*models.MySheet
	err := db.Where("previous_date = ?", previousDate).Find(&mySheets).Error
	return mySheets, err
}
