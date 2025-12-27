package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/linn221/bane/models"
	"github.com/linn221/bane/mystructs"
	"github.com/linn221/bane/utils"
	"gorm.io/gorm"
)

type endpointService struct {
	db           *gorm.DB
	aliasService *aliasService
}

func (s *endpointService) Create(ctx context.Context, input *models.EndpointInput) (*models.Endpoint, error) {
	// Parse the URL to extract all components
	parsedUrl, err := utils.ParseHttpUrl(input.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}

	// Set default for body if not provided
	body := mystructs.VarString{OriginalString: ""}
	if input.Body != nil {
		body = *input.Body
	}

	// Serialize input to JSON for storage
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input to JSON: %w", err)
	}

	// Transform input to endpoint
	endpoint := models.Endpoint{
		Name:        input.Name,
		Description: input.Description,
		ProjectId:   input.ProjectId,
		Https:       parsedUrl.Https,
		Method:      input.Method,
		Domain:      parsedUrl.HttpDomain,
		Path:        parsedUrl.HttpPath,
		Queries:     parsedUrl.HttpQueries,
		Headers:     input.Headers,
		Body:        body,
		Input:       string(inputJSON),
	}

	// Create the endpoint directly
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&endpoint).Error
		if err != nil {
			return err
		}
		// Alias is auto-generated if not provided (handled by CreateAlias)
		err = s.aliasService.CreateAlias(tx, "endpoints", endpoint.Id, "")
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func (s *endpointService) List(ctx context.Context, filter *models.EndpointFilter) ([]*models.Endpoint, error) {
	query := s.db.WithContext(ctx).Model(&models.Endpoint{})

	if filter != nil {
		if filter.Https != nil {
			query = query.Where("http_schema = ?", *filter.Https)
		}

		if filter.Method != "" {
			query = query.Where("http_method = ?", filter.Method)
		}

		if filter.Domain != "" {
			query = query.Where("http_domain = ?", filter.Domain)
		}

		if filter.Search != "" {
			query = query.Where("name LIKE ? OR description LIKE ? OR http_domain LIKE ?",
				"%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
		}
	}

	var results []*models.Endpoint
	err := query.Find(&results).Error
	return results, err
}

func (s *endpointService) Get(ctx context.Context, id *int, alias *string) (*models.Endpoint, error) {
	if id != nil {
		var endpoint models.Endpoint
		err := s.db.WithContext(ctx).First(&endpoint, *id).Error
		return &endpoint, err
	}
	if alias != nil {
		endpointId, err := s.aliasService.GetReferenceId(ctx, *alias)
		if err != nil {
			return nil, err
		}
		var endpoint models.Endpoint
		err = s.db.WithContext(ctx).First(&endpoint, endpointId).Error
		return &endpoint, err
	}
	return nil, gorm.ErrRecordNotFound
}
