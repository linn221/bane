package loaders

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/linn221/bane/models"
	"gorm.io/gorm"
)

// ProgramReader handles efficient loading of Program entities using dataloader
// This reader implements the BatchFunc interface for dataloader to batch multiple
// Program requests into a single database query, reducing N+1 query problems
type ProgramReader struct {
	db *gorm.DB
}

// GetPrograms is the batch function that loads multiple Programs by their IDs
// It receives a slice of Program IDs and returns a slice of dataloader.Result
// This function is called by dataloader when multiple Program requests are batched
func (r *ProgramReader) GetPrograms(ctx context.Context, ids []int) []*dataloader.Result[*models.Program] {
	// Query the database for all Programs with the given IDs
	var results []models.Program
	err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&results).Error
	if err != nil {
		// If there's an error, return error results for all requested IDs
		return handleError[*models.Program](len(ids), err)
	}

	// Create a map for quick lookup of results by ID
	resultMap := make(map[int]*models.Program)
	for _, result := range results {
		// Create a copy to avoid pointer issues
		program := result
		resultMap[program.Id] = &program
	}

	// Create dataloader results in the same order as requested IDs
	loaderResults := make([]*dataloader.Result[*models.Program], 0, len(ids))
	for _, id := range ids {
		if program, exists := resultMap[id]; exists {
			loaderResults = append(loaderResults, &dataloader.Result[*models.Program]{Data: program})
		} else {
			// If Program not found, return nil data (not an error)
			loaderResults = append(loaderResults, &dataloader.Result[*models.Program]{Data: nil})
		}
	}

	return loaderResults
}

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
