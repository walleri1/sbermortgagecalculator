package utils

import (
	"sync"
	"time"

	"sbermortgagecalculator/internal/models"
)

// CacheItem представляет элемент кэша с временем создания
type CacheItem struct {
	Value     models.CachedLoan
	CreatedAt time.Time
}

// Cache представляет кэш с ограничением по размеру и времени жизни
type Cache struct {
	items   map[int64]CacheItem
	mutex   sync.RWMutex
	maxSize int
	ttl     time.Duration
}

// NewCache создает новый экземпляр кэша
func NewCache(maxSize int, ttl time.Duration) *Cache {
	cache := &Cache{
		items:   make(map[int64]CacheItem),
		maxSize: maxSize,
		ttl:     ttl,
	}

	// Запускаем горутину для очистки устаревших элементов
	go cache.cleanup()

	return cache
}

// Store сохраняет значение в кэш
func (c *Cache) Store(key int64, value models.CachedLoan) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Если достигнут максимальный размер, удаляем самый старый элемент
	if len(c.items) >= c.maxSize {
		c.removeOldest()
	}

	c.items[key] = CacheItem{
		Value:     value,
		CreatedAt: time.Now(),
	}
}

// Load загружает значение из кэша
func (c *Cache) Load(key int64) (models.CachedLoan, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, exists := c.items[key]
	if !exists {
		return models.CachedLoan{}, false
	}

	// Проверяем, не устарел ли элемент
	if time.Since(item.CreatedAt) > c.ttl {
		return models.CachedLoan{}, false
	}

	return item.Value, true
}

// Range выполняет функцию для каждого элемента кэша
func (c *Cache) Range(f func(key int64, value models.CachedLoan) bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	for key, item := range c.items {
		if time.Since(item.CreatedAt) <= c.ttl {
			if !f(key, item.Value) {
				break
			}
		}
	}
}

// removeOldest удаляет самый старый элемент из кэша
func (c *Cache) removeOldest() {
	var oldestKey int64
	var oldestTime time.Time

	for key, item := range c.items {
		if oldestTime.IsZero() || item.CreatedAt.Before(oldestTime) {
			oldestKey = key
			oldestTime = item.CreatedAt
		}
	}

	delete(c.items, oldestKey)
}

// cleanup периодически очищает устаревшие элементы
func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.ttl)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		now := time.Now()
		for key, item := range c.items {
			if now.Sub(item.CreatedAt) > c.ttl {
				delete(c.items, key)
			}
		}
		c.mutex.Unlock()
	}
}

// TTL возвращает время жизни кэша
func (c *Cache) TTL() time.Duration {
	return c.ttl
}
