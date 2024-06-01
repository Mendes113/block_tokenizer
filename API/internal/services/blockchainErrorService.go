package services

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// ErrorHandlingService é o serviço responsável por lidar com mensagens de erro do Kafka
type ErrorHandlingService struct {
    errorConsumer *kafka.Consumer
    errChan       chan error
}


// NewErrorHandlingService cria uma nova instância do serviço de manipulação de erros do Kafka
func NewErrorHandlingService(errorConsumer *kafka.Consumer, errChan chan error) *ErrorHandlingService {
    return &ErrorHandlingService{
        errorConsumer: errorConsumer,
        errChan:       errChan,
    }
}


// Start inicia o serviço de manipulação de erros do Kafka
// Start inicia o serviço de manipulação de erros do Kafka
func (s *ErrorHandlingService) Start(consumer *kafka.Consumer) {
    go func() {
        for {
            msg, err := consumer.ReadMessage(-1)
            if err == nil {
                log.Printf("Mensagem recebida do tópico %s: %s\n", *msg.TopicPartition.Topic, string(msg.Value))
            } else {
                log.Printf("Erro ao receber mensagem: %v\n", err)
                s.errChan <- err // envia para o canal
            }
        }
    }()
}



// Stop para o serviço de manipulação de erros do Kafka
func (s *ErrorHandlingService) Stop() {
	s.errorConsumer.Close()
}
