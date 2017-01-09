package main

import (
	"net/http"
	"testing"
	"time"

	krautreporter "github.com/metalmatze/krautreporter-api"
	"github.com/stretchr/testify/assert"
)

func TestFetchParseArticle(t *testing.T) {
	s := &Scraper{client: &http.Client{Timeout: 10 * time.Second}}

	articles := make(map[string]*krautreporter.Article, 5)

	articles["https://krautreporter.de/1705-alt-right-escape"] = &krautreporter.Article{
		Headline: "Alt Right Escape",
		URL:      "https://krautreporter.de/1705-alt-right-escape",
		Excerpt:  "Eine neue rechtsextreme Bewegung nimmt für sich in Anspruch, Donald Trump zum US-Präsidenten gemacht zu haben – durch neue Formen der Propaganda in Sozialen Netzwerken. Zur “Alternativen Rechten” zählen Neonazis, Rassisten und Ultrakonservative ebenso wie Internet-Trolle, die Spaß am Zerstören haben. Droht im Bundestagswahlkampf etwas Ähnliches in Deutschland? Ich habe mit einem Aussteiger gesprochen.",
		Content:  "content",
		AuthorID: 24012,
		Author: &krautreporter.Author{
			ID:   24012,
			Name: "Elisabeth Dietz",
			URL:  "/24012-elisabeth-dietz",
		},
	}

	articles["https://krautreporter.de/1552-hitzefrei"] = &krautreporter.Article{
		Headline: "Hitzefrei!",
		URL:      "https://krautreporter.de/1552-hitzefrei",
		Excerpt:  "Wir rennen normalerweise nicht jedem Kaffeetrend hinterher, aber der schlichte Cold Brew schmeckt so gut, dass Theresa Bäuerlein ihn dringend zum Selbermachen empfiehlt. Zahl der Geräte, die sie für die Zubereitung braucht: null.",
		Content:  "content",
		AuthorID: 1,
		Author: &krautreporter.Author{
			ID:   1,
			Name: "Theresa Bäuerlein",
			URL:  "/1-theresa-bauerlein",
		},
	}

	articles["https://krautreporter.de/1186-hey-fuckers"] = &krautreporter.Article{
		Headline: "Hey, Fuckers!",
		URL:      "https://krautreporter.de/1186-hey-fuckers",
		Excerpt:  "Ein amerikanisches Comedy-Duo bricht mit gleich zwei Vorurteilen: dass Frauen nicht lustig sind und dass es ein Problem ist, wenn sie gerne und viel Sex haben. Ihr Podcast hat einen riesigen Erfolg. Denn er ist wie die HBO-Serie „Girls“ - nur besser.",
		Content:  "content",
		AuthorID: 1,
		Author: &krautreporter.Author{
			ID:   1,
			Name: "Theresa Bäuerlein",
			URL:  "/1-theresa-bauerlein",
		},
	}

	articles["https://krautreporter.de/543-tatort-moskau"] = &krautreporter.Article{
		Headline: "Tatort Moskau",
		URL:      "https://krautreporter.de/543-tatort-moskau",
		Excerpt:  "Der Mord an Boris Nemzow ist, so überraschend es klingen mag, der erste Mord an einem wichtigen Politiker in Russland seit mehr als einem Jahrzehnt. Die Ära Putin ist vielmehr gekennzeichnet durch Attacken auf Menschenrechtler und Journalisten.",
		Content:  "content",
		AuthorID: 17671,
		Author: &krautreporter.Author{
			ID:   17671,
			Name: "Ekaterina Anokhina",
			URL:  "/17671-ekaterina-anokhina",
		},
	}

	articles["https://krautreporter.de/23-ruckkehr-zur-abschreckung"] = &krautreporter.Article{
		Headline: "Rückkehr zur Abschreckung",
		URL:      "https://krautreporter.de/23-ruckkehr-zur-abschreckung",
		Excerpt:  "Innerhalb von neun Minuten einsatzbereit in der Luft: Mit deutschen Kampfflugzeugen auf einem estnischen Flughafen setzt die NATO nach 20 Jahren Entspannung wieder ein Zeichen nach Osten.",
		Content:  "content",
		AuthorID: 25,
		Author: &krautreporter.Author{
			ID:   25,
			Name: "Thomas Wiegold",
			URL:  "/25-thomas-wiegold",
		},
	}

	for url, expected := range articles {
		sa := ScrapeArticle{
			Scraper: s,
			Article: &krautreporter.Article{URL: url},
		}

		doc, err := sa.Fetch()
		assert.NoError(t, err)

		err = sa.Parse(doc)
		assert.NoError(t, err)

		assert.NotZero(t, sa.Article.Content)
		sa.Article.Content = "content"
		assert.Equal(t, expected, sa.Article)
	}
}
