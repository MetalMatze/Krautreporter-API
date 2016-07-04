package service

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"strings"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-kit/kit/log"
)

const mainURL string = "https://krautreporter.de"
const moreURL string = "https://krautreporter.de/articles%s/load_more_navigation_items"

var articleSrcsetRegex = regexp.MustCompile(`(.*) 300w, (.*) 600w, (.*) 1000w, (.*) 2000w`)

func SyncArticles(logger log.Logger) ([]*entity.Article, error) {
	start := time.Now()

	url := mainURL
	articles := []*entity.Article{}
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
				logger.Log("msg", "Error parsing article", "err", err)
			}

			articles = append(articles, &article)
		})

		latestArticle := articles[len(articles)-1]
		url = fmt.Sprintf(moreURL, latestArticle.URL)

		logger.Log("msg", "Synced articles",
			"count", len(articles),
			"duration", time.Since(startNext),
			"next", latestArticle.URL,
		)
	}

	logger.Log("msg", "Synced articles", "count", len(articles), "duration", time.Since(start))

	return articles, nil
}

func parseArticle(s *goquery.Selection) (entity.Article, error) {
	article := entity.Article{}

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

	if srcset, exists := s.Find("img").Attr("srcset"); exists {
		matches := articleSrcsetRegex.FindStringSubmatch(srcset)
		if len(matches) == 5 {
			article.AddImage(entity.Image{Width: 300, Src: matches[1]})
			article.AddImage(entity.Image{Width: 600, Src: matches[2]})
			article.AddImage(entity.Image{Width: 1000, Src: matches[3]})
			article.AddImage(entity.Image{Width: 2000, Src: matches[4]})
		}
	}

	article.ID = id
	article.URL = strings.TrimSpace(url)
	article.Title = strings.TrimSpace(s.Find(".item__title").Text())
	article.Preview = preview
	article.Author = &entity.Author{Name: strings.TrimSpace(s.Find(".meta").Text())}

	return article, nil
}
