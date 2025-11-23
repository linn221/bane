package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type taskService struct {
	db           *gorm.DB
	aliasService *aliasService
}

func (s *taskService) Create(ctx context.Context, input *models.TaskInput) (*models.Task, error) {
	today := utils.Today()
	task := models.Task{
		Title:       input.Title,
		Description: input.Description,
		Status:      models.TaskStatusInProgress,
		Priority:    input.Priority,
		Created:     models.MyDate{Time: today},
	}
	if input.Deadline != nil {
		if task.Deadline.Before(today) {
			return nil, errors.New("deadline is in the past")
		}
		task.Deadline = *input.Deadline
	}
	if input.RemindDate != nil {
		if input.RemindDate.Before(today) {
			return nil, fmt.Errorf("remind date is in the past: %v", input.RemindDate.Time)
		}
		task.RemindDate = *input.RemindDate
	}
	if input.ProjectAlias != "" {
		projectId, err := s.aliasService.GetReferenceId(ctx, input.ProjectAlias)
		if err != nil {
			return nil, err
		}
		task.ProjectId = projectId
	}

	tx := s.db.WithContext(ctx).Begin()
	defer tx.Rollback()

	if err := tx.Create(&task).Error; err != nil {
		return nil, err
	}
	if err := s.aliasService.CreateAlias(tx, "tasks", task.Id, input.Alias); err != nil {
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *taskService) List(ctx context.Context, filter *models.TaskFilter) ([]*models.Task, error) {
	var results []*models.Task
	dbctx := s.db.WithContext(ctx).Model(&models.Task{})
	if filter != nil {
		today := utils.Today()
		if filter.Today {
			dbctx.Where("remind_date = ?", today)
		}
		if filter.Search != "" {
			dbctx.Where("title LIKE ? OR description LIKE ?",
				"%"+filter.Search+"%",
				"%"+filter.Search+"%",
			)
		}
		if filter.Status != nil {
			dbctx.Where("status = ?", filter.Status)
		}
	}
	err := dbctx.Find(&results).Error
	return results, err
}

func (s *taskService) ChangeStatus(ctx context.Context, alias *string, id *int, status models.TaskStatus) (*models.Task, error) {

	task, err := GetRecordByAliasOrId[models.Task](s.db.WithContext(ctx), "tasks", alias, id)
	if err != nil {
		return nil, err
	}
	updates := map[string]any{
		"status": status,
	}
	if status == models.TaskStatusFinished {
		updates["finished_date"] = models.MyDate{Time: utils.Today()}
	}
	if status == models.TaskStatusCancelled {
		updates["cancelled_date"] = models.MyDate{Time: utils.Today()}
	}
	err = s.db.WithContext(ctx).Model(&task).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return GetRecordByAliasOrId[models.Task](s.db.WithContext(ctx), "tasks", alias, id)
}
