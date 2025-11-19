package services

import (
	"context"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type todoService struct {
	GeneralCrud[models.TodoInput, models.Todo]
	db           *gorm.DB
	aliasService *aliasService
}

func newTodoService(db *gorm.DB, aliasService *aliasService) *todoService {
	return &todoService{
		GeneralCrud: GeneralCrud[models.TodoInput, models.Todo]{
			transform: func(input *models.TodoInput) models.Todo {
				result := models.Todo{
					Title:       input.Title,
					Description: input.Description,
					Priority:    input.Priority,
					Status:      models.ToDoStatusInProgress,
					Created:     utils.Today(),
				}
				if input.Status != nil {
					result.Status = *input.Status
				}
				if input.Deadline != nil {
					result.Deadline = input.Deadline.Time
				}
				return result
			},
			updates: func(existing models.Todo, input *models.TodoInput) map[string]any {
				updates := map[string]any{}
				if input.Title != "" {
					updates["Title"] = input.Title
				}
				if input.Description != "" {
					updates["Description"] = input.Description
				}
				if input.Priority != 0 {
					updates["Priority"] = input.Priority
				}
				if input.Status != nil {
					updates["Status"] = *input.Status
				}
				if input.Deadline != nil {
					updates["Deadline"] = input.Deadline.Time
				}
				return updates
			},
			validateWrite: func(db *gorm.DB, input *models.TodoInput, id int) error {
				return input.Validate(db, id)
			},
		},
		db:           db,
		aliasService: aliasService,
	}
}

func (ts *todoService) Create(input *models.TodoInput) (*models.Todo, error) {
	var result *models.Todo
	err := ts.db.Transaction(func(tx *gorm.DB) error {
		var err error
		// Handle project alias lookup
		if input.ProjectAlias != "" {
			// Note: GetId doesn't have context, but we'll keep it for now since Create doesn't have context
			// TODO: Add context to Create method if needed
			projectId, err := ts.aliasService.GetId(input.ProjectAlias)
			if err != nil {
				return err
			}
			// We need to set ProjectId after creation, so we'll do it in a custom way
			result, err = ts.GeneralCrud.Create(tx, input)
			if err != nil {
				return err
			}
			result.ProjectId = projectId
			if err := tx.Save(result).Error; err != nil {
				return err
			}
		} else {
			result, err = ts.GeneralCrud.Create(tx, input)
			if err != nil {
				return err
			}
		}
		// Create alias (will be auto-generated if not provided)
		if err := ts.aliasService.CreateAlias(tx, "todos", result.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (ts *todoService) Get(id *int, alias *string) (*models.Todo, error) {
	if id != nil {
		return ts.GeneralCrud.Get(ts.db, id)
	}
	if alias != nil {
		return ts.GeneralCrud.GetByAlias(context.Background(), ts.db, ts.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ts *todoService) Update(id *int, alias *string, input *models.TodoInput) (*models.Todo, error) {
	if id != nil {
		return ts.GeneralCrud.Update(ts.db, input, id)
	}
	if alias != nil {
		// Note: GetId doesn't have context, but we'll keep it for now since Update doesn't have context
		// TODO: Add context to Update method if needed
		todoId, err := ts.aliasService.GetId(*alias)
		if err != nil {
			return nil, err
		}
		var result *models.Todo
		err = ts.db.Transaction(func(tx *gorm.DB) error {
			result, err = ts.GeneralCrud.Update(tx, input, &todoId)
			if err != nil {
				return err
			}
			// Handle project alias update
			if input.ProjectAlias != "" {
				projectId, err := ts.aliasService.GetId(input.ProjectAlias)
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
				if err := ts.aliasService.CreateAlias(tx, "todos", todoId, input.Alias); err != nil {
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

func (ts *todoService) Delete(id *int, alias *string) (*models.Todo, error) {
	if id != nil {
		return ts.GeneralCrud.Delete(ts.db, id)
	}
	if alias != nil {
		return ts.GeneralCrud.DeleteByAlias(context.Background(), ts.db, ts.aliasService, *alias)
	}
	return nil, gorm.ErrRecordNotFound
}

func (ts *todoService) List(filter *models.TodoFilter) ([]*models.Todo, error) {
	dbctx := ts.db.Model(&models.Todo{})
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
			// Note: GetId doesn't have context, but List doesn't have context either
			// TODO: Add context to List method if needed
			projectId, err := ts.aliasService.GetId(filter.ProjectAlias)
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
