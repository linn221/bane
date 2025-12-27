package loaders

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

type ctxKey string

const (
	loadersKey = ctxKey("dataloaders")
)

// LoadersKey returns the context key for loaders
func LoadersKey() ctxKey {
	return loadersKey
}

// Loaders wrap your data loaders to inject via middleware
type Loaders struct {
	wordAliasLoader     *dataloader.Loader[int, string]
	wordListAliasLoader *dataloader.Loader[int, string]
	endpointAliasLoader *dataloader.Loader[int, string]
	projectAliasLoader  *dataloader.Loader[int, string]
	projectLoader       *dataloader.Loader[int, *models.Project]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *gorm.DB) *Loaders {

	// Create Alias readers for each reference type
	wordAliasReader := &AliasReader{db: conn, referenceType: "words"}
	wordListAliasReader := &AliasReader{db: conn, referenceType: "wordlists"}
	endpointAliasReader := &AliasReader{db: conn, referenceType: "endpoints"}
	projectAliasReader := &AliasReader{db: conn, referenceType: "projects"}
	projectReader := newGenericReader[*models.Project, int](conn,
		func(p *models.Project) int {
			return p.Id
		},
		func(i int) *models.Project {
			return &models.Project{Id: i}
		},
	)

	return &Loaders{
		wordAliasLoader:     wordAliasReader.Loader(),
		wordListAliasLoader: wordListAliasReader.Loader(),
		endpointAliasLoader: endpointAliasReader.Loader(),
		projectAliasLoader:  projectAliasReader.Loader(),
		projectLoader:       projectReader.Loader(),
	}
}

// For retrieves the loaders from the context
// This function is used by loader functions to access the dataloader instances
func For(ctx context.Context) *Loaders {
	return ctx.Value(loadersKey).(*Loaders)
}

// handleError creates array of result with the same error repeated for as many items requested
// This helper function is used when database queries fail to return consistent error results
func handleError[T any](itemsLength int, err error) []*dataloader.Result[T] {
	result := make([]*dataloader.Result[T], itemsLength)
	for i := 0; i < itemsLength; i++ {
		result[i] = &dataloader.Result[T]{Error: err}
	}
	return result
}
