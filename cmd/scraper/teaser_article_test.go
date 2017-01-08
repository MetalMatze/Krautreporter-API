package main

import (
	"encoding/json"
	"testing"
	"time"

	krautreporter "github.com/metalmatze/krautreporter-api"
	"github.com/stretchr/testify/assert"
)

var (
	data = `{
  "articles": [
    {
      "teaser_html": "<div class='article-teaser layout__item'><article class='article card card--hover card--regular js-article-teaser regular-article'>\n<a class='card__link link--camoflaged-hover' href='/737-ich-fuhle-manchmal-diese-ohnmachtige-wut-wenn-ich-zu-den-gefangnissen-fahre'>\n<img\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-src=\"https://krautreporter.imgix.net/system/article/teaser_image/737/gOORuiMBRVZyFLi21t8nBzpLG0mrKckVA_38jwW8v1Y.jpg?ixlib=rails-2.1.2&fit=crop&auto=format&q=30&w=16&h=9\"\n  data-sizes=\"auto\"\n  alt=\"„Ich fühle manchmal diese ohnmächtige Wut, wenn ich zu den Gefängnissen fahre“\"\n  class=\"full-width card__img\"\n/><div class='card__body'>\n<div class='author meta'>\n<div class='author__body'>\n<address class='author__name text-uppercase'>\n<span class='visuallyhidden'>Verfasst von</span>\nPauline Eiferman</address>\n<span class='visuallyhidden'>am</span>\n<time class='author__date author__date--card' datetime='2015-06-15' title='15. Juni 2015'>15. Juni 2015</time>\n</div>\n</div>\n<h2 class='article--title card__title regular-article--title'>\n„Ich fühle manchmal diese ohnmächtige Wut, wenn ich zu den Gefängnissen fahre“</h2>\n</div>\n</a>\n</article></div>\n"
    },
    {
      "teaser_html": "<div class='article-teaser layout__item'><article class='article card card--hover card--regular js-article-teaser regular-article'>\n<a class='card__link link--camoflaged-hover' href='/316-gesucht-das-nachste-aids-virus'>\n<img\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-src=\"https://krautreporter.imgix.net/system/article/teaser_image/316/NsO27uKXZ_tQSbNaq3Y_BL3iw5MkvgVn-wGqFO35TJ0.jpg?ixlib=rails-2.1.2&fit=crop&auto=format&q=30&w=16&h=9\"\n  data-sizes=\"auto\"\n  alt=\"Gesucht: das nächste Aids-Virus \"\n  class=\"full-width card__img\"\n/><div class='card__body'>\n<div class='author meta'>\n<div class='author__body'>\n<address class='author__name text-uppercase'>\n<span class='visuallyhidden'>Verfasst von</span>\nTom Clynes</address>\n<span class='visuallyhidden'>am</span>\n<time class='author__date author__date--card' datetime='2015-02-16' title='16. Februar 2015'>16. Februar 2015</time>\n</div>\n</div>\n<h2 class='article--title card__title regular-article--title'>\nGesucht: das nächste Aids-Virus</h2>\n</div>\n</a>\n</article></div>\n"
    },
    {
      "teaser_html": "<div class='article-teaser layout__item'><article class='article card card--hover card--regular js-article-teaser regular-article'>\n<a class='card__link link--camoflaged-hover' href='/635-ja-wir-haben-abgeklebt'>\n<img\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-src=\"https://krautreporter.imgix.net/system/article/teaser_image/635/LaVKNOaomzorM6oS8mn1GhPGQXOhrD2ok45em6V34QE.jpg?ixlib=rails-2.1.2&fit=crop&auto=format&q=30&w=16&h=9\"\n  data-sizes=\"auto\"\n  alt=\"Ja, wir haben abgeklebt!  \"\n  class=\"full-width card__img\"\n/><div class='card__body'>\n<div class='author meta'>\n<div class='author__body'>\n<address class='author__name text-uppercase'>\n<span class='visuallyhidden'>Verfasst von</span>\nJuliane Schiemenz</address>\n<span class='visuallyhidden'>am</span>\n<time class='author__date author__date--card' datetime='2015-04-24' title='24. April 2015'>24. April 2015</time>\n</div>\n</div>\n<h2 class='article--title card__title regular-article--title'>\nJa, wir haben abgeklebt!</h2>\n</div>\n</a>\n</article></div>\n"
    },
    {
      "teaser_html": "<div class='article-teaser layout__item'><article class='article card card--hover card--regular js-article-teaser regular-article'>\n<a class='card__link link--camoflaged-hover' href='/1702-was-ist-krautreporter'>\n<img\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-src=\"https://krautreporter.imgix.net/system/article/teaser_image/1702/iBs2eREsrYKRRtBKgW7Lud1DSsmewYWCUbJK9SxWrmA.jpg?ixlib=rails-2.1.2&fit=crop&auto=format&q=30&w=16&h=9\"\n  data-sizes=\"auto\"\n  alt=\"Was ist Krautreporter?\"\n  class=\"full-width card__img\"\n/><div class='card__body'>\n<div class='author meta'>\n<div class='author__body'>\n<address class='author__name text-uppercase'>\n<span class='visuallyhidden'>Verfasst von</span>\nKrautreporter</address>\n<span class='visuallyhidden'>am</span>\n<time class='author__date author__date--card' datetime='2017-01-06' title='06. Januar 2017'>06. Januar 2017</time>\n</div>\n</div>\n<h2 class='article--title card__title regular-article--title'>\n30 Tage kostenlos testen!</h2>\n</div>\n</a>\n</article></div>\n"
    },
    {
      "teaser_html": "<div class='article-teaser layout__item'><article class='article card card--hover card--regular js-article-teaser regular-article'>\n<a class='card__link link--camoflaged-hover' href='/796-wie-aus-einem-volksbefreiungsarmee-oberst-die-wahrscheinlich-einflussreichste-frau-chinas-wurde'>\n<img\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-src=\"https://krautreporter.imgix.net/system/article/teaser_image/796/y3e_KCNAmHolgZ0fgiSyj9lZOtfRj8H0VO6OXBTx2OA.jpg?ixlib=rails-2.1.2&fit=crop&auto=format&q=30&w=16&h=9\"\n  data-sizes=\"auto\"\n  alt=\"Wie aus einem Volksbefreiungsarmee-Oberst die wahrscheinlich einflussreichste Frau Chinas wurde\"\n  class=\"full-width card__img\"\n/><div class='card__body'>\n<div class='author meta'>\n<div class='author__body'>\n<address class='author__name text-uppercase'>\n<span class='visuallyhidden'>Verfasst von</span>\nMaximilian Kalkhof</address>\n<span class='visuallyhidden'>am</span>\n<time class='author__date author__date--card' datetime='2015-07-06' title='06. Juli 2015'>06. Juli 2015</time>\n</div>\n</div>\n<h2 class='article--title card__title regular-article--title'>\nWie aus einem Oberst die wahrscheinlich einflussreichste Frau Chinas wurde</h2>\n</div>\n</a>\n</article></div>\n"
    }
  ]
}`

	article1 = &krautreporter.Article{
		ID:       737,
		Ordering: 0,
		Title:    "„Ich fühle manchmal diese ohnmächtige Wut, wenn ich zu den Gefängnissen fahre“",
		Date:     time.Date(2015, 6, 15, 0, 0, 0, 0, time.UTC),
		Preview:  true,
		URL:      "/737-ich-fuhle-manchmal-diese-ohnmachtige-wut-wenn-ich-zu-den-gefangnissen-fahre",
		Images: []krautreporter.Image{
			{
				Width: 1600,
				Src:   "https://krautreporter.imgix.net/system/article/teaser_image/737/gOORuiMBRVZyFLi21t8nBzpLG0mrKckVA_38jwW8v1Y.jpg?w=1600",
			},
		},
	}
	article2 = &krautreporter.Article{
		ID:       316,
		Ordering: 0,
		Title:    "Gesucht: das nächste Aids-Virus",
		Date:     time.Date(2015, 2, 16, 0, 0, 0, 0, time.UTC),
		Preview:  true,
		URL:      "/316-gesucht-das-nachste-aids-virus",
		Images: []krautreporter.Image{
			{
				Width: 1600,
				Src:   "https://krautreporter.imgix.net/system/article/teaser_image/316/NsO27uKXZ_tQSbNaq3Y_BL3iw5MkvgVn-wGqFO35TJ0.jpg?w=1600",
			},
		},
	}
	article3 = &krautreporter.Article{
		ID:       635,
		Ordering: 0,
		Title:    "Ja, wir haben abgeklebt!",
		Date:     time.Date(2015, 4, 24, 0, 0, 0, 0, time.UTC),
		Preview:  true,
		URL:      "/635-ja-wir-haben-abgeklebt",
		Images: []krautreporter.Image{
			{
				Width: 1600,
				Src:   "https://krautreporter.imgix.net/system/article/teaser_image/635/LaVKNOaomzorM6oS8mn1GhPGQXOhrD2ok45em6V34QE.jpg?w=1600",
			},
		},
	}
	article4 = &krautreporter.Article{
		ID:       1702,
		Ordering: 0,
		Title:    "30 Tage kostenlos testen!",
		Date:     time.Date(2017, 1, 6, 0, 0, 0, 0, time.UTC),
		Preview:  true,
		URL:      "/1702-was-ist-krautreporter",
		Images: []krautreporter.Image{
			{
				Width: 1600,
				Src:   "https://krautreporter.imgix.net/system/article/teaser_image/1702/iBs2eREsrYKRRtBKgW7Lud1DSsmewYWCUbJK9SxWrmA.jpg?w=1600",
			},
		},
	}
	article5 = &krautreporter.Article{
		ID:       796,
		Ordering: 0,
		Title:    "Wie aus einem Oberst die wahrscheinlich einflussreichste Frau Chinas wurde",
		Date:     time.Date(2015, 7, 6, 0, 0, 0, 0, time.UTC),
		Preview:  true,
		URL:      "/796-wie-aus-einem-volksbefreiungsarmee-oberst-die-wahrscheinlich-einflussreichste-frau-chinas-wurde",
		Images: []krautreporter.Image{
			{
				Width: 1600,
				Src:   "https://krautreporter.imgix.net/system/article/teaser_image/796/y3e_KCNAmHolgZ0fgiSyj9lZOtfRj8H0VO6OXBTx2OA.jpg?w=1600",
			},
		},
	}
)

func TestTeaserArticle_Parse(t *testing.T) {
	var resp TeaserArticleResponse
	json.Unmarshal([]byte(data), &resp)

	teaser1 := TeaserArticle{HTML: resp.Articles[0].HTML}
	a1, err := teaser1.Parse()
	assert.NoError(t, err)
	assert.Equal(t, article1, a1)

	teaser2 := TeaserArticle{HTML: resp.Articles[1].HTML}
	a2, err := teaser2.Parse()
	assert.NoError(t, err)
	assert.Equal(t, article2, a2)

	teaser3 := TeaserArticle{HTML: resp.Articles[2].HTML}
	a3, err := teaser3.Parse()
	assert.NoError(t, err)
	assert.Equal(t, article3, a3)

	teaser4 := TeaserArticle{HTML: resp.Articles[3].HTML}
	a4, err := teaser4.Parse()
	assert.NoError(t, err)
	assert.Equal(t, article4, a4)

	teaser5 := TeaserArticle{HTML: resp.Articles[4].HTML}
	a5, err := teaser5.Parse()
	assert.NoError(t, err)
	assert.Equal(t, article5, a5)
}
