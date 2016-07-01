package service

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-kit/kit/log"
)

var IDRegex = regexp.MustCompile(`\/(\d*)`)
var authorSrcsetRegex = regexp.MustCompile(`(.*) 50w, (.*) 100w`)

// SyncAuthor finds all author's meta data from the krautreporter homepage
func SyncAuthor(logger log.Logger) ([]entity.Author, error) {
	start := time.Now()

	doc, err := goquery.NewDocument(mainURL)
	if err != nil {
		return nil, err
	}

	authorNodes := doc.Find("#author-list-tab li a")

	var authors []entity.Author
	authorNodes.Each(func(i int, s *goquery.Selection) {
		author := entity.Author{}
		author.URL, _ = s.Attr("href")
		author.ID, _ = strconv.Atoi(IDRegex.FindStringSubmatch(author.URL)[1])
		author.Ordering = authorNodes.Length() - i - 1
		author.Name = strings.TrimSpace(s.Find(".author__name").Text())
		author.Title = strings.TrimSpace(s.Find(".item__title").Text())

		if srcset, exists := s.Find("img").Attr("srcset"); exists {
			matches := authorSrcsetRegex.FindStringSubmatch(srcset)
			if len(matches) == 3 {
				author.AddImage(entity.Image{Width: 130, Src: matches[1]})
				author.AddImage(entity.Image{Width: 260, Src: matches[2]})
			}
		}

		authors = append(authors, author)
	})

	logger.Log("msg", "Synced authors", "count", authorNodes.Length(), "duration", time.Since(start))

	return authors, nil
}
