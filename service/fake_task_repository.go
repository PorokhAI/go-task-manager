package service

import (
	"context"

	"example.com/go-backend/model"
	"example.com/go-backend/repository"
)

type FakeTaskRepository struct {
	Tasks []model.Task

	GetAllCalls int
	CreateCalls int
	UpdateCalls int
	PatchCalls  int
	DeleteCalls int

	CreatedTask model.Task
	DeletedID   int
}

func NewFakeTaskRepository(initialTasks []model.Task) *FakeTaskRepository {
	tasks := make([]model.Task, len(initialTasks))
	copy(tasks, initialTasks)

	return &FakeTaskRepository{
		Tasks: tasks,
	}
}

func (r *FakeTaskRepository) GetAll(ctx context.Context) ([]model.Task, error) {
	r.GetAllCalls++

	tasks := make([]model.Task, len(r.Tasks))
	copy(tasks, r.Tasks)
	return tasks, nil
}

func (r *FakeTaskRepository) Create(ctx context.Context, task model.Task) (model.Task, error) {
	r.CreateCalls++

	r.CreatedTask = task
	r.Tasks = append(r.Tasks, r.CreatedTask)
	return r.CreatedTask, nil
}

func (r *FakeTaskRepository) Update(ctx context.Context, id int, updatedTask model.Task) (model.Task, error) {
	r.UpdateCalls++

	for i, t := range r.Tasks {
		if t.ID == id {
			updatedTask.ID = id
			r.Tasks[i] = updatedTask

			return updatedTask, nil
		}
	}

	return model.Task{}, ErrTaskNotFound
}

func (r *FakeTaskRepository) Patch(ctx context.Context, id int, patch model.TaskPatch) (model.Task, error) {
	r.PatchCalls++

	for i, t := range r.Tasks {
		if t.ID == id {
			if patch.Title != nil {
				t.Title = *patch.Title
			}
			if patch.IsDone != nil {
				t.IsDone = *patch.IsDone
			}

			r.Tasks[i] = t
			return t, nil
		}
	}

	return model.Task{}, ErrTaskNotFound
}

func (r *FakeTaskRepository) Delete(ctx context.Context, id int) error {
	r.DeleteCalls++

	for i, t := range r.Tasks {
		if t.ID == id {
			r.Tasks = append(r.Tasks[:i], r.Tasks[i+1:]...)
			r.DeletedID = id
			return nil
		}
	}

	return ErrTaskNotFound
}

var _ repository.TaskRepository = (*FakeTaskRepository)(nil)
