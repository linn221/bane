package services

import (
	"errors"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"github.com/linn221/bane/validate"
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

func newProgramService(db *gorm.DB, aliasService *aliasService) *programService {
	return &programService{
		GeneralCrud: GeneralCrud[models.ProgramInput, models.Program]{
			transform: func(input *models.ProgramInput) models.Program {
				return models.Program{
					Name:        input.Name,
					Url:         input.URL,
					Description: utils.SafeDeref(input.Description),
					Domain:      input.Domain,
				}
			},
			updates: func(existing models.Program, input *models.ProgramInput) map[string]any {
				return map[string]any{
					"Name":        input.Name,
					"Url":         input.URL,
					"Description": utils.SafeDeref(input.Description),
					"Domain":      input.Domain,
				}
			},
			validateWrite: func(db *gorm.DB, input *models.ProgramInput, id int) error {
				// Check alias uniqueness using AliasService
				if input.Alias != "" {
					existingId, err := aliasService.GetId(input.Alias)
					if err == nil && existingId != id {
						return errors.New("duplicate program alias")
					}
					// If err is gorm.ErrRecordNotFound, alias doesn't exist - that's fine
				}
				return validate.Validate(db,
					validate.NewUniqueRule("programs", "name", input.Name, nil).Except(id).Say("duplicate program name"))
			},
			validateDelete: func(db *gorm.DB, existing models.Program) error {
				return nil
			},
		},
		db:           db,
		aliasService: aliasService,
	}
}

func (ps *programService) Create(input *models.ProgramInput) (*models.Program, error) {
	result, err := ps.GeneralCrud.Create(ps.db, input)
	if err != nil {
		return nil, err
	}
	// Set alias (will be auto-generated if not provided)
	if err := ps.aliasService.SetAlias(string(models.AliasReferenceTypeProgram), result.Id, input.Alias); err != nil {
		return nil, err
	}
	return result, nil
}

func (ps *programService) Get(id *int, alias *string) (*models.Program, error) {
	if id != nil {
		return ps.GeneralCrud.Get(ps.db, id)
	}
	if alias != nil {
		return ps.GeneralCrud.GetByAlias(ps.db, ps.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ps *programService) List() ([]*models.Program, error) {
	var results []*models.Program
	if err := ps.db.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
