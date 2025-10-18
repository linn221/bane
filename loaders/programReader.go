package loaders

import (
	"context"

	"github.com/linn221/bane/models"
)

// GetProgram returns a single Program by ID efficiently using dataloader
// This function should be used in resolvers instead of direct database queries
func GetProgram(ctx context.Context, id int) (*models.Program, error) {
	loaders := For(ctx)
	return loaders.ProgramLoader.Load(ctx, id)()
}

// GetPrograms returns multiple Programs by IDs efficiently using dataloader
// This function should be used when you need to load multiple Programs at once
func GetPrograms(ctx context.Context, ids []int) ([]*models.Program, []error) {
	loaders := For(ctx)
	return loaders.ProgramLoader.LoadMany(ctx, ids)()
}
