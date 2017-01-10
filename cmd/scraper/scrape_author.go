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

func (s Scraper) fetchAuthor(url string) (*goquery.Document, error) {
	resp, err := s.get("authors", url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("scraping %s returned %d", url, resp.StatusCode)
	}

	return goquery.NewDocumentFromResponse(resp)
}

func (s Scraper) parseAuthor(doc *goquery.Document) (*krautreporter.Author, error) {
	author := &krautreporter.Author{}

	authorNode := doc.Find("main .island .author")

	author.Name = strings.TrimSpace(authorNode.Find("h1.author--name").Text())
	author.Biography = strings.TrimSpace(authorNode.Find("p").First().Text())

	smHTML, err := authorNode.Find("p.meta").Html()
	if err != nil {
		return author, err
	}
	author.SocialMedia = strings.TrimSpace(smHTML)

	imgNode := authorNode.Find("img.author__img")
	if imgNode.Length() > 0 {
		srcset, exists := imgNode.Attr("srcset")
		if !exists {
			return author, fmt.Errorf("author has img but no srcset: %s", author.Name)
		}

		matches := srcsetRegex.FindStringSubmatch(srcset)
		if len(matches) != 3 {
			return author, fmt.Errorf("author has img with srcset, but more than 2: %s", author.Name)
		}

		author.AddImage(krautreporter.Image{Width: 170, Src: matches[1]})
		author.AddImage(krautreporter.Image{Width: 340, Src: matches[2]})
	}

	return author, nil
}

func (s Scraper) nextCrawlAuthor(a *krautreporter.Author) {
	constant := 5 * time.Hour
	variable := 30 * time.Minute
	random := rand.Intn(int(variable.Seconds()))

	a.Crawl.Next = time.Now().Add(time.Duration(constant.Seconds() + float64(random)))
}
