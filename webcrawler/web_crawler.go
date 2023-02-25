package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

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

	c.OnHTML("a.lister-page-next", func(e *colly.HTMLElement) {
		nextPage := e.Request.AbsoluteURL(e.Attr("href"))
		fmt.Println(nextPage)
		c.Visit(nextPage)
	})

	c.OnHTML(".mode-detail", func(e *colly.HTMLElement) {
		profileURL := e.ChildAttr("div.lister-item-image > a", "href")
		profileURL = e.Request.AbsoluteURL(profileURL)
		fmt.Println(profileURL)
		infoCollector.Visit(profileURL)
	})

	c.Visit("https://www.imdb.com/search/name/?birth_monthday=12-20")
}
