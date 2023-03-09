package main

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Rodrigo001-dev/golang-intensive/internal/infra/database"
	"github.com/Rodrigo001-dev/golang-intensive/internal/usecases"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/Rodrigo001-dev/golang-intensive/pkg/kafka"
	"github.com/Rodrigo001-dev/golang-intensive/pkg/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	db, err := sql.Open("sqlite3", "./orders.db")

	if err != nil {
		panic(err)
	}

	defer db.Close() // executa tudo e depois executa o Close

	repository := database.NewOrderRepository(db)
	usecase := usecases.CalculateFinalPrice(OrderRepository: repository)

	msgChanKafka := make(chan *ckafka.Message)

	// kafka
	topics := []string("orders")
	servers := "host.docker.internal:9094"
	go kafka.Consumer(topics, servers, msgChanKafka)
	go kafkaWorker(msgChanKafkam, usecase)

	// rabbitmq

	ch, err := rabbitmq.OpenChannle()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	msgRabbitmqChannel := make(chan amqp.Delivery)
	go rabbitmq.Consume(ch, msgRabbitmqChannel)
	rabbitmqWorker(msgRabbitmqChannel, usecase)
}

func kafkaWorker(msgChan chan *ckafka.Message, uc usecases.CalculateFinalPrice) {
	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Value, &OrderInputDTO)

		if err != nil {
			panic(err)
		}

		outputDto, err := uc.Execute(OrderInputDTO)

		if err != nil {
			panic(err)
		}

		fmt.Printf("kafka has process order %s\n", outputDto.ID)
	}
}

func rabbitmqWorker(msgChan chan amqp.Delivery, uc usecases.CalculateFinalPrice) {
	fmt.Println("Rabbitmq worker has started")

	for msg := range msgChan {
		var OrderInputDTO usecase.OrderInputDTO
		err := json.Unmarshal(msg.Body, &OrderInputDTO)

		if err != nil {
			panic(err)
		}

		outputDto, err := uc.Execute(OrderInputDTO)
		if err != nil {
			panic(err)
		}

		msg.Ask(false)
		fmt.Printf("Rabbitmq has processed order %s\n", outputDto)
	}
}