package service

import (
	"errors"
	"log"
	"strconv"

	"fmt"
	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/PuerkitoBio/goquery"
)

const mainURL string = "https://krautreporter.de"
const moreURL string = "https://krautreporter.de/articles%s/load_more_navigation_items"

func CrawlArticles() ([]entity.Article, error) {
	url := mainURL
	articles := []entity.Article{}
	for {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			return nil, err
		}

		selector := "li a"
		if url == mainURL {
			selector = "#article-list-tab " + selector
		}

		// No more articles found, stop crawling and return with found articles
		if doc.Find(selector).Length() <= 0 {
			return articles, nil
		}

		doc.Find(selector).Each(func(i int, s *goquery.Selection) {
			article, err := parseArticle(s)
			if err != nil {
				log.Println(err)
			}

			articles = append(articles, article)
		})

		latestArticle := articles[len(articles)-1]
		url = fmt.Sprintf(moreURL, latestArticle.URL)
		log.Printf("Crawled %d articles, requesting more for %s", len(articles), latestArticle.URL)
	}

	return articles, nil
}

func parseArticle(s *goquery.Selection) (entity.Article, error) {
	url, exists := s.Attr("href")
	if !exists {
		return entity.Article{}, errors.New("Can't find the url for an article")
	}

	id, err := strconv.Atoi(IDRegex.FindStringSubmatch(url)[1])
	if err != nil {
		return entity.Article{}, err
	}

	preview := false
	if s.Find("img").Length() > 1 {
		preview = true
	}

	return entity.Article{
		ID:      id,
		URL:     url,
		Title:   s.Find(".item__title").Text(),
		Preview: preview,
		Author: &entity.Author{
			Name: s.Find(".meta").Text(),
		},
	}, nil
}
