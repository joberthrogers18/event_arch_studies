package main

import (
	"github.com/gin-gonic/gin"
	config "github.com/joberthrogers18/webcrawler_queue/producer/config"
	"github.com/joberthrogers18/webcrawler_queue/producer/routes"
)

func main() {
	conn, publisher := config.InitializeRabbitMq()

	defer conn.Close()
	defer publisher.Close()

	r := gin.Default()

	routes.GetRoutesCrawler(r, publisher)

	err := r.Run(":8080")

	if err != nil {
		panic("[Error] failed to started Gin server due to: " + err.Error())
	}
}
