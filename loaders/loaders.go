package loaders

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/linn221/bane/app"
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
	RIdLoader              *dataloader.Loader[app.Reference, int]
	WordAliasLoader        *dataloader.Loader[int, string]
	WordListAliasLoader    *dataloader.Loader[int, string]
	EndpointAliasLoader    *dataloader.Loader[int, string]
	MySheetAliasLoader     *dataloader.Loader[int, string]
	ProjectAliasLoader     *dataloader.Loader[int, string]
	TaskAliasLoader        *dataloader.Loader[int, string]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *gorm.DB) *Loaders {
	ridReader := &RIdReader{}

	// Create Alias readers for each reference type
	wordAliasReader := &AliasReader{db: conn, referenceType: "words"}
	wordListAliasReader := &AliasReader{db: conn, referenceType: "wordlists"}
	endpointAliasReader := &AliasReader{db: conn, referenceType: "endpoints"}
	mySheetAliasReader := &AliasReader{db: conn, referenceType: "my_sheets"}
	projectAliasReader := &AliasReader{db: conn, referenceType: "projects"}
	taskAliasReader := &AliasReader{db: conn, referenceType: "tasks"}

	return &Loaders{
		RIdLoader:          dataloader.NewBatchedLoader(ridReader.GetRIds, dataloader.WithWait[app.Reference, int](time.Millisecond)),
		WordAliasLoader:    dataloader.NewBatchedLoader(wordAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		WordListAliasLoader: dataloader.NewBatchedLoader(wordListAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		EndpointAliasLoader: dataloader.NewBatchedLoader(endpointAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		MySheetAliasLoader: dataloader.NewBatchedLoader(mySheetAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		ProjectAliasLoader: dataloader.NewBatchedLoader(projectAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		TaskAliasLoader:    dataloader.NewBatchedLoader(taskAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
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
