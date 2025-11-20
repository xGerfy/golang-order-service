package repository_test

import (
	"fmt"
	"testing"

	"order-service/internal/models"
	"order-service/internal/repository"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		cache := repository.NewCache(2)

		order1 := &models.Order{OrderUID: "order1"}
		order2 := &models.Order{OrderUID: "order2"}
		order3 := &models.Order{OrderUID: "order3"}

		// Test Set and Get
		cache.Set(order1)
		cache.Set(order2)

		result, exists := cache.Get("order1")
		assert.True(t, exists)
		assert.Equal(t, order1, result)

		result, exists = cache.Get("order2")
		assert.True(t, exists)
		assert.Equal(t, order2, result)

		// Test capacity - adding third item should evict first
		cache.Set(order3)

		_, exists = cache.Get("order1")
		assert.False(t, exists, "order1 should be evicted")

		result, exists = cache.Get("order2")
		assert.True(t, exists)
		assert.Equal(t, order2, result)

		result, exists = cache.Get("order3")
		assert.True(t, exists)
		assert.Equal(t, order3, result)
	})

	t.Run("update existing order", func(t *testing.T) {
		cache := repository.NewCache(2)

		order1 := &models.Order{OrderUID: "order1", TrackNumber: "old"}
		order1Updated := &models.Order{OrderUID: "order1", TrackNumber: "new"}

		cache.Set(order1)
		cache.Set(order1Updated)

		result, exists := cache.Get("order1")
		assert.True(t, exists)
		assert.Equal(t, "new", result.TrackNumber)
	})

	t.Run("get moves to front", func(t *testing.T) {
		cache := repository.NewCache(3)

		order1 := &models.Order{OrderUID: "order1"}
		order2 := &models.Order{OrderUID: "order2"}
		order3 := &models.Order{OrderUID: "order3"}
		order4 := &models.Order{OrderUID: "order4"}

		cache.Set(order1)
		cache.Set(order2)
		cache.Set(order3)

		// Access order1 to move it to front
		_, exists := cache.Get("order1")
		assert.True(t, exists)

		// Add order4 - should evict order2 (least recently used)
		cache.Set(order4)

		_, exists = cache.Get("order2")
		assert.False(t, exists, "order2 should be evicted")

		// order1, order3, order4 should exist
		_, exists = cache.Get("order1")
		assert.True(t, exists)
		_, exists = cache.Get("order3")
		assert.True(t, exists)
		_, exists = cache.Get("order4")
		assert.True(t, exists)
	})

	t.Run("restore cache", func(t *testing.T) {
		cache := repository.NewCache(2)

		orders := map[string]*models.Order{
			"order1": {OrderUID: "order1"},
			"order2": {OrderUID: "order2"},
			"order3": {OrderUID: "order3"}, // This should be evicted due to capacity
		}

		cache.Restore(orders)

		assert.Equal(t, 2, cache.Size())

		_, exists := cache.Get("order3")
		assert.False(t, exists, "order3 should be evicted due to capacity")
	})

	t.Run("get all", func(t *testing.T) {
		cache := repository.NewCache(3)

		order1 := &models.Order{OrderUID: "order1"}
		order2 := &models.Order{OrderUID: "order2"}

		cache.Set(order1)
		cache.Set(order2)

		allOrders := cache.GetAll()
		assert.Equal(t, 2, len(allOrders))
		assert.Equal(t, order1, allOrders["order1"])
		assert.Equal(t, order2, allOrders["order2"])
	})

	t.Run("size", func(t *testing.T) {
		cache := repository.NewCache(3)

		assert.Equal(t, 0, cache.Size())

		cache.Set(&models.Order{OrderUID: "order1"})
		assert.Equal(t, 1, cache.Size())

		cache.Set(&models.Order{OrderUID: "order2"})
		assert.Equal(t, 2, cache.Size())
	})
}

func TestLRUCache_EdgeCases(t *testing.T) {
	t.Run("zero capacity", func(t *testing.T) {
		cache := repository.NewCache(0)
		// Вместо проверки неэкспортируемого поля, проверяем поведение
		cache.Set(&models.Order{OrderUID: "order1"})
		cache.Set(&models.Order{OrderUID: "order2"})
		// Должен работать без паники
		assert.True(t, cache.Size() > 0)
	})

	t.Run("negative capacity", func(t *testing.T) {
		cache := repository.NewCache(-5)
		// Вместо проверки неэкспортируемого поля, проверяем поведение
		cache.Set(&models.Order{OrderUID: "order1"})
		assert.Equal(t, 1, cache.Size())
	})

	t.Run("concurrent access", func(t *testing.T) {
		cache := repository.NewCache(100)

		// Run multiple goroutines to test concurrent access
		done := make(chan bool)

		// Writer goroutine
		go func() {
			for i := range 100 {
				orderUID := fmt.Sprintf("order%d", i)
				cache.Set(&models.Order{OrderUID: orderUID})
			}
			done <- true
		}()

		// Reader goroutine
		go func() {
			for i := range 100 {
				orderUID := fmt.Sprintf("order%d", i)
				cache.Get(orderUID)
			}
			done <- true
		}()

		// Wait for both goroutines
		<-done
		<-done

		// Should not panic and maintain integrity
		assert.True(t, cache.Size() <= 100)
	})
}
