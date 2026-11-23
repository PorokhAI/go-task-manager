package handler

import (
	"context"

	"example.com/go-backend/model"
)

type TaskService interface {
	GetAll(ctx context.Context) ([]model.Task, error)
	Create(ctx context.Context, task model.Task) (model.Task, error)
	Update(ctx context.Context, id int, task model.Task) (model.Task, error)
	Patch(ctx context.Context, id int, patch model.TaskPatch) (model.Task, error)
	Delete(ctx context.Context, id int) error
}
