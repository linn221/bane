package services

import (
	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"github.com/linn221/bane/validate"
	"gorm.io/gorm"
)

var ProgramCrud = GeneralCrud[models.NewProgram, models.Program]{
	Transform: func(input models.NewProgram) models.Program {
		return models.Program{
			Name:        input.Name,
			Url:         input.URL,
			Description: utils.SafeDeref(input.Description),
			Domain:      input.Domain,
		}
	},
	Updates: func(existing models.Program, input models.NewProgram) map[string]any {
		return map[string]any{
			"Name":        input.Name,
			"Url":         input.URL,
			"Description": utils.SafeDeref(input.Description),
			"Domain":      input.Domain,
		}
	},

	ValidateWrite: func(db *gorm.DB, input models.NewProgram, id int) error {
		return validate.Validate(db, validate.NewUniqueRule("programs", "name", input.Name, nil).Except(id).Say("duplicate program name"))

	},
	ValidateDelete: func(db *gorm.DB, existing models.Program) error {
		return nil
	},
}
