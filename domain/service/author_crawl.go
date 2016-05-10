package service

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/PuerkitoBio/goquery"
)

var IDRegex = regexp.MustCompile(`\/(\d*)--`)

func CrawlAuthor() ([]entity.Author, error) {
	doc, err := goquery.NewDocument("https://krautreporter.de")
	if err != nil {
		return nil, err
	}

	authorNodes := doc.Find("#author-list-tab li a")
	log.Printf("Found %d authors", authorNodes.Length())

	authors := []entity.Author{}
	authorNodes.Each(func(i int, s *goquery.Selection) {
		author := entity.Author{}
		author.URL, _ = s.Attr("href")
		author.ID, _ = strconv.Atoi(IDRegex.FindStringSubmatch(author.URL)[1])
		author.Ordering = authorNodes.Length() - i - 1
		author.Name = strings.TrimSpace(s.Find(".author__name").Text())
		author.Title = s.Find(".item__title").Text()

		authors = append(authors, author)
	})

	return authors, nil
}
