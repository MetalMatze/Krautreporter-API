package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	krautreporter "github.com/metalmatze/krautreporter-api"
	"github.com/metalmatze/krautreporter-api/repository"
	"github.com/urfave/cli"
)

const (
	indexInterval = 5 * time.Minute
	crawlInterval = indexInterval
)

var (
	idRegex = regexp.MustCompile(`\/(\d*)`)
)

// Scrape interface makes sure implementations are usable by the pipeline
type Scrape interface {
	Type() string
	Fetch() (*goquery.Document, error)
	Parse(*goquery.Document) error
	Save() error
}

// Scraper knows the host to scrape and connects to the database
type Scraper struct {
	host       string
	client     *http.Client
	Repository *repository.Repository
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

	return s.runCrawl()
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
		go s.scrape(scrapeChan)
	}

	for {
		if err := s.index(); err != nil {
			return err
		}

		authors, err := s.Repository.FindOutdatedAuthors()
		if err != nil {
			return err
		}

		for _, a := range authors {
			scrapeChan <- &ScrapeAuthor{
				Scraper: s,
				Author:  a,
			}
		}

		articles, err := s.Repository.FindOutdatedArticles()
		if err != nil {
			return err
		}

		for _, a := range articles {
			scrapeChan <- &ScrapeArticle{
				Scraper: s,
				Article: a,
			}
		}

		time.Sleep(crawlInterval)
	}
}

func (s *Scraper) scrape(scrapeChan <-chan Scrape) {
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
	var teaserArticles []TeaserArticle

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

		teaserArticles = append(teaserArticles, articleResp.Articles...)
		page++
	}

	articles := []*krautreporter.Article{}
	var successCount, errorCount float64
	for i, ta := range teaserArticles {
		if len(ta.HTML) == 0 {
			errorCount++
			continue
		}

		article, err := ta.Parse()
		if err != nil {
			log.Println(err)
			errorCount++
			continue
		}

		article.Ordering = len(teaserArticles) - i - 1

		if article.Crawl == nil {
			article.Crawl = &krautreporter.Crawl{Next: time.Now()}
		}

		articles = append(articles, article)
		successCount++
	}

	indexArticleGauge.WithLabelValues("success").Set(successCount)
	indexArticleGauge.WithLabelValues("error").Set(errorCount)

	if err := s.Repository.SaveAllArticles(articles); err != nil {
		return err
	}
	log.Printf("Saved %d indexed articles", len(articles))

	return nil
}
