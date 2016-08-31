package main

import (
	"encoding/json"
	"errors"
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
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/urfave/cli"
)

const (
	indexInterval = time.Minute
	crawlInterval = 5 * time.Minute
)

var (
	idRegex     = regexp.MustCompile(`\/(\d*)`)
	srcsetRegex = regexp.MustCompile(`(.*) 300w, (.*) 600w, (.*) 1000w, (.*) 2000w`)
)

func main() {
	db, err := gorm.Open("postgres", "postgres://postgres:postgres@localhost:54320/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	c := Scraper{
		host: os.Getenv("HOST"),
		db:   db,
	}

	app := cli.NewApp()
	app.Name = "scraper"
	app.Usage = "Index & crawl krautreporter.de"

	app.Commands = []cli.Command{
		{
			Name:   "index",
			Usage:  "Index all articles to start crawling",
			Action: c.indexCommand,
		},
		{
			Name:   "crawl",
			Usage:  "Crawl to get all missing data",
			Action: c.crawlCommand,
		},
	}

	app.Run(os.Args)
}

type Scraper struct {
	host string
	db   *gorm.DB
}

func (scraper Scraper) indexCommand(c *cli.Context) error {
	for {
		if err := scraper.index(); err != nil {
			return err
		}
		time.Sleep(indexInterval)
	}

	return nil
}

func (scraper Scraper) crawlCommand(c *cli.Context) error {

	articleChan := make(chan *entity.Article, 1024)

	for i := 0; i < 10; i++ {
		go func() {
			for a := range articleChan {
				scraper.crawlArticle(a)
			}
		}()
	}

	go func() {
		for {
			var crawls []*entity.Crawl
			scraper.db.Where("next < ?", time.Now()).Where("crawlable_type = ?", "articles").Order("next").Find(&crawls)

			var IDs []int
			for _, c := range crawls {
				IDs = append(IDs, c.CrawlableID)
			}

			var articles []*entity.Article
			scraper.db.Preload("Crawl").Where(IDs).Find(&articles)

			for _, a := range articles {
				articleChan <- a
			}

			time.Sleep(crawlInterval)
		}
	}()

	scraper.indexCommand(c) // blocks and crawling works in a goroutine

	return nil
}

func (scraper Scraper) index() error {
	var articles []TeaserArticle

	page := 1
	for {
		url := scraper.host + "/api/articles"
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

	if err := scraper.SaveArticles(articles); err != nil {
		log.Println(err)
	}

	return nil
}

// SaveArticles takes a slice of TeaserArticles and saves them to the db
func (scraper Scraper) SaveArticles(articles []TeaserArticle) error {
	tx := scraper.db.Begin()
	for i, a := range articles {
		article, err := a.Parse()
		if err != nil {
			log.Println(err)
			continue
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
	idMatches := idRegex.FindStringSubmatch(URL)
	if len(idMatches) != 2 {
		return nil, fmt.Errorf("couldn't find id in %s", URL)
	}

	id, err := strconv.Atoi(idMatches[1])
	if err != nil {
		return nil, fmt.Errorf("couldn't find id in %s", URL)
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

func (scraper Scraper) crawlArticle(a *entity.Article) {
	log.Println(scraper.host + a.URL)
	doc, err := goquery.NewDocument(scraper.host + a.URL)
	if err != nil {
		log.Println(err)
	}

	articleNode := doc.Find("main article.article")
	contentNode := articleNode.Find(".article-content")

	if articleNode.Length() == 0 {
		log.Printf("article %s has no content", a.URL)
	}

	contentHTML, err := contentNode.Html()
	if err != nil {
		log.Println(err)
	}

	a.Headline = strings.TrimSpace(articleNode.Find(".article--title").Text())
	a.Excerpt = strings.TrimSpace(contentNode.Find(".article--teaser").Text())
	a.Content = strings.TrimSpace(contentHTML)

	authorNode := articleNode.Find(".author .author--link")
	authorURL, _ := authorNode.Attr("href")
	authorName := strings.TrimSpace(authorNode.Text())

	// ID
	authorID, err := strconv.Atoi(idRegex.FindStringSubmatch(authorURL)[1])
	if err != nil {
		log.Printf("couldn't parse id for author %s\n", authorURL)
	}

	author := entity.Author{
		ID:   authorID,
		Name: authorName,
		URL:  authorURL,
	}
	scraper.db.Preload("Images").Preload("Crawl").FirstOrCreate(&author)

	if author.Crawl.ID == 0 {
		author.Crawl = entity.Crawl{Next: time.Now()}
	}

	scraper.db.Save(&author)

	a.AuthorID = author.ID
	scraper.db.Save(&a)
}
