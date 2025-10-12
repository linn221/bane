package services

import (
	"context"
	"fmt"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type TagService struct {
	DB *gorm.DB
}

func (s *TagService) createWithTx(tx *gorm.DB, input *models.NewTag) (uint, error) {
	result := models.Tag{
		Name:        input.Name,
		Description: input.Description,
		Alias:       input.Alias,
		Priority:    input.Priority,
	}

	if err := tx.Create(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to create tag: %v", err)
	}
	return uint(result.Id), nil
}

// Create creates a new tag and returns its Id
func (s *TagService) Create(ctx context.Context, tag *models.NewTag) (*models.Tag, error) {
	id, err := s.createWithTx(s.DB.WithContext(ctx), tag)
	if err != nil {
		return nil, err
	}
	return first[models.Tag](s.DB.WithContext(ctx), id)
}

// List retrieves all tags
func (s *TagService) List(ctx context.Context, search *string) ([]*models.Tag, error) {
	var tags []*models.Tag
	dbctx := s.DB.WithContext(ctx).Model(&models.Tag{})
	if search != nil {
		dbctx.Where("name LIKE ?", utils.SurroundPercentages(*search))
	}

	if err := dbctx.Order("priority DESC, name ASC").Find(&tags).Error; err != nil {
		return nil, fmt.Errorf("failed to list tags: %v", err)
	}
	return tags, nil
}

// Get retrieves a tag by ID
func (s *TagService) Get(ctx context.Context, id uint) (*models.Tag, error) {
	var input models.Tag
	if err := s.DB.WithContext(ctx).First(&input, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("tag not found")
		}
		return nil, fmt.Errorf("failed to get tag: %v", err)
	}
	return &input, nil
}

// // Update updates an existing tag
// func (s *TagService) Update(ctx context.Context, id int, tag *models.Tag) (int, error) {
// 	// Check if tag exists
// 	var existingTag models.Tag
// 	if err := s.DB.WithContext(ctx).First(&existingTag, id).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return 0, fmt.Errorf("tag not found")
// 		}
// 		return 0, fmt.Errorf("failed to find tag: %v", err)
// 	}

// 	// Update fields
// 	existingTag.Name = tag.Name
// 	existingTag.Priority = tag.Priority

// 	if err := s.DB.WithContext(ctx).Save(&existingTag).Error; err != nil {
// 		return 0, fmt.Errorf("failed to update tag: %v", err)
// 	}

// 	return existingTag.Id, nil
// }

// // ConnectTagToReference connects a tag to a reference (polymorphic relationship)
// func (s *TagService) ConnectTagToReference(ctx context.Context, tagId int, referenceType string, referenceId int) error {
// 	// Check if tag exists
// 	var tag models.Tag
// 	if err := s.DB.WithContext(ctx).First(&tag, tagId).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return fmt.Errorf("tag not found")
// 		}
// 		return fmt.Errorf("failed to find tag: %v", err)
// 	}

// 	// Check if connection already exists
// 	var existingRef models.Taggable
// 	err := s.DB.WithContext(ctx).Where("tag_id = ? AND taggable_type = ? AND taggable_id = ?",
// 		tagId, referenceType, referenceId).First(&existingRef).Error
// 	if err == nil {
// 		return fmt.Errorf("tag is already connected to this reference")
// 	}

// 	// Create the connection
// 	tagRef := models.Taggable{
// 		TagID:        tagId,
// 		TaggableType: referenceType,
// 		TaggableID:   referenceId,
// 	}

// 	if err := s.DB.WithContext(ctx).Create(&tagRef).Error; err != nil {
// 		return fmt.Errorf("failed to connect tag to reference: %v", err)
// 	}

// 	return nil
// }

// // DisconnectTagFromReference disconnects a tag from a reference
// func (s *TagService) DisconnectTagFromReference(ctx context.Context, tagId int, referenceType string, referenceId int) error {
// 	result := s.DB.WithContext(ctx).Where("tag_id = ? AND taggable_type = ? AND taggable_id = ?",
// 		tagId, referenceType, referenceId).Delete(&models.Taggable{})

// 	if result.Error != nil {
// 		return fmt.Errorf("failed to disconnect tag from reference: %v", result.Error)
// 	}

// 	if result.RowsAffected == 0 {
// 		return fmt.Errorf("tag connection not found")
// 	}

// 	return nil
// }

// // GetTagsForReference gets all tags for a specific reference
// func (s *TagService) GetTagsForReference(ctx context.Context, referenceType string, referenceId int) ([]*models.Tag, error) {
// 	var tags []*models.Tag
// 	if err := s.DB.WithContext(ctx).
// 		Joins("JOIN taggables ON tags.id = taggables.tag_id").
// 		Where("taggables.taggable_type = ? AND taggables.taggable_id = ?", referenceType, referenceId).
// 		Order("tags.priority DESC, tags.name ASC").
// 		Find(&tags).Error; err != nil {
// 		return nil, fmt.Errorf("failed to get tags for reference: %v", err)
// 	}
// 	return tags, nil
// }

// // ConnectTagsToReference connects multiple tags to a reference
// func (s *TagService) ConnectTagsToReference(ctx context.Context, tagIds []int, referenceType string, referenceId int) error {
// 	for _, tagId := range tagIds {
// 		if err := s.ConnectTagToReference(ctx, tagId, referenceType, referenceId); err != nil {
// 			// If tag is already connected, continue with other tags
// 			if err.Error() == "tag is already connected to this reference" {
// 				continue
// 			}
// 			return err
// 		}
// 	}
// 	return nil
// }
