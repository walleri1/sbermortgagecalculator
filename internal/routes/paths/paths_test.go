package paths

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"sbermortgagecalculator/internal/models"
)

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
	loanCache = sync.Map{}

	req := httptest.NewRequest(http.MethodGet, "/cache", nil)
	rec := httptest.NewRecorder()

	GetCachedLoans(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, but got %d", http.StatusNotFound, rec.Code)
	}

	expectedBody := `{"error":"empty cache"}` + "\n"
	if rec.Body.String() != expectedBody {
		t.Errorf("Expected body %q, but got %q", expectedBody, rec.Body.String())
	}
}

func TestGetCachedLoans_NonEmptyCache(t *testing.T) {
	loanCache = sync.Map{}
	loan1 := models.CachedLoan{ID: 0}
	loan2 := models.CachedLoan{ID: 1}
	loanCache.Store(loan1.ID, loan1)
	loanCache.Store(loan2.ID, loan2)

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

	expectedLoans := []models.CachedLoan{loan1, loan2}
	if len(loans) != len(expectedLoans) {
		t.Fatalf("Expected %d loans, but got %d", len(expectedLoans), len(loans))
	}
	for i, loan := range loans {
		if loan != expectedLoans[i] {
			t.Errorf("Expected loan %v, but got %v", expectedLoans[i], loan)
		}
	}
}

func TestExecuteLoanCalculation_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/execute", nil)
	rec := httptest.NewRecorder()

	ExecuteLoanCalculation(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, but got %d", http.StatusMethodNotAllowed, rec.Code)
	}

	expected := `{"error":"Only POST method is allowed"}` + "\n"
	if rec.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestExecuteLoanCalculation_ReadBodyError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/execute", nil)
	rec := httptest.NewRecorder()

	ExecuteLoanCalculation(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExecuteLoanCalculation_InvalidJSON(t *testing.T) {
	body := bytes.NewBufferString("invalid JSON")
	req := httptest.NewRequest(http.MethodPost, "/execute", body)
	rec := httptest.NewRecorder()

	ExecuteLoanCalculation(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, rec.Code)
	}

	expected := `{"error":"Invalid JSON format"}` + "\n"
	if rec.Body.String() != expected {
		t.Errorf("Expected body %q, got %q", expected, rec.Body.String())
	}
}

func TestExecuteLoanCalculation_CalculationError(t *testing.T) {
	request := models.LoanRequest{
		LoanParams: models.LoanParams{},
	}
	body, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	ExecuteLoanCalculation(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, but got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestExecuteLoanCalculation_Success(t *testing.T) {
	request := models.LoanRequest{
		LoanParams: models.LoanParams{
			ObjectCost:     5000000,
			InitialPayment: 1000000,
			Months:         240,
		},
		Program: models.Program{Salary: true},
	}
	body, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewReader(body))
	rec := httptest.NewRecorder()
	ExecuteLoanCalculation(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, but got %d", http.StatusOK, rec.Code)
	}

	var response models.LoanResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse JSON response: %v", err)
	}

	if response.Result.Aggregates.LoanSum != 4000000 {
		t.Errorf("Expected loan amount 4000000, but got %d", response.Result.Aggregates.LoanSum)
	}

	if response.Result.Aggregates.Rate != 8 {
		t.Errorf("Expected rate amount 8, but got %d", response.Result.Aggregates.Rate)
	}

}
