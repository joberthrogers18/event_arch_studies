package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"

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
	BirthData string
	Bio       string
	TopMovies []movie
}

func main() {
	crawler()
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

	// c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
	// 	fmt.Println("============== Aquiiiiiiiii =================")
	// 	nextPage := e.Request.AbsoluteURL(e.Attr("href"))
	// 	fmt.Println(nextPage)
	// 	c.Visit(nextPage)
	// })

	infoCollector.OnHTML("section.sc-781acd7f-0", func(e *colly.HTMLElement) {
		fmt.Println(e.ChildText("h1"))
		tmpProfile := star{}
		tmpProfile.Name = e.ChildText("h1 > span.sc-afe43def-1.fDTGTb")
		tmpProfile.Photo = e.ChildAttr("img.ipc-image", "src")
		tmpProfile.Bio = e.ChildText(".ipc-html-content--baseAlt > .ipc-html-content-inner-div")
		// tmpProfile.JobTitle = e.ChildText("#name-job-categories > a > span.itemprop")
		birthDayStr := e.ChildText("span.sc-dec7a8b-2.haviXP:nth-child(2)")
		// tmpProfile.BirthData = e.ChildTexts("div.sc-dec7a8b-1")
		var middleStr int64 = int64(math.Round(float64(len(birthDayStr) / 2)))
		fmt.Println(birthDayStr[middleStr:])

		// tmpProfile.Bio = strings.TrimSpace(e.ChildText("#name-bio-text > div.name-trivia-bio-text > div.inline"))

		// e.ForEach("div.knownfor-title", func(_ int, kf *colly.HTMLElement) {
		// 	tmpMovie := movie{}
		// 	tmpMovie.Title = kf.ChildText("div.knownfor-title-role > a.knownfor-ellipsis")
		// 	tmpMovie.Year = kf.ChildText("div.knownfor-year > span.knownfor-ellipsis")
		// 	tmpProfile.TopMovies = append(tmpProfile.TopMovies, tmpMovie)
		// })

		js, err := json.MarshalIndent(tmpProfile, "", "    ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(js))
	})

	c.Visit("https://www.imdb.com/search/name/?birth_monthday=12-20")
	// infoCollector.Visit("https://www.imdb.com/search/name/?birth_monthday=12-20")
}
