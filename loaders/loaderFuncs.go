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

func GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	loaders := For(ctx)
	return loaders.taskByIdLoader.Load(ctx, id)()
}

func GetChildrenTasks(ctx context.Context, id int) ([]*models.Task, error) {
	loaders := For(ctx)
	return loaders.childrenTasksLoader.Load(ctx, id)()
}

func GetMySheetLabelById(ctx context.Context, id int) (*models.MySheetLabel, error) {
	loaders := For(ctx)
	return loaders.mySheetLabelByIdLoader.Load(ctx, id)()
}
