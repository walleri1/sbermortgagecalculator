// Package calculator ...
package calculator

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"sbermortgagecalculator/internal/models"
)

const (
	SalaryRate   = 8
	MilitaryRate = 9
	BaseRate     = 10
)

// Errors for validation
var (
	ErrNoProgramSelected      = errors.New("choose program")
	ErrMultiplePrograms       = errors.New("choose only 1 program")
	ErrInitialPaymentTooLow   = errors.New("the initial payment should be more than or equal to 20% of the object cost")
	ErrMonthsShouldBePositive = errors.New("loan term in months should be a positive number")
	ErrLoanSumZeroOrNegative  = errors.New("loan sum must be greater than zero")
	ErrCalculationError       = errors.New("calculation error: denominator is zero")
)

// Calculate computes the loan parameters (rate, loan amount, monthly payment, overpayment, etc.).
func Calculate(request models.LoanRequest) (*models.Aggregates, error) {
	// Validate input
	if err := validateRequest(request); err != nil {
		return nil, err
	}

	// Determine the interest rate
	rate, err := selectRate(request.Program)
	if err != nil {
		return nil, err
	}

	// Convert inputs to decimal
	objectCost := decimal.NewFromInt(int64(request.Params.ObjectCost))
	initialPayment := decimal.NewFromInt(int64(request.Params.InitialPayment))
	loanSum := objectCost.Sub(initialPayment) // Loan sum = ObjectCost - InitialPayment

	if loanSum.LessThanOrEqual(decimal.Zero) {
		return nil, ErrMonthsShouldBePositive
	}

	// Ensure the number of months is positive
	if request.Params.Months <= 0 {
		return nil, ErrMonthsShouldBePositive
	}
	months := decimal.NewFromInt(int64(request.Params.Months))

	// Monthly interest rate in decimal form: rate / 100 / 12
	rateDecimal := decimal.NewFromInt(int64(rate))
	monthlyRate := rateDecimal.Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))

	// Calculate the monthly payment (annuity formula)
	monthlyPayment, err := CalculateMonthlyPayment(loanSum, monthlyRate, months)
	if err != nil {
		return nil, err
	}

	// Total payment over the loan period
	totalPayment := monthlyPayment.Mul(months)

	// Overpayment is the total payment minus the loan amount
	overpayment := totalPayment.Sub(loanSum)

	// Calculate the last payment date
	lastPaymentDate := time.Now().AddDate(0, int(months.IntPart()), 0).Format("2006-01-02")

	// Return aggregates
	return &models.Aggregates{
		Rate:            int(rate),
		LoanSum:         int(loanSum.IntPart()),
		MonthlyPayment:  int(monthlyPayment.IntPart()),
		Overpayment:     int(overpayment.IntPart()),
		LastPaymentDate: lastPaymentDate,
	}, nil
}

// CalculateMonthlyPayment computes the monthly payment using the annuity formula.
func CalculateMonthlyPayment(loanSum, monthlyRate, months decimal.Decimal) (decimal.Decimal, error) {
	// Formula: P = S * (r * (1 + r)^n) / ((1 + r)^n - 1)
	// Where:
	// S = loanSum, r = monthlyRate, n = months

	if monthlyRate.IsZero() {
		// If the interest rate is 0%, the monthly payment is simply the loan amount divided by months
		return loanSum.Div(months), nil
	}

	one := decimal.NewFromInt(1)

	// (1 + r)^n
	compoundRate := one.Add(monthlyRate).Pow(months)

	// Calculate monthly payment
	numerator := loanSum.Mul(monthlyRate).Mul(compoundRate)
	denominator := compoundRate.Sub(one)

	if denominator.IsZero() {
		return decimal.Zero, ErrCalculationError
	}

	return numerator.Div(denominator), nil
}

// selectRate determines the loan's interest rate based on the selected program.
func selectRate(program models.Program) (int, error) {
	if program.Salary && !program.Military && !program.Base {
		return SalaryRate, nil
	}
	if program.Military && !program.Salary && !program.Base {
		return MilitaryRate, nil
	}
	if program.Base && !program.Salary && !program.Military {
		return BaseRate, nil
	}

	// Validate program selection
	if !program.Salary && !program.Military && !program.Base {
		return 0, ErrNoProgramSelected
	}
	return 0, ErrMultiplePrograms
}

// validateRequest validates the loan request parameters. Ensures initial payment, programs, and loan terms are valid.
func validateRequest(request models.LoanRequest) error {
	minInitialPayment := decimal.NewFromInt(int64(request.Params.ObjectCost)).Mul(decimal.NewFromFloat(0.2))
	initialPayment := decimal.NewFromInt(int64(request.Params.InitialPayment))

	if initialPayment.LessThan(minInitialPayment) {
		return ErrInitialPaymentTooLow
	}
	return nil
}
