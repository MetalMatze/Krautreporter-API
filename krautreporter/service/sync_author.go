package service

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/PuerkitoBio/goquery"
	"github.com/gollection/gollection/log"
)

var IDRegex = regexp.MustCompile(`\/(\d*)--`)
var SrcsetRegex = regexp.MustCompile(`(.*) 50w, (.*) 100w`)

func SyncAuthor(log log.Logger) ([]entity.Author, error) {
	start := time.Now()

	doc, err := goquery.NewDocument("https://krautreporter.de")
	if err != nil {
		return nil, err
	}

	authorNodes := doc.Find("#author-list-tab li a")

	authors := []entity.Author{}
	authorNodes.Each(func(i int, s *goquery.Selection) {
		author := entity.Author{}
		author.URL, _ = s.Attr("href")
		author.ID, _ = strconv.Atoi(IDRegex.FindStringSubmatch(author.URL)[1])
		author.Ordering = authorNodes.Length() - i - 1
		author.Name = strings.TrimSpace(s.Find(".author__name").Text())
		author.Title = s.Find(".item__title").Text()

		if srcset, exists := s.Find("img").Attr("srcset"); exists {
			matches := SrcsetRegex.FindStringSubmatch(srcset)
			if len(matches) == 3 {
				author.Images = append(author.Images, entity.Image{Width: 50, Src: matches[1]})
				author.Images = append(author.Images, entity.Image{Width: 100, Src: matches[2]})
			}
		}

		authors = append(authors, author)
	})

	log.Info("Synced authors", "count", authorNodes.Length(), "duration", time.Since(start))

	return authors, nil
}
