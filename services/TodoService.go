package services

import (
	"context"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type todoService struct {
	GeneralCrud[models.TodoInput, models.Todo]
	db           *gorm.DB
	aliasService *aliasService
}

func (s *todoService) Create(ctx context.Context, input *models.TodoInput) (*models.Todo, error) {
	var result *models.Todo
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		// Handle project alias lookup
		if input.ProjectAlias != "" {
			projectId, err := s.aliasService.GetReferenceId(ctx, input.ProjectAlias)
			if err != nil {
				return err
			}
			// We need to set ProjectId after creation, so we'll do it in a custom way
			result, err = s.GeneralCrud.Create(tx, input)
			if err != nil {
				return err
			}
			result.ProjectId = projectId
			if err := tx.Save(result).Error; err != nil {
				return err
			}
		} else {
			result, err = s.GeneralCrud.Create(tx, input)
			if err != nil {
				return err
			}
		}
		// Create alias (will be auto-generated if not provided)
		if err := s.aliasService.CreateAlias(tx, "todos", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *todoService) Get(ctx context.Context, id *int, alias *string) (*models.Todo, error) {
	if id != nil {
		var result models.Todo
		err := s.db.WithContext(ctx).First(&result, *id).Error
		return &result, err
	}
	if alias != nil {
		return s.GeneralCrud.GetByAlias(ctx, s.db.WithContext(ctx), s.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *todoService) Update(ctx context.Context, id *int, alias *string, input *models.TodoInput) (*models.Todo, error) {
	if id != nil {
		return s.GeneralCrud.Update(s.db.WithContext(ctx), input, id)
	}
	if alias != nil {
		todoId, err := s.aliasService.GetReferenceId(ctx, *alias)
		if err != nil {
			return nil, err
		}
		var result *models.Todo
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			result, err = s.GeneralCrud.Update(tx, input, &todoId)
			if err != nil {
				return err
			}
			// Handle project alias update
			if input.ProjectAlias != "" {
				projectId, err := s.aliasService.GetReferenceId(ctx, input.ProjectAlias)
				if err != nil {
					return err
				}
				result.ProjectId = projectId
				if err := tx.Save(result).Error; err != nil {
					return err
				}
			}
			// Create alias if provided
			if input.Alias != "" {
				if err := s.aliasService.CreateAlias(tx, "todos", todoId, input.Alias); err != nil {
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

func (s *todoService) List(ctx context.Context, filter *models.TodoFilter) ([]*models.Todo, error) {
	dbctx := s.db.WithContext(ctx).Model(&models.Todo{})
	if filter != nil {
		if filter.Title != "" {
			dbctx = dbctx.Where("title LIKE ?", "%"+filter.Title+"%")
		}
		if filter.Status != "" {
			dbctx = dbctx.Where("status = ?", filter.Status)
		}
		if filter.ProjectId != 0 {
			dbctx = dbctx.Where("project_id = ?", filter.ProjectId)
		}
		if filter.ProjectAlias != "" {
			projectId, err := s.aliasService.GetReferenceId(ctx, filter.ProjectAlias)
			if err == nil {
				dbctx = dbctx.Where("project_id = ?", projectId)
			}
		}
		if filter.ParentId != 0 {
			dbctx = dbctx.Where("parent_id = ?", filter.ParentId)
		}
		if filter.PriorityMin != 0 {
			dbctx = dbctx.Where("priority >= ?", filter.PriorityMin)
		}
		if filter.PriorityMax != 0 {
			dbctx = dbctx.Where("priority <= ?", filter.PriorityMax)
		}
		if filter.DeadlineFrom != nil && !filter.DeadlineFrom.Time.IsZero() {
			dbctx = dbctx.Where("deadline >= ?", filter.DeadlineFrom.Time)
		}
		if filter.DeadlineTo != nil && !filter.DeadlineTo.Time.IsZero() {
			dbctx = dbctx.Where("deadline <= ?", filter.DeadlineTo.Time)
		}
		if filter.Search != "" {
			dbctx = dbctx.Where("title LIKE ? OR description LIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
		}
	}
	var results []*models.Todo
	err := dbctx.Find(&results).Error
	return results, err
}
