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
      "teaser_html": "\u003cdiv class='article-teaser layout__item'\u003e\u003carticle class='article card card--hover card--regular js-article-teaser regular-article'\u003e\n\u003ca class='card__link link--camoflaged-hover' href='/737-ich-fuhle-manchmal-diese-ohnmachtige-wut-wenn-ich-zu-den-gefangnissen-fahre'\u003e\n\u003cimg\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-path=\"/system/article/teaser_image/737/gOORuiMBRVZyFLi21t8nBzpLG0mrKckVA_38jwW8v1Y.jpg\"\n  ix-params='{\"fit\":\"crop\",\"auto\":\"format\",\"q\":\"30\",\"w\":\"16\",\"h\":\"9\"}'\n  data-sizes=\"auto\"\n  alt=\"„Ich fühle manchmal diese ohnmächtige Wut, wenn ich zu den Gefängnissen fahre“\"\n  class=\"full-width card__img\"\n/\u003e\u003cdiv class='card__body'\u003e\n\u003cdiv class='author meta'\u003e\n\u003cdiv class='author__body'\u003e\n\u003caddress class='author__name text-uppercase'\u003e\n\u003cspan class='visuallyhidden'\u003eVerfasst von\u003c/span\u003e\nPauline Eiferman\u003c/address\u003e\n\u003cspan class='visuallyhidden'\u003eam\u003c/span\u003e\n\u003ctime class='author__date author__date--card' datetime='2015-06-15' title='15. Juni 2015'\u003e15. Juni 2015\u003c/time\u003e\n\u003c/div\u003e\n\u003c/div\u003e\n\u003ch2 class='article--title card__title regular-article--title'\u003e\n„Ich fühle manchmal diese ohnmächtige Wut, wenn ich zu den Gefängnissen fahre“\u003c/h2\u003e\n\u003c/div\u003e\n\u003c/a\u003e\n\u003c/article\u003e\u003c/div\u003e\n"
    },
    {
      "teaser_html": "\u003cdiv class='article-teaser layout__item'\u003e\u003carticle class='article card card--hover card--regular js-article-teaser regular-article'\u003e\n\u003ca class='card__link link--camoflaged-hover' href='/316-gesucht-das-nachste-aids-virus'\u003e\n\u003cimg\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-path=\"/system/article/teaser_image/316/NsO27uKXZ_tQSbNaq3Y_BL3iw5MkvgVn-wGqFO35TJ0.jpg\"\n  ix-params='{\"fit\":\"crop\",\"auto\":\"format\",\"q\":\"30\",\"w\":\"16\",\"h\":\"9\"}'\n  data-sizes=\"auto\"\n  alt=\"Gesucht: das nächste Aids-Virus \"\n  class=\"full-width card__img\"\n/\u003e\u003cdiv class='card__body'\u003e\n\u003cdiv class='author meta'\u003e\n\u003cdiv class='author__body'\u003e\n\u003caddress class='author__name text-uppercase'\u003e\n\u003cspan class='visuallyhidden'\u003eVerfasst von\u003c/span\u003e\nTom Clynes\u003c/address\u003e\n\u003cspan class='visuallyhidden'\u003eam\u003c/span\u003e\n\u003ctime class='author__date author__date--card' datetime='2015-02-16' title='16. Februar 2015'\u003e16. Februar 2015\u003c/time\u003e\n\u003c/div\u003e\n\u003c/div\u003e\n\u003ch2 class='article--title card__title regular-article--title'\u003e\nGesucht: das nächste Aids-Virus\u003c/h2\u003e\n\u003c/div\u003e\n\u003c/a\u003e\n\u003c/article\u003e\u003c/div\u003e\n"
    },
    {
      "teaser_html": "\u003cdiv class='article-teaser layout__item'\u003e\u003carticle class='article card card--hover card--regular js-article-teaser regular-article'\u003e\n\u003ca class='card__link link--camoflaged-hover' href='/635-ja-wir-haben-abgeklebt'\u003e\n\u003cimg\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-path=\"/system/article/teaser_image/635/LaVKNOaomzorM6oS8mn1GhPGQXOhrD2ok45em6V34QE.jpg\"\n  ix-params='{\"fit\":\"crop\",\"auto\":\"format\",\"q\":\"30\",\"w\":\"16\",\"h\":\"9\"}'\n  data-sizes=\"auto\"\n  alt=\"Ja, wir haben abgeklebt!  \"\n  class=\"full-width card__img\"\n/\u003e\u003cdiv class='card__body'\u003e\n\u003cdiv class='author meta'\u003e\n\u003cdiv class='author__body'\u003e\n\u003caddress class='author__name text-uppercase'\u003e\n\u003cspan class='visuallyhidden'\u003eVerfasst von\u003c/span\u003e\nJuliane Schiemenz\u003c/address\u003e\n\u003cspan class='visuallyhidden'\u003eam\u003c/span\u003e\n\u003ctime class='author__date author__date--card' datetime='2015-04-24' title='24. April 2015'\u003e24. April 2015\u003c/time\u003e\n\u003c/div\u003e\n\u003c/div\u003e\n\u003ch2 class='article--title card__title regular-article--title'\u003e\nJa, wir haben abgeklebt!\u003c/h2\u003e\n\u003c/div\u003e\n\u003c/a\u003e\n\u003c/article\u003e\u003c/div\u003e\n"
    },
    {
      "teaser_html": "\u003cdiv class='article-teaser layout__item'\u003e\u003carticle class='article card card--hover card--regular js-article-teaser regular-article'\u003e\n\u003ca class='card__link link--camoflaged-hover' href='/1702-was-ist-krautreporter'\u003e\n\u003cimg\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-path=\"/system/article/teaser_image/1702/iBs2eREsrYKRRtBKgW7Lud1DSsmewYWCUbJK9SxWrmA.jpg\"\n  ix-params='{\"fit\":\"crop\",\"auto\":\"format\",\"q\":\"30\",\"w\":\"16\",\"h\":\"9\"}'\n  data-sizes=\"auto\"\n  alt=\"Was ist Krautreporter?\"\n  class=\"full-width card__img\"\n/\u003e\u003cdiv class='card__body'\u003e\n\u003cdiv class='author meta'\u003e\n\u003cdiv class='author__body'\u003e\n\u003caddress class='author__name text-uppercase'\u003e\n\u003cspan class='visuallyhidden'\u003eVerfasst von\u003c/span\u003e\nKrautreporter\u003c/address\u003e\n\u003cspan class='visuallyhidden'\u003eam\u003c/span\u003e\n\u003ctime class='author__date author__date--card' datetime='2017-01-06' title='06. Januar 2017'\u003e06. Januar 2017\u003c/time\u003e\n\u003c/div\u003e\n\u003c/div\u003e\n\u003ch2 class='article--title card__title regular-article--title'\u003e\n30 Tage kostenlos testen!\u003c/h2\u003e\n\u003c/div\u003e\n\u003c/a\u003e\n\u003c/article\u003e\u003c/div\u003e\n"
    },
    {
      "teaser_html": "\u003cdiv class='article-teaser layout__item'\u003e\u003carticle class='article card card--hover card--regular js-article-teaser regular-article'\u003e\n\u003ca class='card__link link--camoflaged-hover' href='/796-wie-aus-einem-volksbefreiungsarmee-oberst-die-wahrscheinlich-einflussreichste-frau-chinas-wurde'\u003e\n\u003cimg\n  src=\"data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==\"\n  ix-path=\"/system/article/teaser_image/796/y3e_KCNAmHolgZ0fgiSyj9lZOtfRj8H0VO6OXBTx2OA.jpg\"\n  ix-params='{\"fit\":\"crop\",\"auto\":\"format\",\"q\":\"30\",\"w\":\"16\",\"h\":\"9\"}'\n  data-sizes=\"auto\"\n  alt=\"Wie aus einem Volksbefreiungsarmee-Oberst die wahrscheinlich einflussreichste Frau Chinas wurde\"\n  class=\"full-width card__img\"\n/\u003e\u003cdiv class='card__body'\u003e\n\u003cdiv class='author meta'\u003e\n\u003cdiv class='author__body'\u003e\n\u003caddress class='author__name text-uppercase'\u003e\n\u003cspan class='visuallyhidden'\u003eVerfasst von\u003c/span\u003e\nMaximilian Kalkhof\u003c/address\u003e\n\u003cspan class='visuallyhidden'\u003eam\u003c/span\u003e\n\u003ctime class='author__date author__date--card' datetime='2015-07-06' title='06. Juli 2015'\u003e06. Juli 2015\u003c/time\u003e\n\u003c/div\u003e\n\u003c/div\u003e\n\u003ch2 class='article--title card__title regular-article--title'\u003e\nWie aus einem Oberst die wahrscheinlich einflussreichste Frau Chinas wurde\u003c/h2\u003e\n\u003c/div\u003e\n\u003c/a\u003e\n\u003c/article\u003e\u003c/div\u003e\n"
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
