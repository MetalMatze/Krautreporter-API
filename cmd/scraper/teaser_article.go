package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	krautreporter "github.com/metalmatze/krautreporter-api"
)

var (
	ixSrcRegex = regexp.MustCompile(`.*\?`) // Match everything to the first ?
)

type (
	// TeaserArticle is just the teaser part of an article
	TeaserArticle struct {
		HTML string `json:"teaser_html"`
	}
	// TeaserArticleResponse is a http json response with TeaserArticles
	TeaserArticleResponse struct {
		Articles []TeaserArticle `json:"articles"`
	}
)

// Parse a TeaserArticle and return every data for an Article
func (ta TeaserArticle) Parse() (*krautreporter.Article, error) {
	article := &krautreporter.Article{}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(ta.HTML))
	if err != nil {
		return article, err
	}

	cardNode := doc.Find("article a")

	// Title
	title := cardNode.Find("h2").Text()
	title = strings.TrimPrefix(title, "\\n")
	article.Title = strings.TrimSpace(title)

	// URL
	URL, exists := cardNode.Attr("href")
	if !exists {
		return article, fmt.Errorf("card node has no href attr: %s", article.Title)
	}
	article.URL = URL

	// ID
	idMatches := idRegex.FindStringSubmatch(URL)
	if len(idMatches) != 2 {
		return article, fmt.Errorf("couldn't find id in %s", URL)
	}

	id, err := strconv.Atoi(idMatches[1])
	if err != nil {
		return article, fmt.Errorf("couldn't find id in %s", URL)
	}
	article.ID = id

	// Date
	dateString, exists := cardNode.Find("time").Attr("datetime")
	if !exists {
		return article, fmt.Errorf("no date for %d found", id)
	}
	// TODO: Parse with Europe/Berlin location timezone
	date, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return article, fmt.Errorf("couldn't parse date for %d: %s", id, err)
	}
	article.Date = date

	imageNode := doc.Find("img")
	if imageNode.Length() > 0 { // preview available if img node exists
		article.Preview = true

		ixSrc, exists := imageNode.Attr("ix-src")
		if !exists {
			return article, fmt.Errorf("article img has no ix-src attr")
		}

		src := ixSrcRegex.FindString(ixSrc)
		article.AddImage(krautreporter.Image{
			Width: 1600,
			Src:   src + "w=1600",
		})
	}

	return article, nil
}
