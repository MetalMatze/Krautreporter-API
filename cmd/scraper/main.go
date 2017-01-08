package main

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/urfave/cli"
)

// Config for the scraper binary
type Config struct {
	DSN  string
	Host string
}

var (
	// BuildCommit is the git commit the binary will be build upon
	BuildCommit string

	indexCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "krautreporter_index_total",
			Help: "How often the s indexed the site",
		},
	)
	indexArticleGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "krautreporter_index_article_total",
			Help: "How often articles are successfully scraped",
		},
		[]string{"status"},
	)
	crawlCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "krautreporter_crawls_total",
			Help: "How often authors & articles are crawled",
		},
		[]string{"type", "status"},
	)
	scrapeTimeHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "krautreporter_scrape_time_seconds",
			Help:    "Scrape time in seconds",
			Buckets: []float64{0.1, 0.5, 1.0, 2.0, 3.0, 5.0, 7.0, 10.0, 15.0, 20.0, 30.0},
		},
		[]string{"type"},
	)
)

func init() {
	rand.Seed(time.Now().UnixNano())
	prometheus.MustRegister(indexCounter, indexArticleGauge, crawlCounter, scrapeTimeHistogram)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err)
	}

	config := &Config{
		DSN:  os.Getenv("DSN"),
		Host: os.Getenv("HOST"),
	}

	db, err := gorm.Open("postgres", config.DSN)
	if err != nil {
		log.Fatal(err)
	}

	c := Scraper{
		host: config.Host,
		db:   db,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	app := cli.NewApp()
	app.Name = "krautreporter-scraper"
	app.Usage = "Index & crawl krautreporter.de"

	app.Commands = []cli.Command{
		{
			Name:   "index",
			Usage:  "Index all articles to start crawling",
			Action: c.ActionIndex,
		},
		{
			Name:   "crawl",
			Usage:  "Crawl to get all missing data",
			Action: c.ActionCrawl,
		},
	}

	go func() {
		http.Handle("/metrics", prometheus.Handler())
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	app.Run(os.Args)
}
