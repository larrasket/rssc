package main

import (
	"fmt"
	"net/http"

	_ "embed"

	"github.com/larrasket/rssc"
)

//go:embed README.html
var README string

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(README))
	})

	http.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
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
		ftype := q.Get("t")
		entries, err := rssc.FilterFeeds(&prams)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		feeds, err := rssc.GenerateFeeds(entries, ftype)
		if err != nil {
			w.Write([]byte(err.Error()))
		}
		if ftype == "json" {
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.Header().Set("Content-Type", "application/xml")
		}
		w.Write([]byte(feeds))
	})

	fmt.Println("Listening on :8080...")
	http.ListenAndServe(":8080", nil)

}
