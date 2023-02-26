package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
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

func main() {
	crawler()
}

func getDatesBirthAndDeath(data string) []string {
	regCompile := regexp.MustCompile(`[A-Z][a-z]+\s[0-9]{1,2}\,\s[0-9]{4}`)
	regex := regCompile.FindAllString(data, -1)

	return regex
}

func crawler() {
	c := colly.NewCollector(
		colly.AllowedDomains("imdb.com", "www.imdb.com"),
	)

	infoCollector := c.Clone()

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL.String())
	})

	infoCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting Profile URL: ", r.URL.String())
	})

	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileURL := e.ChildAttr("div.lister-item-image > a", "href")
		profileURL = e.Request.AbsoluteURL(profileURL)
		fmt.Println(profileURL)
		infoCollector.Visit(profileURL)
	})

	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println(nextPage)
		c.Visit(nextPage)
	})

	infoCollector.OnHTML(".ipc-page-wrapper.ipc-page-wrapper--base", func(e *colly.HTMLElement) {
		tmpProfile := star{}
		tmpProfile.Name = e.ChildText("h1 > span.sc-afe43def-1.fDTGTb")
		tmpProfile.Photo = e.ChildAttr("img.ipc-image", "src")
		tmpProfile.Bio = strings.TrimSpace(e.ChildText(".ipc-html-content--baseAlt > .ipc-html-content-inner-div"))
		birthDayStr := e.ChildText("span.sc-dec7a8b-2.haviXP:nth-child(2)")
		var middleStr int64 = int64(math.Round(float64(len(birthDayStr) / 2)))
		var datesActor []string = getDatesBirthAndDeath(birthDayStr[middleStr:])

		tmpProfile.BirthDate = datesActor[0]

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
	})

	c.Visit("https://www.imdb.com/search/name/?birth_monthday=12-20")
	// infoCollector.Visit("https://www.imdb.com/search/name/?birth_monthday=12-20")
}
