package loaders

import (
	"context"

	"github.com/linn221/bane/models"
)

// GetNotesForProgram returns Notes for a single Program ID efficiently using dataloader
// This function should be used in the Program.Notes resolver instead of direct database queries
func GetNotesForProgram(ctx context.Context, programId int) ([]*models.Note, error) {
	loaders := For(ctx)
	return loaders.NotesLoaderForProgram.Load(ctx, programId)()
}

// GetNotesForPrograms returns Notes for multiple Program IDs efficiently using dataloader
// This function should be used when you need to load Notes for multiple Programs at once
func GetNotesForPrograms(ctx context.Context, programIds []int) ([][]*models.Note, []error) {
	loaders := For(ctx)
	return loaders.NotesLoaderForProgram.LoadMany(ctx, programIds)()
}
