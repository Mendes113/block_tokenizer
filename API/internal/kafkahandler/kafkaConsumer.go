package kafka

import (
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "log"
)

var consumer *kafka.Consumer

// InitConsumer inicializa o consumidor Kafka
func InitConsumer(topic string) {
    var err error
    consumer, err = kafka.NewConsumer(&kafka.ConfigMap{
        "bootstrap.servers":  "localhost:9092",
        "group.id":           "blockchain_data",
        "auto.offset.reset": "earliest",
    })
    if err != nil {
        log.Fatalf("Falha ao criar o consumidor Kafka: %v", err)
    }

    err = consumer.SubscribeTopics([]string{topic}, nil)
    if err != nil {
        log.Fatalf("Falha ao subscrever tópico: %v", err)
    }

    go func() {
        for {
            msg, err := consumer.ReadMessage(-1)
            if err == nil {
                log.Printf("Mensagem recebida do tópico %s: %s\n", *msg.TopicPartition.Topic, string(msg.Value))
            } else {
                log.Printf("Erro ao receber mensagem: %v\n", err)
            }
        }
    }()
}

// CloseConsumer fecha o consumidor Kafka
func CloseConsumer() {
    consumer.Close()
}
