package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"order-service/internal/models"
	"order-service/internal/repository"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	repo     repository.OrderRepository
	cache    repository.OrderCache
	validate *validator.Validate
}

func New(repo repository.OrderRepository, cache repository.OrderCache) *Service {
	validate := validator.New(validator.WithRequiredStructEnabled())

	svc := &Service{
		repo:     repo,
		cache:    cache,
		validate: validate,
	}

	// Запись в кэш из бд при загрузке
	svc.restoreCache()

	return svc
}

// ProcessOrder обрабатывает заказ из Kafka
func (s *Service) ProcessOrder(order *models.Order) error {
	ctx := context.Background()

	// ВАЛИДАЦИЯ перед сохранением
	if err := s.validateOrder(order); err != nil {
		return fmt.Errorf("валидация заказа failed: %w", err)
	}

	// Сохранение в бд
	if err := s.repo.SaveOrder(ctx, order); err != nil {
		return err
	}

	// Обновление кэша
	s.cache.Set(order)

	log.Printf("Заказ %s успешно обработан и сохранен", order.OrderUID)
	return nil
}

// ProcessOrderFromJSON обрабатывает сырые JSON данные из Kafka
func (s *Service) ProcessOrderFromJSON(data []byte) error {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	return s.ProcessOrder(&order)
}

func (s *Service) GetOrder(ctx context.Context, orderUID string) (*models.Order, error) {
	// Взять из кэша (быстрая операция - не нужен context)
	if order, exists := s.cache.Get(orderUID); exists {
		return order, nil
	}

	// Запрос к БД (медленная операция - нужен context)
	order, err := s.repo.GetOrder(ctx, orderUID) // Передаем context в репозиторий
	if err != nil {
		return nil, err
	}

	s.cache.Set(order)
	return order, nil
}

func (s *Service) restoreCache() {
	ctx := context.Background()
	log.Println("Получение кэша из базы данных...")

	orders, err := s.repo.GetAllOrders(ctx)
	if err != nil {
		log.Printf("Ошибка получения кэша: %v", err)
		return
	}

	s.cache.Restore(orders)
	log.Printf("Кэш загружен, записей: %d", len(orders))
}

func (s *Service) HealthCheck(ctx context.Context) error {
	return s.repo.HealthCheck(ctx) // Передаем context в репозиторий
}

// validateOrder валидирует заказ с помощью go-validator
func (s *Service) validateOrder(order *models.Order) error {
	// Базовая валидация структуры
	if err := s.validate.Struct(order); err != nil {
		return fmt.Errorf("структурная валидация: %w", err)
	}

	// Дополнительная бизнес-логика валидации
	return s.businessValidation(order)
}

// businessValidation дополнительная проверка бизнес-правил
func (s *Service) businessValidation(order *models.Order) error {
	// Проверка суммы товаров
	itemsTotal := 0
	for _, item := range order.Items {
		itemsTotal += item.TotalPrice
	}

	if itemsTotal != order.Payment.GoodsTotal {
		return fmt.Errorf("несоответствие сумм: goods_total=%d, сумма товаров=%d",
			order.Payment.GoodsTotal, itemsTotal)
	}

	// Проверка даты (не может быть из будущего)
	if order.DateCreated.After(time.Now().Add(24 * time.Hour)) {
		return fmt.Errorf("дата создания заказа из будущего: %s", order.DateCreated)
	}

	return nil
}
