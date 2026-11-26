package service

import (
	"context"
	"errors"
	"testing"

	"example.com/go-backend/model"
)

func TestTaskServiceCreateTrimsTitle(t *testing.T) {
	repo := NewFakeTaskRepository([]model.Task{})

	service := NewTaskService(repo)

	newTask := model.Task{
		Title:  "    Изучить Go    ",
		IsDone: false,
	}
	createdTask, err := service.Create(context.Background(), newTask)
	if err != nil {
		t.Fatalf("Ожидалось успешное создание задачи, получено %v", err)
	}
	if repo.CreateCalls != 1 {
		t.Errorf("Ожидался 1 вызов Create, получено %d", repo.CreateCalls)
	}
	expectedTitle := "Изучить Go"
	if createdTask.Title != expectedTitle {
		t.Errorf("Неожиданный Title, ожидалось %q, получено %q", expectedTitle, createdTask.Title)
	}
	if repo.CreatedTask.Title != expectedTitle {
		t.Errorf("Неожиданный Title, ожидалось %q, получено %q", expectedTitle, repo.CreatedTask.Title)
	}
}

func TestTaskServiceDeleteNotFound(t *testing.T) {
	repo := NewFakeTaskRepository([]model.Task{})

	service := NewTaskService(repo)

	err := service.Delete(context.Background(), 999)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("Ожидалась ошибка %v, получено: %v", ErrTaskNotFound, err)
	}
	if repo.DeleteCalls != 1 {
		t.Errorf("Ожидался 1 вызов Delete, получено %d", repo.DeleteCalls)
	}

}

func TestTaskServiceDeleteSuccess(t *testing.T) {
	repo := NewFakeTaskRepository([]model.Task{
		{ID: 1, Title: "Задача 1"},
	})

	service := NewTaskService(repo)

	err := service.Delete(context.Background(), 1)
	if err != nil {
		t.Fatalf("Ожидалось успешное удаление задачи, получено %v", err)
	}
	if repo.DeleteCalls != 1 {
		t.Errorf("Ожидался 1 вызов Delete, получено %d", repo.DeleteCalls)
	}
	if repo.DeletedID != 1 {
		t.Errorf("Ожидалась 1 в поле DeletedID, получено %d", repo.DeletedID)
	}
	if len(repo.Tasks) != 0 {
		t.Errorf("Не ожидалось задач в срезе, получено %d задач", len(repo.Tasks))
	}
}

func TestTaskServicePatchTitleOnly(t *testing.T) {
	repo := NewFakeTaskRepository([]model.Task{
		{
			ID:     1,
			Title:  "Старая задача",
			IsDone: true,
		},
	})

	service := NewTaskService(repo)
	expectedTitle := "Новое название"
	patchedTask, err := service.Patch(context.Background(), 1, model.TaskPatch{Title: &expectedTitle})

	if err != nil {
		t.Fatalf("Ожидалось успешное обновление задачи, получено %v", err)
	}

	if patchedTask.Title != expectedTitle {
		t.Errorf("Неожиданный Title, ожидалось %q, получено %q", expectedTitle, patchedTask.Title)
	}

	if !patchedTask.IsDone {
		t.Errorf("Ожидалось IsDone = true")
	}

	if len(repo.Tasks) != 1 {
		t.Fatalf("Ожидалась 1 задача, получено %d", len(repo.Tasks))
	}

	if repo.Tasks[0].Title != expectedTitle {
		t.Errorf("В репозитории ожидалось название %q, получено %q",
			expectedTitle, repo.Tasks[0].Title)
	}

	if !repo.Tasks[0].IsDone {
		t.Errorf("В репозитории ожидалось IsDone = true")
	}

	if repo.PatchCalls != 1 {
		t.Errorf("Ожидался 1 вызов Patch, получено %d", repo.PatchCalls)
	}
}

