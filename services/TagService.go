package services

import (
	"context"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type tagService struct {
	db           *gorm.DB
	aliasService *aliasService
}

func (s *tagService) Validate(ctx context.Context, input *models.TagInput) error {
	panic("ss")
}

func (s *tagService) Create(ctx context.Context, input *models.TagInput) (*models.Tag, error) {
	if err := s.Validate(ctx, input); err != nil {
		return nil, err
	}
	var tag models.Tag
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		tag = models.Tag{
			Name:        input.Name,
			Description: input.Description,
			Priority:    input.Priority,
		}
		if err := tx.Create(&tag).Error; err != nil {
			return err
		}
		// Create alias (will be auto-generated if not provided)
		if err := s.aliasService.CreateAlias(tx, "tags", tag.Id, input.Alias); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &tag, nil
}

func (s *tagService) Get(ctx context.Context, alias string) (*models.Tag, error) {
	return first[models.Tag](ctx, s.db, s.aliasService, alias)
}

// func (s *tagService) Update(id *int, alias *string, input *models.TagInput) (*models.Tag, error) {
// 	if id != nil {
// 		return s.GeneralCrud.Update(s.db, input, id)
// 	}
// 	if alias != nil {
// 		return s.GeneralCrud.UpdateByAlias(s.db, input, *alias)
// 	}
// 	return nil, gorm.ErrRecordNotFound
// }

func (s *tagService) Delete(ctx context.Context, alias string) (*models.Tag, error) {
	tag, err := s.Get(ctx, alias)
	if err != nil {
		return nil, err
	}
	if err := s.db.WithContext(ctx).Delete(&tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}
