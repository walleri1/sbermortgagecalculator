// Package cache implements an in-memory caching layer.
package cache

import (
	"sort"
	"sync"

	model "sbermortgagecalculator/internal/models"
)

// CachedLoanStore is a thread-safe storage for CachedLoan
type CachedLoanStore struct {
	mu    sync.RWMutex
	store map[int64]model.CachedLoan
}

// NewCachedLoanStore creates a new storage for CachedLoan
func NewCachedLoanStore() *CachedLoanStore {
	return &CachedLoanStore{
		store: make(map[int64]model.CachedLoan),
	}
}

// Add adds a new CachedLoan to the repository
func (s *CachedLoanStore) Add(cachedLoan model.CachedLoan) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.store[cachedLoan.ID] = cachedLoan
}

// Get retrieves CachedLoan from storage by ID
func (s *CachedLoanStore) Get(id int64) (model.CachedLoan, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	cachedLoan, exists := s.store[id]
	return cachedLoan, exists
}

// Remove removes CachedLoan from storage by ID
func (s *CachedLoanStore) Remove(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.store, id)
}

// Exists checks whether a CachedLoan exists with the given ID
func (s *CachedLoanStore) Exists(id int64) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.store[id]
	return exists
}

// GetSortedLoans returns all loans sorted by ID
func (s *CachedLoanStore) GetSortedLoans() []model.CachedLoan {
	s.mu.RLock()
	defer s.mu.RUnlock()

	loans := make([]model.CachedLoan, 0, len(s.store))
	for _, loan := range s.store {
		loans = append(loans, loan)
	}

	sort.Slice(loans, func(i, j int) bool {
		return loans[i].ID < loans[j].ID
	})

	return loans
}
