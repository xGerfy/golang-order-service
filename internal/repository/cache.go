package repository

import (
	"container/list"
	"sync"

	"order-service/internal/models"
)

// Проверяем, что Cache реализует интерфейс repository.Cache
var _ OrderCache = (*LRUCache)(nil)

type LRUItem struct {
	order   *models.Order
	element *list.Element
}

type LRUCache struct {
	mu       sync.RWMutex
	orders   map[string]*LRUItem
	list     *list.List
	capacity int
}

func NewCache(capacity int) *LRUCache {
	if capacity <= 0 {
		capacity = 1000 // дефолтный размер
	}

	return &LRUCache{
		orders:   make(map[string]*LRUItem),
		list:     list.New(),
		capacity: capacity,
	}
}

func (c *LRUCache) Set(order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если уже существует - обновляем и перемещаем в начало
	if item, exists := c.orders[order.OrderUID]; exists {
		item.order = order
		c.list.MoveToFront(item.element)
		return
	}

	// Если достигли capacity - удаляем самый старый (инвалидация!)
	if c.list.Len() >= c.capacity {
		oldest := c.list.Back()
		if oldest != nil {
			oldestItem := oldest.Value.(*LRUItem)
			delete(c.orders, oldestItem.order.OrderUID)
			c.list.Remove(oldest)
		}
	}

	// Добавляем новый элемент
	element := c.list.PushFront(&LRUItem{order: order})
	c.orders[order.OrderUID] = &LRUItem{
		order:   order,
		element: element,
	}
}

func (c *LRUCache) Get(orderUID string) (*models.Order, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, exists := c.orders[orderUID]
	if !exists {
		return nil, false
	}

	// Перемещаем в начало (последний использованный)
	c.list.MoveToFront(item.element)
	return item.order, true
}

func (c *LRUCache) GetAll() map[string]*models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]*models.Order)
	for key, item := range c.orders {
		result[key] = item.order
	}

	return result
}

func (c *LRUCache) Restore(orders map[string]*models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.orders = make(map[string]*LRUItem)
	c.list = list.New()

	// Добавляем элементы пока не достигнем capacity
	for _, order := range orders {
		// Если достигли capacity - выходим
		if c.list.Len() >= c.capacity {
			break
		}

		element := c.list.PushFront(&LRUItem{order: order})
		c.orders[order.OrderUID] = &LRUItem{
			order:   order,
			element: element,
		}
	}
}

func (c *LRUCache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.orders)
}
