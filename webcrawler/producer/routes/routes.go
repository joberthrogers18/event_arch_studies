package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wagslane/go-rabbitmq"
)

func getRoutesCrawler(r *gin.Engine, publisher *rabbitmq.Publisher) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.GET("/", controllers.webCrawler)
}
