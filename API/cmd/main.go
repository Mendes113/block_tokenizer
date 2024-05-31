package main

import (
	"blockchain.api/internal/handlers"
	"blockchain.api/internal/services"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	mykafka "blockchain.api/internal/kafkahandler"
)

func main() {
	app := fiber.New()

	// Configurar o produtor Kafka
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	mykafka.InitConsumer("blockchain_data")	

	// Inicializar os serviços da blockchain com o produtor e o consumidor
	blockchainService := services.NewBlockchainService(producer)


	// Definir rotas
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/health", handlers.HealthCheckHandler)
	app.Post("/block", handlers.BlockAddHandler(blockchainService))

	// Escutar na porta 3000
	app.Listen(":3000")
}
