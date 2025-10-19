package services

import (
	"github.com/linn221/bane/app"
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type endpointService struct {
	GeneralCrud[models.NewEndpoint, models.Endpoint]
}

var EndpointService = endpointService{
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
}

func (es *endpointService) CreateEndpoint(app *app.App, db *gorm.DB, input *models.NewEndpoint) (*models.Endpoint, error) {
	// Find the program by alias
	var program models.Program
	err := db.Where("alias = ?", input.ProgramAlias).First(&program).Error
	if err != nil {
		return nil, err
	}

	// Set the program ID
	endpoint := es.transform(input)
	endpoint.ProgramId = program.Id

	// Create the endpoint directly
	err = db.Create(&endpoint).Error
	if err != nil {
		return nil, err
	}

	return &endpoint, nil
}

func (es *endpointService) ListEndpoints(app *app.App, db *gorm.DB, filter *models.EndpointFilter) ([]*models.Endpoint, error) {
	query := db.Model(&models.Endpoint{})

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
			query = query.Where("name ILIKE ? OR description ILIKE ? OR http_domain ILIKE ?",
				"%"+filter.Search+"%", "%"+filter.Search+"%", "%"+filter.Search+"%")
		}
	}

	var results []*models.Endpoint
	err := query.Find(&results).Error
	return results, err
}

func (es *endpointService) GetEndpointByID(app *app.App, db *gorm.DB, id int) (*models.Endpoint, error) {
	var endpoint models.Endpoint
	err := db.Preload("Program").First(&endpoint, id).Error
	return &endpoint, err
}

func (es *endpointService) UpdateEndpoint(app *app.App, db *gorm.DB, id int, input *models.NewEndpoint) (*models.Endpoint, error) {
	// Find the program by alias if provided
	if input.ProgramAlias != "" {
		var program models.Program
		err := db.Where("alias = ?", input.ProgramAlias).First(&program).Error
		if err != nil {
			return nil, err
		}

		// Update the program ID
		updates := map[string]interface{}{
			"program_id": program.Id,
		}
		err = db.Model(&models.Endpoint{}).Where("id = ?", id).Updates(updates).Error
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

	err := db.Model(&models.Endpoint{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}

	return es.GetEndpointByID(app, db, id)
}

func (es *endpointService) DeleteEndpoint(app *app.App, db *gorm.DB, id int) (*models.Endpoint, error) {
	var endpoint models.Endpoint
	err := db.First(&endpoint, id).Error
	if err != nil {
		return nil, err
	}

	err = db.Delete(&endpoint).Error
	return &endpoint, err
}
