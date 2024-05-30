package main

import (
	"blockchain.api/internal/handlers"
	"blockchain.api/internal/services"
	"github.com/confluentinc/confluent-kafka-go/kafka" //using fiber
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	 producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	 if err != nil {
		 panic(err)
	 }
 
	 defer producer.Close()
 
	 blockchainService := services.NewBlockchainService(producer)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/health", handlers.HealthCheckHandler)
	app.Post("/block", handlers.BlockAddHandler(blockchainService))

	app.Listen(":3000")
}

