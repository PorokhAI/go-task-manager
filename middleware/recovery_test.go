package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPanicRecovery(t *testing.T) {
	response := httptest.NewRecorder()
	recorder := &statusRecorder{
		ResponseWriter: response,
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks",
		nil,
	)

	handler := RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("тестовая паника")
	}))

	handler.ServeHTTP(recorder, req)

	if recorder.status != http.StatusInternalServerError {
		t.Errorf(
			"ожидался статус %d, получен %d",
			http.StatusInternalServerError,
			recorder.status,
		)
	}

	if !strings.Contains(response.Body.String(), "Внутренняя ошибка сервера") {
		t.Errorf("неожиданное тело ответа: %q", response.Body.String())
	}
}

func TestOKRequest(t *testing.T) {
	response := httptest.NewRecorder()
	recorder := &statusRecorder{
		ResponseWriter: response,
	}

	req := httptest.NewRequest(
		http.MethodGet,
		"/tasks",
		nil,
	)

	handler := RecoveryMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, err := w.Write([]byte("hello"))
		if err != nil {
			t.Fatalf("ошибка записи тела: %v", err)
		}

		if n != 5 {
			t.Errorf("ожидалось записать 5 байт, записано %d", n)
		}

		w.WriteHeader(http.StatusOK)
	}))

	handler.ServeHTTP(recorder, req)

	if recorder.status != http.StatusOK {
		t.Errorf("ожидался статус %d, получен %d", http.StatusOK, recorder.status)
	}
	if response.Body.String() != "hello" {
		t.Errorf(
			"ожидалось тело %q, получено %q",
			"hello",
			response.Body.String(),
		)
	}
}
