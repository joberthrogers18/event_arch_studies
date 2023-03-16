package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/streadway/amqp"
)

type movie struct {
	Title string
	Year  string
}

type star struct {
	Name  string
	Photo string
	// JobTitle  string
	BirthDate string
	DeathDate string
	Bio       string
	TopMovies []movie
}

var wg sync.WaitGroup

func getDatesBirthAndDeath(data string) []string {
	regCompile := regexp.MustCompile(`[A-Z][a-z]+\s[0-9]{1,2}\,\s[0-9]{4}`)
	regex := regCompile.FindAllString(data, -1)

	return regex
}

func crawler(url string, channel *amqp.Channel) {
	defer wg.Done()

	collyInstMain := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	infoCollectorInst := collyInstMain.Clone()

	collyInstMain.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	infoCollectorInst.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Profile URL: ", r.URL.String())
	})

	collyInstMain.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileURL := e.ChildAttr("div.lister-item-image > a", "href")
		profileURL = e.Request.AbsoluteURL(profileURL)
		fmt.Println(profileURL)
		infoCollectorInst.Visit(profileURL)
	})

	infoCollectorInst.OnHTML(".ipc-page-wrapper.ipc-page-wrapper--base", func(e *colly.HTMLElement) {
		tmpProfile := star{}
		tmpProfile.Name = e.ChildText("h1 > span.sc-afe43def-1.fDTGTb")
		tmpProfile.Photo = e.ChildAttr("img.ipc-image", "src")
		tmpProfile.Bio = strings.TrimSpace(e.ChildText(".ipc-html-content--baseAlt > .ipc-html-content-inner-div"))
		birthDayStr := e.ChildText("span.sc-dec7a8b-2.haviXP:nth-child(2)")
		var middleStr int64 = int64(math.Round(float64(len(birthDayStr) / 2)))
		var datesActor []string = getDatesBirthAndDeath(birthDayStr[middleStr:])

		tmpProfile.BirthDate = "-"
		if len(datesActor) > 0 {
			tmpProfile.BirthDate = datesActor[0]
		}

		tmpProfile.DeathDate = "-"
		if len(datesActor) > 1 {
			tmpProfile.DeathDate = datesActor[1]
		}

		e.ForEach("div.ipc-list-card--span.ipc-list-card--border-line", func(_ int, kf *colly.HTMLElement) {
			tmpMovie := movie{}
			tmpMovie.Title = kf.ChildText("div.ipc-primary-image-list-card__content-top > a.ipc-primary-image-list-card__title")
			tmpMovie.Year = kf.ChildText("div.ipc-primary-image-list-card__content-bottom > ul > li > span.ipc-primary-image-list-card__secondary-text")
			tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
		})

		js, err := json.MarshalIndent(tmpProfile, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))

		err = channel.Publish(
			"",
			"testing",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(string(js)),
			},
		)

		if err != nil {
			fmt.Println("Not Published in queue")
		}
	})

	collyInstMain.Visit(url)
}

func startWebCrawler(channel *amqp.Channel) {

	collyInst := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	collyInst.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
		wg.Add(1)
		go crawler(r.URL.String(), channel)
	})

	collyInst.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		var nextPage string = e.Request.AbsoluteURL(e.Attr("href"))
		collyInst.Visit(nextPage)
	})

	collyInst.Visit("https://www.imdb.com/search/name/?birth_monthday=12-20")
	wg.Wait()
}

func initializeRabbitMq() *amqp.Channel {
	fmt.Println("RabbitMQ: Getting started")

	connection, err := amqp.Dial("amqp:guest:guest@localhost:5672/")

	if err != nil {
		panic(err)
	}

	defer connection.Close()

	fmt.Print("Successfully connected to RabbitMQ instance \n\n")

	channel, err := connection.Channel()
	if err != nil {
		panic(err)
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"testing",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	fmt.Println("Queue status:", queue)

	return channel
}

func main() {
	channel := initializeRabbitMq()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		fmt.Println("A new webCrawler was started in go routine")
		go startWebCrawler(channel)
	})

	err := r.Run()

	if err != nil {
		panic("[Error] failed to started Gin server due to: " + err.Error())
	}
}
