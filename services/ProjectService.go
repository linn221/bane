package services

import (
	"context"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type projectService struct {
	db           *gorm.DB
	aliasService *aliasService
}

func (s *projectService) Create(ctx context.Context, input *models.ProjectInput) (*models.Project, error) {
	project := models.Project{
		Name:        input.Name,
		Description: input.Description,
	}

	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		err := tx.Create(&project).Error
		if err != nil {
			return err
		}
		// Create alias (will be auto-generated if not provided)
		if err := s.aliasService.CreateAlias(tx, "projects", project.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (s *projectService) Get(ctx context.Context, id *int, alias *string) (*models.Project, error) {
	if id != nil {
		var result models.Project
		err := s.db.WithContext(ctx).First(&result, *id).Error
		return &result, err
	}
	if alias != nil {
		projectId, err := s.aliasService.GetReferenceId(ctx, *alias)
		if err != nil {
			return nil, err
		}
		var result models.Project
		err = s.db.WithContext(ctx).First(&result, projectId).Error
		return &result, err
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
