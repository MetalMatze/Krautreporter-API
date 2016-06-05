package service

import (
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/PuerkitoBio/goquery"
	"log"
	"strings"
	"time"
)

func CrawlArticle(a entity.Article) (entity.Article, error) {
	doc, err := goquery.NewDocument("https://krautreporter.de" + a.URL)
	if err != nil {
		log.Println("Failed to fetch %s", a.URL)
		return a, err
	}

	node := doc.Find("main article.article.article--full")
	nodeHeader := node.Find("header.article__header")
	nodeContent := node.Find(".article__content")

	date, err := time.Parse("02.01.2006", nodeHeader.Find("h2.meta").Text())
	if err != nil {
		return a, err
	}

	excerpt := strings.TrimSpace(nodeContent.Find("h2.gamma").Text())

	html, err := nodeContent.Html()
	if err != nil {
		return a, err
	}

	a.Date = date
	a.Headline = strings.TrimSpace(nodeHeader.Find("h1.article__title").Text())
	a.Excerpt = excerpt
	a.Content = html

	return a, nil
}
