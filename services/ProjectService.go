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

func (s *projectService) Create(ctx context.Context, input *models.ProjectInput) (*models.Project, error) {
	var result *models.Project
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.GeneralCrud.Create(tx, input)
		if err != nil {
			return err
		}
		// Create alias (will be auto-generated if not provided)
		if err := s.aliasService.CreateAlias(tx, "projects", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *projectService) Get(ctx context.Context, id *int, alias *string) (*models.Project, error) {
	if id != nil {
		var result models.Project
		err := s.db.WithContext(ctx).First(&result, *id).Error
		return &result, err
	}
	if alias != nil {
		return s.GeneralCrud.GetByAlias(ctx, s.db.WithContext(ctx), s.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *projectService) Update(ctx context.Context, id *int, alias *string, input *models.ProjectInput) (*models.Project, error) {
	if id != nil {
		return s.GeneralCrud.Update(s.db.WithContext(ctx), input, id)
	}
	if alias != nil {
		projectId, err := s.aliasService.GetReferenceId(ctx, *alias)
		if err != nil {
			return nil, err
		}
		var result *models.Project
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			result, err = s.GeneralCrud.Update(tx, input, &projectId)
			if err != nil {
				return err
			}
			// Create alias if provided
			if input.Alias != "" {
				if err := s.aliasService.CreateAlias(tx, "projects", projectId, input.Alias); err != nil {
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

func (s *projectService) List(ctx context.Context, filter *models.ProjectFilter) ([]*models.Project, error) {
	dbctx := s.db.WithContext(ctx).Model(&models.Project{})
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
