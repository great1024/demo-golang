package main

import (
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

// 博客爬虫
func main() {
	domain := ""
	run(domain)
}

func run(domain string) {
	res, err := http.Get(domain)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	print(doc)
	//allUrl := doc.Find("a")
	//for i := 0; i < len(allUrl); i++ {
	//
	//}
}
