package commands

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/MetalMatze/Krautreporter-API/domain/entity"
	"github.com/MetalMatze/Krautreporter-API/domain/interactor"
	"github.com/PuerkitoBio/goquery"
	"github.com/codegangsta/cli"
)

func CrawlCommand(authorInteractor interactor.AuthorInteractor) cli.Command {
	return cli.Command{
		Name:  "crawl",
		Usage: "Display an inspiring quote",
		Subcommands: []cli.Command{{
			Name:  "authors",
			Usage: "Crawl all authors from krautreporter.de",
			Action: func(context *cli.Context) {
				doc, err := goquery.NewDocument("https://krautreporter.de")
				if err != nil {
					log.Fatal(err)
				}

				authorNodes := doc.Find("#author-list-tab li a")
				log.Printf("Found %d authors, start parsing and saving", authorNodes.Length())

				authors := []entity.Author{}

				idRegex := regexp.MustCompile(`\/(\d*)--`)
				authorNodes.Each(func(i int, s *goquery.Selection) {
					author := entity.Author{}
					author.URL, _ = s.Attr("href")
					author.ID, _ = strconv.Atoi(idRegex.FindStringSubmatch(author.URL)[1])
					author.Ordering = authorNodes.Length() - i - 1
					author.Name = strings.TrimSpace(s.Find(".author__name").Text())
					author.Title = s.Find(".item__title").Text()

					authors = append(authors, author)
				})

				authorInteractor.SaveAll(authors)
			},
		}},
	}
}
