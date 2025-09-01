package service

import (
	"log"
	"order-service/internal/cache"
	"order-service/internal/database"
	"order-service/internal/models"
)

type Service struct {
	db    *database.Postgres
	cache *cache.Cache
}

func New(db *database.Postgres, cache *cache.Cache) *Service {
	svc := &Service{
		db:    db,
		cache: cache,
	}

	// Запись в кэш из бд при загрузке
	svc.restoreCache()

	return svc
}

func (s *Service) ProcessOrder(order *models.Order) error {
	// Сохранение в бд
	if err := s.db.SaveOrder(order); err != nil {
		return err
	}

	// Обновление кэша
	s.cache.Set(order)

	return nil
}

func (s *Service) GetOrder(orderUID string) (*models.Order, error) {
	// Взять из кэша
	if order, exists := s.cache.Get(orderUID); exists {
		return order, nil
	}

	// Если нет в кэе берем из бж
	order, err := s.db.GetOrder(orderUID)
	if err != nil {
		return nil, err
	}

	// Обновление кэша
	s.cache.Set(order)

	return order, nil
}

func (s *Service) restoreCache() {
	log.Println("Получение кэша из базы данных...")

	orders, err := s.db.GetAllOrders()
	if err != nil {
		log.Printf("Ошибка получения кэша: %v", err)
		return
	}

	s.cache.Restore(orders)
	log.Printf("Кэш загружен с %d записями", len(orders))
}

func (s *Service) HealthCheck() error {
	return s.db.HealthCheck()
}
