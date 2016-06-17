package cli

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter"
	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/icrowley/fake"
	"github.com/urfave/cli"
)

const authorsCount int = 100
const articlesCount int = 1000

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func SeedCommand(kr *krautreporter.Krautreporter) cli.Command {
	return cli.Command{
		Name:  "seed",
		Usage: "Seed the database with some test data",
		Action: func(c *cli.Context) error {
			start := time.Now()

			var authors []entity.Author
			for i := 1; i <= authorsCount; i++ {
				author := entity.Author{
					ID:          i,
					Ordering:    authorsCount - i,
					Name:        fake.FullName(),
					Title:       fake.JobTitle(),
					URL:         "",
					Biography:   fake.Paragraphs(),
					SocialMedia: fake.Paragraphs(),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				authors = append(authors, author)
			}
			kr.CrawlInteractor.SaveAuthors(authors)

			var articles []entity.Article
			for i := 1; i <= articlesCount; i++ {
				article := entity.Article{
					ID:        i,
					Ordering:  articlesCount - i,
					Title:     fake.Title(),
					Headline:  fake.Title(),
					Date:      time.Now(),
					Preview:   false,
					URL:       fmt.Sprintf("%d-%s", i, fake.Title()),
					Excerpt:   fake.Paragraph(),
					Content:   fake.ParagraphsN(r.Intn(15) + 5),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
					Author:    &authors[r.Intn(authorsCount)],
				}

				articles = append(articles, article)
			}
			kr.CrawlInteractor.SaveArticles(articles)
			fmt.Printf("%+v\n", articles[len(articles)-1])

			kr.Log.Info("Database is seeded", "duration", time.Since(start))

			return nil
		},
	}
}
