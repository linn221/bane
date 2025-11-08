package services

import (
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type endpointService struct {
	GeneralCrud[models.EndpointInput, models.Endpoint]
	db           *gorm.DB
	deducer      Deducer
	aliasService *aliasService
}

func newEndpointService(db *gorm.DB, deducer Deducer, aliasService *aliasService) *endpointService {
	return &endpointService{
		GeneralCrud: GeneralCrud[models.EndpointInput, models.Endpoint]{
			transform: func(input *models.EndpointInput) models.Endpoint {
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
		db:           db,
		deducer:      deducer,
		aliasService: aliasService,
	}
}

func (es *endpointService) Create(input *models.EndpointInput) (*models.Endpoint, error) {
	// Find the program by alias using AliasService
	programId, err := es.aliasService.GetId(input.ProgramAlias)
	if err != nil {
		return nil, err
	}

	// Set the program ID
	endpoint := es.transform(input)
	endpoint.ProgramId = programId

	// Create the endpoint directly
	err = es.db.Create(&endpoint).Error
	if err != nil {
		return nil, err
	}

	// Set alias (will be auto-generated if not provided)
	if err := es.aliasService.SetAlias(string(models.AliasReferenceTypeEndpoint), endpoint.Id, input.Alias); err != nil {
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
		endpointId, err := es.aliasService.GetId(*alias)
		if err != nil {
			return nil, err
		}
		var endpoint models.Endpoint
		err = es.db.Preload("Program").First(&endpoint, endpointId).Error
		return &endpoint, err
	}
	return nil, gorm.ErrRecordNotFound
}
