package services

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type projectService struct {
	GeneralCrud[models.ProjectInput, models.Project]
	db           *gorm.DB
	aliasService *aliasService
}

func newProjectService(db *gorm.DB, aliasService *aliasService) *projectService {
	return &projectService{
		GeneralCrud: GeneralCrud[models.ProjectInput, models.Project]{
			transform: func(input *models.ProjectInput) models.Project {
				return models.Project{
					Name:        input.Name,
					Description: input.Description,
				}
			},
			updates: func(existing models.Project, input *models.ProjectInput) map[string]any {
				updates := map[string]any{}
				if input.Name != "" {
					updates["Name"] = input.Name
				}
				if input.Description != "" {
					updates["Description"] = input.Description
				}
				return updates
			},
			validateWrite: func(db *gorm.DB, input *models.ProjectInput, id int) error {
				return input.Validate(db, id)
			},
		},
		db:           db,
		aliasService: aliasService,
	}
}

func (ps *projectService) Create(input *models.ProjectInput) (*models.Project, error) {
	result, err := ps.GeneralCrud.Create(ps.db, input)
	if err != nil {
		return nil, err
	}
	// Set alias (will be auto-generated if not provided)
	if err := ps.aliasService.SetAlias(string(models.AliasReferenceTypeProject), result.Id, input.Alias); err != nil {
		return nil, err
	}
	return result, nil
}

func (ps *projectService) Get(id *int, alias *string) (*models.Project, error) {
	if id != nil {
		return ps.GeneralCrud.Get(ps.db, id)
	}
	if alias != nil {
		return ps.GeneralCrud.GetByAlias(ps.db, ps.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ps *projectService) Update(id *int, alias *string, input *models.ProjectInput) (*models.Project, error) {
	if id != nil {
		return ps.GeneralCrud.Update(ps.db, input, id)
	}
	if alias != nil {
		projectId, err := ps.aliasService.GetId(*alias)
		if err != nil {
			return nil, err
		}
		result, err := ps.GeneralCrud.Update(ps.db, input, &projectId)
		if err != nil {
			return nil, err
		}
		// Set alias if provided
		if input.Alias != "" {
			if err := ps.aliasService.SetAlias(string(models.AliasReferenceTypeProject), projectId, input.Alias); err != nil {
				return nil, err
			}
		}
		return result, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (ps *projectService) Delete(id *int, alias *string) (*models.Project, error) {
	if id != nil {
		return ps.GeneralCrud.Delete(ps.db, id)
	}
	if alias != nil {
		return ps.GeneralCrud.DeleteByAlias(ps.db, ps.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ps *projectService) List(filter *models.ProjectFilter) ([]*models.Project, error) {
	dbctx := ps.db.Model(&models.Project{})
	if filter != nil {
		if filter.Name != "" {
			dbctx = dbctx.Where("name LIKE ?", "%"+filter.Name+"%")
		}
		if filter.Search != "" {
			dbctx = dbctx.Where("name LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
		}
	}
	var results []*models.Project
	err := dbctx.Find(&results).Error
	return results, err
}
