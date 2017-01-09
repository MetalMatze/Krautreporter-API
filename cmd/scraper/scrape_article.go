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

type ScrapeArticle struct {
	Scraper *Scraper
	Article *krautreporter.Article
}

func (sa *ScrapeArticle) Type() string {
	return "articles"
}

func (sa *ScrapeArticle) Fetch() (*goquery.Document, error) {
	resp, err := sa.Scraper.get("articles", sa.Scraper.host+sa.Article.URL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("scraping %s returned %d", sa.Article.URL, resp.StatusCode)
	}

	return goquery.NewDocumentFromResponse(resp)
}

func (sa *ScrapeArticle) Parse(doc *goquery.Document) error {
	articleNode := doc.Find("main article.article")
	contentNode := articleNode.Find(".article-content")

	if articleNode.Length() == 0 {
		log.Printf("article %s has no content") // TODO log with ID
	}

	content, err := contentNode.Html()
	if err != nil {
		return err
	}
	sa.Article.Content = strings.TrimSpace(content)

	sa.Article.Headline = strings.TrimSpace(articleNode.Find("h1.article--title").Text())
	sa.Article.Excerpt = strings.TrimSpace(contentNode.Find(".article--teaser").Text())

	if err = sa.parseAuthor(articleNode.Find(".author .author--link")); err != nil {
		return err
	}

	return nil
}

func (sa *ScrapeArticle) parseAuthor(node *goquery.Selection) error {
	if sa.Article.Author == nil {
		sa.Article.Author = &krautreporter.Author{}
	}

	sa.Article.Author.Name = strings.TrimSpace(node.Text())

	// URL
	authorURL, exists := node.Attr("href")
	if !exists {
		return fmt.Errorf("author link doesn't exist for %s", sa.Article.Author.ID) // TODO log with ID
	}
	sa.Article.Author.URL = authorURL

	// ID
	idMatches := idRegex.FindStringSubmatch(authorURL)
	if len(idMatches) != 2 {
		return fmt.Errorf("couldn't parse article's author id, article: %s, author: %s", "", authorURL)
	}

	authorID, err := strconv.Atoi(idMatches[1])
	if err != nil {
		return fmt.Errorf("couldn't parse article's author id, article: %s, author: %s", "", authorURL)
	}
	sa.Article.Author.ID = authorID
	sa.Article.AuthorID = authorID

	return nil
}

func (sa *ScrapeArticle) nextCrawl() {
	constant := 5 * time.Hour
	variable := 30 * time.Minute
	random := rand.Intn(int(variable.Seconds()))

	sa.Article.Crawl.Next = time.Now().Add(time.Duration(constant.Seconds() + float64(random)))
}

func (sa *ScrapeArticle) Save() error {
	return nil
}
