package repository

import (
	"context"
	"sync"

	"example.com/go-backend/model"
)

type InMemoryTaskRepository struct {
	mu     sync.Mutex
	tasks  []model.Task
	nextID int
}

func NewInMemoryTaskRepository(initialTasks []model.Task) *InMemoryTaskRepository {
	tasks := make([]model.Task, len(initialTasks))
	copy(tasks, initialTasks)

	nextID := 1
	for i := range tasks {
		if tasks[i].ID >= nextID {
			nextID = tasks[i].ID + 1
		}
	}

	return &InMemoryTaskRepository{
		tasks:  tasks,
		nextID: nextID,
	}
}

func (r *InMemoryTaskRepository) GetAll(_ context.Context) ([]model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasks := make([]model.Task, len(r.tasks))
	copy(tasks, r.tasks)
	return tasks, nil
}

func (r *InMemoryTaskRepository) Create(_ context.Context, task model.Task) (model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	r.nextID++
	r.tasks = append(r.tasks, task)

	return task, nil
}

func (r *InMemoryTaskRepository) Update(_ context.Context, id int, updatedTask model.Task) (model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			updatedTask.ID = id
			r.tasks[i] = updatedTask
			return updatedTask, nil
		}
	}

	return model.Task{}, ErrNotFound
}

func (r *InMemoryTaskRepository) Patch(_ context.Context, id int, patch model.TaskPatch) (model.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			if patch.Title != nil {
				task.Title = *patch.Title
			}
			if patch.IsDone != nil {
				task.IsDone = *patch.IsDone
			}

			r.tasks[i] = task
			return task, nil
		}
	}

	return model.Task{}, ErrNotFound
}

func (r *InMemoryTaskRepository) Delete(_ context.Context, id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, task := range r.tasks {
		if task.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}

var _ TaskRepository = (*InMemoryTaskRepository)(nil)
