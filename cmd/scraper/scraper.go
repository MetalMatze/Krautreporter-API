package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
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

type Scrape interface {
	Type() string
	Fetch() (*goquery.Document, error)
	Parse(*goquery.Document) error
	Save() error
}

// Scraper knows the host to scrape and connects to the database
type Scraper struct {
	db *gorm.DB

	host   string
	client *http.Client
}

func (s *Scraper) get(types, url string) (*http.Response, error) {
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
func (s *Scraper) ActionIndex(c *cli.Context) error {
	log.Printf("starting, commit: %s\n", BuildCommit)

	return s.runIndex()
}

// ActionCrawl runs and endless loop of crawls and indexes in a certain interval
func (s *Scraper) ActionCrawl(c *cli.Context) error {
	log.Printf("starting, commit: %s\n", BuildCommit)

	//go s.runIndex()
	go s.runCrawl()

	select {}
}

func (s *Scraper) runIndex() error {
	for {
		if err := s.index(); err != nil {
			return err
		}
		indexCounter.Inc()
		time.Sleep(indexInterval)
	}
}

func (s *Scraper) runCrawl() error {
	scrapeChan := make(chan Scrape, 2000)

	for i := 0; i < 10; i++ {
		go s.Scrape(scrapeChan)
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
			scrapeChan <- &ScrapeArticle{
				Scraper: s,
				Article: a,
			}
		}

		//// authors
		//s.db.Where("next < ?", time.Now()).
		//	Where("crawlable_type = ?", "authors").
		//	Order("next").
		//	Find(&crawls)
		//
		//for _, c := range crawls {
		//	IDs = append(IDs, c.CrawlableID)
		//}
		//
		//var authors []*krautreporter.Author
		//s.db.Preload("Crawl").
		//	Where(IDs).
		//	Find(&authors)
		//
		//for _, a := range authors {
		//	authorChan <- a
		//}

		time.Sleep(crawlInterval)
	}
}

func (s *Scraper) Scrape(scrapeChan <-chan Scrape) {
	for scrape := range scrapeChan {
		doc, err := scrape.Fetch()
		if err != nil {
			crawlCounter.WithLabelValues(scrape.Type(), "error").Inc()
			log.Println(err)
		}

		if err := scrape.Parse(doc); err != nil {
			crawlCounter.WithLabelValues(scrape.Type(), "error").Inc()
			log.Println(err)
		}

		if err := scrape.Save(); err != nil {
			crawlCounter.WithLabelValues(scrape.Type(), "error").Inc()
			log.Println(err)
		}

		crawlCounter.WithLabelValues(scrape.Type(), "success").Inc()
	}
}

func (s *Scraper) index() error {
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

	//if err := s.SaveArticles(articles); err != nil {
	//	log.Println(err)
	//}

	return nil
}
