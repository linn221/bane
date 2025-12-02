package loaders

import (
	"context"

	"github.com/linn221/bane/models"
)

func GetProject(ctx context.Context, id int) (*models.Project, error) {
	loaders := For(ctx)
	return loaders.projectLoader.Load(ctx, id)()
}

func GetTasksByProjectId(ctx context.Context, projectId int) ([]*models.Task, error) {
	loaders := For(ctx)
	return loaders.tasksByProjectIdLoader.Load(ctx, projectId)()
}
