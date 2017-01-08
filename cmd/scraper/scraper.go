package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	krautreporter "github.com/metalmatze/krautreporter-api"
	"github.com/urfave/cli"
)

const (
	indexInterval = 5 * time.Minute
	crawlInterval = indexInterval
)

var (
	idRegex = regexp.MustCompile(`\/(\d*)`)
)

// Scraper knows the host to scrape and connects to the database
type Scraper struct {
	db *gorm.DB

	host   string
	client *http.Client
}

func (s Scraper) get(types, url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "github.com/metalmatze/krautreporter-api/1.0.0")

	start := time.Now()
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}

	duration := time.Since(start).Seconds()
	scrapeTimeHistogram.WithLabelValues(types).Observe(duration)

	log.Println(resp.Status, url)

	return resp, err
}

// ActionIndex runs and endless loop of indexing in a certain interval
func (s Scraper) ActionIndex(c *cli.Context) error {
	log.Printf("starting, commit: %s\n", BuildCommit)

	return s.runIndex()
}

// ActionCrawl runs and endless loop of crawls and indexes in a certain interval
func (s Scraper) ActionCrawl(c *cli.Context) error {
	log.Printf("starting, commit: %s\n", BuildCommit)

	//go s.runIndex()
	go s.runCrawl()

	select {}
}

func (s Scraper) runIndex() error {
	for {
		if err := s.index(); err != nil {
			return err
		}
		indexCounter.Inc()
		time.Sleep(indexInterval)
	}
}

func (s Scraper) runCrawl() error {
	articleChan := make(chan *krautreporter.Article, 1024)
	authorChan := make(chan *krautreporter.Author, 1024)

	for i := 0; i < 10; i++ {
		go func() {
			for a := range articleChan {
				if err := s.scrapeArticle(a); err != nil {
					log.Println(err)
					crawlCounter.WithLabelValues("articles", "error").Inc()
				}
				crawlCounter.WithLabelValues("articles", "success").Inc()
			}
		}()
	}

	for i := 0; i < 10; i++ {
		go func() {
			for a := range authorChan {
				doc, err := s.fetchAuthor(s.host + a.URL)
				if err != nil {
					log.Println(err)
					crawlCounter.WithLabelValues("authors", "error").Inc()
					continue
				}

				author, err := s.parseAuthor(doc)
				if err != nil {
					log.Println(err)
					crawlCounter.WithLabelValues("authors", "error").Inc()
					continue
				}

				s.nextCrawlAuthor(author)

				//TODO: Save the new author

				crawlCounter.WithLabelValues("authors", "success").Inc()
			}
		}()
	}

	for {
		// articles
		var crawls []*krautreporter.Crawl
		s.db.Where("next < ?", time.Now()).
			Where("crawlable_type = ?", "articles").
			Order("next").
			Find(&crawls)

		var IDs []int
		for _, c := range crawls {
			IDs = append(IDs, c.CrawlableID)
		}

		var articles []*krautreporter.Article
		s.db.Preload("Crawl").
			Where(IDs).
			Find(&articles)

		for _, a := range articles {
			articleChan <- a
		}

		// authors
		s.db.Where("next < ?", time.Now()).
			Where("crawlable_type = ?", "authors").
			Order("next").
			Find(&crawls)

		for _, c := range crawls {
			IDs = append(IDs, c.CrawlableID)
		}

		var authors []*krautreporter.Author
		s.db.Preload("Crawl").
			Where(IDs).
			Find(&authors)

		for _, a := range authors {
			authorChan <- a
		}

		time.Sleep(crawlInterval)
	}
}

func (s Scraper) index() error {
	var articles []TeaserArticle

	page := 1
	for {
		url := s.host + "/api/articles"
		if page > 1 {
			url = fmt.Sprintf("%s?page=%d", url, page)
		}

		resp, err := s.get("teaser_articles", url)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			log.Printf("request for %s returned %d", url, resp.StatusCode)
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(err)
			}
			log.Println(string(body))
			resp.Body.Close()

			continue
		}

		var articleResp TeaserArticleResponse
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&articleResp)
		if err != nil {
			return err
		}
		resp.Body.Close()

		if len(articleResp.Articles) == 0 {
			break
		}

		articles = append(articles, articleResp.Articles...)
		page++
	}

	if err := s.SaveArticles(articles); err != nil {
		log.Println(err)
	}

	return nil
}
