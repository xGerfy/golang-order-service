package main

import (
	"log"

	"order-service/internal/config"
	"order-service/internal/kafka"
	"order-service/internal/repository"
	"order-service/internal/service"
)

func setupComponents(cfg *config.Config) (*repository.DB, *service.Service, *kafka.Consumer) {
	db, err := repository.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к бд: %v", err)
	}

	cache := repository.NewCache(cfg.CacheCapacity)
	svc := service.New(db, cache)
	kafkaConsumer := kafka.New(cfg.KafkaBrokers, cfg.KafkaTopic, cfg.KafkaGroupID, svc)

	return db, svc, kafkaConsumer
}
