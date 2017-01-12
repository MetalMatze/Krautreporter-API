package main

import (
	"net/http"
	"testing"
	"time"

	krautreporter "github.com/metalmatze/krautreporter-api"
	"github.com/stretchr/testify/assert"
)

func TestFetchParseAuthor(t *testing.T) {
	s := &Scraper{client: &http.Client{Timeout: 10 * time.Second}}

	authors := make(map[string]*krautreporter.Author, 4)

	authors["https://krautreporter.de/6-sebastian-esser"] = &krautreporter.Author{
		Name:        "Sebastian Esser",
		URL:         "https://krautreporter.de/6-sebastian-esser",
		Biography:   "Sebastian Esser, Jahrgang 1976, arbeitete als Politik- und Medienredakteur und gründete zwei Journalismus-Startups. Er lebt in Berlin.",
		SocialMedia: "Sebastian Esser\nauf\n<a class=\"link link--camoflaged text-bold twitter--link\" target=\"_blank\" href=\"https://twitter.com/@sebastianesser\">TWITTER</a> | <a class=\"link link--camoflaged text-bold facebook--link\" target=\"_blank\" href=\"https://www.facebook.com/https://www.facebook.com/sebastian.esser\">FACEBOOK</a> | <a class=\"link link--camoflaged text-bold xing--link\" target=\"_blank\" href=\"https://www.xing.com/profile/\">XING</a> | <a class=\"link link--camoflaged text-bold linkedin--link\" target=\"_blank\" href=\"https://www.linkedin.com/in/\">LINKEDIN</a> | <a class=\"link link--camoflaged text-bold homepage--link\" target=\"_blank\" href=\"http://sebastian-esser.de\">HOMEPAGE</a>",
		Images: []*krautreporter.Image{{
			Width: 170,
			Src:   "/system/user/profile_image/6/title_26Voo_0sItW8-BkHh872x9wzbJIJh5WN1SwBroOcBS4.png",
		}, {
			Width: 340,
			Src:   "/system/user/profile_image/6/title_retina_26Voo_0sItW8-BkHh872x9wzbJIJh5WN1SwBroOcBS4.png",
		}},
	}

	authors["https://krautreporter.de/13-tilo-jung"] = &krautreporter.Author{
		Name:        "Tilo Jung",
		URL:         "https://krautreporter.de/13-tilo-jung",
		Biography:   "Tilo Jung",
		SocialMedia: "Tilo Jung\nauf\n<a class=\"link link--camoflaged text-bold twitter--link\" target=\"_blank\" href=\"https://twitter.com/@TiloJung\">TWITTER</a> | <a class=\"link link--camoflaged text-bold facebook--link\" target=\"_blank\" href=\"https://www.facebook.com/tilo.jung\">FACEBOOK</a> | <a class=\"link link--camoflaged text-bold xing--link\" target=\"_blank\" href=\"https://www.xing.com/profile/\">XING</a> | <a class=\"link link--camoflaged text-bold linkedin--link\" target=\"_blank\" href=\"https://www.linkedin.com/in/\">LINKEDIN</a> | <a class=\"link link--camoflaged text-bold homepage--link\" target=\"_blank\" href=\"http://www.jungundnaiv.de\">HOMEPAGE</a>",
		Images: []*krautreporter.Image{{
			Width: 170,
			Src:   "/system/user/profile_image/13/title_KUcC_EJQiDrfoGvsdsSz0SjfNhTEouNaONYGGC1vaMw.png",
		}, {
			Width: 340,
			Src:   "/system/user/profile_image/13/title_retina_KUcC_EJQiDrfoGvsdsSz0SjfNhTEouNaONYGGC1vaMw.png",
		}},
	}

	authors["https://krautreporter.de/175-juliane-wiedemeier"] = &krautreporter.Author{
		Name:        "Juliane Wiedemeier",
		URL:         "https://krautreporter.de/175-juliane-wiedemeier",
		Biography:   "Juliane Wiedemeier ist freie Journalistin in Berlin und Gründungsredakteurin der Prenzlauer Berg Nachrichten.",
		SocialMedia: "Juliane Wiedemeier\nauf\n<a class=\"link link--camoflaged text-bold twitter--link\" target=\"_blank\" href=\"https://twitter.com/\">TWITTER</a> | <a class=\"link link--camoflaged text-bold facebook--link\" target=\"_blank\" href=\"https://www.facebook.com/\">FACEBOOK</a> | <a class=\"link link--camoflaged text-bold xing--link\" target=\"_blank\" href=\"https://www.xing.com/profile/\">XING</a> | <a class=\"link link--camoflaged text-bold linkedin--link\" target=\"_blank\" href=\"https://www.linkedin.com/in/\">LINKEDIN</a> | <a class=\"link link--camoflaged text-bold homepage--link\" target=\"_blank\" href=\"\">HOMEPAGE</a>",
		Images: []*krautreporter.Image{{
			Width: 170,
			Src:   "/system/user/profile_image/175/title_pfd5hoLcRM6j_VXzfbugAudHJHvT1xVaJEceuWjT0DI.png",
		}, {
			Width: 340,
			Src:   "/system/user/profile_image/175/title_retina_pfd5hoLcRM6j_VXzfbugAudHJHvT1xVaJEceuWjT0DI.png",
		}},
	}

	authors["https://krautreporter.de/24012-elisabeth-dietz"] = &krautreporter.Author{
		Name:        "Elisabeth Dietz",
		URL:         "https://krautreporter.de/24012-elisabeth-dietz",
		Biography:   "",
		SocialMedia: "",
		Images: []*krautreporter.Image{{
			Width: 170,
			Src:   "/system/user/profile_image/24012/title_GE1OGrRtJPg2AXjf8eiwFDDpEk8bVW8yip6wl9333zA.jpg",
		}, {
			Width: 340,
			Src:   "/system/user/profile_image/24012/title_retina_GE1OGrRtJPg2AXjf8eiwFDDpEk8bVW8yip6wl9333zA.jpg",
		}},
	}

	for url, expected := range authors {
		sa := ScrapeAuthor{
			Scraper: s,
			Author: &krautreporter.Author{
				URL: url,
			},
		}

		doc, err := sa.Fetch()
		assert.NoError(t, err)

		err = sa.Parse(doc)
		assert.NoError(t, err)
		assert.Equal(t, expected, sa.Author)
	}
}
