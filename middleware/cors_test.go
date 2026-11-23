package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var allowedOrigins = []string{"http://localhost:3000"}

func TestCORSAllowedOrigin(t *testing.T) {
	response := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
	}), allowedOrigins)

	handler.ServeHTTP(response, req)

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("ожидался разрешённый origin, получен %q", got)
	}

	if response.Code != http.StatusOK {
		t.Errorf("ожидался статус %d, получен %d", http.StatusOK, response.Code)
	}

	if got := response.Body.String(); got != "hello" {
		t.Errorf("ожидалось тело %q, получено %q", "hello", got)
	}
}

func TestCORSOptionsRequest(t *testing.T) {
	response := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodOptions, "/tasks", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "GET")
	req.Header.Set("Access-Control-Request-Headers", "content-type")

	nextCalled := false

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	}), allowedOrigins)

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusNoContent {
		t.Errorf("ожидался статус %d, получен %d", http.StatusNoContent, response.Code)
	}

	if response.Body.Len() != 0 {
		t.Errorf("ожидалось пустое тело, получено %q", response.Body.String())
	}

	if nextCalled {
		t.Error("следующий обработчик не должен вызываться для OPTIONS")
	}

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "http://localhost:3000" {
		t.Errorf("ожидался разрешённый origin, получен %q", got)
	}
}

func TestCORSOKRequest(t *testing.T) {
	response := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)

	handler := CORSMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello"))
	}), allowedOrigins)

	handler.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("ожидался статус %d, получен %d", http.StatusOK, response.Code)
	}

	if got := response.Body.String(); got != "hello" {
		t.Errorf("ожидалось тело %q, получено %q", "hello", got)
	}
}

func TestCORSNotAllowedOrigin(t *testing.T) {
	response := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
	req.Header.Set("Origin", "http://localhost:3000")

	handler := CORSMiddleware(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("hello"))
		}),
		[]string{"http://localhost:8888"},
	)

	handler.ServeHTTP(response, req)

	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf("для запрещённого origin заголовок должен отсутствовать, получен %q", got)
	}

	if response.Code != http.StatusOK {
		t.Errorf("ожидался статус %d, получен %d", http.StatusOK, response.Code)
	}

	if got := response.Body.String(); got != "hello" {
		t.Errorf("ожидалось тело %q, получено %q", "hello", got)
	}
}

func TestCORSOriginHeaders(t *testing.T) {
	tests := []struct {
		name           string
		requestOrigin  string
		allowedOrigins []string
		wantHeader     string
	}{
		{
			name:           "разрешённый origin",
			requestOrigin:  "http://localhost:3000",
			allowedOrigins: []string{"http://localhost:3000", "http://localhost:8888"},
			wantHeader:     "http://localhost:3000",
		},
		{
			name:           "запрещённый origin",
			requestOrigin:  "http://localhost:5000",
			allowedOrigins: []string{"http://localhost:3000", "http://localhost:8888"},
			wantHeader:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/tasks", nil)
			req.Header.Set("Origin", tt.requestOrigin)

			handler := CORSMiddleware(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				}),
				tt.allowedOrigins,
			)

			handler.ServeHTTP(response, req)

			got := response.Header().Get("Access-Control-Allow-Origin")
			if got != tt.wantHeader {
				t.Errorf(
					"Access-Control-Allow-Origin: ожидалось %q, получено %q",
					tt.wantHeader,
					got,
				)
			}
		})
	}
}

func TestCORSWithoutOrigin(t *testing.T) {
	response := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/tasks", nil)

	nextCalled := false

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	})

	handler := CORSMiddleware(
		nextHandler,
		[]string{
			"http://localhost:3000",
			"http://localhost:8888",
		},
	)

	handler.ServeHTTP(response, req)

	// Проверяем, что CORS-заголовок не установлен.
	if got := response.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Errorf(
			"ожидался пустой Access-Control-Allow-Origin, получено %q",
			got,
		)
	}

	// Проверяем, что запрос дошёл до следующего обработчика.
	if !nextCalled {
		t.Error("ожидалось, что запрос будет передан следующему обработчику")
	}

	// Проверяем, что обработчик вернул 200 OK.
	if response.Code != http.StatusOK {
		t.Errorf(
			"ожидался статус %d, получен %d",
			http.StatusOK,
			response.Code,
		)
	}
}
