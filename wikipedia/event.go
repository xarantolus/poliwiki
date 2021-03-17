package wikipedia

import (
	"net/url"
	"path"
	"strconv"
)

type Event struct {
	Type string `json:"type"`

	Meta struct {
		URI string `json:"uri"`
	} `json:"meta"`

	// Article title
	Title string `json:"title"`

	// Why the change was made
	Comment string `json:"comment"`

	Timestamp int `json:"timestamp"`

	// Could be interesting for IP
	User string `json:"user"`

	// Should probably be filtered out
	Bot bool `json:"bot"`

	// Contains the old & new ID
	Revision Revision `json:"revision"`

	// Wiki name, e.g. "dewiki" or "enwiki"
	Wiki string `json:"wiki"`
}

type Revision struct {
	Old int `json:"old"`
	New int `json:"new"`
}

// DiffURL returns the URL for seeing the difference between two versions of an article
func (e *Event) DiffURL() (us string, ok bool) {
	pageSlug := path.Base(e.Meta.URI)
	if pageSlug == "" {
		return
	}

	var u = url.URL{
		Scheme: "https",
		Host:   "de.wikipedia.org",
		Path:   "/w/index.php",
	}

	var q = make(url.Values, 4)

	q.Set("title", pageSlug)
	q.Set("diff", strconv.Itoa(e.Revision.New))
	q.Set("oldid", strconv.Itoa(e.Revision.Old))

	// diffonly removes the article text from the page, no need to load it
	q.Set("diffonly", "yes")

	u.RawQuery = q.Encode()

	return u.String(), true
}
