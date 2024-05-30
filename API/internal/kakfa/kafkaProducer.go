package kafka

import (
    "github.com/confluentinc/confluent-kafka-go/kafka"
    "log"
)

var producer *kafka.Producer

// InitProducer inicializa o produtor Kafka
func InitProducer() {
    var err error
    producer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
    if err != nil {
        log.Fatalf("Failed to create Kafka producer: %v", err)
    }

    go func() {
        for e := range producer.Events() {
            switch ev := e.(type) {
            case *kafka.Message:
                if ev.TopicPartition.Error != nil {
                    log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
                } else {
                    log.Printf("Delivered message to topic %s [%d] at offset %v\n",
                        *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset)
                }
            }
        }
    }()
}

// SendMessage envia uma mensagem para um tópico Kafka
func SendMessage(topic string, message []byte) {
    producer.Produce(&kafka.Message{
        TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
        Value:          message,
    }, nil)
}

// CloseProducer fecha o produtor Kafka
func CloseProducer() {
    producer.Close()
}
