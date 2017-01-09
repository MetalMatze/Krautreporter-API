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

func (s Scraper) fetchArticle(url string) (*goquery.Document, error) {
	resp, err := s.get("articles", url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("scraping %s returned %d", url, resp.StatusCode)
	}

	return goquery.NewDocumentFromResponse(resp)
}

func (s Scraper) parseArticle(doc *goquery.Document) (*krautreporter.Article, error) {
	article := &krautreporter.Article{}

	articleNode := doc.Find("main article.article")
	contentNode := articleNode.Find(".article-content")

	if articleNode.Length() == 0 {
		log.Printf("article %s has no content") // TODO log with ID
	}

	content, err := contentNode.Html()
	if err != nil {
		return article, err
	}
	article.Content = strings.TrimSpace(content)

	article.Headline = strings.TrimSpace(articleNode.Find("h1.article--title").Text())
	article.Excerpt = strings.TrimSpace(contentNode.Find(".article--teaser").Text())

	// Author
	author, err := s.parseArticleAuthor(articleNode.Find(".author .author--link"))
	if err != nil {
		return article, err
	}
	article.Author = author
	article.AuthorID = author.ID

	return article, nil
}

func (s Scraper) parseArticleAuthor(node *goquery.Selection) (*krautreporter.Author, error) {
	author := &krautreporter.Author{}

	author.Name = strings.TrimSpace(node.Text())

	// URL
	authorURL, exists := node.Attr("href")
	if !exists {
		return author, fmt.Errorf("author link doesn't exist for %s") // TODO log with ID
	}
	author.URL = authorURL

	// ID
	idMatches := idRegex.FindStringSubmatch(authorURL)
	if len(idMatches) != 2 {
		return author, fmt.Errorf("couldn't parse article's author id, article: %s, author: %s", "", authorURL)
	}

	authorID, err := strconv.Atoi(idMatches[1])
	if err != nil {
		return author, fmt.Errorf("couldn't parse article's author id, article: %s, author: %s", "", authorURL)
	}
	author.ID = authorID

	return author, nil
}

func (s Scraper) nextCrawlArticle(a *krautreporter.Article) {
	constant := 5 * time.Hour
	variable := 30 * time.Minute
	random := rand.Intn(int(variable.Seconds()))

	a.Crawl.Next = time.Now().Add(time.Duration(constant.Seconds() + float64(random)))
}
