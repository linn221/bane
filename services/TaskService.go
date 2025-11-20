package services

import (
	"context"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type taskService struct {
	GeneralCrud[models.TaskInput, models.Task]
	db           *gorm.DB
	aliasService *aliasService
}

func (s *taskService) Create(ctx context.Context, input *models.TaskInput) (*models.Task, error) {
	var result *models.Task
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

func (s *taskService) Get(ctx context.Context, id *int, alias *string) (*models.Task, error) {
	if id != nil {
		var result models.Task
		err := s.db.WithContext(ctx).First(&result, *id).Error
		return &result, err
	}
	if alias != nil {
		return s.GeneralCrud.GetByAlias(ctx, s.db.WithContext(ctx), s.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (s *taskService) Update(ctx context.Context, id *int, alias *string, input *models.TaskInput) (*models.Task, error) {
	if id != nil {
		return s.GeneralCrud.Update(s.db.WithContext(ctx), input, id)
	}
	if alias != nil {
		taskId, err := s.aliasService.GetReferenceId(ctx, *alias)
		if err != nil {
			return nil, err
		}
		var result *models.Task
		err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			result, err = s.GeneralCrud.Update(tx, input, &taskId)
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
				if err := s.aliasService.CreateAlias(tx, "todos", taskId, input.Alias); err != nil {
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

func (s *taskService) List(ctx context.Context) ([]*models.Task, error) {
	var results []*models.Task
	err := s.db.WithContext(ctx).Model(&models.Task{}).Find(&results).Error
	return results, err
}

func (s *taskService) Cancel(ctx context.Context, id *int, alias *string) (*models.Task, error) {
	var taskId int
	if id != nil {
		taskId = *id
	} else if alias != nil {
		var err error
		taskId, err = s.aliasService.GetReferenceId(ctx, *alias)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	var result models.Task
	err := s.db.WithContext(ctx).First(&result, taskId).Error
	if err != nil {
		return nil, err
	}

	result.Status = models.TaskStatusCancelled
	result.CancelledDate = models.MyDate{Time: utils.Today()}
	err = s.db.WithContext(ctx).Save(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *taskService) Finish(ctx context.Context, id *int, alias *string) (*models.Task, error) {
	var taskId int
	if id != nil {
		taskId = *id
	} else if alias != nil {
		var err error
		taskId, err = s.aliasService.GetReferenceId(ctx, *alias)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	var result models.Task
	err := s.db.WithContext(ctx).First(&result, taskId).Error
	if err != nil {
		return nil, err
	}

	result.Status = models.TaskStatusFinished
	result.FinishedDate = models.MyDate{Time: utils.Today()}
	err = s.db.WithContext(ctx).Save(&result).Error
	if err != nil {
		return nil, err
	}

	return &result, nil
}

