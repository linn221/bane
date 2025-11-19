package loaders

import (
	"context"
	"time"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/linn221/bane/app"
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
	ProgramLoader          *dataloader.Loader[int, *models.Program]
	NotesLoaderForProgram  *dataloader.Loader[int, []*models.Note]
	RIdLoader              *dataloader.Loader[app.Reference, int]
	TagAliasLoader         *dataloader.Loader[int, string]
	ProgramAliasLoader     *dataloader.Loader[int, string]
	WordAliasLoader        *dataloader.Loader[int, string]
	WordListAliasLoader    *dataloader.Loader[int, string]
	EndpointAliasLoader    *dataloader.Loader[int, string]
	MemorySheetAliasLoader *dataloader.Loader[int, string]
	MySheetAliasLoader     *dataloader.Loader[int, string]
	ProjectAliasLoader     *dataloader.Loader[int, string]
	TodoAliasLoader        *dataloader.Loader[int, string]
}

// NewLoaders instantiates data loaders for the middleware
func NewLoaders(conn *gorm.DB, deducer *app.Deducer) *Loaders {
	// Create Program and Notes readers for efficient loading
	programReader := &ProgramReader{db: conn}
	notesReader := &NotesReader{db: conn, referenceType: "programs"}
	ridReader := &RIdReader{deducer: deducer}

	// Create Alias readers for each reference type
	tagAliasReader := &AliasReader{db: conn, referenceType: "tags"}
	programAliasReader := &AliasReader{db: conn, referenceType: "programs"}
	wordAliasReader := &AliasReader{db: conn, referenceType: "words"}
	wordListAliasReader := &AliasReader{db: conn, referenceType: "wordlists"}
	endpointAliasReader := &AliasReader{db: conn, referenceType: "endpoints"}
	memorySheetAliasReader := &AliasReader{db: conn, referenceType: "memory_sheets"}
	mySheetAliasReader := &AliasReader{db: conn, referenceType: "my_sheets"}
	projectAliasReader := &AliasReader{db: conn, referenceType: "projects"}
	todoAliasReader := &AliasReader{db: conn, referenceType: "todos"}

	return &Loaders{
		ProgramLoader:          dataloader.NewBatchedLoader(programReader.GetPrograms, dataloader.WithWait[int, *models.Program](time.Millisecond)),
		NotesLoaderForProgram:  dataloader.NewBatchedLoader(notesReader.GetNotes, dataloader.WithWait[int, []*models.Note](time.Millisecond)),
		RIdLoader:              dataloader.NewBatchedLoader(ridReader.GetRIds, dataloader.WithWait[app.Reference, int](time.Millisecond)),
		TagAliasLoader:         dataloader.NewBatchedLoader(tagAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		ProgramAliasLoader:     dataloader.NewBatchedLoader(programAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		WordAliasLoader:        dataloader.NewBatchedLoader(wordAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		WordListAliasLoader:    dataloader.NewBatchedLoader(wordListAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		EndpointAliasLoader:    dataloader.NewBatchedLoader(endpointAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		MemorySheetAliasLoader: dataloader.NewBatchedLoader(memorySheetAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		MySheetAliasLoader:     dataloader.NewBatchedLoader(mySheetAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		ProjectAliasLoader:     dataloader.NewBatchedLoader(projectAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
		TodoAliasLoader:        dataloader.NewBatchedLoader(todoAliasReader.GetAliases, dataloader.WithWait[int, string](time.Millisecond)),
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
