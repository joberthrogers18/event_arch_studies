package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"
	"sync"

	"github.com/gocolly/colly"
	"github.com/wagslane/go-rabbitmq"
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

var OK_STATUS = 200

func GetDatesBirthAndDeath(data string) []string {
	regCompile := regexp.MustCompile(`[A-Z][a-z]+\s[0-9]{1,2}\,\s[0-9]{4}`)
	regex := regCompile.FindAllString(data, -1)

	return regex
}

func GetCrawler(url string, publisher *rabbitmq.Publisher) {
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
		var datesActor []string = GetDatesBirthAndDeath(birthDayStr[middleStr:])

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

		err = publisher.PublishWithContext(
			context.Background(),
			[]byte(js),
			[]string{"my_routing_key"},
			rabbitmq.WithPublishOptionsContentType("application/json"),
			rabbitmq.WithPublishOptionsMandatory,
			rabbitmq.WithPublishOptionsPersistentDelivery,
			rabbitmq.WithPublishOptionsExchange("events"),
		)

		if err != nil {
			log.Println(err)
		}

		if err != nil {
			fmt.Println("Not Published in queue", err)
		}
	})

	collyInstMain.Visit(url)
}

func StartWebCrawler(publisher *rabbitmq.Publisher) {

	collyInst := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	collyInst.OnRequest(func(r *colly.Request) {
		fmt.Println(r.URL.String())
		wg.Add(1)
		go GetCrawler(r.URL.String(), publisher)
	})

	collyInst.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		var nextPage string = e.Request.AbsoluteURL(e.Attr("href"))
		collyInst.Visit(nextPage)
	})

	collyInst.Visit("https://www.imdb.com/search/name/?birth_monthday=12-20")
	wg.Wait()
}
