package utils

import (
	"encoding/json"
	"net/http"

	"sbermortgagecalculator/internal/models"
)

// ErrorResponse представляет структуру ответа с ошибкой
type ErrorResponse struct {
	Error string `json:"error"`
}

// WriteJSONError записывает JSON-ответ с ошибкой
func WriteJSONError(w http.ResponseWriter, errorMessage string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errorMessage})
}

// ValidationError представляет ошибку валидации
type ValidationError struct {
	Field   string
	Message string
}

// ValidateLoanRequest проверяет корректность запроса на кредит
func ValidateLoanRequest(request models.LoanRequest) []ValidationError {
	var errors []ValidationError

	// Проверка стоимости объекта
	if request.ObjectCost <= 0 {
		errors = append(errors, ValidationError{
			Field:   "object_cost",
			Message: "стоимость объекта должна быть положительным числом",
		})
	}

	// Проверка первоначального взноса
	if request.InitialPayment < 0 {
		errors = append(errors, ValidationError{
			Field:   "initial_payment",
			Message: "первоначальный взнос не может быть отрицательным",
		})
	}

	// Проверка срока кредита
	if request.Months <= 0 {
		errors = append(errors, ValidationError{
			Field:   "months",
			Message: "срок кредита должен быть положительным числом",
		})
	}

	// Проверка программы кредитования
	programsSelected := 0
	if request.Program.Salary {
		programsSelected++
	}
	if request.Program.Military {
		programsSelected++
	}
	if request.Program.Base {
		programsSelected++
	}

	if programsSelected == 0 {
		errors = append(errors, ValidationError{
			Field:   "program",
			Message: "необходимо выбрать программу кредитования",
		})
	} else if programsSelected > 1 {
		errors = append(errors, ValidationError{
			Field:   "program",
			Message: "можно выбрать только одну программу кредитования",
		})
	}

	// Проверка первоначального взноса (20% от стоимости)
	if request.ObjectCost > 0 && request.InitialPayment > 0 {
		minInitialPayment := float64(request.ObjectCost) * 0.2
		if float64(request.InitialPayment) < minInitialPayment {
			errors = append(errors, ValidationError{
				Field:   "initial_payment",
				Message: "первоначальный взнос должен быть не менее 20% от стоимости объекта",
			})
		}
	}

	return errors
}
