package service

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/PuerkitoBio/goquery"
	"github.com/gollection/gollection/log"
)

const mainURL string = "https://krautreporter.de"
const moreURL string = "https://krautreporter.de/articles%s/load_more_navigation_items"

func SyncArticles(log log.Logger) ([]entity.Article, error) {
	start := time.Now()

	url := mainURL
	articles := []entity.Article{}
	for {
		startNext := time.Now()
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
				log.Warn("Error parsing article", "err", err)
			}

			articles = append(articles, article)
		})

		latestArticle := articles[len(articles)-1]
		url = fmt.Sprintf(moreURL, latestArticle.URL)

		log.Debug("Synced articles",
			"count", len(articles),
			"duration", time.Since(startNext),
			"next", latestArticle.URL,
		)
	}

	log.Info("Synced articles", "count", len(articles), "duration", time.Since(start))

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
	if s.Find("img").Length() > 0 {
		preview = true
	}

	return entity.Article{
		ID:      id,
		URL:     strings.TrimSpace(url),
		Title:   strings.TrimSpace(s.Find(".item__title").Text()),
		Preview: preview,
		Author: &entity.Author{
			Name: s.Find(".meta").Text(),
		},
	}, nil
}
