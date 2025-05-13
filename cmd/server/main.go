package main

import (
	"log"

	"bff-service/internal/api"
)

func main() {
	// Инициализация и запуск REST сервера
	server := api.NewServer()
	log.Println("REST server starting on :8081")
	if err := server.Start(":8081"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
