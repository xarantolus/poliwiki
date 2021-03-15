package wikipedia

import (
	"net/url"
	"path"
	"strconv"
)

type Event struct {
	Type string `json:"type"`

	Meta struct {
		URL string `json:"url"`
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

func (e *Event) DiffURL() (us string, ok bool) {
	pageSlug := path.Base(e.Meta.URL)
	if pageSlug == "" {
		return
	}

	var u = url.URL{
		Scheme: "https",
		Host:   "de.wikipedia.org",
		Path:   "/w/index.php",
	}

	u.Query().Set("title", pageSlug)
	u.Query().Set("diff", strconv.Itoa(e.Revision.New))
	u.Query().Set("oldid", strconv.Itoa(e.Revision.Old))
	u.Query().Set("diffonly", "yes")

	return u.String(), true
}

type Revision struct {
	Old int `json:"old"`
	New int `json:"new"`
}
