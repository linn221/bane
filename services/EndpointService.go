package services

import (
	"context"

	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type endpointService struct {
	db           *gorm.DB
	aliasService *aliasService
}

func (s *endpointService) Create(ctx context.Context, input *models.EndpointInput) (*models.Endpoint, error) {
	// Transform input to endpoint
	endpoint := models.Endpoint{
		Name:        input.Name,
		Description: input.Description,
		HttpSchema:  input.HttpSchema,
		HttpMethod:  input.HttpMethod,
		HttpDomain:  input.HttpDomain,
		HttpPath:    input.HttpPath,
		HttpQueries: input.HttpQueries,
		HttpHeaders: input.HttpHeaders,
		HttpCookies: input.HttpCookies,
		HttpBody:    input.HttpBody,
	}

	// Create the endpoint directly
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Create(&endpoint).Error
		if err != nil {
			return err
		}
		err = s.aliasService.CreateAlias(tx, "endpoints", endpoint.Id, input.Alias)
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
		if filter.HttpSchema != "" {
			query = query.Where("http_schema = ?", filter.HttpSchema)
		}

		if filter.HttpMethod != "" {
			query = query.Where("http_method = ?", filter.HttpMethod)
		}

		if filter.HttpDomain != "" {
			query = query.Where("http_domain = ?", filter.HttpDomain)
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
