package services

import (
	"context"

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
	var result *models.Project
	err := ps.db.Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = ps.GeneralCrud.Create(tx, input)
		if err != nil {
			return err
		}
		// Create alias (will be auto-generated if not provided)
		if err := ps.aliasService.CreateAlias(tx, "projects", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ps *projectService) Get(id *int, alias *string) (*models.Project, error) {
	if id != nil {
		return ps.GeneralCrud.Get(ps.db, id)
	}
	if alias != nil {
		return ps.GeneralCrud.GetByAlias(context.Background(), ps.db, ps.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ps *projectService) Update(id *int, alias *string, input *models.ProjectInput) (*models.Project, error) {
	if id != nil {
		return ps.GeneralCrud.Update(ps.db, input, id)
	}
	if alias != nil {
		// Note: GetId doesn't have context, but we'll keep it for now since Update doesn't have context
		// TODO: Add context to Update method if needed
		projectId, err := ps.aliasService.GetId(*alias)
		if err != nil {
			return nil, err
		}
		var result *models.Project
		err = ps.db.Transaction(func(tx *gorm.DB) error {
			result, err = ps.GeneralCrud.Update(tx, input, &projectId)
			if err != nil {
				return err
			}
			// Create alias if provided
			if input.Alias != "" {
				if err := ps.aliasService.CreateAlias(tx, "projects", projectId, input.Alias); err != nil {
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

func (ps *projectService) Delete(id *int, alias *string) (*models.Project, error) {
	if id != nil {
		return ps.GeneralCrud.Delete(ps.db, id)
	}
	if alias != nil {
		return ps.GeneralCrud.DeleteByAlias(context.Background(), ps.db, ps.aliasService, *alias)
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
