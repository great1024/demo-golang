package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)

func main() {
	c := colly.NewCollector()
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		fmt.Println("Visiting", e.DOM)
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit("https://amazon.com")
}
