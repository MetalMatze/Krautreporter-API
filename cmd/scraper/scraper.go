package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	idRegex            = regexp.MustCompile(`\/(\d*)`)
	articleSrcsetRegex = regexp.MustCompile(`(.*) 300w, (.*) 600w, (.*) 1000w, (.*) 2000w`)
	authorSrcsetRegex  = regexp.MustCompile(`(.*) 170w, (.*) 340w`)
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

func (s Scraper) ActionIndex(c *cli.Context) error {
	log.Printf("starting, commit: %s\n", BuildCommit)

	return s.runIndex()
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

func (s Scraper) ActionCrawl(c *cli.Context) error {
	log.Printf("starting, commit: %s\n", BuildCommit)

	go s.runCrawl()
	go s.runIndex()

	select {}
}

func (s Scraper) runCrawl() error {
	articleChan := make(chan *krautreporter.Article, 1024)
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

	authorChan := make(chan *krautreporter.Author, 1024)
	for i := 0; i < 10; i++ {
		go func() {
			for a := range authorChan {
				if err := s.scrapeAuthor(a); err != nil {
					crawlCounter.WithLabelValues("authors", "error").Inc()
				}
				crawlCounter.WithLabelValues("authors", "success").Inc()
			}
		}()
	}

	for {
		// articles
		var crawls []*krautreporter.Crawl
		s.db.Where("next < ?", time.Now()).Where("crawlable_type = ?", "articles").Order("next").Find(&crawls)

		var IDs []int
		for _, c := range crawls {
			IDs = append(IDs, c.CrawlableID)
		}

		var articles []*krautreporter.Article
		s.db.Preload("Crawl").Where(IDs).Find(&articles)

		for _, a := range articles {
			articleChan <- a
		}

		// authors
		s.db.Where("next < ?", time.Now()).Where("crawlable_type = ?", "authors").Order("next").Find(&crawls)

		for _, c := range crawls {
			IDs = append(IDs, c.CrawlableID)
		}

		var authors []*krautreporter.Author
		s.db.Preload("Crawl").Where(IDs).Find(&authors)

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

// SaveArticles takes a slice of TeaserArticles and saves them to the db
func (s Scraper) SaveArticles(articles []TeaserArticle) error {
	var successCount, errorCount float64
	tx := s.db.Begin()
	for i, a := range articles {
		article, err := a.Parse()
		if err != nil {
			log.Println(err)
			errorCount++
			continue
		}

		article.Ordering = len(articles) - i - 1

		if article.Crawl.ID == 0 {
			article.Crawl = krautreporter.Crawl{Next: time.Now()}
		}

		tx.Preload("Images").Preload("Crawl").FirstOrCreate(&article)
		tx.Save(&article)
		successCount++
	}
	tx.Commit()

	indexArticleGauge.WithLabelValues("success").Set(successCount)
	indexArticleGauge.WithLabelValues("error").Set(errorCount)

	return nil
}

// ParseArticleImages takes a string with srcset and returns a slice of Images
func ParseArticleImages(srcset string) ([]krautreporter.Image, error) {
	var images []krautreporter.Image

	matches := articleSrcsetRegex.FindStringSubmatch(srcset)
	if len(matches) == 5 {
		images = append(images, krautreporter.Image{Width: 300, Src: matches[1]})
		images = append(images, krautreporter.Image{Width: 600, Src: matches[2]})
		images = append(images, krautreporter.Image{Width: 1000, Src: matches[3]})
		images = append(images, krautreporter.Image{Width: 2000, Src: matches[4]})
	}

	return images, nil
}

func (s Scraper) scrapeArticle(a *krautreporter.Article) error {
	resp, err := s.get("articles", s.host+a.URL)
	if err != nil {
		return err
	}

	log.Println(resp.Status, s.host+a.URL)

	if resp.StatusCode != 200 {
		return fmt.Errorf("scraping %s returned %d", a.URL, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	articleNode := doc.Find("main article.article")
	contentNode := articleNode.Find(".article-content")

	if articleNode.Length() == 0 {
		log.Printf("article %s has no content", a.URL)
	}

	contentHTML, err := contentNode.Html()
	if err != nil {
		return err
	}

	a.Headline = strings.TrimSpace(articleNode.Find(".article--title").Text())
	a.Excerpt = strings.TrimSpace(contentNode.Find(".article--teaser").Text())
	a.Content = strings.TrimSpace(contentHTML)

	authorNode := articleNode.Find(".author .author--link")
	authorURL, exists := authorNode.Attr("href")
	authorName := strings.TrimSpace(authorNode.Text())

	if !exists {
		return fmt.Errorf("author link doesn't exist for %s", a.URL)
	}

	idMatches := idRegex.FindStringSubmatch(authorURL)
	if len(idMatches) != 2 {
		log.Printf("couldn't parse article's author id, article: %s, author: %s\n", a.URL, authorURL)
	}

	// ID
	authorID, err := strconv.Atoi(idMatches[1])
	if err != nil {
		log.Printf("couldn't parse article's author id, article: %s, author: %s\n", a.URL, authorURL)
	}

	author := krautreporter.Author{
		ID:   authorID,
		Name: authorName,
		URL:  authorURL,
	}
	s.db.Preload("Images").Preload("Crawl").FirstOrCreate(&author)

	if author.Crawl.ID == 0 {
		author.Crawl = krautreporter.Crawl{Next: time.Now()}
	}

	s.db.Save(&author)

	a.Crawl.Next = time.Now().Add(time.Duration(float64(rand.Intn(18000))+30*time.Minute.Seconds()) * time.Second)
	a.AuthorID = author.ID
	s.db.Save(&a)

	return nil
}

func (s Scraper) scrapeAuthor(a *krautreporter.Author) error {
	resp, err := s.get("authors", s.host+a.URL)
	if err != nil {
		return err
	}
	log.Println(resp.Status, s.host+a.URL)

	if resp.StatusCode != 200 {
		return fmt.Errorf("scraping %s returned %d", a.URL, resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	s.db.Preload("Images").Preload("Crawl").FirstOrCreate(&a)

	authorNode := doc.Find("main .island .author")
	imageNode := authorNode.Find("img.author__img")

	a.Biography = strings.TrimSpace(authorNode.Find("p").First().Text())

	html, err := authorNode.Find("p.meta").Html()
	if err != nil {
		return err
	}

	a.SocialMedia = strings.TrimSpace(html)
	var images []krautreporter.Image
	if imageNode.Length() > 0 {
		srcset, _ := imageNode.Attr("srcset")
		images, err = ParseAuthorImages(srcset)
		if err != nil {
			return err
		}
	}

	for _, i := range images {
		a.AddImage(i)
	}

	a.Crawl.Next = time.Now().Add(time.Duration(float64(rand.Intn(18000))+30*time.Minute.Seconds()) * time.Second)
	s.db.Save(&a)

	return nil
}

// ParseAuthorImages takes a string with srcset and returns a slice of Images
func ParseAuthorImages(srcset string) ([]krautreporter.Image, error) {
	var images []krautreporter.Image

	matches := authorSrcsetRegex.FindStringSubmatch(srcset)
	if len(matches) == 3 {
		images = append(images, krautreporter.Image{Width: 170, Src: matches[1]})
		images = append(images, krautreporter.Image{Width: 340, Src: matches[2]})
	}

	return images, nil
}
