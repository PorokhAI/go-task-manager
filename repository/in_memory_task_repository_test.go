package repository

import (
	"context"
	"testing"

	"example.com/go-backend/model"
)

func TestInMemoryTaskRepositoryNextID(t *testing.T) {
	repo := NewInMemoryTaskRepository([]model.Task{
		{ID: 1, Title: "Задача 1"},
		{ID: 3, Title: "Задача 3"},
	})

	created, err := repo.Create(context.Background(), model.Task{
		Title: "Новая задача",
	})
	if err != nil {
		t.Fatalf("Ошибка при создании задачи: %v", err)
	}

	if created.ID != 4 {
		t.Errorf("ожидался ID 4, получен %d", created.ID)
	}
}

func TestInMemoryTaskRepositoryNextIDAfterDelete(t *testing.T) {
	repo := NewInMemoryTaskRepository([]model.Task{
		{ID: 1, Title: "Задача 1"},
		{ID: 2, Title: "Задача 2"},
	})

	err := repo.Delete(context.Background(), 2)
	if err != nil {
		t.Fatalf("Ошибка при удалении задачи: %v", err)
	}

	created, err := repo.Create(context.Background(), model.Task{
		Title: "Новая задача",
	})
	if err != nil {
		t.Fatalf("Ошибка при создании задачи: %v", err)
	}

	if created.ID != 3 {
		t.Errorf("ожидался ID 3 после удаления, получен %d", created.ID)
	}
}
