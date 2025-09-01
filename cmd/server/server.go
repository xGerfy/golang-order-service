package main

import (
	"log"
	"net/http"
	"order-service/internal/handler"
	"order-service/internal/service"

	"github.com/gorilla/mux"
)

func setupRouter(svc *service.Service) *mux.Router {
	h := handler.New(svc)
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/order/{id}", h.GetOrder).Methods("GET")
	router.HandleFunc("/health", h.HealthCheck).Methods("GET")

	// Web interface
	router.HandleFunc("/", h.ServeWebInterface).Methods("GET")
	router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("web/static"))),
	)

	return router
}

func startHTTPServer(addr string, router *mux.Router) *http.Server {
	server := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("HTTP сервер запущен на %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка HTTP сервера: %v", err)
		}
	}()

	return server
}
