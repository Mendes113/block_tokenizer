package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func main() {
	// Use a URL do Ganache
	var ganacheUrl = "http://127.0.0.1:8545"
	client, err := ethclient.DialContext(context.Background(), ganacheUrl)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a connection")
	defer client.Close()

	// Endereço do contrato
	contractAddress := common.HexToAddress("0xD7ACd2a9FD159E69Bb102A1ca21C9a3e3A5F771B")

	// ABI do contrato (baseado no ABI fornecido)
	contractABI := `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"userAddress","type":"address"},{"indexed":false,"internalType":"string","name":"message","type":"string"}],"name":"Notification","type":"event"}]`

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatal("Failed to parse ABI:", err)
	}

	// Iniciar consumidor Kafka
	consumer := InitConsumer("blockchain_data")
	defer consumer.Close()


	godotenv.Load()
	PKEY := os.Getenv("PRIVATE_KEY")
	privateKeyHex := PKEY  // Chave privada do Ganache

	// Converter a chave privada hexadecimal em uma estrutura *ecdsa.PrivateKey
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("Failed to convert private key:", err)
	}

	for {
		select {
		case <-time.After(10 * time.Second):
			processLogs(client, parsedABI, contractAddress)
		case msg := <-consumeKafkaMessages(consumer):
			// Processar a mensagem Kafka
			novoBlocoDeDados := string(msg.Value)
			fmt.Printf("Mensagem recebida do tópico %s: %s\n", *msg.TopicPartition.Topic, novoBlocoDeDados)
			err = adicionarBlocoDeDados(client, privateKey, contractAddress, novoBlocoDeDados)
			if err != nil {
				log.Println("Failed to add new data block to contract:", err)
			} else {
				log.Println("New data block added to contract successfully")
			}
		}
	}
}

func processLogs(client *ethclient.Client, parsedABI abi.ABI, contractAddress common.Address) {
	// Obter logs de eventos
	logs, err := client.FilterLogs(context.Background(), ethereum.FilterQuery{
		Addresses: []common.Address{contractAddress},
	})
	if err != nil {
		log.Println("Failed to get logs:", err)
		return
	}

	// Processar logs de eventos
	for _, log := range logs {
		var eventMap = make(map[string]interface{})
		err := parsedABI.UnpackIntoMap(eventMap, "Notification", log.Data)
		if err != nil {
			fmt.Print("Failed to decode event log:", err)
			continue
		}
		var event struct {
			UserAddress common.Address
			Message     string
		}
		event.UserAddress = eventMap["userAddress"].(common.Address)
		event.Message = eventMap["message"].(string)
		fmt.Printf("Nova notificação recebida: UserAddress=%s, Message=%s\n", event.UserAddress.Hex(), event.Message)
	}
}

func adicionarBlocoDeDados(client *ethclient.Client, privateKey *ecdsa.PrivateKey, contratoEndereco common.Address, dados string) error {
	contratoAbi := `[{"constant":false,"inputs":[{"name":"_dados","type":"string"}],"name":"adicionarBlocoDeDados","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"anonymous":false,"inputs":[{"indexed":true,"name":"quem","type":"address"},{"indexed":false,"name":"dados","type":"string"}],"name":"NovoBlocoDeDadosAdicionado","type":"event"}]`
	parsedABI, err := abi.JSON(strings.NewReader(contratoAbi))
	if err != nil {
		return fmt.Errorf("falha ao fazer o parsing da ABI do contrato: %v", err)
	}

	// Obter o ID do método 'adicionarBlocoDeDados'
	metodo, ok := parsedABI.Methods["adicionarBlocoDeDados"]
	if !ok {
		return fmt.Errorf("método adicionarBlocoDeDados não encontrado na ABI do contrato")
	}

	// Codificar os argumentos do método
	dado, err := parsedABI.Pack(metodo.Name, dados)
	if err != nil {
		return fmt.Errorf("falha ao empacotar os dados para o método adicionarBlocoDeDados: %v", err)
	}

	// Obter o nonce
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("falha ao obter o nonce da conta: %v", err)
	}

	// Definir o gasPrice
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("falha ao obter o preço do gas: %v", err)
	}

	// Definir o limite de gas
	gasLimit := uint64(300000) // ajuste conforme necessário

	// Construir a transação
	tx := types.NewTransaction(nonce, contratoEndereco, big.NewInt(0), gasLimit, gasPrice, dado)

	// Assinar a transação
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		return fmt.Errorf("falha ao assinar a transação: %v", err)
	}

	// Enviar a transação
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("falha ao enviar a transação: %v", err)
	}

	return nil
}

func InitConsumer(topic string) *kafka.Consumer {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
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

	return consumer
}

func consumeKafkaMessages(consumer *kafka.Consumer) <-chan *kafka.Message {
	msgChan := make(chan *kafka.Message)

	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err == nil {
				log.Printf("Mensagem recebida do tópico %s: %s\n", *msg.TopicPartition.Topic, string(msg.Value))
				msgChan <- msg
			} else {
				log.Printf("Erro ao receber mensagem: %v\n", err)
			}
		}
	}()

	return msgChan
}
