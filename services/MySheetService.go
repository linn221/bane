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

func getNextRemindingDate(currentDate time.Time, index int) time.Time {
	var intervals = [...]int{1, 1, 5, 7, 14, 19, 30, 50, 60, 80, 100, 120}
	i := min(index, len(intervals)-1)
	return currentDate.AddDate(0, 0, intervals[i])
}

func (s *mySheetService) CreateLabel(ctx context.Context, input *models.MySheetLabelInput) (*models.MySheetLabel, error) {
	label := models.MySheetLabel{
		Name: input.Name,
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&label).Error; err != nil {
			return err
		}
		if err := s.aliasService.CreateAlias(tx, "my_sheet_labels", label.Id, input.Alias); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &label, nil
}

func (s *mySheetService) ListMySheetLabels(ctx context.Context, name *string) ([]*models.MySheetLabel, error) {
	var labels []*models.MySheetLabel
	var err error
	if name != nil {
		err = s.db.WithContext(ctx).Where("name LIKE ?", "%"+*name+"%").Find(&labels).Error
	} else {
		err = s.db.WithContext(ctx).Find(&labels).Error
	}
	return labels, err
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
	if input.Label != nil {
		labelId, err := s.aliasService.GetReferenceId(ctx, *input.Label)
		if err != nil {
			return nil, err
		}
		mysheet.LabelId = labelId
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
		if filter.Label != nil {
			labelId, err := s.aliasService.GetReferenceId(ctx, *filter.Label)
			if err != nil {
				return nil, err
			}
			dbctx.Where("label_id = ?", labelId)
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
		nextDate := getNextRemindingDate(currentDate.Time, existing.Index+1)
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
	currentSheets, err := getMySheetsByPreviousDate(s.db.WithContext(ctx), currentDate)
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
