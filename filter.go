package rssc

import (
	"regexp"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/mmcdole/gofeed"
	"github.com/pkg/errors"
)

type FilterPrams struct {
	URL,
	AuthorRegex,
	ContentRegex,
	TitleRegex,
	DescriptionRegex string
	DotNet bool // whether to use .NET Regex Engineg
}

// var tmpf FilterPrams = FilterPrams{TitleRegex: `\bAnalysis\b`}
// FilterFeed("https://hnrss.org/newest", &tmpf)

// Filter `url` based on `prams` if `prams` is nil (or all filters are empty)
// return the original feeds
func FilterFeeds(prams *FilterPrams) (*gofeed.Feed, error) {
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
	if prams.DotNet {
		return withDotNETRegex(feeds, prams)
	}
	return withRE2(feeds, prams)
}

// Implements filtering for Google's RE2 standards
func withRE2(feeds *gofeed.Feed, prams *FilterPrams) (*gofeed.Feed, error) {
	var authorF, contentF, titleF, DescriptionF *regexp.Regexp

	// validate filters
	authorF, err := validateRE2(prams.AuthorRegex)
	if err != nil {
		return nil, err
	}

	contentF, err = validateRE2(prams.ContentRegex)
	if err != nil {
		return nil, err
	}

	titleF, err = validateRE2(prams.TitleRegex)
	if err != nil {
		return nil, err
	}

	DescriptionF, err = validateRE2(prams.DescriptionRegex)
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

// Implements filtering for Microsoft's .NET regex engine
// TODO refactor
func withDotNETRegex(feeds *gofeed.Feed, prams *FilterPrams) (*gofeed.Feed, error) {
	var authorF, contentF, titleF, DescriptionF *regexp2.Regexp

	// validate filters
	authorF, err := validateNET(prams.AuthorRegex)
	if err != nil {
		return nil, err
	}

	contentF, err = validateNET(prams.ContentRegex)
	if err != nil {
		return nil, err
	}

	titleF, err = validateNET(prams.TitleRegex)
	if err != nil {
		return nil, err
	}

	DescriptionF, err = validateNET(prams.DescriptionRegex)
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
				p, err := authorF.MatchString(a.Name)
				if err != nil {
					return nil, errors.Wrapf(err, "timeout on expression %s",
						prams.AuthorRegex)
				}
				if p {
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
		if contentF != nil {
			p, err := contentF.MatchString(i.Content)
			if err != nil {
				return nil, errors.Wrapf(err, "timeout on expression %s",
					prams.ContentRegex)
			}
			if p {
				fitlerd = append(fitlerd, i)
				continue
			}
		}
		// search title
		if titleF != nil {
			p, err := titleF.MatchString(i.Title)
			if err != nil {
				return nil, errors.Wrapf(err, "timeout on expression %s",
					prams.TitleRegex)
			}
			if p {
				fitlerd = append(fitlerd, i)
				continue
			}
		}
		// search description
		if DescriptionF != nil {
			p, err := DescriptionF.MatchString(i.Description)
			if err != nil {
				return nil, errors.Wrapf(err, "timeout on expression %s",
					prams.DescriptionRegex)
			}
			if p {
				fitlerd = append(fitlerd, i)
			}
		}
	}
	feeds.Items = fitlerd
	return feeds, nil
}

func validateRE2(reg string) (*regexp.Regexp, error) {
	if len(reg) == 0 {
		return nil, nil
	}
	r, err := regexp.Compile(reg)
	if err != nil {
		return nil, errors.Wrapf(err, "malformatted regex: %s", r)
	}
	return r, nil
}

func validateNET(reg string) (*regexp2.Regexp, error) {
	if len(reg) == 0 {
		return nil, nil
	}
	r, err := regexp2.Compile(reg, regexp2.None)
	if err != nil {
		return nil, errors.Wrapf(err, "malformatted regex: %s", r)
	}
	return r, nil
}

func init() {
	regexp2.SetTimeoutCheckPeriod(time.Second * 3)
}
