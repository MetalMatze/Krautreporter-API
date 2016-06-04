package service

import (
	"log"
	"strings"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/PuerkitoBio/goquery"
)

func CrawlAuthor(a entity.Author) (entity.Author, error) {
	doc, err := goquery.NewDocument("https://krautreporter.de" + a.URL)
	if err != nil {
		log.Println("Failed to fetch %s", a.URL)
		return entity.Author{}, err
	}

	doc.Find("header.article__header").Each(func(i int, s *goquery.Selection) {
		a.Biography = strings.TrimSpace(s.Find(".author__bio").Text())

		html, err := s.Find("#author-page--media-links").Html()
		if err == nil {
			a.SocialMedia = html
		}
	})

	a.Crawl.Next = time.Now().Add(6 * time.Hour)

	return a, nil
}
