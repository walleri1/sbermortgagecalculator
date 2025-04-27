package paths

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sbermortgagecalculator/internal/models"
)

func TestExecuteLoanCalculation(t *testing.T) {
	// Инициализируем кэш
	InitCache(10, time.Hour)

	tests := []struct {
		name           string
		request        models.LoanRequest
		expectedStatus int
		checkCache     bool
	}{
		{
			name: "Валидный запрос",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     1000000,
					InitialPayment: 200000,
					Months:         120,
				},
				Program: models.Program{
					Salary: true,
				},
			},
			expectedStatus: http.StatusOK,
			checkCache:     true,
		},
		{
			name: "Невалидный запрос - отрицательная стоимость",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     -1000000,
					InitialPayment: 200000,
					Months:         120,
				},
				Program: models.Program{
					Salary: true,
				},
			},
			expectedStatus: http.StatusBadRequest,
			checkCache:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем тестовый запрос
			body, _ := json.Marshal(tt.request)
			req := httptest.NewRequest("POST", "/execute", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			// Создаем ResponseWriter
			rr := httptest.NewRecorder()

			// Вызываем обработчик
			ExecuteLoanCalculation(rr, req)

			// Проверяем статус ответа
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedStatus, status)
			}

			// Если запрос успешный, проверяем кэш
			if tt.checkCache && tt.expectedStatus == http.StatusOK {
				var response models.LoanResponse
				if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
					t.Fatalf("Ошибка разбора ответа: %v", err)
				}

				// Проверяем, что результат сохранен в кэше
				if cached, ok := loanCache.Load(0); !ok {
					t.Error("Ожидалось, что результат будет в кэше")
				} else if cached.ID != 0 {
					t.Errorf("Ожидался ID 0, получен %d", cached.ID)
				}
			}
		})
	}
}

func TestExecuteLoanCalculation_InvalidContentType(t *testing.T) {
	req := httptest.NewRequest("POST", "/execute", nil)
	req.Header.Set("Content-Type", "text/plain")
	rr := httptest.NewRecorder()

	ExecuteLoanCalculation(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("Ожидался статус %d, получен %d", http.StatusBadRequest, status)
	}
}
