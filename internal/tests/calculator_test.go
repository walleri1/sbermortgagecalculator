package calculator

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"sbermortgagecalculator/internal/calculator"
	"sbermortgagecalculator/internal/models"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name            string
		request         models.LoanRequest
		expectedRate    int
		expectedLoan    int
		expectedPayment int
		expectErr       error
	}{
		{
			name: "Valid corporate program",
			request: models.LoanRequest{
				Params: models.LoanParams{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
				},
				Program: models.Program{Salary: true},
			},
			expectedRate:    8,
			expectedLoan:    4000000,
			expectedPayment: 33458, // примерный платеж
			expectErr:       nil,
		},
		{
			name: "Valid military program",
			request: models.LoanRequest{
				Params: models.LoanParams{
					ObjectCost:     3000000,
					InitialPayment: 600000,
					Months:         180,
				},
				Program: models.Program{Military: true},
			},
			expectedRate:    9,
			expectedLoan:    2400000,
			expectedPayment: 24933, // примерный платеж
			expectErr:       nil,
		},
		{
			name: "Invalid program selection",
			request: models.LoanRequest{
				Params: models.LoanParams{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
				},
				Program: models.Program{}, // не выбрана программа
			},
			expectErr: calculator.ErrNoProgramSelected,
		},
		{
			name: "Low initial payment",
			request: models.LoanRequest{
				Params: models.LoanParams{
					ObjectCost:     6000000,
					InitialPayment: 500000, // меньше 20% от стоимости дома
					Months:         240,
				},
				Program: models.Program{Salary: true},
			},
			expectErr: calculator.ErrInitialPaymentTooLow,
		},
		{
			name: "Zero Months",
			request: models.LoanRequest{
				Params: models.LoanParams{
					ObjectCost:     4000000,
					InitialPayment: 800000,
					Months:         0, // количество месяцев равно 0
				},
				Program: models.Program{Base: true},
			},
			expectErr: calculator.ErrMonthsShouldBePositive,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calculator.Calculate(tc.request)
			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedRate, result.Rate)
			assert.Equal(t, tc.expectedLoan, result.LoanSum)
			assert.Equal(t, tc.expectedPayment, result.MonthlyPayment)
		})
	}
}

func TestCalculateMonthlyPayment(t *testing.T) {
	tests := []struct {
		name            string
		loanSum         decimal.Decimal
		monthlyRate     decimal.Decimal
		months          decimal.Decimal
		expectedPayment decimal.Decimal
		expectErr       error
	}{
		{
			name:            "Standard case",
			loanSum:         decimal.NewFromInt(4000000),
			monthlyRate:     decimal.NewFromFloat(0.006666), // 8% annually
			months:          decimal.NewFromInt(240),
			expectedPayment: decimal.NewFromInt(33458),
			expectErr:       nil,
		},
		{
			name:            "Zero interest rate",
			loanSum:         decimal.NewFromInt(4000000),
			monthlyRate:     decimal.Zero, // No interest rate
			months:          decimal.NewFromInt(240),
			expectedPayment: decimal.NewFromInt(16666),
			expectErr:       nil,
		},
		{
			name:            "Zero months",
			loanSum:         decimal.NewFromInt(4000000),
			monthlyRate:     decimal.NewFromFloat(0.006666),
			months:          decimal.Zero, // Invalid months
			expectedPayment: decimal.Zero,
			expectErr:       calculator.ErrCalculationError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payment, err := calculator.CalculateMonthlyPayment(tc.loanSum, tc.monthlyRate, tc.months)

			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr)
				return
			}

			// Проверяем отсутствие ошибок
			assert.NoError(t, err)

			// Преобразование в float64
			expectedPaymentFloat, ok1 := tc.expectedPayment.Float64()
			assert.True(t, ok1, "Expected payment conversion to float64 failed")

			paymentFloat, ok2 := payment.Float64()
			assert.True(t, ok2, "Payment conversion to float64 failed")

			// Сравнение с допустимой погрешностью
			assert.InEpsilon(t, expectedPaymentFloat, paymentFloat, 0.01, "Expected %f, got %f", expectedPaymentFloat, paymentFloat)
		})
	}
}
