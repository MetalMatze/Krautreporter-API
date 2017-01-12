package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	krautreporter "github.com/metalmatze/krautreporter-api"
)

var (
	srcsetRegex = regexp.MustCompile(`(.*) 170w, (.*) 340w`)
)

// ScrapeAuthor implementes the Scrape interface to scrape one specific author
type ScrapeAuthor struct {
	Scraper *Scraper
	Author  *krautreporter.Author
}

// Type returns a string representing the type of the Scrape interface implementation
func (sa *ScrapeAuthor) Type() string {
	return "authors"
}

// Fetch an author and return a goquery.Document with its content
func (sa *ScrapeAuthor) Fetch() (*goquery.Document, error) {
	resp, err := sa.Scraper.get("authors", sa.Author.URL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("scraping %s returned %d", sa.Author.URL, resp.StatusCode)
	}

	return goquery.NewDocumentFromResponse(resp)
}

// Parse a goquery.Document into the given author
func (sa *ScrapeAuthor) Parse(doc *goquery.Document) error {
	authorNode := doc.Find("main .island .author")

	sa.Author.Name = strings.TrimSpace(authorNode.Find("h1.author--name").Text())
	sa.Author.Biography = strings.TrimSpace(authorNode.Find("p").First().Text())

	smHTML, err := authorNode.Find("p.meta").Html()
	if err != nil {
		return err
	}
	sa.Author.SocialMedia = strings.TrimSpace(smHTML)

	imgNode := authorNode.Find("img.author__img")
	if imgNode.Length() > 0 {
		srcset, exists := imgNode.Attr("srcset")
		if !exists {
			return fmt.Errorf("author has img but no srcset: %s", sa.Author.Name)
		}

		matches := srcsetRegex.FindStringSubmatch(srcset)
		if len(matches) != 3 {
			return fmt.Errorf("author has img with srcset, but more than 2: %s", sa.Author.Name)
		}

		sa.Author.AddImage(&krautreporter.Image{Width: 170, Src: matches[1]})
		sa.Author.AddImage(&krautreporter.Image{Width: 340, Src: matches[2]})
	}

	return nil
}

// Save the updated author after fetching & parsing
func (sa *ScrapeAuthor) Save() error {
	sa.Author.Crawl.NextRandom()

	return sa.Scraper.Repository.SaveAuthor(sa.Author)
}
