// Package calculator with sets of functions for calculating the credit burden of credit programs
package calculator

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"

	"sbermortgagecalculator/internal/models"
)

const (
	CorporateRate = 8  // corporate program
	MilitaryRate  = 9  // military program
	BaseRate      = 10 // base program
)

// Errors for validation
var (
	ErrNoProgramSelected      = errors.New("choose program")
	ErrMultiplePrograms       = errors.New("choose only 1 program")
	ErrInitialPaymentTooLow   = errors.New("the initial payment should be more than or equal to 20% of the object cost")
	ErrMonthsShouldBePositive = errors.New("loan term in months should be a positive number")
	ErrLoanSumZeroOrNegative  = errors.New("loan sum must be greater than zero")
	ErrCalculationError       = errors.New("undefined behavior: division by zero")
)

// CalculateMortgageAggregates computes the loan parameters (rate, loan amount, monthly payment, overpayment, etc.)
func CalculateMortgageAggregates(request models.LoanRequest) (*models.Aggregates, error) {
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
	loanSum := objectCost.Sub(initialPayment)

	// Making sure that the borrower needs the money
	if loanSum.LessThanOrEqual(decimal.Zero) {
		return nil, ErrLoanSumZeroOrNegative
	}

	// Ensure the number of loanMonths is positive
	loanMonths := decimal.NewFromInt(int64(request.Params.Months))
	if loanMonths.LessThanOrEqual(decimal.Zero) {
		return nil, ErrMonthsShouldBePositive
	}

	// Monthly interest rate in decimal form: rate / 100 / 12
	monthlyRate := decimal.NewFromInt(int64(rate)).Div(decimal.NewFromInt(100)).Div(decimal.NewFromInt(12))

	// Calculate the monthly payment (annuity formula - docs example_golang.xlsx)
	monthlyPayment, err := calculateMonthlyPayment(loanSum, monthlyRate, loanMonths)
	if err != nil {
		return nil, err
	}

	// Total amount of payments for the entire loan period with body and interest
	totalPayment := monthlyPayment.Mul(loanMonths)

	// Interest for using the bank's money
	overpayment := totalPayment.Sub(loanSum)

	// Last payment date
	lastPaymentDate := time.Now().AddDate(0, int(loanMonths.IntPart()), 0).Format("2006-01-02")

	return &models.Aggregates{
		Rate:            int8(rate),
		LoanSum:         loanSum.IntPart(),
		MonthlyPayment:  int32(monthlyPayment.IntPart()),
		Overpayment:     overpayment.IntPart(),
		LastPaymentDate: lastPaymentDate,
	}, nil
}

// calculateMonthlyPayment computes the monthly payment using the annuity formula
func calculateMonthlyPayment(loanSum, monthlyRate, months decimal.Decimal) (decimal.Decimal, error) {
	// Formula: P = S * (G * (1 + G)^T) / ((1 + G)^T - 1)
	// Where:
	// S = loanSum, G = monthlyRate, T = months

	// When there is nothing to borrow
	if loanSum.LessThanOrEqual(decimal.Zero) || months.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, ErrCalculationError
	}

	// If the interest rate is 0% (Muslim mortgage), then the monthly payment is simply the loan amount divided by the number of months
	if monthlyRate.IsZero() {
		return loanSum.Div(months), nil
	}

	// (1 + G)^T
	compoundRate := decimal.NewFromInt(1).Add(monthlyRate).Pow(months)

	// Calculate monthly payment
	numerator := loanSum.Mul(monthlyRate).Mul(compoundRate) //S * (G * (1 + G)^T)
	denominator := compoundRate.Sub(decimal.NewFromInt(1))  //((1 + G)^T - 1)

	// ((1 + G)^T - 1) == 0
	if denominator.IsZero() {
		return decimal.Zero, ErrCalculationError
	}

	return numerator.Div(denominator), nil
}

// selectRate determines the loan's interest rate based on the selected program
func selectRate(program models.Program) (int, error) {
	if program.Salary && !program.Military && !program.Base {
		return CorporateRate, nil
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
