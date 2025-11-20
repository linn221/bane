package loaders

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

// AliasReader handles efficient loading of aliases for entities using dataloader
// This reader implements the BatchFunc interface for dataloader to batch multiple
// alias requests into a single database query, reducing N+1 query problems
type AliasReader struct {
	db            *gorm.DB
	referenceType string
}

// GetAliases is the batch function that loads multiple aliases by their reference IDs
// It receives a slice of reference IDs and returns a slice of dataloader.Result containing strings
// This function is called by dataloader when multiple alias requests are batched
func (r *AliasReader) GetAliases(ctx context.Context, referenceIds []int) []*dataloader.Result[string] {
	// Query the database for all aliases with the given reference IDs and type
	var aliases []models.Alias
	err := r.db.WithContext(ctx).
		Where("reference_id IN ? AND reference_type = ?", referenceIds, r.referenceType).
		Find(&aliases).Error

	if err != nil {
		// If there's an error, return error results for all requested IDs
		return handleError[string](len(referenceIds), err)
	}

	// Create a map for quick lookup of aliases by reference ID
	aliasMap := make(map[int]string)
	for _, alias := range aliases {
		aliasMap[alias.ReferenceId] = alias.Name
	}

	// Create dataloader results in the same order as requested IDs
	loaderResults := make([]*dataloader.Result[string], 0, len(referenceIds))
	for _, id := range referenceIds {
		if alias, exists := aliasMap[id]; exists {
			loaderResults = append(loaderResults, &dataloader.Result[string]{Data: alias})
		} else {
			// If alias not found, return empty string (not an error)
			loaderResults = append(loaderResults, &dataloader.Result[string]{Data: ""})
		}
	}

	return loaderResults
}

// GetTagAlias returns a single alias for a Tag by ID efficiently using dataloader
func GetTagAlias(ctx context.Context, id int) (string, error) {
	loaders := For(ctx)
	return loaders.TagAliasLoader.Load(ctx, id)()
}

// GetProgramAlias returns a single alias for a Program by ID efficiently using dataloader
func GetProgramAlias(ctx context.Context, id int) (string, error) {
	loaders := For(ctx)
	return loaders.ProgramAliasLoader.Load(ctx, id)()
}

// GetWordAlias returns a single alias for a Word by ID efficiently using dataloader
func GetWordAlias(ctx context.Context, id int) (string, error) {
	loaders := For(ctx)
	return loaders.WordAliasLoader.Load(ctx, id)()
}

// GetWordListAlias returns a single alias for a WordList by ID efficiently using dataloader
func GetWordListAlias(ctx context.Context, id int) (string, error) {
	loaders := For(ctx)
	return loaders.WordListAliasLoader.Load(ctx, id)()
}

// GetEndpointAlias returns a single alias for an Endpoint by ID efficiently using dataloader
func GetEndpointAlias(ctx context.Context, id int) (string, error) {
	loaders := For(ctx)
	return loaders.EndpointAliasLoader.Load(ctx, id)()
}

// GetMemorySheetAlias returns a single alias for a MemorySheet by ID efficiently using dataloader
func GetMemorySheetAlias(ctx context.Context, id int) (string, error) {
	loaders := For(ctx)
	return loaders.MemorySheetAliasLoader.Load(ctx, id)()
}

// GetMySheetAlias returns a single alias for a MySheet by ID efficiently using dataloader
func GetMySheetAlias(ctx context.Context, id int) (string, error) {
	loaders := For(ctx)
	return loaders.MySheetAliasLoader.Load(ctx, id)()
}

// GetProjectAlias returns a single alias for a Project by ID efficiently using dataloader
func GetProjectAlias(ctx context.Context, id int) (string, error) {
	loaders := For(ctx)
	return loaders.ProjectAliasLoader.Load(ctx, id)()
}

// GetTaskAlias returns a single alias for a Task by ID efficiently using dataloader
func GetTaskAlias(ctx context.Context, id int) (string, error) {
	loaders := For(ctx)
	return loaders.TaskAliasLoader.Load(ctx, id)()
}
