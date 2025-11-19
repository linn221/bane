package services

import (
	"context"
	"time"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type mySheetService struct {
	GeneralCrud[models.MySheetInput, models.MySheet]
	db           *gorm.DB
	aliasService *aliasService
}

func (s *mySheetService) Create(ctx context.Context, input *models.MySheetInput) (*models.MySheet, error) {
	var result *models.MySheet
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.GeneralCrud.Create(tx, input)
		if err != nil {
			return err
		}
		// Create alias (will be auto-generated if not provided)
		if err := s.aliasService.CreateAlias(tx, "my_sheets", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *mySheetService) Get(ctx context.Context, id *int, alias *string) (*models.MySheet, error) {
	if id != nil {
		var result models.MySheet
		err := s.db.WithContext(ctx).First(&result, *id).Error
		return &result, err
	}
	if alias != nil {
		return s.GeneralCrud.GetByAlias(ctx, s.db.WithContext(ctx), s.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *mySheetService) Update(ctx context.Context, id *int, alias *string, input *models.MySheetInput) (*models.MySheet, error) {
	if id != nil {
		return s.GeneralCrud.Update(s.db.WithContext(ctx), input, id)
	}
	if alias != nil {
		mySheetId, err := s.aliasService.GetReferenceId(ctx, *alias)
		if err != nil {
			return nil, err
		}
		var result *models.MySheet
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			result, err = s.GeneralCrud.Update(tx, input, &mySheetId)
			if err != nil {
				return err
			}
			// Create alias if provided
			if input.Alias != "" {
				if err := s.aliasService.CreateAlias(tx, "my_sheets", mySheetId, input.Alias); err != nil {
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
	for _, nSheet := range nextSheets {
		id := nSheet.Id
		_, err := s.GeneralCrud.Update(tx, &models.MySheetInput{UpdateNextDate: true}, &id)
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
