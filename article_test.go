package krautreporter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestArticle_AddImage(t *testing.T) {
	a := Article{}
	assert.Nil(t, a.Images)
	assert.Len(t, a.Images, 0)

	i1 := &Image{Width: 130, Src: "/foobar.jpg"}
	a.AddImage(i1)
	assert.Len(t, a.Images, 1)
	assert.Equal(t, i1, a.Images[0])

	i2 := &Image{Width: 260, Src: "/foobaz.jpg"}
	a.AddImage(i2)
	assert.Len(t, a.Images, 2)
	assert.Equal(t, i2, a.Images[1])

	i3 := &Image{Width: 130, Src: "/baz.jpg"}
	a.AddImage(i3)
	assert.Len(t, a.Images, 2)
	assert.Equal(t, i3, a.Images[0])
}

func TestArticle_NextCrawl(t *testing.T) {
	now := time.Now()

	a := Article{}
	assert.Nil(t, a.Crawl)

	a.NextCrawl(&Crawl{Next: now})
	assert.Equal(t, 0, a.Crawl.ID)
	assert.Equal(t, now, a.Crawl.Next)

	a.NextCrawl(&Crawl{ID: 1234, Next: now})
	assert.Equal(t, 1234, a.Crawl.ID)
	assert.Equal(t, now, a.Crawl.Next)

	now = time.Now()
	a.NextCrawl(&Crawl{Next: now})
	assert.Equal(t, 1234, a.Crawl.ID)
	assert.Equal(t, now, a.Crawl.Next)

	a.NextCrawl(&Crawl{ID: 12345, Next: now})
	assert.Equal(t, 12345, a.Crawl.ID)
	assert.Equal(t, now, a.Crawl.Next)
}
