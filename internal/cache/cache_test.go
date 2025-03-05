package cache

import (
	"reflect"
	"testing"

	"sbermortgagecalculator/internal/models"
)

func TestCachedLoanStore_AddAndGet(t *testing.T) {
	store := NewCachedLoanStore()

	loan := models.CachedLoan{
		ID: 1,
		CalculationResult: models.CalculationResult{
			Params: models.LoanParams{
				ObjectCost:     5000000,
				InitialPayment: 1000000,
				Months:         240,
			},
			Program: models.Program{
				Salary: true,
				Base:   true,
			},
			Aggregates: models.Aggregates{
				Rate:           8,
				LoanSum:        4000000,
				MonthlyPayment: 30000,
				Overpayment:    1200000,
			},
		},
	}

	store.Add(loan)

	retrievedLoan, exists := store.Get(loan.ID)
	if !exists {
		t.Fatalf("Failed to retrieve loan with ID %d", loan.ID)
	}

	if !reflect.DeepEqual(loan, retrievedLoan) {
		t.Errorf("Retrieved loan does not match:\nExpected: %+v\nGot: %+v", loan, retrievedLoan)
	}
}

func TestCachedLoanStore_Remove(t *testing.T) {
	store := NewCachedLoanStore()

	loan := models.CachedLoan{
		ID: 1,
		CalculationResult: models.CalculationResult{
			Params: models.LoanParams{
				ObjectCost:     5000000,
				InitialPayment: 1000000,
				Months:         240,
			},
		},
	}

	store.Add(loan)
	store.Remove(loan.ID)

	_, exists := store.Get(loan.ID)
	if exists {
		t.Fatalf("Loan with ID %d was not removed from the store", loan.ID)
	}
}

func TestCachedLoanStore_Exists(t *testing.T) {
	store := NewCachedLoanStore()

	loan := models.CachedLoan{
		ID: 1,
	}

	store.Add(loan)
	if !store.Exists(loan.ID) {
		t.Fatalf("Loan with ID %d should exist in the store", loan.ID)
	}

	store.Remove(loan.ID)
	if store.Exists(loan.ID) {
		t.Fatalf("Loan with ID %d should not exist in the store", loan.ID)
	}
}

func TestCachedLoanStore_GetSortedLoans(t *testing.T) {
	store := NewCachedLoanStore()

	loans := []models.CachedLoan{
		{ID: 3},
		{ID: 1},
		{ID: 2},
	}

	for _, loan := range loans {
		store.Add(loan)
	}

	expectedOrder := []int{1, 2, 3}
	sortedLoans := store.GetSortedLoans()

	if len(sortedLoans) != len(loans) {
		t.Fatalf("Expected %d loans, but got %d", len(loans), len(sortedLoans))
	}

	for i, loan := range sortedLoans {
		if loan.ID != expectedOrder[i] {
			t.Errorf("Expected loan ID %d at index %d, but got %d", expectedOrder[i], i, loan.ID)
		}
	}
}

func TestCachedLoanStore_Concurrency(t *testing.T) {
	store := NewCachedLoanStore()

	const numGoroutines = 100

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			store.Add(models.CachedLoan{ID: i})
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			_, _ = store.Get(i)
		}(i)
	}

	for i := 0; i < numGoroutines; i++ {
		go func(i int) {
			store.Remove(i)
		}(i)
	}
}
