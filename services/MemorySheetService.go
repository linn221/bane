package services

import (
	"context"
	"time"

	"github.com/linn221/bane/models"
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

func (s *memorySheetService) Create(ctx context.Context, input *models.MemorySheetInput) (*models.MemorySheet, error) {
	var result *models.MemorySheet
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.GeneralCrud.Create(tx, input)
		if err != nil {
			return err
		}
		// Create alias (will be auto-generated if not provided)
		if err := s.aliasService.CreateAlias(tx, "memory_sheets", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *memorySheetService) Get(ctx context.Context, id *int) (*models.MemorySheet, error) {
	if id == nil {
		return nil, gorm.ErrRecordNotFound
	}
	var result models.MemorySheet
	err := s.db.WithContext(ctx).First(&result, *id).Error
	return &result, err
}

func (s *memorySheetService) Update(ctx context.Context, id *int, alias *string, input *models.MemorySheetInput) (*models.MemorySheet, error) {
	if id != nil {
		return s.GeneralCrud.Update(s.db.WithContext(ctx), input, id)
	}
	if alias != nil {
		memorySheetId, err := s.aliasService.GetReferenceId(ctx, *alias)
		if err != nil {
			return nil, err
		}
		var result *models.MemorySheet
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			result, err = s.GeneralCrud.Update(tx, input, &memorySheetId)
			if err != nil {
				return err
			}
			// Create alias if provided
			if input.Alias != "" {
				if err := s.aliasService.CreateAlias(tx, "memory_sheets", memorySheetId, input.Alias); err != nil {
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

func (s *memorySheetService) GetTodayNotes(ctx context.Context, currentDate time.Time) ([]*models.MemorySheet, error) {
	nextSheets, err := getMemorySheetsByNextDate(s.db.WithContext(ctx), currentDate)
	if err != nil {
		return nil, err
	}
	tx := s.db.WithContext(ctx).Begin()
	defer tx.Rollback()
	for _, nSheet := range nextSheets {
		id := nSheet.Id
		_, err := s.GeneralCrud.Update(tx, &models.MemorySheetInput{UpdateNextDate: true}, &id)
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
