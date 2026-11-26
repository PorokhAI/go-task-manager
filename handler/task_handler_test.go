package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"example.com/go-backend/model"
	"example.com/go-backend/repository"
	"example.com/go-backend/service"
)

func newTestTaskHandler() *TaskHandler {
	repo := repository.NewInMemoryTaskRepository(nil)
	taskService := service.NewTaskService(repo)

	return NewTaskHandler(taskService)
}

func TestTaskHandlerCreate(t *testing.T) {
	handler := newTestTaskHandler()

	body := `{"title":"Изучить Go","isDone":false}`
	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusCreated,
			res.StatusCode,
		)
	}

	if contentType := res.Header.Get("Content-Type"); contentType != "application/json; charset=utf-8" {
		t.Errorf("ожидался Content-Type application/json; charset=utf-8, получен %q", contentType)
	}

	var task model.Task
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		t.Fatalf("не удалось декодировать ответ: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("ожидался ID 1, получен %d", task.ID)
	}

	if task.Title != "Изучить Go" {
		t.Errorf("неожиданный Title: %q", task.Title)
	}

	if task.IsDone {
		t.Error("ожидалось IsDone = false")
	}
}

func TestTaskHandlerPostEmptyTitle(t *testing.T) {
	handler := newTestTaskHandler()

	body := `{"title":"   ","isDone":false}`
	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusBadRequest,
			res.StatusCode,
		)
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("не удалось прочитать тело ответа: %v", err)
	}

	expectedMessage := "Название задачи не может быть пустым"
	if !strings.Contains(string(responseBody), expectedMessage) {
		t.Errorf(
			"ожидалось сообщение %q, получено %q",
			expectedMessage,
			string(responseBody),
		)
	}
}

func TestTaskHandlerPostUnknownField(t *testing.T) {
	handler := newTestTaskHandler()

	body := `{"title":"Изучить GO","unknown":123}`
	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusBadRequest,
			res.StatusCode,
		)
	}

	tasks, err := handler.Service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("ожидалось 0 задач, получено %d", len(tasks))
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("не удалось прочитать тело ответа: %v", err)
	}

	expectedMessage := "Некорректный JSON"
	if !strings.Contains(string(responseBody), expectedMessage) {
		t.Errorf(
			"ожидалось сообщение %q, получено %q",
			expectedMessage,
			string(responseBody),
		)
	}
}

func TestTaskHandlerPut(t *testing.T) {
	handler := newTestTaskHandler()

	initialTask, err := handler.Service.Create(context.Background(), model.Task{
		Title:  "Исходная задача",
		IsDone: false,
	})

	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	body := `{"title":"Обновлённая задача","isDone":true}`
	req := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/tasks/%d", initialTask.ID),
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()
	req.SetPathValue("id", strconv.Itoa(initialTask.ID))

	handler.Update(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	if contentType := res.Header.Get("Content-Type"); contentType != "application/json; charset=utf-8" {
		t.Errorf("ожидался Content-Type application/json; charset=utf-8, получен %q", contentType)
	}

	var task model.Task
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		t.Fatalf("не удалось декодировать ответ: %v", err)
	}

	if task.ID != initialTask.ID {
		t.Errorf("ожидался ID %d, получен %d", initialTask.ID, task.ID)
	}

	if task.Title != "Обновлённая задача" {
		t.Errorf("неожиданный Title: %q", task.Title)
	}

	if !task.IsDone {
		t.Error("ожидалось IsDone = true")
	}
}

