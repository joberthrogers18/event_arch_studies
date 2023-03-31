package config

import (
	"fmt"

	"github.com/wagslane/go-rabbitmq"
)

func initializeRabbitMq() (*rabbitmq.Conn, *rabbitmq.Publisher) {
	conn, err := rabbitmq.NewConn(
		"amqp://rabbitmq:rabbitmq@rabbitmq1/",
		rabbitmq.WithConnectionOptionsLogging,
	)

	if err != nil {
		fmt.Println(err)
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName("events"),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)

	if err != nil {
		fmt.Println(err)
	}

	return conn, publisher
}
