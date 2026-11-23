package service

import (
	"context"
	"errors"
	"strings"

	"example.com/go-backend/model"
	"example.com/go-backend/repository"
)

var (
	ErrEmptyTitle   = errors.New("название задачи не может быть пустым")
	ErrTaskNotFound = errors.New("задача не найдена")
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (s *TaskService) GetAll(ctx context.Context) ([]model.Task, error) {
	return s.repo.GetAll(ctx)
}

func (s *TaskService) Create(ctx context.Context, task model.Task) (model.Task, error) {
	title := strings.TrimSpace(task.Title)
	if title == "" {
		return model.Task{}, ErrEmptyTitle
	}

	task.Title = title
	return s.repo.Create(ctx, task)
}

func (s *TaskService) Update(ctx context.Context, id int, task model.Task) (model.Task, error) {
	title := strings.TrimSpace(task.Title)
	if title == "" {
		return model.Task{}, ErrEmptyTitle
	}

	task.Title = title

	updatedTask, err := s.repo.Update(ctx, id, task)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	return updatedTask, nil
}

func (s *TaskService) Patch(ctx context.Context, id int, patch model.TaskPatch) (model.Task, error) {
	if patch.Title != nil {
		title := strings.TrimSpace(*patch.Title)
		if title == "" {
			return model.Task{}, ErrEmptyTitle
		}

		patch.Title = &title
	}

	updatedTask, err := s.repo.Patch(ctx, id, patch)
	if errors.Is(err, repository.ErrNotFound) {
		return model.Task{}, ErrTaskNotFound
	}
	if err != nil {
		return model.Task{}, err
	}

	return updatedTask, nil
}

func (s *TaskService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrTaskNotFound
	}
	if err != nil {
		return err
	}

	return nil
}
