package main

import (
	"log"
	"order-service/internal/cache"
	"order-service/internal/config"
	"order-service/internal/database"
	"order-service/internal/kafka"
	"order-service/internal/service"
)

func setupComponents(cfg *config.Config) (*database.Postgres, *service.Service, *kafka.Consumer) {
	db, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к бд: %v", err)
	}

	cache := cache.New()
	svc := service.New(db, cache)
	kafkaConsumer := kafka.New(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID, svc)

	return db, svc, kafkaConsumer
}
