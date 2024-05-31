package main

import (
	"blockchain.api/internal/database"
	"blockchain.api/internal/handlers"
	mykafka "blockchain.api/internal/kafkahandler"
	"blockchain.api/internal/services"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Inicializar Fiber
	app := fiber.New()

	// Abrir conexão com o MongoDB
	client, err := database.OpenConnection()
	if err != nil {
		log.Fatalf("Erro ao abrir a conexão com o MongoDB: %v", err)
	}
	defer database.Close(client)

	// Configurar o produtor Kafka
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Erro ao configurar o produtor Kafka: %v", err)
	}
	defer producer.Close()

	// Inicializar os serviços da blockchain com o produtor e o cliente MongoDB
	msgChan := make(chan *kafka.Message)
	consumer := mykafka.InitConsumer("blockchain_data", msgChan)

	blockchainService := services.NewBlockchainService(producer, client)
	blockchainDataService := services.NewBlockchainDataService(consumer, client, msgChan)
	_ = blockchainDataService

	// Definir rotas
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/health", handlers.HealthCheckHandler)
	app.Post("/block", handlers.BlockAddHandler(blockchainService))

	// Manejar o encerramento gracioso
	go func() {
		// Escutar na porta 3000
		if err := app.Listen(":3000"); err != nil {
			log.Fatalf("Erro ao iniciar o servidor Fiber: %v", err)
		}
	}()

	// Esperar por sinais de encerramento para fechar o aplicativo corretamente
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
	<-sigint

	// Fechar o aplicativo Fiber
	if err := app.Shutdown(); err != nil {
		log.Fatalf("Erro ao fechar o servidor Fiber: %v", err)
	}

	// Fechar o consumidor Kafka
	mykafka.CloseConsumer()
}
