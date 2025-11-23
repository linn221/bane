package services

import (
	"context"
	"time"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type mySheetService struct {
	db           *gorm.DB
	aliasService *aliasService
}

func (s *mySheetService) Create(ctx context.Context, input *models.MySheetInput) (*models.MySheet, error) {
	created := utils.Today()
	if input.Date != nil {
		created = input.Date.Time
	}
	mysheet := models.MySheet{
		Title:   input.Title,
		Body:    input.Body,
		Created: models.MyDate{Time: created},
		Index:   0,
	}
	mysheet.PreviousDate = mysheet.Created
	mysheet.NextDate = models.MyDate{Time: mysheet.Created.AddDate(0, 0, 1)}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&mysheet).Error; err != nil {
			return err
		}
		if err := s.aliasService.CreateAlias(tx, "my_sheets", mysheet.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &mysheet, nil
}

func (s *mySheetService) List(ctx context.Context, filter *models.MySheetFilter) ([]*models.MySheet, error) {
	dbctx := s.db.WithContext(ctx).Model(&models.MySheet{})
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

func (s *mySheetService) GetTodaySheets(ctx context.Context, currentDate time.Time) ([]*models.MySheet, error) {
	nextSheets, err := getMySheetsByNextDate(s.db.WithContext(ctx), currentDate)
	if err != nil {
		return nil, err
	}
	tx := s.db.WithContext(ctx).Begin()
	defer tx.Rollback()
	for _, existing := range nextSheets {
		updates := make(map[string]any)
		currentDate := existing.NextDate
		nextDate := GetNextDate(currentDate.Time, existing.Index+1)
		// NextDate has moved for the sheet
		updates["PreviousDate"] = currentDate
		updates["NextDate"] = models.MyDate{Time: nextDate}
		updates["Index"] = existing.Index + 1
		err := tx.Model(&existing).Updates(updates).Error
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
