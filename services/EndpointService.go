package services

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type endpointService struct {
	GeneralCrud[models.NewEndpoint, models.Endpoint]
	db      *gorm.DB
	deducer Deducer
}

func newEndpointService(db *gorm.DB, deducer Deducer) *endpointService {
	return &endpointService{
		GeneralCrud: GeneralCrud[models.NewEndpoint, models.Endpoint]{
			transform: func(input *models.NewEndpoint) models.Endpoint {
				return models.Endpoint{
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
			},
		},
		db:      db,
		deducer: deducer,
	}
}

func (es *endpointService) Create(input *models.NewEndpoint) (*models.Endpoint, error) {
	// Find the program by alias
	var program models.Program
	err := es.db.Where("alias = ?", input.ProgramAlias).First(&program).Error
	if err != nil {
		return nil, err
	}

	// Set the program ID
	endpoint := es.transform(input)
	endpoint.ProgramId = program.Id

	// Create the endpoint directly
	err = es.db.Create(&endpoint).Error
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func (es *endpointService) List(filter *models.EndpointFilter) ([]*models.Endpoint, error) {
	query := es.db.Model(&models.Endpoint{})

	if filter != nil {
		if filter.ProgramAlias != "" {
			// Join with programs table to filter by program alias
			query = query.Joins("JOIN programs ON endpoints.program_id = programs.id").
				Where("programs.alias = ?", filter.ProgramAlias)
		}

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

func (es *endpointService) Get(id *int, alias *string) (*models.Endpoint, error) {
	if id != nil {
		var endpoint models.Endpoint
		err := es.db.Preload("Program").First(&endpoint, *id).Error
		return &endpoint, err
	}
	if alias != nil {
		var endpoint models.Endpoint
		err := es.db.Preload("Program").Where("alias = ?", *alias).First(&endpoint).Error
		return &endpoint, err
	}
	return nil, gorm.ErrRecordNotFound
}

func (es *endpointService) Update(id *int, alias *string, input *models.NewEndpoint) (*models.Endpoint, error) {
	var endpointId int
	if id != nil {
		endpointId = *id
	} else if alias != nil {
		var endpoint models.Endpoint
		err := es.db.Where("alias = ?", *alias).First(&endpoint).Error
		if err != nil {
			return nil, err
		}
		endpointId = endpoint.Id
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	// Find the program by alias if provided
	if input.ProgramAlias != "" {
		var program models.Program
		err := es.db.Where("alias = ?", input.ProgramAlias).First(&program).Error
		if err != nil {
			return nil, err
		}

		// Update the program ID
		updates := map[string]interface{}{
			"program_id": program.Id,
		}
		err = es.db.Model(&models.Endpoint{}).Where("id = ?", endpointId).Updates(updates).Error
		if err != nil {
			return nil, err
		}
	}

	// Update other fields
	updates := map[string]interface{}{
		"name":         input.Name,
		"description":  input.Description,
		"http_schema":  input.HttpSchema,
		"http_method":  input.HttpMethod,
		"http_domain":  input.HttpDomain,
		"http_path":    input.HttpPath,
		"http_queries": input.HttpQueries,
		"http_headers": input.HttpHeaders,
		"http_cookies": input.HttpCookies,
		"http_body":    input.HttpBody,
	}

	err := es.db.Model(&models.Endpoint{}).Where("id = ?", endpointId).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return es.Get(&endpointId, nil)
}

func (es *endpointService) Patch(id *int, alias *string, input *models.PatchEndpoint) (*models.Endpoint, error) {
	var endpointId int
	if id != nil {
		endpointId = *id
	} else if alias != nil {
		var endpoint models.Endpoint
		err := es.db.Where("alias = ?", *alias).First(&endpoint).Error
		if err != nil {
			return nil, err
		}
		endpointId = endpoint.Id
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	updates := make(map[string]any)

	// Handle program alias separately
	if input.ProgramAlias != nil && *input.ProgramAlias != "" {
		var program models.Program
		err := es.db.Where("alias = ?", *input.ProgramAlias).First(&program).Error
		if err != nil {
			return nil, err
		}
		updates["program_id"] = program.Id
	}

	// Add other fields if they are provided
	if input.Name != nil && *input.Name != "" {
		updates["name"] = *input.Name
	}
	if input.Alias != nil && *input.Alias != "" {
		updates["alias"] = *input.Alias
	}
	if input.Description != nil && *input.Description != "" {
		updates["description"] = *input.Description
	}
	if input.HttpSchema != nil {
		updates["http_schema"] = *input.HttpSchema
	}
	if input.HttpMethod != nil {
		updates["http_method"] = *input.HttpMethod
	}
	if input.HttpDomain != nil && *input.HttpDomain != "" {
		updates["http_domain"] = *input.HttpDomain
	}
	if input.HttpPort != nil {
		updates["http_port"] = *input.HttpPort
	}
	if input.HttpTimeout != nil {
		updates["http_timeout"] = *input.HttpTimeout
	}
	if input.HttpFollowRedirects != nil {
		updates["http_follow_redirects"] = *input.HttpFollowRedirects
	}
	// Check custom types using their IsZero methods
	if input.HttpPath != nil && !input.HttpPath.IsZero() {
		updates["http_path"] = *input.HttpPath
	}
	if input.HttpQueries != nil && !input.HttpQueries.IsZero() {
		updates["http_queries"] = *input.HttpQueries
	}
	if input.HttpHeaders != nil && !input.HttpHeaders.IsZero() {
		updates["http_headers"] = *input.HttpHeaders
	}
	if input.HttpCookies != nil && !input.HttpCookies.IsZero() {
		updates["http_cookies"] = *input.HttpCookies
	}
	if input.HttpBody != nil && !input.HttpBody.IsZero() {
		updates["http_body"] = *input.HttpBody
	}

	// Use the GeneralCrud Patch method
	return es.GeneralCrud.Patch(es.db, updates, &endpointId)
}

func (es *endpointService) Delete(id *int, alias *string) (*models.Endpoint, error) {
	var endpointId int
	if id != nil {
		endpointId = *id
	} else if alias != nil {
		var endpoint models.Endpoint
		err := es.db.Where("alias = ?", *alias).First(&endpoint).Error
		if err != nil {
			return nil, err
		}
		endpointId = endpoint.Id
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	var endpoint models.Endpoint
	err := es.db.First(&endpoint, endpointId).Error
	if err != nil {
		return nil, err
	}

	err = es.db.Delete(&endpoint).Error
	return &endpoint, err
}
