package main

import (
	"log"

	"github.com/wagslane/go-rabbitmq"
)

func main() {
	conn, err := rabbitmq.NewConn(
		"amqp://guest:guest@rabbitmq:5672/",
		rabbitmq.WithConnectionOptionsLogging,
	)

	if err != nil {
		log.Fatal("Error start connection amqp: ", err)
	}

	defer conn.Close()

	consumer, err := rabbitmq.NewConsumer(
		conn,
		func(d rabbitmq.Delivery) (action rabbitmq.Action) {
			log.Printf("consumed: %v", string(d.Body))

			return rabbitmq.Ack
		},
		"hello",
		rabbitmq.WithConsumerOptionsRoutingKey("my_routing_key"),
		rabbitmq.WithConsumerOptionsExchangeName("events"),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
	)

	if err != nil {
		log.Fatal(err)
	}

	defer consumer.Close()
}
