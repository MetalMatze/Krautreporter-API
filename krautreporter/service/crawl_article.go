package service

import (
	"log"
	"strings"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/PuerkitoBio/goquery"
)

func CrawlArticle(a *entity.Article) error {
	doc, err := goquery.NewDocument(mainURL + a.URL)
	if err != nil {
		log.Println("Failed to fetch %s", a.URL)
		return err
	}

	node := doc.Find("main article.article.article--full")
	nodeHeader := node.Find("header.article__header")
	nodeContent := node.Find(".article__content")

	date, err := time.Parse("02.01.2006", nodeHeader.Find("h2.meta").Text())
	if err != nil {
		return err
	}

	excerpt := strings.TrimSpace(nodeContent.Find("h2.gamma").Text())

	html, err := nodeContent.Html()
	if err != nil {
		return err
	}

	a.Date = date
	a.Headline = strings.TrimSpace(nodeHeader.Find("h1.article__title").Text())
	a.Excerpt = excerpt
	a.Content = html

	a.Crawl.Next = time.Now().Add(6 * time.Hour)

	return nil
}
