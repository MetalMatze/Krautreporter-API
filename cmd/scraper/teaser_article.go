package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/metalmatze/krautreporter-api/entity"
)

type (
	// TeaserArticle is just the teaser part of an article
	TeaserArticle struct {
		TeaserHTML string `json:"teaser_html"`
	}
	// TeaserArticleResponse is a http json response with TeaserArticles
	TeaserArticleResponse struct {
		Articles []TeaserArticle `json:"articles"`
	}
)

// Parse a TeaserArticle and return every data for an Article
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
		images, err = ParseArticleImages(srcset)
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
