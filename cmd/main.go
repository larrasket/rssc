package main

import (
	"fmt"
	"net/http"
	"time"

	_ "embed"

	"github.com/larrasket/rssc"
	"github.com/patrickmn/go-cache"
)

//go:embed README.html
var README string

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(README))
	})

	c := cache.New(5*time.Second, 20*time.Second)
	http.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		if req, found := c.Get(q.Encode()); found {
			res := req.(string)
			w.Write([]byte(res))
			return
		}

		prams := rssc.FilterPrams{AuthorRegex: q.Get("authorf"),
			ContentRegex:     q.Get("contentf"),
			TitleRegex:       q.Get("titlef"),
			DescriptionRegex: q.Get("descriptionf"),
			URL:              q.Get("src"),
			DotNet: func(b string) bool {
				if b == "1" {
					return true
				}
				return false
			}(q.Get("net")),
		}

		entries, err := rssc.FilterFeeds(&prams)
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		ftype := q.Get("t")
		feeds, err := rssc.GenerateFeeds(entries, ftype)
		if err != nil {
			w.Write([]byte(err.Error()))
		}

		if ftype == "json" {
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.Header().Set("Content-Type", "application/xml")
		}

		c.Set(q.Encode(), feeds, cache.DefaultExpiration)
		w.Write([]byte(feeds))
	})

	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)

}
