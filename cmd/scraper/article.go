package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	krautreporter "github.com/metalmatze/krautreporter-api"
)

func (s Scraper) scrapeArticle(a *krautreporter.Article) error {
	resp, err := s.get("articles", s.host+a.URL)
	if err != nil {
		return err
	}

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
		return fmt.Errorf("couldn't parse article's author id, article: %s, author: %s", a.URL, authorURL)
	}

	// ID
	authorID, err := strconv.Atoi(idMatches[1])
	if err != nil {
		return fmt.Errorf("couldn't parse article's author id, article: %s, author: %s", a.URL, authorURL)
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

	fmt.Printf("%+v\n", a.Images)

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
