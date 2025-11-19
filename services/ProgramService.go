package services

import (
	"context"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

func ListPrograms(db *gorm.DB, search *string) ([]*models.AllProgram, error) {
	var programs []*models.Program
	query := db.Model(&models.Program{})

	if search != nil && *search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ? OR domain ILIKE ?",
			"%"+*search+"%", "%"+*search+"%", "%"+*search+"%")
	}

	err := query.Find(&programs).Error
	if err != nil {
		return nil, err
	}

	// Convert to AllProgram type
	// Note: Alias needs to be fetched from AliasService - this function needs to be updated
	// For now, we'll leave it empty or fetch aliases separately
	var allPrograms []*models.AllProgram
	for _, program := range programs {
		allPrograms = append(allPrograms, &models.AllProgram{
			ID:          program.Id,
			Alias:       "", // TODO: Fetch from AliasService
			Name:        program.Name,
			Description: &program.Description,
			Domain:      program.Domain,
			URL:         program.Url,
		})
	}
	return allPrograms, nil

}

type programService struct {
	GeneralCrud[models.ProgramInput, models.Program]
	db           *gorm.DB
	aliasService *aliasService
}

func (s *programService) Create(ctx context.Context, input *models.ProgramInput) (*models.Program, error) {
	var result *models.Program
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.GeneralCrud.Create(tx, input)
		if err != nil {
			return err
		}
		// Create alias (will be auto-generated if not provided)
		if err := s.aliasService.CreateAlias(tx, "programs", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *programService) Get(ctx context.Context, id *int, alias *string) (*models.Program, error) {
	if id != nil {
		var result models.Program
		err := s.db.WithContext(ctx).First(&result, *id).Error
		return &result, err
	}
	if alias != nil {
		return s.GeneralCrud.GetByAlias(ctx, s.db.WithContext(ctx), s.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *programService) List(ctx context.Context) ([]*models.Program, error) {
	var results []*models.Program
	if err := s.db.WithContext(ctx).Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
