package paths

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sbermortgagecalculator/internal/models"
)

func TestGetCachedResult(t *testing.T) {
	// Инициализируем кэш
	InitCache(10, time.Hour)

	// Добавляем тестовые данные в кэш
	testResponse := models.CachedLoan{
		ID: 1,
		CalculationResult: models.CalculationResult{
			Params: models.LoanParams{
				ObjectCost:     1000000,
				InitialPayment: 200000,
				Months:         120,
			},
			Program: models.Program{
				Salary: true,
			},
			Aggregates: models.Aggregates{
				MonthlyPayment:  10000,
				LoanSum:         1200000,
				Overpayment:     200000,
				Rate:            7,
				LastPaymentDate: "2024-03-06",
			},
		},
		CreatedAt: time.Now(),
	}
	loanCache.Store(1, testResponse)

	tests := []struct {
		name           string
		id             string
		expectedStatus int
	}{
		{
			name:           "Существующий ID",
			id:             "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Несуществующий ID",
			id:             "999",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Невалидный ID",
			id:             "invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/cache/"+tt.id, nil)
			rr := httptest.NewRecorder()

			GetCachedResult(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedStatus, status)
			}
		})
	}
}

func TestGetCachedResult_MethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest("POST", "/cache/1", nil)
	rr := httptest.NewRecorder()

	GetCachedResult(rr, req)

	if status := rr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusMethodNotAllowed, status)
	}
}

func TestGetCachedResult_ExpiredCache(t *testing.T) {
	// Инициализируем кэш с очень коротким TTL
	InitCache(10, time.Millisecond)

	// Добавляем тестовые данные в кэш
	testResponse := models.CachedLoan{
		ID: 2,
		CalculationResult: models.CalculationResult{
			Params: models.LoanParams{
				ObjectCost:     1000000,
				InitialPayment: 200000,
				Months:         120,
			},
			Program: models.Program{
				Salary: true,
			},
			Aggregates: models.Aggregates{
				MonthlyPayment:  10000,
				LoanSum:         1200000,
				Overpayment:     200000,
				Rate:            7,
				LastPaymentDate: "2024-03-06",
			},
		},
	}
	loanCache.Store(2, testResponse)

	// Ждем, пока кэш истечет
	time.Sleep(2 * time.Millisecond)

	req := httptest.NewRequest("GET", "/cache/2", nil)
	rr := httptest.NewRecorder()

	GetCachedResult(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusNotFound, status)
	}
}
