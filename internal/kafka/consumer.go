package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"order-service/internal/models"
	"order-service/internal/service"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	service *service.Service
}

func New(brokers []string, topic string, groupID string, svc *service.Service) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
	})

	return &Consumer{
		reader:  reader,
		service: svc,
	}
}

func (c *Consumer) Start() {
	log.Println("Запуск косюмера кафки...")

	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка чтения сообщения: %v", err)
			continue
		}

		go c.processMessage(msg.Value)
	}
}

func (c *Consumer) processMessage(data []byte) {
	var order models.Order
	if err := json.Unmarshal(data, &order); err != nil {
		log.Printf("Ошибка преобразования сообщения: %v", err)
		return
	}

	if err := c.service.ProcessOrder(&order); err != nil {
		log.Printf("Ошибка обработки заказа %s: %v", order.OrderUID, err)
		return
	}

	log.Printf("Успешная обработка заказа: %s", order.OrderUID)
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
