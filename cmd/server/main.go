package main

import (
	"log"

	"order-service/internal/config"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация компонентов
	db, svc, kafkaConsumer := setupComponents(cfg)
	defer db.Close()
	defer kafkaConsumer.Close()

	// Запуск Kafka консюмер в фоне
	go kafkaConsumer.Start()

	// Настройка и запуск HTTP сервера
	router := setupRouter(svc)
	server := startHTTPServer(cfg.HTTP_ADDR, router)

	// Ожидание сигнала завершения
	waitForShutdown(server)
}
