package handler

import "net/http"

func NewRouter(taskHandler *TaskHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /tasks", taskHandler.GetAll)
	mux.HandleFunc("POST /tasks", taskHandler.Create)
	mux.HandleFunc("PUT /tasks/{id}", taskHandler.Update)
	mux.HandleFunc("PATCH /tasks/{id}", taskHandler.Patch)
	mux.HandleFunc("DELETE /tasks/{id}", taskHandler.Delete)
	mux.HandleFunc("GET /health", taskHandler.Health)

	return mux
}
