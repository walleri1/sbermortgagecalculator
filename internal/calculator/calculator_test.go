package calculator

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"

	"sbermortgagecalculator/internal/models"
)

func TestCalculateMortgageAggregates(t *testing.T) {
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
				LoanParams: models.LoanParams{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
				},
				Program: models.Program{Salary: true},
			},
			expectedRate:    CorporateRate,
			expectedLoan:    4000000,
			expectedPayment: 33457,
			expectErr:       nil,
		},
		{
			name: "Valid military program",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     3000000,
					InitialPayment: 600000,
					Months:         180,
				},
				Program: models.Program{Military: true},
			},
			expectedRate:    MilitaryRate,
			expectedLoan:    2400000,
			expectedPayment: 24342,
			expectErr:       nil,
		},
		{
			name: "Valid base program",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     3000000,
					InitialPayment: 600000,
					Months:         180,
				},
				Program: models.Program{Base: true},
			},
			expectedRate:    BaseRate,
			expectedLoan:    2400000,
			expectedPayment: 25790,
			expectErr:       nil,
		},
		{
			name: "Invalid no program selection",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
				},
				Program: models.Program{},
			},
			expectErr: ErrNoProgramSelected,
		},
		{
			name: "Invalid two program selection: salary, base",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
				},
				Program: models.Program{
					Salary: true,
					Base:   true,
				},
			},
			expectErr: ErrMultiplePrograms,
		},
		{
			name: "Invalid two program selection: salary, military",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
				},
				Program: models.Program{
					Salary:   true,
					Military: true,
				},
			},
			expectErr: ErrMultiplePrograms,
		},
		{
			name: "Invalid two program selection: base, military",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
				},
				Program: models.Program{
					Base:     true,
					Military: true,
				},
			},
			expectErr: ErrMultiplePrograms,
		},
		{
			name: "Low initial payment",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     6000000,
					InitialPayment: 500000,
					Months:         240,
				},
				Program: models.Program{Salary: true},
			},
			expectErr: ErrInitialPaymentTooLow,
		},
		{
			name: "Lack of loan amount",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     6000000,
					InitialPayment: 6000000,
					Months:         240,
				},
				Program: models.Program{Salary: true},
			},
			expectErr: ErrLoanSumZeroOrNegative,
		},
		{
			name: "Zero Months",
			request: models.LoanRequest{
				LoanParams: models.LoanParams{
					ObjectCost:     4000000,
					InitialPayment: 800000,
					Months:         0,
				},
				Program: models.Program{Base: true},
			},
			expectErr: ErrMonthsShouldBePositive,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CalculateMortgageAggregates(tc.request)
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

func TestCalculateMortgageAggregatesCached(t *testing.T) {
	request := models.LoanRequest{
		LoanParams: models.LoanParams{
			ObjectCost:     5000000,
			InitialPayment: 1000000,
			Months:         240,
		},
		Program: models.Program{Salary: true},
	}
	resultFirst, err := CalculateMortgageAggregates(request)
	assert.NoError(t, err)
	resutlSecond, err := CalculateMortgageAggregates(request)
	assert.NoError(t, err)
	assert.Equal(t, resultFirst, resutlSecond)
}

func TestCalculateMonthlyPayment(t *testing.T) {
	tests := []struct {
		name           string
		loanSum        decimal.Decimal
		monthlyRate    decimal.Decimal
		months         decimal.Decimal
		expectedResult decimal.Decimal
		expectErr      error
	}{
		{
			name:           "Standard case with 8% annual interest rate",
			loanSum:        decimal.NewFromInt(4000000),
			monthlyRate:    decimal.NewFromInt32(CorporateRate).Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12)),
			months:         decimal.NewFromInt(20).Mul(decimal.NewFromInt(12)),
			expectedResult: decimal.RequireFromString("33457.6"),
			expectErr:      nil,
		},
		{
			name:           "Zero interest rate",
			loanSum:        decimal.NewFromInt(4000000),
			monthlyRate:    decimal.Zero,
			months:         decimal.NewFromInt(20).Mul(decimal.NewFromInt(12)),
			expectedResult: decimal.NewFromFloat(16666.66).Round(2),
			expectErr:      nil,
		},
		{
			name:           "1 month loan with 12% annual interest rate",
			loanSum:        decimal.NewFromInt(1000000),
			monthlyRate:    decimal.NewFromFloat(0.01),
			months:         decimal.NewFromInt(1),
			expectedResult: decimal.RequireFromString("1010000.00"),
			expectErr:      nil,
		},
		{
			name:           "Zero months (invalid input)",
			loanSum:        decimal.NewFromInt(4000000),
			monthlyRate:    decimal.NewFromInt32(CorporateRate).Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12)),
			months:         decimal.Zero,
			expectedResult: decimal.Zero,
			expectErr:      ErrCalculationError,
		},
		{
			name:           "Edge case with 0% and 0 months",
			loanSum:        decimal.NewFromInt(4000000),
			monthlyRate:    decimal.Zero,
			months:         decimal.NewFromInt(0),
			expectedResult: decimal.Zero,
			expectErr:      ErrCalculationError,
		},
		{
			name:           "Negative loan amount",
			loanSum:        decimal.NewFromInt(-1000000),
			monthlyRate:    decimal.NewFromInt(1).Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12)),
			months:         decimal.NewFromInt(12),
			expectedResult: decimal.Zero,
			expectErr:      ErrCalculationError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := calculateMonthlyPayment(tc.loanSum, tc.monthlyRate, tc.months)

			if tc.expectErr != nil {
				assert.ErrorIs(t, err, tc.expectErr, "Expected error did not occur")
				return
			}

			assert.NoError(t, err, "Unexpected error occurred")

			delta := decimal.NewFromFloat(0.01)
			diff := result.Sub(tc.expectedResult).Abs()
			assert.True(t, diff.LessThanOrEqual(delta),
				"Expected result: %s, got: %s (diff: %s)", tc.expectedResult.String(), result.String(), diff.String())
		})
	}
}
