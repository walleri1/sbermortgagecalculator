package paths

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sbermortgagecalculator/internal/models"
	"sync"
	"testing"
)

func Test_GetCachedLoans_EmptyCache(t *testing.T) {
	loanCache = sync.Map{}

	req := httptest.NewRequest(http.MethodGet, "/cache", nil)
	w := httptest.NewRecorder()

	GetCachedLoans(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unable to parse response: %v", err)
	}
	if resp["error"] != "empty cache" {
		t.Errorf("expected error 'empty cache', got %v", resp["error"])
	}
}

func Test_GetCachedLoans_WithData(t *testing.T) {
	loanCache = sync.Map{}
	cachedLoan := models.CachedLoan{
		ID: 1,
		CalculationResult: models.CalculationResult{
			Aggregates: models.Aggregates{
				LastPaymentDate: "2024-12-31",
				LoanSum:         1000000,
				Overpayment:     200000,
				MonthlyPayment:  80000,
				Rate:            7,
			},
		},
	}
	loanCache.Store(1, cachedLoan)

	req := httptest.NewRequest(http.MethodGet, "/cache", nil)
	w := httptest.NewRecorder()

	GetCachedLoans(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
	}

	var loans []models.CachedLoan
	if err := json.Unmarshal(w.Body.Bytes(), &loans); err != nil {
		t.Fatalf("unable to parse response: %v", err)
	}
	if len(loans) != 1 {
		t.Errorf("expected 1 loan, got %d", len(loans))
	}
	if loans[0].ID != 1 {
		t.Errorf("expected loan ID 1, got %d", loans[0].ID)
	}
}

func Test_ExecuteLoanCalculation_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/execute", nil)
	w := httptest.NewRecorder()

	ExecuteLoanCalculation(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unable to parse response: %v", err)
	}
	if resp["error"] != "Only POST method is allowed" {
		t.Errorf("expected error 'Only POST method is allowed', got %v", resp["error"])
	}
}

func Test_GetCachedLoans_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/cache", nil)
	w := httptest.NewRecorder()

	GetCachedLoans(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected status code %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}

	var resp map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	if err != nil {
		t.Fatalf("unable to parse response body as JSON: %v", err)
	}

	expectedErrorMessage := "Only GET method is allowed"
	if resp["error"] != expectedErrorMessage {
		t.Errorf("expected error message '%s', got '%s'", expectedErrorMessage, resp["error"])
	}
}

func Test_ExecuteLoanCalculation_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewReader([]byte("invalid json")))
	w := httptest.NewRecorder()

	ExecuteLoanCalculation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unable to parse response: %v", err)
	}
	if resp["error"] != "Invalid JSON format" {
		t.Errorf("expected error 'Invalid JSON format', got %v", resp["error"])
	}
}

func Test_ExecuteLoanCalculation_Success(t *testing.T) {
	request := models.LoanRequest{
		LoanParams: models.LoanParams{
			ObjectCost:     2000000,
			InitialPayment: 300000,
			Months:         240,
		},
		Program: models.Program{
			Salary: true,
		},
	}
	body, _ := json.Marshal(request)

	req := httptest.NewRequest(http.MethodPost, "/execute", bytes.NewReader(body))
	w := httptest.NewRecorder()

	ExecuteLoanCalculation(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var resp models.LoanResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unable to parse response: %v", err)
	}
}
