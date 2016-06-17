package cli

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/MetalMatze/Krautreporter-API/krautreporter/entity"
	"github.com/gollection/gollection"
	"github.com/icrowley/fake"
	"github.com/urfave/cli"
)

const authorsCount int = 100
const articlesCount int = 1000

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func SeedCommand(g *gollection.Gollection) cli.Command {
	return cli.Command{
		Name:  "seed",
		Usage: "Seed the database with some test data",
		Action: func(c *cli.Context) error {
			start := time.Now()

			tx := g.DB.Begin()
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
				tx.Create(&author)
			}

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
				tx.Create(&article)
			}

			tx.Commit()

			g.Log.Info("Database is seeded", "duration", time.Since(start))

			return nil
		},
	}
}
