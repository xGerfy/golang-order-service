package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"order-service/internal/models"

	"github.com/segmentio/kafka-go"
)

func main() {
	// Адрес брокера Kafka
	brokers := []string{"localhost:9092"}

	// Топик, куда отправляются заказы
	topic := "orders"

	// Создаем экземпляр нового writer'а
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})

	// Формируем тестовый заказ
	testOrder := models.Order{
		OrderUID:    "testorder",
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: models.Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: models.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []models.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}

	// Преобразуем заказ в JSON
	jsonData, err := json.Marshal(testOrder)
	if err != nil {
		log.Fatalln("Ошибка преобразования в JSON:", err)
		os.Exit(1)
	}

	// Отправляем сообщение
	err = writer.WriteMessages(context.Background(), kafka.Message{
		Value: jsonData,
	})
	if err != nil {
		log.Fatalln("Ошибка отправки сообщения:", err)
		os.Exit(1)
	}

	// Закрываем writer
	err = writer.Close()
	if err != nil {
		log.Fatalln("Ошибка закрытия writer:", err)
		os.Exit(1)
	}

	log.Println("Сообщение успешно отправлено!")
}
