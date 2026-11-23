package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"example.com/go-backend/model"
	"example.com/go-backend/service"
)

type TaskHandler struct {
	Service TaskService
}

func NewTaskHandler(service TaskService) *TaskHandler {
	return &TaskHandler{
		Service: service,
	}
}

func decodeJSON(r *http.Request, dst any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		if err == nil {
			return errors.New("в запросе должен быть один JSON-объект")
		}
		return err
	}

	return nil
}

func (h *TaskHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.Service.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Ошибка получения списка задач", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		fmt.Println("Ошибка сериализации", err)
	}
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	newTask := model.Task{}

	if err := decodeJSON(r, &newTask); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	createdTask, err := h.Service.Create(r.Context(), newTask)
	if errors.Is(err, service.ErrEmptyTitle) {
		http.Error(
			w,
			"Название задачи не может быть пустым",
			http.StatusBadRequest,
		)
		return
	}
	if err != nil {
		http.Error(w, "Ошибка создания задачи", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(createdTask); err != nil {
		fmt.Println("Ошибка отправки ответа:", err)
	}
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	updatedTask := model.Task{}

	taskId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || taskId <= 0 {
		http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
		return
	}

	if err := decodeJSON(r, &updatedTask); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	updatedTask, updateError := h.Service.Update(r.Context(), taskId, updatedTask)
	if errors.Is(updateError, service.ErrEmptyTitle) {
		http.Error(w, "Название задачи не может быть пустым",
			http.StatusBadRequest)
		return
	}

	if errors.Is(updateError, service.ErrTaskNotFound) {
		http.Error(w, "Задача не найдена",
			http.StatusNotFound)
		return
	}

	if updateError != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		fmt.Println("Ошибка отправки ответа", err)
	}

}

func (h *TaskHandler) Patch(w http.ResponseWriter, r *http.Request) {
	patch := model.TaskPatch{}

	taskId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || taskId <= 0 {
		http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
		return
	}

	if err := decodeJSON(r, &patch); err != nil {
		http.Error(w, "Некорректный JSON", http.StatusBadRequest)
		return
	}

	if patch.Title == nil && patch.IsDone == nil {
		http.Error(w, "Нет полей для обновления", http.StatusBadRequest)
		return
	}

	updatedTask, patchError := h.Service.Patch(r.Context(), taskId, patch)
	if errors.Is(patchError, service.ErrEmptyTitle) {
		http.Error(w, "Название задачи не может быть пустым",
			http.StatusBadRequest)
		return
	}

	if errors.Is(patchError, service.ErrTaskNotFound) {
		http.Error(w, "Задача не найдена",
			http.StatusNotFound)
		return
	}

	if patchError != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
		fmt.Println("Ошибка отправки ответа", err)
	}

}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	taskId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || taskId <= 0 {
		http.Error(w, "Некорректный ID задачи", http.StatusBadRequest)
		return
	}

	err = h.Service.Delete(r.Context(), taskId)

	if errors.Is(err, service.ErrTaskNotFound) {
		http.Error(w, "Задача не найдена", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
