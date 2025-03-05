// Package models contains data structures and database interaction logic.
package models

// LoanParams stores the user's request parameters.
type LoanParams struct {
	ObjectCost     int `json:"object_cost"`     // Cost object.
	InitialPayment int `json:"initial_payment"` // Initial payment.
	Months         int `json:"months"`          // Loan term in months.
}

// Program describes the selected loan program.
type Program struct {
	Salary   bool `json:"salary,omitempty"`   // Corporate program.
	Military bool `json:"military,omitempty"` // Military program.
	Base     bool `json:"base,omitempty"`     // Base program.
}

// Aggregates describes the results of loan calculations.
type Aggregates struct {
	LastPaymentDate string `json:"last_payment_date"` // Last payment dates.
	LoanSum         int    `json:"loan_sum"`          // Credit amount.
	Overpayment     int    `json:"overpayment"`       // Overpayment for the entire period.
	MonthlyPayment  int    `json:"monthly_payment"`   // Monthly payment.
	Rate            int    `json:"rate"`              // Annual interest rate.
}

// LoanRequest is a structure representing a JSON request.
type LoanRequest struct {
	Params  LoanParams `json:"params"`
	Program Program    `json:"program"`
}

// CalculationResult combines a query and a calculation result.
type CalculationResult struct {
	Aggregates Aggregates `json:"aggregates"`
	Params     LoanParams `json:"params"`
	Program    Program    `json:"program"`
}

// LoanResponse structure for the response.
type LoanResponse struct {
	Result CalculationResult `json:"result"`
}

// CachedLoan is a structure for storing data in a cache.
type CachedLoan struct {
	Result CalculationResult `json:"result"`
	ID     int               `json:"id"`
}