func TestTaskServicePatchIsDoneOnly(t *testing.T) {
	expectedTitle := "Изучить GO"
	repo := NewFakeTaskRepository([]model.Task{
		{
			ID:     1,
			Title:  expectedTitle,
			IsDone: false,
		},
	})

	service := NewTaskService(repo)
	expectedIsDone := true
	patchedTask, err := service.Patch(context.Background(), 1, model.TaskPatch{IsDone: &expectedIsDone})

	if err != nil {
		t.Fatalf("Ожидалось успешное обновление задачи, получено %v", err)
	}

	if patchedTask.Title != expectedTitle {
		t.Errorf("Неожиданный Title, ожидалось %q, получено %q", expectedTitle, patchedTask.Title)
	}

	if !patchedTask.IsDone {
		t.Errorf("Ожидалось IsDone = true")
	}

	if len(repo.Tasks) != 1 {
		t.Fatalf("Ожидалась 1 задача, получено %d", len(repo.Tasks))
	}

	if repo.Tasks[0].Title != expectedTitle {
		t.Errorf("В репозитории ожидалось название %q, получено %q",
			expectedTitle, repo.Tasks[0].Title)
	}

	if !repo.Tasks[0].IsDone {
		t.Errorf("В репозитории ожидалось IsDone = true")
	}

	if repo.PatchCalls != 1 {
		t.Errorf("Ожидался 1 вызов Patch, получено %d", repo.PatchCalls)
	}
}

func TestTaskServicePatchIsDoneFalse(t *testing.T) {
	expectedTitle := "Изучить GO"
	repo := NewFakeTaskRepository([]model.Task{
		{
			ID:     1,
			Title:  expectedTitle,
			IsDone: true,
		},
	})

	service := NewTaskService(repo)
	expectedIsDone := false
	patchedTask, err := service.Patch(context.Background(), 1, model.TaskPatch{IsDone: &expectedIsDone})

	if err != nil {
		t.Fatalf("Ожидалось успешное обновление задачи, получено %v", err)
	}

	if patchedTask.Title != expectedTitle {
		t.Errorf("Неожиданный Title, ожидалось %q, получено %q", expectedTitle, patchedTask.Title)
	}

	if patchedTask.IsDone {
		t.Errorf("Ожидалось IsDone = false")
	}

	if len(repo.Tasks) != 1 {
		t.Fatalf("Ожидалась 1 задача, получено %d", len(repo.Tasks))
	}

	if repo.Tasks[0].Title != expectedTitle {
		t.Errorf("В репозитории ожидалось название %q, получено %q",
			expectedTitle, repo.Tasks[0].Title)
	}

	if repo.Tasks[0].IsDone {
		t.Errorf("В репозитории ожидалось IsDone = false")
	}

	if repo.PatchCalls != 1 {
		t.Errorf("Ожидался 1 вызов Patch, получено %d", repo.PatchCalls)
	}
}

func TestTaskServicePatchNotFound(t *testing.T) {
	repo := NewFakeTaskRepository([]model.Task{})

	service := NewTaskService(repo)
	newTitle := "Новая задача"
	_, err := service.Patch(context.Background(), 999, model.TaskPatch{Title: &newTitle})
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("Ожидалась ошибка %v, получено: %v", ErrTaskNotFound, err)
	}
	if repo.PatchCalls != 1 {
		t.Errorf("Ожидался 1 вызов Patch, получено %d", repo.PatchCalls)
	}
}
func TestTaskServiceUpdateEmptyTitle(t *testing.T) {
	expectedTitle := "Старая задача"
	repo := NewFakeTaskRepository([]model.Task{
		{
			ID:     1,
			Title:  expectedTitle,
			IsDone: false,
		},
	})

	service := NewTaskService(repo)
	_, err := service.Update(context.Background(), 1, model.Task{Title: "   ", IsDone: true})
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("Ожидалась ошибка %v, получено: %v", ErrEmptyTitle, err)
	}
	if repo.Tasks[0].Title != expectedTitle {
		t.Errorf("В репозитории ожидалось название %q, получено %q",
			expectedTitle, repo.Tasks[0].Title)
	}

	if repo.Tasks[0].IsDone {
		t.Errorf("В репозитории ожидалось IsDone = false")
	}
	if repo.UpdateCalls != 0 {
		t.Errorf("Ожидалось 0 вызовов Update, получено %d", repo.UpdateCalls)
	}
}

func TestTaskServiceUpdateNotFound(t *testing.T) {
	repo := NewFakeTaskRepository([]model.Task{})

	service := NewTaskService(repo)
	_, err := service.Update(context.Background(), 999, model.Task{Title: "Новая задача", IsDone: true})
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("Ожидалась ошибка %v, получено: %v", ErrTaskNotFound, err)
	}

	if repo.UpdateCalls != 1 {
		t.Errorf("Ожидался 1 вызов Update, получено %d", repo.UpdateCalls)
	}
}

