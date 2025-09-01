package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBrokers []string
	KafkaTopic   string
	KafkaGroupID string
	DatabaseURL  string
	HTTP_ADDR    string
}

func Load() (*Config, error) {
	// Загружаем .env файл, но игнорируем ошибку если файла нет
	if err := godotenv.Load(); err != nil {
		// Проверяем, что ошибка именно из-за отсутствия файла, а не других проблем
		if !os.IsNotExist(err) {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
		// Продолжаем работу, так как переменные могут быть установлены в окружении
	}

	return &Config{
		KafkaBrokers: strings.Split(os.Getenv("KAFKA_BROKERS"), ","),
		KafkaTopic:   os.Getenv("KAFKA_TOPIC"),
		KafkaGroupID: os.Getenv("KAFKA_GROUP_ID"),
		DatabaseURL:  os.Getenv("DB_URL"),
		HTTP_ADDR:    os.Getenv("HTTP_ADDR"),
	}, nil
}
