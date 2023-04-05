package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/joberthrogers18/webcrawler_queue/producer/controllers"
	"github.com/wagslane/go-rabbitmq"
)

func GetRoutesCrawler(r *gin.Engine, publisher *rabbitmq.Publisher) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.GET("/crawler", func(c *gin.Context) {
		controllers.WebCrawlerGet(c, publisher)
		c.JSON(200, gin.H{
			"message": "webCrawler started succefully",
		})
	})
}
