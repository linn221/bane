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

func newTagService(db *gorm.DB, aliasService *aliasService) *tagService {
	return &tagService{db: db, aliasService: aliasService}
}

func (ts *tagService) Validate(ctx context.Context, input *models.NewTag) error {
	panic("ss")
}

func (ts *tagService) Create(ctx context.Context, input *models.NewTag) (*models.Tag, error) {
	if err := ts.Validate(ctx, input); err != nil {
		return nil, err
	}
	tag := models.Tag{
		Name:        input.Name,
		Description: input.Description,
		Priority:    input.Priority,
	}
	if err := ts.db.WithContext(ctx).Create(&tag).Error; err != nil {
		return nil, err
	}
	// Set alias if provided
	if input.Alias != "" {
		if err := ts.aliasService.SetAlias(string(models.AliasReferenceTypeTag), tag.Id, input.Alias); err != nil {
			return nil, err
		}
	}
	return &tag, nil
}

func (ts *tagService) Get(ctx context.Context, alias string) (*models.Tag, error) {
	return first[models.Tag](ts.db.WithContext(ctx), ts.aliasService, alias)
}

// func (ts *tagService) Update(id *int, alias *string, input *models.NewTag) (*models.Tag, error) {
// 	if id != nil {
// 		return ts.GeneralCrud.Update(ts.db, input, id)
// 	}
// 	if alias != nil {
// 		return ts.GeneralCrud.UpdateByAlias(ts.db, input, *alias)
// 	}
// 	return nil, gorm.ErrRecordNotFound
// }

func (ts *tagService) Delete(ctx context.Context, alias string) (*models.Tag, error) {
	tag, err := ts.Get(ctx, alias)
	if err != nil {
		return nil, err
	}
	if err := ts.db.WithContext(ctx).Delete(&tag).Error; err != nil {
		return nil, err
	}

	return tag, nil
}
