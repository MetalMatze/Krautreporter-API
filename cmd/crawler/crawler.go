package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-errors/errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/urfave/cli"
)

var (
	IDRegex = regexp.MustCompile(`\/(v2\/)?(\d*)`)
)

func main() {
	db, err := gorm.Open("postgres", "postgres://postgres:postgres@localhost:54320/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	c := Crawler{
		host: os.Getenv("HOST"),
		db:   db,
	}

	app := cli.NewApp()
	app.Name = "crawler"
	app.Usage = "Crawls krautreporter.de and saves to db"

	app.Commands = []cli.Command{{
		Name:   "index",
		Usage:  "Only crawls to build an index",
		Action: c.indexCommand,
	}}

	app.Run(os.Args)
}

type Crawler struct {
	host    string
	db      *gorm.DB
	metrics crawlerMetrics
}

func (crawler Crawler) indexCommand(c *cli.Context) error {
	var articles []TeaserArticle

	for {
		a, err := crawler.indexHomepage()
		if err != nil {
			return err
		}
		articles = append(articles, a...)

		page := 2
		for {
			url := fmt.Sprintf(crawler.host+"/api/articles?page=%d", page)
			log.Println(url)

			resp, err := http.Get(url)
			if err != nil {
				return err
			}
			if resp.StatusCode != 200 {
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

		time.Sleep(20 * time.Second)
	}

	return nil
}

func (crawler Crawler) indexHomepage() ([]TeaserArticle, error) {
	var articles []TeaserArticle

	log.Println(crawler.host)
	doc, err := goquery.NewDocument(crawler.host)
	if err != nil {
		return articles, err
	}
	// TODO: Need to check if really 200 OK

	layoutNode := doc.Find(".content .layout .layout__item")
	layoutNode.Each(func(i int, s *goquery.Selection) {
		if i == 9 { // the 10th node is the more button, don't include it
			return
		}

		html, err := s.Html()
		if err != nil {
			log.Println(err)
		}
		articles = append(articles, TeaserArticle{TeaserHTML: html})
	})

	return articles, nil
}

func (crawler Crawler) SaveArticles(articles []TeaserArticle) error {
	tx := crawler.db.Begin()
	for i, a := range articles {
		article, err := a.Parse()
		if err != nil {
			return err
		}
		article.Ordering = len(articles) - i - 1

		if tx := tx.FirstOrCreate(&article); tx.Error != nil {
			return err
		}
	}
	tx.Commit()

	return nil
}

type (
	TeaserArticle struct {
		TeaserHTML string `json:"teaser_html"`
	}
	TeaserArticleResponse struct {
		Articles []TeaserArticle `json:"articles"`
	}
)

func (ta TeaserArticle) Parse() (*entity.Article, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(ta.TeaserHTML))
	if err != nil {
		return nil, err
	}

	cardNode := doc.Find("article a")
	authorNode := cardNode.Find(".author__body")
	imageNode := cardNode.Find("img.card__img")

	// Title
	title := strings.TrimSpace(cardNode.Find("h2").Text())

	// URL
	URL, exists := cardNode.Attr("href")
	if !exists {
		return nil, errors.New("URL attr doesn't exist")
	}

	// ID
	id, err := strconv.Atoi(IDRegex.FindStringSubmatch(URL)[2])
	if err != nil {
		return nil, fmt.Errorf("couldn't parse id for %d:  ID attr doesn't exist", id)
	}

	// Date
	dateString, exists := authorNode.Find("time").Attr("datetime")
	if !exists {
		return nil, fmt.Errorf("no date for %d found", id)
	}
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse date for %d: %s", id, err)
	}

	// Preview
	var preview bool
	if imageNode.Length() > 0 { // preview available if img node exists
		preview = true

		imageHTML, err := imageNode.Html()
		if err != nil {
			return nil, err
		}

		fmt.Println(imageHTML)
	}

	return &entity.Article{
		ID:      id,
		Title:   title,
		Date:    date,
		Preview: preview,
		URL:     URL,
	}, nil
}
