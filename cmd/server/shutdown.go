package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func waitForShutdown(server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	log.Println("Получен сигнал выключения...")

	if err := server.Close(); err != nil {
		log.Printf("Ошибка во время остановки сервера: %v", err)
	}

	log.Println("Сервер остановлен")
}
