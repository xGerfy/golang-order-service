package repository

import (
	"context"

	"order-service/internal/models"
)

// OrderRepository интерфейс для работы с заказами в БД
type OrderRepository interface {
	SaveOrder(ctx context.Context, order *models.Order) error
	GetOrder(ctx context.Context, orderUID string) (*models.Order, error)
	GetAllOrders(ctx context.Context) (map[string]*models.Order, error)
	HealthCheck(ctx context.Context) error
	Close()
}

// Cache интерфейс для кэша
type OrderCache interface {
	Set(order *models.Order)
	Get(orderUID string) (*models.Order, bool)
	GetAll() map[string]*models.Order
	Restore(orders map[string]*models.Order)
	Size() int
}