func TestTaskServiceCreateValidation(t *testing.T) {
	tests := []struct {
		name      string
		title     string
		wantError error
	}{
		{
			name:      "empty title",
			title:     "",
			wantError: ErrEmptyTitle,
		},
		{
			name:      "whitespace title",
			title:     "   ",
			wantError: ErrEmptyTitle,
		},
		{
			name:      "valid title",
			title:     "Изучить Go",
			wantError: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := NewFakeTaskRepository([]model.Task{})

			service := NewTaskService(repo)
			_, err := service.Create(context.Background(), model.Task{Title: tt.title, IsDone: false})
			if !errors.Is(err, tt.wantError) {
				t.Fatalf("Ожидалась ошибка %v, получено: %v", tt.wantError, err)
			}
			if tt.wantError != nil && repo.CreateCalls != 0 {
				t.Errorf("Ожидалось 0 вызовов Create, получено %d", repo.CreateCalls)
			}
			if tt.wantError == nil && repo.CreateCalls != 1 {
				t.Errorf("Ожидался 1 вызов Create, получено %d", repo.CreateCalls)
			}
		})
	}
}

func TestTaskServiceUpdateSuccess(t *testing.T) {
	expectedTitle := "Обновленная задача"
	repo := NewFakeTaskRepository([]model.Task{
		{
			ID:     1,
			Title:  "Старая задача",
			IsDone: false,
		},
	})

	service := NewTaskService(repo)
	updatedTask, err := service.Update(context.Background(), 1, model.Task{IsDone: true, Title: expectedTitle})

	if err != nil {
		t.Fatalf("Ожидалось успешное обновление задачи, получено %v", err)
	}

	if updatedTask.Title != expectedTitle {
		t.Errorf("Неожиданный Title, ожидалось %q, получено %q", expectedTitle, updatedTask.Title)
	}

	if !updatedTask.IsDone {
		t.Errorf("Ожидалось IsDone = true")
	}

	if updatedTask.ID != 1 {
		t.Errorf("Ожидался ID = 1, получено %d", updatedTask.ID)
	}

	if len(repo.Tasks) != 1 {
		t.Fatalf("Ожидалась 1 задача, получено %d", len(repo.Tasks))
	}

	if repo.Tasks[0].Title != expectedTitle {
		t.Errorf("В репозитории ожидалось название %q, получено %q",
			expectedTitle, repo.Tasks[0].Title)
	}

	if !repo.Tasks[0].IsDone {
		t.Errorf("В репозитории ожидалось IsDone = true")
	}

	if repo.Tasks[0].ID != 1 {
		t.Errorf("В репозитории ожидался ID = 1, получено %d", repo.Tasks[0].ID)
	}

	if repo.UpdateCalls != 1 {
		t.Errorf("Ожидался 1 вызов Update, получено %d", repo.UpdateCalls)
	}
}

func TestTaskServicePatchEmptyTitle(t *testing.T) {
	expectedTitle := "Старая задача"

	repo := NewFakeTaskRepository([]model.Task{
		{
			ID:     1,
			Title:  expectedTitle,
			IsDone: false,
		},
	})

	service := NewTaskService(repo)
	title := "   "

	_, err := service.Patch(context.Background(), 1, model.TaskPatch{
		Title: &title,
	})

	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("Ожидалась ошибка %v, получено: %v", ErrEmptyTitle, err)
	}

	if repo.PatchCalls != 0 {
		t.Errorf("Ожидалось 0 вызовов Patch, получено %d", repo.PatchCalls)
	}

	if len(repo.Tasks) != 1 {
		t.Fatalf("Ожидалась 1 задача, получено %d", len(repo.Tasks))
	}

	if repo.Tasks[0].Title != expectedTitle {
		t.Errorf("Ожидалось название %q, получено %q",
			expectedTitle, repo.Tasks[0].Title)
	}

	if repo.Tasks[0].IsDone {
		t.Errorf("Ожидалось IsDone = false")
	}
}

