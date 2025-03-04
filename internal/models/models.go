// Package models contains data structures and database interaction logic.
package models

// LoanParams stores the user's request parameters
type LoanParams struct {
	ObjectCost     int64 `json:"object_cost"`     // Cost object
	InitialPayment int64 `json:"initial_payment"` // Initial payment
	Months         int16 `json:"months"`          // Loan term in months
}

// Program describes the selected loan program
type Program struct {
	Salary   bool `json:"salary,omitempty"`   // Corporate program
	Military bool `json:"military,omitempty"` // Military program
	Base     bool `json:"base,omitempty"`     // Base program
}

// Aggregates describes the results of loan calculations
type Aggregates struct {
	LoanSum         int64  `json:"loan_sum"`          // Credit amount
	Overpayment     int64  `json:"overpayment"`       // Overpayment for the entire period
	MonthlyPayment  int32  `json:"monthly_payment"`   // Monthly payment
	Rate            int8   `json:"rate"`              // Annual interest rate
	LastPaymentDate string `json:"last_payment_date"` // Last payment date
}

// LoanRequest is a structure representing a JSON request
type LoanRequest struct {
	Params  LoanParams `json:"params"`
	Program Program    `json:"program"`
}

// CalculationResult combines a query and a calculation result
type CalculationResult struct {
	Params     LoanParams `json:"params"`
	Program    Program    `json:"program"`
	Aggregates Aggregates `json:"aggregates"`
}

// LoanResponse structure for the response
type LoanResponse struct {
	Result CalculationResult `json:"result"`
}

// CachedLoan is a structure for storing data in a cache
type CachedLoan struct {
	ID     int64             `json:"id"`
	Result CalculationResult `json:"result"`
}
