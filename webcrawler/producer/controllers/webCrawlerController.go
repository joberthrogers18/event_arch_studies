package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	services "github.com/joberthrogers18/webcrawler_queue/producer/services"
	"github.com/joberthrogers18/webcrawler_queue/producer/utils"
	"github.com/wagslane/go-rabbitmq"
)

func WebCrawlerGet(c *gin.Context, publisher *rabbitmq.Publisher) {
	fmt.Println("A new webCrawler was started in go routine")
	go services.StartWebCrawler(publisher)
	c.JSON(utils.OK_STATUS, gin.H{
		"message": "the crawler has just been started",
	})
}