func TestTaskServicePatchTrimsTitle(t *testing.T) {
	expectedTitle := "Новое название"

	repo := NewFakeTaskRepository([]model.Task{
		{
			ID:     1,
			Title:  "Старая задача",
			IsDone: false,
		},
	})

	service := NewTaskService(repo)
	title := "   " + expectedTitle + "   "

	updatedTask, err := service.Patch(context.Background(), 1, model.TaskPatch{
		Title: &title,
	})

	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}

	if updatedTask.Title != expectedTitle {
		t.Errorf("Ожидалось название %q, получено %q",
			expectedTitle, updatedTask.Title)
	}

	if repo.Tasks[0].Title != expectedTitle {
		t.Errorf("В репозитории ожидалось название %q, получено %q",
			expectedTitle, repo.Tasks[0].Title)
	}

	if repo.PatchCalls != 1 {
		t.Errorf("Ожидался 1 вызов Patch, получено %d", repo.PatchCalls)
	}
}

func TestTaskServiceUpdateTrimsTitle(t *testing.T) {
	expectedTitle := "Обновлённая задача"

	repo := NewFakeTaskRepository([]model.Task{
		{
			ID:     1,
			Title:  "Старая задача",
			IsDone: false,
		},
	})

	service := NewTaskService(repo)

	updatedTask, err := service.Update(context.Background(), 1, model.Task{
		Title:  "   " + expectedTitle + "   ",
		IsDone: true,
	})

	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}

	if updatedTask.Title != expectedTitle {
		t.Errorf("Ожидалось название %q, получено %q",
			expectedTitle, updatedTask.Title)
	}

	if repo.Tasks[0].Title != expectedTitle {
		t.Errorf("В репозитории ожидалось название %q, получено %q",
			expectedTitle, repo.Tasks[0].Title)
	}

	if repo.UpdateCalls != 1 {
		t.Errorf("Ожидался 1 вызов Update, получено %d", repo.UpdateCalls)
	}
}

func TestTaskServiceGetAll(t *testing.T) {
	expectedTasks := []model.Task{
		{ID: 1, Title: "Изучить Go", IsDone: false},
		{ID: 2, Title: "Написать тесты", IsDone: true},
	}

	repo := NewFakeTaskRepository(expectedTasks)
	service := NewTaskService(repo)

	tasks, err := service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}

	if len(tasks) != len(expectedTasks) {
		t.Fatalf("Ожидалось %d задач, получено %d",
			len(expectedTasks), len(tasks))
	}

	for i, expected := range expectedTasks {
		if tasks[i] != expected {
			t.Errorf("Задача с индексом %d: ожидалось %+v, получено %+v",
				i, expected, tasks[i])
		}
	}

	if repo.GetAllCalls != 1 {
		t.Errorf("Ожидался 1 вызов GetAll, получено %d", repo.GetAllCalls)
	}
}

func TestTaskServiceCreateInvalidPriority(t *testing.T) {
	repo := NewFakeTaskRepository([]model.Task{})
	service := NewTaskService(repo)

	_, err := service.Create(context.Background(), model.Task{
		Title:    "Learn Go",
		Priority: 3,
	})

	if !errors.Is(err, ErrInvalidPriority) {
		t.Fatalf("expected %v, got %v", ErrInvalidPriority, err)
	}

	if repo.CreateCalls != 0 {
		t.Errorf("expected 0 Create calls, got %d", repo.CreateCalls)
	}
}

func TestTaskServiceUpdateInvalidPriority(t *testing.T) {
	repo := NewFakeTaskRepository([]model.Task{
		{ID: 1, Title: "Learn Go"},
	})
	service := NewTaskService(repo)

	_, err := service.Update(context.Background(), 1, model.Task{
		Title:    "Learn Go",
		Priority: -1,
	})

	if !errors.Is(err, ErrInvalidPriority) {
		t.Fatalf("expected %v, got %v", ErrInvalidPriority, err)
	}

	if repo.UpdateCalls != 0 {
		t.Errorf("expected 0 Update calls, got %d", repo.UpdateCalls)
	}
}

func TestTaskServicePatchInvalidPriority(t *testing.T) {
	repo := NewFakeTaskRepository([]model.Task{
		{ID: 1, Title: "Learn Go"},
	})
	service := NewTaskService(repo)

	priority := 5

	_, err := service.Patch(context.Background(), 1, model.TaskPatch{
		Priority: &priority,
	})

	if !errors.Is(err, ErrInvalidPriority) {
		t.Fatalf("expected %v, got %v", ErrInvalidPriority, err)
	}

	if repo.PatchCalls != 0 {
		t.Errorf("expected 0 Patch calls, got %d", repo.PatchCalls)
	}
}