func TestTaskHandlerPatch(t *testing.T) {
	handler := newTestTaskHandler()

	initialTask, err := handler.Service.Create(context.Background(), model.Task{
		Title:  "Изучить GO",
		IsDone: false,
	})

	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	body := `{"isDone":true}`
	req := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/tasks/%d", initialTask.ID),
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()
	req.SetPathValue("id", strconv.Itoa(initialTask.ID))

	handler.Patch(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusOK,
			res.StatusCode,
		)
	}

	if contentType := res.Header.Get("Content-Type"); contentType != "application/json; charset=utf-8" {
		t.Errorf("ожидался Content-Type application/json; charset=utf-8, получен %q", contentType)
	}

	var task model.Task
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		t.Fatalf("не удалось декодировать ответ: %v", err)
	}

	if task.ID != initialTask.ID {
		t.Errorf("ожидался ID %d, получен %d", initialTask.ID, task.ID)
	}

	if task.Title != "Изучить GO" {
		t.Errorf("неожиданный Title: %q", task.Title)
	}

	if !task.IsDone {
		t.Error("ожидалось IsDone = true")
	}
}

func TestTaskHandlerPutNotFound(t *testing.T) {
	handler := newTestTaskHandler()

	body := `{"title":"Обновлённая задача","isDone":true}`
	req := httptest.NewRequest(
		http.MethodPut,
		"/tasks/999",
		strings.NewReader(body),
	)
	req.SetPathValue("id", "999")

	rec := httptest.NewRecorder()
	handler.Update(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusNotFound,
			res.StatusCode,
		)
	}
}

func TestTaskHandlerPatchEmptyBodyObject(t *testing.T) {
	handler := newTestTaskHandler()

	req := httptest.NewRequest(
		http.MethodPatch,
		"/tasks/1",
		strings.NewReader(`{}`),
	)
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()
	handler.Patch(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusBadRequest {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusBadRequest,
			res.StatusCode,
		)
	}
}

func TestRouterPostTask(t *testing.T) {
	taskHandler := newTestTaskHandler()
	router := NewRouter(taskHandler)

	body := `{"title":"Проверка маршрута","isDone":false}`
	req := httptest.NewRequest(
		http.MethodPost,
		"/tasks",
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusCreated,
			res.StatusCode,
		)
	}

	var task model.Task
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		t.Fatalf("не удалось декодировать ответ: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("ожидался ID 1, получен %d", task.ID)
	}

	if task.Title != "Проверка маршрута" {
		t.Errorf(
			"ожидался Title %q, получен %q",
			"Проверка маршрута",
			task.Title,
		)
	}
}

func TestRouterPutTask(t *testing.T) {
	taskHandler := newTestTaskHandler()
	router := NewRouter(taskHandler)

	createdTask, err := taskHandler.Service.Create(context.Background(), model.Task{
		Title:  "Исходная задача",
		IsDone: false,
	})
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	body := `{"title":"Обновлённая задача","isDone":true}`
	req := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/tasks/%d", createdTask.ID),
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("ожидался статус 200, получен %d", res.StatusCode)
	}

	var task model.Task
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		t.Fatalf("не удалось декодировать ответ: %v", err)
	}

	if task.Title != "Обновлённая задача" || !task.IsDone {
		t.Errorf("задача не обновлена корректно: %+v", task)
	}
}

func TestRouterPatchTask(t *testing.T) {
	taskHandler := newTestTaskHandler()
	router := NewRouter(taskHandler)

	createdTask, err := taskHandler.Service.Create(context.Background(), model.Task{
		Title:  "Изучить Go",
		IsDone: false,
	})
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	body := `{"isDone":true}`
	req := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/tasks/%d", createdTask.ID),
		strings.NewReader(body),
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("ожидался статус 200, получен %d", res.StatusCode)
	}

	var task model.Task
	if err := json.NewDecoder(res.Body).Decode(&task); err != nil {
		t.Fatalf("не удалось декодировать ответ: %v", err)
	}

	if task.Title != "Изучить Go" {
		t.Errorf("PATCH изменил Title: %q", task.Title)
	}

	if !task.IsDone {
		t.Error("ожидалось IsDone = true")
	}
}

