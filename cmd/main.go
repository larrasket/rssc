package main

import (
	"fmt"
	"net/http"

	"github.com/larrasket/rssc"
)

func main() {
	//AuthorRegex, ContentRegex, TitleRegex, DescriptionRegex string
	http.HandleFunc("/rss", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		prams := rssc.FilterPrams{AuthorRegex: q.Get("authorf"),
			ContentRegex:     q.Get("contentf"),
			TitleRegex:       q.Get("titlef"),
			DescriptionRegex: q.Get("descriptionf"),
			URL:              q.Get("source"),
		}
		ftype := q.Get("ftype")
		entries, err := rssc.FilterFeed(&prams)
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
