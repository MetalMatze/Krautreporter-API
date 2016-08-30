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
	idRegex     = regexp.MustCompile(`\/(v2\/)?(\d*)`)
	srcsetRegex = regexp.MustCompile(`(.*) 300w, (.*) 600w, (.*) 1000w, (.*) 2000w`)
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
	host string
	db   *gorm.DB
}

func (crawler Crawler) indexCommand(c *cli.Context) error {
	var articles []TeaserArticle

	for {
		page := 1
		for {
			url := crawler.host + "/api/articles"
			if page > 1 {
				url = fmt.Sprintf("%s?page=%d", url, page)
			}

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

		if err := crawler.SaveArticles(articles); err != nil {
			log.Println(err)
		}

		time.Sleep(15 * time.Second)
	}

	return nil
}

// SaveArticles takes a slice of TeaserArticles and saves them to the db
func (crawler Crawler) SaveArticles(articles []TeaserArticle) error {
	tx := crawler.db.Begin()
	for i, a := range articles {
		article, err := a.Parse()
		if err != nil {
			return err
		}

		article.Ordering = len(articles) - i - 1

		if article.Crawl.ID == 0 {
			article.Crawl = entity.Crawl{Next: time.Now()}
		}

		tx.Preload("Images").Preload("Crawl").FirstOrCreate(&article)
		tx.Save(&article)
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
	id, err := strconv.Atoi(idRegex.FindStringSubmatch(URL)[2])
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
	var images []entity.Image
	if imageNode.Length() > 0 { // preview available if img node exists
		preview = true

		srcset, _ := imageNode.Attr("srcset")
		images, err = ParseImages(srcset)
		if err != nil {
			return nil, err
		}
	}

	article := &entity.Article{
		ID:      id,
		Title:   title,
		Date:    date,
		Preview: preview,
		URL:     URL,
	}

	for _, i := range images {
		article.AddImage(i)
	}

	return article, nil
}

// ParseImages takes a string with srcset and returns a slice of Images
func ParseImages(srcset string) ([]entity.Image, error) {
	var images []entity.Image

	matches := srcsetRegex.FindStringSubmatch(srcset)
	if len(matches) == 5 {
		images = append(images, entity.Image{Width: 300, Src: matches[1]})
		images = append(images, entity.Image{Width: 600, Src: matches[2]})
		images = append(images, entity.Image{Width: 1000, Src: matches[3]})
		images = append(images, entity.Image{Width: 2000, Src: matches[4]})
	}

	return images, nil
}
