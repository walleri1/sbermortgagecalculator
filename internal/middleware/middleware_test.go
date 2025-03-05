package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestLoggingMiddleware(t *testing.T) {
	var buffer bytes.Buffer
	log.SetOutput(&buffer)

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(50 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	handlerToTest := LoggingMiddleware(nextHandler)

	req := httptest.NewRequest("GET", "http://sber.com/test", nil)
	recorder := httptest.NewRecorder()

	handlerToTest.ServeHTTP(recorder, req)

	if status := recorder.Result().StatusCode; status != http.StatusOK {
		t.Fatalf("the status was expected to be %d, but received %d", http.StatusOK, status)
	}

	expectedBody := "OK"
	if body := recorder.Body.String(); body != expectedBody {
		t.Fatalf("the response body was expected to be '%s', but we received '%s'", expectedBody, body)
	}

	logOutput := buffer.String()
	if !strings.Contains(logOutput, "status_code: 200") {
		t.Fatalf("the log should contain 'status_code: 200', but the contents of the log: '%s'", logOutput)
	}
	if !strings.Contains(logOutput, "duration:") {
		t.Fatalf("the log should contain 'duration:', but the contents of the log: '%s'", logOutput)
	}
}
