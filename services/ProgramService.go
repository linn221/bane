package services

import (
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
	var allPrograms []*models.AllProgram
	for _, program := range programs {
		allPrograms = append(allPrograms, &models.AllProgram{
			ID:          program.Id,
			Alias:       program.Alias,
			Name:        program.Name,
			Description: &program.Description,
			Domain:      program.Domain,
			URL:         program.Url,
		})
	}
	return allPrograms, nil

}

type programService struct {
	GeneralCrud[models.NewProgram, models.Program]
}

var ProgramService = programService{
	GeneralCrud: GeneralCrud[models.NewProgram, models.Program]{
		transform: func(input *models.NewProgram) models.Program {
			return models.Program{
				Alias:       input.Alias,
				Name:        input.Name,
				Url:         input.URL,
				Description: utils.SafeDeref(input.Description),
				Domain:      input.Domain,
			}
		},
		updates: func(existing models.Program, input *models.NewProgram) map[string]any {
			return map[string]any{
				"Alias":       input.Alias,
				"Name":        input.Name,
				"Url":         input.URL,
				"Description": utils.SafeDeref(input.Description),
				"Domain":      input.Domain,
			}
		},
		validateWrite: func(db *gorm.DB, input *models.NewProgram, id int) error {
			return validate.Validate(db,
				validate.NewUniqueRule("programs", "name", input.Name, nil).Except(id).Say("duplicate program name"),
				validate.NewUniqueRule("programs", "alias", input.Alias, nil).Except(id).Say("duplicate program alias"))
		},
		validateDelete: func(db *gorm.DB, existing models.Program) error {
			return nil
		},
	},
}

func (crud *programService) ListPrograms(db *gorm.DB) ([]*models.Program, error) {
	var results []*models.Program
	if err := db.Find(&results).Error; err != nil {
		return nil, err
	}

	return results, nil
}
