package krautreporter

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCrawl_NextRandom(t *testing.T) {
	now := time.Now()
	c := Crawl{Next: now}
	assert.Equal(t, now, c.Next)

	c.NextRandom()
	assert.True(t, time.Now().Before(c.Next))
	fmt.Printf("%v < %v", time.Now(), c.Next)

}