func TestRouterDeleteTask(t *testing.T) {
	taskHandler := newTestTaskHandler()
	router := NewRouter(taskHandler)

	createdTask, err := taskHandler.Service.Create(context.Background(), model.Task{
		Title:  "Удалить задачу",
		IsDone: false,
	})
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	req := httptest.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("/tasks/%d", createdTask.ID),
		nil,
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("ожидался статус 204, получен %d", res.StatusCode)
	}

	tasks, err := taskHandler.Service.GetAll(context.Background())
	if err != nil {
		t.Fatalf("Не ожидалась ошибка, получено: %v", err)
	}
	if len(tasks) != 0 {
		t.Errorf("ожидалось 0 задач после удаления, получено %d", len(tasks))
	}
}

func TestRouterMethodNotAllowed(t *testing.T) {
	taskHandler := newTestTaskHandler()
	router := NewRouter(taskHandler)

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks/1",
		nil,
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusMethodNotAllowed,
			res.StatusCode,
		)
	}
}

func TestRouterNotFound(t *testing.T) {
	taskHandler := newTestTaskHandler()
	router := NewRouter(taskHandler)

	req := httptest.NewRequest(
		http.MethodGet,
		"/unknown",
		nil,
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusNotFound,
			res.StatusCode,
		)
	}
}

func TestRouterDeleteTaskNotFound(t *testing.T) {
	taskHandler := newTestTaskHandler()
	router := NewRouter(taskHandler)

	req := httptest.NewRequest(
		http.MethodDelete,
		"/tasks/999",
		nil,
	)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Fatalf(
			"ожидался статус %d, получен %d",
			http.StatusNotFound,
			res.StatusCode,
		)
	}
}

func TestTaskHandlerCreateInvalidPriority(t *testing.T) {
	handler := newTestTaskHandler()

	body := `{"title":"Изучить Go","priority":3}`
	req := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ожидался статус 400, получен %d", rec.Code)
	}
}

func TestTaskHandlerUpdateInvalidPriority(t *testing.T) {
	handler := newTestTaskHandler()

	task, err := handler.Service.Create(context.Background(), model.Task{
		Title: "Изучить Go",
	})
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	body := `{"title":"Изучить Go","priority":-1}`
	req := httptest.NewRequest(
		http.MethodPut,
		fmt.Sprintf("/tasks/%d", task.ID),
		strings.NewReader(body),
	)
	req.SetPathValue("id", strconv.Itoa(task.ID))

	rec := httptest.NewRecorder()
	handler.Update(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ожидался статус 400, получен %d", rec.Code)
	}
}

func TestTaskHandlerPatchInvalidPriority(t *testing.T) {
	handler := newTestTaskHandler()

	task, err := handler.Service.Create(context.Background(), model.Task{
		Title: "Изучить Go",
	})
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	body := `{"priority":3}`
	req := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/tasks/%d", task.ID),
		strings.NewReader(body),
	)
	req.SetPathValue("id", strconv.Itoa(task.ID))

	rec := httptest.NewRecorder()
	handler.Patch(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("ожидался статус 400, получен %d", rec.Code)
	}
}

func TestTaskHandlerPatchOnlyPriority(t *testing.T) {
	handler := newTestTaskHandler()

	task, err := handler.Service.Create(context.Background(), model.Task{
		Title: "Изучить Go",
	})
	if err != nil {
		t.Fatalf("не удалось создать задачу: %v", err)
	}

	body := `{"priority":2}`
	req := httptest.NewRequest(
		http.MethodPatch,
		fmt.Sprintf("/tasks/%d", task.ID),
		strings.NewReader(body),
	)
	req.SetPathValue("id", strconv.Itoa(task.ID))

	rec := httptest.NewRecorder()
	handler.Patch(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("ожидался статус 200, получен %d: %s", rec.Code, rec.Body.String())
	}

	var updatedTask model.Task
	if err := json.NewDecoder(rec.Body).Decode(&updatedTask); err != nil {
		t.Fatalf("не удалось декодировать ответ: %v", err)
	}

	if updatedTask.Priority != 2 {
		t.Errorf("ожидался priority 2, получен %d", updatedTask.Priority)
	}
}
