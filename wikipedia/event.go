package wikipedia

type Event struct {
	Type string `json:"type"`

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
