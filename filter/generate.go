package filter

import (
	"fmt"
	"time"

	generator "github.com/gorilla/feeds"
	"github.com/mmcdole/gofeed"
)

// Geerate a feed string of type t (a string should be json, rss, or atom). If t
// is malformatted or invalid, it returns XML rss feeds
func GenerateFeeds(feeds *gofeed.Feed, t string) (string, error) {

	// TODO unlike mmcdole/gofeed, does not support multiple authors. I think it
	// has other limitation as well, I will consider writing a better XML
	// generator.  Or even better, considering not changing the original feeds format
	if feeds == nil {
		return "", fmt.Errorf(" no feeds provided")
	}

	author := &generator.Author{}
	if len(feeds.Authors) > 0 {
		author.Name = feeds.Authors[0].Name
		author.Email = feeds.Authors[0].Email
	}
	var created time.Time
	if feeds.PublishedParsed != nil {
		created = *feeds.PublishedParsed
	}
	feed := &generator.Feed{
		Title:       feeds.Title,
		Link:        &generator.Link{Href: feeds.Link},
		Description: feeds.Description,
		Author:      author,
		Created:     created,
	}
	gfeeds := []*generator.Item{}
	for _, i := range feeds.Items {
		if len(feeds.Authors) > 0 {
			author.Name = feeds.Authors[0].Name
			author.Email = feeds.Authors[0].Email
		} else {
			author = nil
		}
		gfeeds = append(gfeeds, &generator.Item{
			Title:       i.Title,
			Link:        &generator.Link{Href: i.Link},
			Description: i.Description,
			Author:      author,
			Content:     i.Content,
			Created:     *i.PublishedParsed})
	}
	feed.Items = gfeeds
	if t == "json" {
		return feed.ToJSON()
	}
	if t == "atom" {
		return feed.ToAtom()
	}
	return feed.ToRss()
}
