package utils

import (
	"sync"
	"testing"
	"time"

	"sbermortgagecalculator/internal/models"
)

func TestCache_StoreAndLoad(t *testing.T) {
	cache := NewCache(10, time.Hour)
	loan := models.CachedLoan{
		ID: 1,
		CalculationResult: models.CalculationResult{
			Params: models.LoanParams{
				ObjectCost:     1000000,
				InitialPayment: 200000,
				Months:         120,
			},
			Program: models.Program{
				Salary: true,
			},
			Aggregates: models.Aggregates{
				Rate:            8,
				LoanSum:         800000,
				MonthlyPayment:  10000,
				Overpayment:     400000,
				LastPaymentDate: "2033-01-01",
			},
		},
	}

	// Тест сохранения и загрузки
	cache.Store(1, loan)
	if loaded, ok := cache.Load(1); !ok {
		t.Error("Ожидалось, что элемент будет найден в кэше")
	} else if loaded.ID != loan.ID {
		t.Errorf("Ожидался ID %d, получен %d", loan.ID, loaded.ID)
	}

	// Тест загрузки несуществующего элемента
	if _, ok := cache.Load(2); ok {
		t.Error("Ожидалось, что элемент не будет найден в кэше")
	}
}

func TestCache_MaxSize(t *testing.T) {
	cache := NewCache(2, time.Hour)
	loan := models.CachedLoan{ID: 1}

	// Заполняем кэш до предела
	cache.Store(1, loan)
	cache.Store(2, loan)

	// Добавляем еще один элемент
	cache.Store(3, loan)

	// Проверяем, что первый элемент был удален
	if _, ok := cache.Load(1); ok {
		t.Error("Ожидалось, что самый старый элемент будет удален")
	}

	// Проверяем, что новые элементы на месте
	if _, ok := cache.Load(2); !ok {
		t.Error("Ожидалось, что элемент 2 будет в кэше")
	}
	if _, ok := cache.Load(3); !ok {
		t.Error("Ожидалось, что элемент 3 будет в кэше")
	}
}

func TestCache_TTL(t *testing.T) {
	cache := NewCache(10, 100*time.Millisecond)
	loan := models.CachedLoan{ID: 1}

	cache.Store(1, loan)

	// Проверяем, что элемент доступен сразу после сохранения
	if _, ok := cache.Load(1); !ok {
		t.Error("Ожидалось, что элемент будет доступен сразу после сохранения")
	}

	// Ждем, пока элемент устареет
	time.Sleep(200 * time.Millisecond)

	// Проверяем, что элемент больше не доступен
	if _, ok := cache.Load(1); ok {
		t.Error("Ожидалось, что элемент будет удален после истечения TTL")
	}
}

func TestCache_Range(t *testing.T) {
	cache := NewCache(10, time.Hour)
	expectedCount := 3

	// Добавляем несколько элементов
	for i := int64(0); i < int64(expectedCount); i++ {
		value := models.CachedLoan{
			ID: i,
			CalculationResult: models.CalculationResult{
				Params: models.LoanParams{
					ObjectCost:     1000000,
					InitialPayment: 200000,
					Months:         120,
				},
			},
		}
		cache.Store(i, value)
	}

	// Подсчитываем элементы через Range
	count := 0
	cache.Range(func(key int64, value models.CachedLoan) bool {
		count++
		return true
	})

	if count != expectedCount {
		t.Errorf("Expected %d elements, but got %d", expectedCount, count)
	}
}

func TestCache_Store(t *testing.T) {
	cache := NewCache(10, time.Hour)
	key := int64(1)
	value := models.CachedLoan{
		ID: key,
		CalculationResult: models.CalculationResult{
			Params: models.LoanParams{
				ObjectCost:     1000000,
				InitialPayment: 200000,
				Months:         120,
			},
			Program: models.Program{
				Salary: true,
			},
			Aggregates: models.Aggregates{
				Rate:            8,
				LoanSum:         800000,
				MonthlyPayment:  10000,
				Overpayment:     400000,
				LastPaymentDate: "2033-01-01",
			},
		},
	}

	cache.Store(key, value)

	if val, ok := cache.Load(key); !ok || val.ID != value.ID {
		t.Errorf("Expected value %v for key %v, but got %v", value, key, val)
	}
}

func TestCache_Load_NotFound(t *testing.T) {
	cache := NewCache(10, time.Hour)
	key := int64(999)

	if _, ok := cache.Load(key); ok {
		t.Error("Expected key not found, but it was found")
	}
}

func TestCache_Load_Expired(t *testing.T) {
	cache := NewCache(10, time.Millisecond)
	key := int64(1)
	value := models.CachedLoan{
		ID: key,
		CalculationResult: models.CalculationResult{
			Params: models.LoanParams{
				ObjectCost:     1000000,
				InitialPayment: 200000,
				Months:         120,
			},
		},
	}

	cache.Store(key, value)
	time.Sleep(2 * time.Millisecond)

	if _, ok := cache.Load(key); ok {
		t.Error("Expected expired key not found, but it was found")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	cache := NewCache(100, time.Hour)
	var wg sync.WaitGroup
	numGoroutines := 10
	numOperations := 100

	// Запускаем горутины для записи
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := int64(id*100 + j)
				value := models.CachedLoan{
					ID: key,
					CalculationResult: models.CalculationResult{
						Params: models.LoanParams{
							ObjectCost:     1000000,
							InitialPayment: 200000,
							Months:         120,
						},
					},
				}
				cache.Store(key, value)
			}
		}(i)
	}

	// Запускаем горутины для чтения
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				key := int64(id*100 + j)
				cache.Load(key)
			}
		}(i)
	}

	wg.Wait()
}

func TestCache_Size(t *testing.T) {
	maxSize := 2
	cache := NewCache(maxSize, time.Hour)

	// Добавляем элементы до максимального размера и один сверх
	for i := 0; i <= maxSize; i++ {
		key := int64(i)
		value := models.CachedLoan{
			ID: key,
			CalculationResult: models.CalculationResult{
				Params: models.LoanParams{
					ObjectCost:     1000000,
					InitialPayment: 200000,
					Months:         120,
				},
			},
		}
		cache.Store(key, value)
	}

	// Проверяем, что некоторые старые значения были удалены
	found := 0
	for i := 0; i <= maxSize; i++ {
		if _, ok := cache.Load(int64(i)); ok {
			found++
		}
	}

	if found > maxSize {
		t.Errorf("Cache contains %d elements, which exceeds max size %d", found, maxSize)
	}
}

func TestCache_Expiration(t *testing.T) {
	cache := NewCache(10, 50*time.Millisecond)
	key := int64(1)
	value := models.CachedLoan{
		ID: key,
		CalculationResult: models.CalculationResult{
			Params: models.LoanParams{
				ObjectCost:     1000000,
				InitialPayment: 200000,
				Months:         120,
			},
		},
	}

	cache.Store(key, value)

	// Проверяем, что значение доступно сразу после сохранения
	if _, ok := cache.Load(key); !ok {
		t.Error("Значение должно быть доступно сразу после сохранения")
	}

	// Ждем, пока значение устареет
	time.Sleep(100 * time.Millisecond)

	// Проверяем, что значение больше недоступно
	if _, ok := cache.Load(key); ok {
		t.Error("Значение должно быть удалено после истечения TTL")
	}
}
