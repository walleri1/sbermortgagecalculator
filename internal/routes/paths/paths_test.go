package paths

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"sbermortgagecalculator/internal/models"
)

func init() {
	InitCache(10, time.Hour)
}

func TestGetLoansFromSyncMap(t *testing.T) {
	var m sync.Map

	loan1 := models.CachedLoan{ID: 0}
	loan2 := models.CachedLoan{ID: 1}
	m.Store(0, loan1)
	m.Store(1, loan2)
	m.Store(2, "invalid_type")

	loans := getLoansFromSyncMap(&m)

	// Ожидаемый результат
	expectedLoans := []models.CachedLoan{loan1, loan2}

	// Проверка длины результата
	if len(loans) != len(expectedLoans) {
		t.Errorf("Expected %d loans, but got %d", len(expectedLoans), len(loans))
	}

	// Проверка содержимого результата
	for i, loan := range loans {
		if loan.ID != expectedLoans[i].ID {
			t.Errorf("Expected loan ID %d, but got %d", expectedLoans[i].ID, loan.ID)
		}
	}
}

func TestWriteJSONResponse(t *testing.T) {
	recorder := httptest.NewRecorder()

	data := map[string]string{"message": "success"}
	statusCode := http.StatusOK

	writeJSONResponse(recorder, data, statusCode)

	if recorder.Code != statusCode {
		t.Errorf("Expected status code %d, but got %d", statusCode, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', but got '%s'", contentType)
	}

	var responseBody map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if responseBody["message"] != data["message"] {
		t.Errorf("Expected response message '%s', but got '%s'", data["message"], responseBody["message"])
	}
}

func TestWriteJSONResponse_Error(t *testing.T) {
	recorder := httptest.NewRecorder()

	data := make(chan int)

	statusCode := http.StatusOK

	writeJSONResponse(recorder, data, statusCode)

	if recorder.Code != statusCode {
		t.Errorf("Expected status code %d, but got %d", statusCode, recorder.Code)
	}

	if recorder.Code == http.StatusInternalServerError {
		t.Errorf("Expected JSON serialization error, but got: %d", recorder.Code)
	}
}

func TestWriteJSONError(t *testing.T) {
	recorder := httptest.NewRecorder()

	errorMessage := "test error"
	statusCode := http.StatusBadRequest

	writeJSONError(recorder, errorMessage, statusCode)

	if recorder.Code != statusCode {
		t.Errorf("Expected status code %d, but got %d", statusCode, recorder.Code)
	}

	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', but got '%s'", contentType)
	}

	var responseBody map[string]string
	if err := json.NewDecoder(recorder.Body).Decode(&responseBody); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if responseBody["error"] != errorMessage {
		t.Errorf("Expected error message '%s', but got '%s'", errorMessage, responseBody["error"])
	}
}

func TestGetCachedLoans_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/cache", nil)
	rec := httptest.NewRecorder()

	GetCachedLoans(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, but got %d", http.StatusMethodNotAllowed, rec.Code)
	}

	expectedBody := `{"error":"Only GET method is allowed"}` + "\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, but got %q", expectedBody, rec.Body.String())
	}
}

func TestGetCachedLoans_EmptyCache(t *testing.T) {
	// Очищаем кэш перед тестом
	InitCache(10, time.Hour)

	req := httptest.NewRequest(http.MethodGet, "/cache", nil)
	rec := httptest.NewRecorder()

	GetCachedLoans(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, rec.Code)
	}

	expectedBody := `{"error":"empty cache"}` + "\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, but got %q", expectedBody, rec.Body.String())
	}
}

func TestGetCachedLoans_NonEmptyCache(t *testing.T) {
	// Очищаем кэш перед тестом
	InitCache(10, time.Hour)

	loan1 := models.CachedLoan{ID: 1}
	loan2 := models.CachedLoan{ID: 2}
	loanCache.Store(int64(loan1.ID), loan1)
	loanCache.Store(int64(loan2.ID), loan2)

	req := httptest.NewRequest(http.MethodGet, "/cache", nil)
	rec := httptest.NewRecorder()

	GetCachedLoans(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, rec.Code)
	}

	var loans []models.CachedLoan
	if err := json.Unmarshal(rec.Body.Bytes(), &loans); err != nil {
		t.Fatalf("Failed to decode JSON response: %v", err)
	}

	if len(loans) != 2 {
		t.Fatalf("Expected %d loans, but got %d", 2, len(loans))
	}

	// Проверяем, что все ожидаемые займы присутствуют
	found := make(map[int64]bool)
	for _, loan := range loans {
		found[loan.ID] = true
	}

	if !found[1] || !found[2] {
		t.Error("Not all expected loans were found in the response")
	}
}

func TestExecuteLoanCalculation_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer([]byte(`invalid json`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ExecuteLoanCalculation(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, rec.Code)
	}

	expectedBody := `{"error":"Invalid JSON format"}` + "\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, got %q", expectedBody, rec.Body.String())
	}
}

func TestExecuteLoanCalculation_Success(t *testing.T) {
	// Очищаем кэш перед тестом
	InitCache(10, time.Hour)

	request := models.LoanRequest{
		LoanParams: models.LoanParams{
			ObjectCost:     5000000,
			InitialPayment: 1000000,
			Months:         240,
		},
		Program: models.Program{
			Salary: true,
		},
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ExecuteLoanCalculation(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, rec.Code)
	}

	var response models.LoanResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	// Проверяем основные параметры ответа
	if response.Result.Params.ObjectCost != request.LoanParams.ObjectCost {
		t.Errorf("Expected object cost %d, but got %d", request.LoanParams.ObjectCost, response.Result.Params.ObjectCost)
	}

	if response.Result.Params.InitialPayment != request.LoanParams.InitialPayment {
		t.Errorf("Expected initial payment %d, but got %d", request.LoanParams.InitialPayment, response.Result.Params.InitialPayment)
	}

	if response.Result.Params.Months != request.LoanParams.Months {
		t.Errorf("Expected months %d, but got %d", request.LoanParams.Months, response.Result.Params.Months)
	}

	if response.Result.Program.Salary != request.Program.Salary {
		t.Errorf("Expected salary program %v, but got %v", request.Program.Salary, response.Result.Program.Salary)
	}
}

func TestExecuteLoanCalculation_ReadBodyError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/execute", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ExecuteLoanCalculation(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExecuteLoanCalculation_CalculationError(t *testing.T) {
	request := models.LoanRequest{
		LoanParams: models.LoanParams{
			ObjectCost:     -1000000, // Отрицательная стоимость вызовет ошибку
			InitialPayment: 200000,
			Months:         120,
		},
	}

	body, _ := json.Marshal(request)
	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	ExecuteLoanCalculation(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, rec.Code)
	}
}
