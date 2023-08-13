package rssc

import (
	"regexp"

	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

type FilterPrams struct {
	URL, AuthorRegex, ContentRegex, TitleRegex, DescriptionRegex string
}

// var tmpf FilterPrams = FilterPrams{TitleRegex: `\bAnalysis\b`}
// FilterFeed("https://hnrss.org/newest", &tmpf)

// Filter `url` based on `prams`, if `prams` is nil (or all filters are empty)
// return the original feeds
func FilterFeed(prams *FilterPrams) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feeds, err := fp.ParseURL(prams.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't read feeds %s ", prams.URL)
	}

	if prams == nil ||
		len(prams.AuthorRegex)+len(prams.ContentRegex)+
			len(prams.TitleRegex)+len(prams.DescriptionRegex) == 0 {
		return feeds, nil
	}

	var authorF, contentF, titleF, DescriptionF *regexp.Regexp

	// validate filters
	authorF, err = validateRegex(prams.AuthorRegex)
	if err != nil {
		return nil, err
	}

	contentF, err = validateRegex(prams.ContentRegex)
	if err != nil {
		return nil, err
	}

	titleF, err = validateRegex(prams.TitleRegex)
	if err != nil {
		return nil, err
	}

	DescriptionF, err = validateRegex(prams.DescriptionRegex)
	if err != nil {
		return nil, err
	}

	// I considered using the new Slices.Delete function instead of building new
	// a new feeds; however "Delete is O(len(s)-j)" (from current go docs), I
	// think we can save some overhead by doing it by hand.
	fitlerd := []*gofeed.Item{}
	for _, i := range feeds.Items {
		// search authors
		if authorF != nil {
			match := false
			for _, a := range i.Authors {
				if authorF.MatchString(a.Name) {
					fitlerd = append(fitlerd, i)
					match = true
					break
				}
			}
			if match {
				continue
			}
		}
		// search content
		if contentF != nil && contentF.MatchString(i.Content) {
			fitlerd = append(fitlerd, i)
			continue
		}
		// search title
		if titleF != nil && titleF.MatchString(i.Title) {
			fitlerd = append(fitlerd, i)
			continue
		}
		// search description
		if DescriptionF != nil && DescriptionF.MatchString(i.Description) {
			fitlerd = append(fitlerd, i)
		}
	}
	feeds.Items = fitlerd
	return feeds, nil

}

func validateRegex(reg string) (*regexp.Regexp, error) {
	if len(reg) == 0 {
		return nil, nil
	}
	r, err := regexp.Compile(reg)
	if err != nil {
		return nil, errors.Wrapf(err, "malformatted regex: %s", r)
	}
	return r, nil
}
