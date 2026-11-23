package service_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"testing"

	"example.com/go-backend/model"
	"example.com/go-backend/repository"
	"example.com/go-backend/service"
)

func newInMemoryTaskService() *service.TaskService {
	repo := repository.NewInMemoryTaskRepository(nil)
	return service.NewTaskService(repo)
}

func TestTaskServiceCreate(t *testing.T) {
	taskService := newInMemoryTaskService()

	created, err := taskService.Create(context.Background(), model.Task{
		Title:  "Изучить Go",
		IsDone: false,
	})
	if err != nil {
		t.Fatalf("не ожидалась ошибка создания: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("ожидался ID 1, получен %d", created.ID)
	}
	if created.Title != "Изучить Go" {
		t.Errorf("ожидался заголовок %q, получен %q",
			"Изучить Go", created.Title)
	}
	if created.IsDone {
		t.Error("ожидалось, что задача будет не выполнена")
	}
}

func TestTaskServiceCreateTrimsTitle(t *testing.T) {
	taskService := newInMemoryTaskService()

	created, err := taskService.Create(context.Background(), model.Task{
		Title: "  Изучить Go  ",
	})
	if err != nil {
		t.Fatalf("не ожидалась ошибка создания: %v", err)
	}

	if created.Title != "Изучить Go" {
		t.Errorf("ожидался очищенный заголовок %q, получен %q",
			"Изучить Go", created.Title)
	}
}

func TestTaskServiceCreateEmptyTitle(t *testing.T) {
	taskService := newInMemoryTaskService()

	for _, title := range []string{"", "   ", "\t\n"} {
		_, err := taskService.Create(context.Background(), model.Task{Title: title})

		if !errors.Is(err, service.ErrEmptyTitle) {
			t.Errorf("для заголовка %q ожидалась ErrEmptyTitle, получено %v",
				title, err)
		}
	}

	tasks, err := taskService.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(tasks) != 0 {
		t.Error("при ошибке валидации задачи не должны создаваться")
	}
}

func TestTaskServiceUpdate(t *testing.T) {
	taskService := newInMemoryTaskService()

	created, err := taskService.Create(context.Background(), model.Task{
		Title: "Old title",
	})
	if err != nil {
		t.Fatal(err)
	}

	updatedTask, err := taskService.Update(context.Background(), created.ID, model.Task{
		Title:  "Новая задача",
		IsDone: true,
	})
	if err != nil {
		t.Fatalf("не ожидалась ошибка обновления: %v", err)
	}

	if updatedTask.ID != created.ID {
		t.Errorf("ожидался ID %d, получен %d", created.ID, updatedTask.ID)
	}
	if updatedTask.Title != "Новая задача" {
		t.Errorf("неверный заголовок: %q", updatedTask.Title)
	}
	if !updatedTask.IsDone {
		t.Error("ожидалось, что задача будет выполнена")
	}

	_, err = taskService.Update(context.Background(), 999, model.Task{
		Title: "Несуществующая задача",
	})
	if !errors.Is(err, service.ErrTaskNotFound) {
		t.Errorf("ожидалась ErrTaskNotFound, получено %v", err)
	}
}

func TestTaskServiceUpdateEmptyTitleDoesNotModifyTask(t *testing.T) {
	taskService := newInMemoryTaskService()

	created, err := taskService.Create(context.Background(), model.Task{
		Title:  "Исходный заголовок",
		IsDone: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, title := range []string{"", "   "} {
		_, err := taskService.Update(context.Background(), created.ID, model.Task{
			Title:  title,
			IsDone: false,
		})

		if !errors.Is(err, service.ErrEmptyTitle) {
			t.Errorf("ожидалась ErrEmptyTitle, получено %v", err)
		}
	}

	tasks, err := taskService.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("ожидалась 1 задача, получено %d", len(tasks))
	}
	if tasks[0].Title != "Исходный заголовок" || !tasks[0].IsDone {
		t.Error("задача не должна измениться после ошибки валидации")
	}
}

func TestTaskServicePatch(t *testing.T) {
	taskService := newInMemoryTaskService()

	created, err := taskService.Create(context.Background(), model.Task{
		Title:  "Old title",
		IsDone: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	newTitle := "Новая задача"
	isDone := false

	updatedTask, err := taskService.Patch(context.Background(), created.ID, model.TaskPatch{
		Title:  &newTitle,
		IsDone: &isDone,
	})
	if err != nil {
		t.Fatalf("не ожидалась ошибка PATCH: %v", err)
	}

	if updatedTask.Title != newTitle {
		t.Errorf("неверный заголовок: %q", updatedTask.Title)
	}
	if updatedTask.IsDone {
		t.Error("ожидалось, что задача будет не выполнена")
	}

	anotherTitle := "  Только заголовок  "

	updatedTask, err = taskService.Patch(context.Background(), created.ID, model.TaskPatch{
		Title: &anotherTitle,
	})
	if err != nil {
		t.Fatalf("не ожидалась ошибка PATCH заголовка: %v", err)
	}

	if updatedTask.Title != strings.TrimSpace(anotherTitle) {
		t.Errorf("ожидался заголовок %q, получен %q",
			strings.TrimSpace(anotherTitle), updatedTask.Title)
	}
	if updatedTask.IsDone {
		t.Error("IsDone должен остаться false")
	}

	_, err = taskService.Patch(context.Background(), 999, model.TaskPatch{Title: &newTitle})
	if !errors.Is(err, service.ErrTaskNotFound) {
		t.Errorf("ожидалась ErrTaskNotFound, получено %v", err)
	}
}

func TestTaskServicePatchEmptyTitleDoesNotModifyTask(t *testing.T) {
	taskService := newInMemoryTaskService()

	created, err := taskService.Create(context.Background(), model.Task{
		Title:  "Исходный заголовок",
		IsDone: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, title := range []string{"", "   "} {
		_, err := taskService.Patch(context.Background(), created.ID, model.TaskPatch{
			Title: &title,
		})

		if !errors.Is(err, service.ErrEmptyTitle) {
			t.Errorf("ожидалась ErrEmptyTitle, получено %v", err)
		}
	}

	tasks, err := taskService.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("ожидалась 1 задача, получено %d", len(tasks))
	}
	if tasks[0].Title != "Исходный заголовок" || !tasks[0].IsDone {
		t.Error("задача не должна измениться после ошибки валидации")
	}
}

func TestTaskServicePatchOnlyIsDone(t *testing.T) {
	taskService := newInMemoryTaskService()

	created, err := taskService.Create(context.Background(), model.Task{
		Title:  "Исходный заголовок",
		IsDone: false,
	})
	if err != nil {
		t.Fatal(err)
	}

	isDone := true
	updatedTask, err := taskService.Patch(context.Background(), created.ID, model.TaskPatch{
		IsDone: &isDone,
	})
	if err != nil {
		t.Fatalf("не ожидалась ошибка PATCH: %v", err)
	}

	if updatedTask.Title != "Исходный заголовок" {
		t.Errorf("заголовок не должен измениться: %q", updatedTask.Title)
	}
	if !updatedTask.IsDone {
		t.Error("ожидалось, что IsDone станет true")
	}
}

func TestTaskServiceDelete(t *testing.T) {
	taskService := newInMemoryTaskService()

	_, err := taskService.Create(context.Background(), model.Task{Title: "Задача 1"})
	if err != nil {
		t.Fatal(err)
	}
	created, err := taskService.Create(context.Background(), model.Task{Title: "Задача 2"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = taskService.Create(context.Background(), model.Task{Title: "Задача 3"})
	if err != nil {
		t.Fatal(err)
	}

	err = taskService.Delete(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("ожидалось успешное удаление задачи, получено: %v", err)
	}

	tasks, err := taskService.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(tasks) != 2 {
		t.Fatalf("ожидалось 2 задачи, получено %d", len(tasks))
	}

	for _, task := range tasks {
		if task.ID == created.ID {
			t.Error("удалённая задача всё ещё присутствует")
		}
	}

	err = taskService.Delete(context.Background(), created.ID)
	if !errors.Is(err, service.ErrTaskNotFound) {
		t.Errorf("ожидалась ошибка ErrTaskNotFound, получено: %v", err)
	}

	err = taskService.Delete(context.Background(), 999)
	if !errors.Is(err, service.ErrTaskNotFound) {
		t.Errorf("ожидалась ошибка ErrTaskNotFound, получено: %v", err)
	}
}

func TestTaskServiceConcurrentCreate(t *testing.T) {
	const count = 100

	taskService := newInMemoryTaskService()
	var wg sync.WaitGroup

	for i := 0; i < count; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			_, err := taskService.Create(context.Background(), model.Task{
				Title: fmt.Sprintf("Задача %d", i),
			})
			if err != nil {
				t.Errorf("ошибка создания задачи %d: %v", i, err)
			}
		}(i)
	}

	wg.Wait()

	tasks, err := taskService.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(tasks) != count {
		t.Fatalf("ожидалось %d задач, получено %d", count, len(tasks))
	}

	seenIDs := make(map[int]bool)

	for _, task := range tasks {
		if task.ID <= 0 {
			t.Errorf("ожидался положительный ID, получен %d", task.ID)
		}
		if seenIDs[task.ID] {
			t.Errorf("обнаружен повторяющийся ID: %d", task.ID)
		}
		seenIDs[task.ID] = true
	}
}
