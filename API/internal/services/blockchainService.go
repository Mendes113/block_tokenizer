package services

import (
	"log"
	"blockchain.api/internal/database"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.mongodb.org/mongo-driver/mongo"
)

// BlockchainService é a estrutura que contém os métodos relacionados à blockchain
type BlockchainService struct {
	producer *kafka.Producer
	client   *mongo.Client
}

// BlockchainDataService é a estrutura que contém os métodos relacionados ao consumo de dados da blockchain
type BlockchainDataService struct {
	consumer *kafka.Consumer
	client   *mongo.Client
	msgChan  chan *kafka.Message
}

// NewBlockchainService cria uma nova instância do serviço BlockchainService
func NewBlockchainService(producer *kafka.Producer, client *mongo.Client) *BlockchainService {
	return &BlockchainService{
		producer: producer,
		client:   client,
	}
}

// NewBlockchainDataService cria uma nova instância do serviço BlockchainDataService
func NewBlockchainDataService(consumer *kafka.Consumer, client *mongo.Client, msgChan chan *kafka.Message) *BlockchainDataService {
	service := &BlockchainDataService{
		consumer: consumer,
		client:   client,
		msgChan:  msgChan,
	}

	go service.processMessages()

	return service
}

// AddBlock envia uma mensagem para o Kafka para adicionar um bloco à blockchain
func (s *BlockchainService) AddBlock(data string) error {
	topic := "blockchain"
	s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(data),
	}, nil)

	log.Printf("Sent message to topic %s: %s\n", topic, data)
	return nil
}

// GetMessageAndSave salva uma mensagem recebida no MongoDB
func (s *BlockchainDataService) GetMessageAndSave(message MessageData) {
	database.SaveCollection(s.client, "blockchain_data", message)
}

// processMessages processa mensagens recebidas do Kafka
func (s *BlockchainDataService) processMessages() {
	for msg := range s.msgChan {
		messageData := MessageData{Message: string(msg.Value)}
		s.GetMessageAndSave(messageData)
	}
}



type MessageData struct {
	Message string `json:"message"`
}

