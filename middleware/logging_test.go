package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusRecorderWriteHeaderNotFound(t *testing.T) {
	recorder := &statusRecorder{
		ResponseWriter: httptest.NewRecorder(),
	}
	recorder.WriteHeader(http.StatusNotFound)

	if recorder.status != http.StatusNotFound {
		t.Errorf("ожидался статус %d, получен %d", http.StatusNotFound, recorder.status)
	}
}

func TestStatusRecorderWriteHeaderDoubleCall(t *testing.T) {
	recorder := &statusRecorder{
		ResponseWriter: httptest.NewRecorder(),
	}
	recorder.WriteHeader(http.StatusNotFound)
	recorder.WriteHeader(http.StatusOK)

	if recorder.status != http.StatusNotFound {
		t.Errorf("ожидался статус %d, получен %d", http.StatusNotFound, recorder.status)
	}
}

func TestStatusRecorderWrite(t *testing.T) {
	response := httptest.NewRecorder()

	recorder := &statusRecorder{
		ResponseWriter: response,
	}

	n, err := recorder.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("ошибка записи тела: %v", err)
	}

	if n != 5 {
		t.Errorf("ожидалось записать 5 байт, записано %d", n)
	}

	if recorder.status != http.StatusOK {
		t.Errorf(
			"ожидался статус %d, получен %d",
			http.StatusOK,
			recorder.status,
		)
	}

	if response.Body.String() != "hello" {
		t.Errorf(
			"ожидалось тело %q, получено %q",
			"hello",
			response.Body.String(),
		)
	}
}
