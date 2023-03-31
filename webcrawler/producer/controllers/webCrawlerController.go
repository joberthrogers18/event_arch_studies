package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/joberthrogers18/webcrawler_queue/producer/utils"
)

func webCrawlerGet(c *gin.Context) {
	fmt.Println("A new webCrawler was started in go routine")
	// go startWebCrawler(publisher)
	c.JSON(utils.OK_STATUS, gin.H{
		"message": "the crawler has just been started",
	})
}
